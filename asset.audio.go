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
	am  *AssetManager
	ctx *audio.Context // lazy — created on first real play

	musicName   string
	musicPlayer *audio.Player

	sfxCache  map[string][]byte // resampled PCM, keyed by Asset.Name
	sfxActive []*audio.Player   // in-flight one-shots; pruned each PlaySound

	failed map[string]bool // assets that errored once — don't retry every frame
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		sfxCache: make(map[string][]byte),
		failed:   make(map[string]bool),
	}
}

// attach wires the audio player to the asset registry. Called from
// NewGame after AssetManager is constructed.
func (a *AudioPlayer) attach(am *AssetManager) { a.am = am }

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

// resolveAsset returns the registered Asset for name, or false if it's
// unknown, the wrong kind, or already on the failure list.
func (a *AudioPlayer) resolveAsset(name string) (Asset, bool) {
	if name == "" || a.am == nil || a.failed[name] {
		return Asset{}, false
	}
	asset, ok := a.am.Get(name)
	if !ok || asset.Kind != AssetAudio || asset.Path == "" {
		return Asset{}, false
	}
	return asset, true
}

func (a *AudioPlayer) PlayMusic(name string) {
	if a.musicName == name && a.musicPlayer != nil && a.musicPlayer.IsPlaying() {
		return
	}
	a.stopMusic()
	a.musicName = name

	asset, ok := a.resolveAsset(name)
	if !ok {
		logf("audio.PlayMusic %q: no asset", name)
		return
	}
	ctx := a.ensureContext()
	stream, length, err := decodeAudioStream(asset.Path, audioSampleRate)
	if err != nil {
		logf("audio.PlayMusic %q: %v", name, err)
		a.failed[name] = true
		return
	}
	loop := audio.NewInfiniteLoop(stream, length)
	p, err := audio.NewPlayer(ctx, loop)
	if err != nil {
		logf("audio.PlayMusic %q: player: %v", name, err)
		a.failed[name] = true
		return
	}
	a.musicPlayer = p
	p.Play()
}

func (a *AudioPlayer) StopMusic() {
	if a.musicName == "" && a.musicPlayer == nil {
		return
	}
	logf("audio.StopMusic (was %q)", a.musicName)
	a.stopMusic()
	a.musicName = ""
}

func (a *AudioPlayer) stopMusic() {
	if a.musicPlayer == nil {
		return
	}
	_ = a.musicPlayer.Close()
	a.musicPlayer = nil
}

func (a *AudioPlayer) PlaySound(name string) {
	asset, ok := a.resolveAsset(name)
	if !ok {
		logf("audio.PlaySound %q: no asset", name)
		return
	}
	ctx := a.ensureContext()

	pcm, ok := a.sfxCache[name]
	if !ok {
		data, err := decodeAudioBytes(asset.Path, audioSampleRate)
		if err != nil {
			logf("audio.PlaySound %q: %v", name, err)
			a.failed[name] = true
			return
		}
		a.sfxCache[name] = data
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
