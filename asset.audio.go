package inkwell

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// audioSampleRate is the playback rate of the shared audio context. WAV,
// OGG and MP3 streams are resampled to this rate on decode.
const audioSampleRate = 48000

// AudioPlayer plays music and one-shot sound effects through Ebiten's
// audio subsystem. Music is streamed and looped infinitely; sound effects
// are decoded once, cached as raw PCM, and replayed via
// NewPlayerFromBytes for low-latency triggering.
//
// AudioPlayer is fail-soft by design: missing asset files, decode errors
// and a missing audio device all degrade to no-ops with a debug log
// message — the game keeps running. The action constructors
// (PlayMusic/PlaySound/StopMusic) call through here unchanged.
type AudioPlayer struct {
	ctx *audio.Context // lazy — created on first real play

	music       *Asset
	musicPlayer *audio.Player

	sfxCache  map[*Asset][]byte // resampled PCM, one per asset
	sfxActive []*audio.Player   // in-flight one-shots; pruned each PlaySound

	failed map[*Asset]bool // assets that errored once — don't retry every frame
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		sfxCache: make(map[*Asset][]byte),
		failed:   make(map[*Asset]bool),
	}
}

func (a *AudioPlayer) ensureContext() *audio.Context {
	if a.ctx != nil {
		return a.ctx
	}
	// audio.NewContext panics if called twice, so check for an existing
	// singleton first — this lets a test or a host app pre-create one.
	if c := audio.CurrentContext(); c != nil {
		a.ctx = c
		return c
	}
	a.ctx = audio.NewContext(audioSampleRate)
	return a.ctx
}

// playable reports whether an asset is something the player can start:
// present, audio, with a path, and not already on the failure list.
func (a *AudioPlayer) playable(asset *Asset) bool {
	return asset != nil && !a.failed[asset] && asset.Kind == AssetAudio && asset.Path != ""
}

func assetName(asset *Asset) string {
	if asset == nil {
		return "<nil>"
	}
	return asset.Name
}

func (a *AudioPlayer) PlayMusic(asset *Asset) {
	if a.music == asset && a.musicPlayer != nil && a.musicPlayer.IsPlaying() {
		return
	}
	a.stopMusic()
	a.music = asset

	if !a.playable(asset) {
		logf("audio.PlayMusic %s: no asset", assetName(asset))
		return
	}
	ctx := a.ensureContext()
	stream, length, err := decodeAudioStream(asset.Path, audioSampleRate)
	if err != nil {
		logf("audio.PlayMusic %q: %v", asset.Name, err)
		a.failed[asset] = true
		return
	}
	loop := audio.NewInfiniteLoop(stream, length)
	p, err := audio.NewPlayer(ctx, loop)
	if err != nil {
		logf("audio.PlayMusic %q: player: %v", asset.Name, err)
		a.failed[asset] = true
		return
	}
	a.musicPlayer = p
	p.Play()
}

func (a *AudioPlayer) StopMusic() {
	if a.music == nil && a.musicPlayer == nil {
		return
	}
	logf("audio.StopMusic (was %s)", assetName(a.music))
	a.stopMusic()
	a.music = nil
}

func (a *AudioPlayer) stopMusic() {
	if a.musicPlayer == nil {
		return
	}
	_ = a.musicPlayer.Close()
	a.musicPlayer = nil
}

func (a *AudioPlayer) PlaySound(asset *Asset) {
	if !a.playable(asset) {
		logf("audio.PlaySound %s: no asset", assetName(asset))
		return
	}
	ctx := a.ensureContext()

	pcm, ok := a.sfxCache[asset]
	if !ok {
		data, err := decodeAudioBytes(asset.Path, audioSampleRate)
		if err != nil {
			logf("audio.PlaySound %q: %v", asset.Name, err)
			a.failed[asset] = true
			return
		}
		a.sfxCache[asset] = data
		pcm = data
	}
	p := audio.NewPlayerFromBytes(ctx, pcm)
	p.Play()
	a.sfxActive = append(a.sfxActive, p)
	a.pruneSfx()
}

func (a *AudioPlayer) pruneSfx() {
	n := 0
	for _, p := range a.sfxActive {
		if p.IsPlaying() {
			a.sfxActive[n] = p
			n++
		} else {
			_ = p.Close()
		}
	}
	for i := n; i < len(a.sfxActive); i++ {
		a.sfxActive[i] = nil
	}
	a.sfxActive = a.sfxActive[:n]
}

// audioStream is the common interface satisfied by wav/vorbis/mp3 Stream
// types — io.ReadSeeker plus a known total length for loop bookkeeping.
type audioStream interface {
	io.ReadSeeker
	Length() int64
}

func decodeAudioStream(path string, sampleRate int) (audioStream, int64, error) {
	f, err := openAsset(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, 0, err
	}
	r := bytes.NewReader(data)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".wav":
		s, err := wav.DecodeWithSampleRate(sampleRate, r)
		if err != nil {
			return nil, 0, err
		}
		return s, s.Length(), nil
	case ".ogg":
		s, err := vorbis.DecodeWithSampleRate(sampleRate, r)
		if err != nil {
			return nil, 0, err
		}
		return s, s.Length(), nil
	case ".mp3":
		s, err := mp3.DecodeWithSampleRate(sampleRate, r)
		if err != nil {
			return nil, 0, err
		}
		return s, s.Length(), nil
	default:
		return nil, 0, fmt.Errorf("unsupported audio extension %q", filepath.Ext(path))
	}
}

// decodeAudioBytes decodes a file to raw int16 stereo PCM at the given
// sample rate. Used for SFX so NewPlayerFromBytes can replay them
// without re-decoding on every trigger.
func decodeAudioBytes(path string, sampleRate int) ([]byte, error) {
	s, _, err := decodeAudioStream(path, sampleRate)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(s)
}
