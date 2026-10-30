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

// audioFrameBytes is one frame of Ebiten's PCM: 16-bit stereo.
const audioFrameBytes = 4

// AudioPlayer plays music, recorded lines and one-shot sound effects
// through Ebiten's audio subsystem. Music is streamed, looping when its
// Music says so; a Voice is decoded and played once, and can be stopped
// when the line is skipped; sound effects are decoded once, cached as raw
// PCM, and replayed via NewPlayerFromBytes for low-latency triggering.
//
// AudioPlayer is fail-soft by design: missing files, decode errors and a
// missing audio device all degrade to no-ops with a debug log message — the
// game keeps running.
type AudioPlayer struct {
	ctx *audio.Context // lazy — created on first real play

	music       *Music
	musicPlayer *audio.Player

	voicePlayer *audio.Player

	sfxCache  map[string][]byte // resampled PCM, keyed by Source
	sfxActive []*audio.Player   // in-flight one-shots; pruned each PlaySound

	failed map[string]bool // sources that errored once — don't retry every frame
}

// NewAudioPlayer returns an AudioPlayer with nothing playing. The audio
// context is not created until the first thing is played.
func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		sfxCache: make(map[string][]byte),
		failed:   make(map[string]bool),
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

// playable reports whether something can be started: present, with a
// source, and not already on the failure list.
func (a *AudioPlayer) playable(audible Audible) bool {
	if audible == nil {
		return false
	}
	source := audible.Source()
	return source != "" && !a.failed[source]
}

func sourceName(audible Audible) string {
	if audible == nil {
		return "<nil>"
	}
	return audible.Source()
}

// PlayMusic starts music in place of the current track, looping it when
// its Loop says so. The same track already playing carries on untouched; nil
// or a file that cannot be played stops the current track and plays
// nothing.
func (a *AudioPlayer) PlayMusic(music *Music) {
	if a.music == music && a.musicPlayer != nil && a.musicPlayer.IsPlaying() {
		return
	}
	a.stopMusic()
	a.music = music

	if music == nil || !a.playable(music) {
		logf("audio.PlayMusic %s: no source", sourceName(music))
		return
	}
	ctx := a.ensureContext()
	stream, length, err := decodeAudioStream(music.Path, audioSampleRate)
	if err != nil {
		logf("audio.PlayMusic %q: %v", music.Path, err)
		a.failed[music.Path] = true
		return
	}
	var source io.Reader = stream
	if music.Loop {
		source = audio.NewInfiniteLoop(stream, length)
	}
	p, err := audio.NewPlayer(ctx, source)
	if err != nil {
		logf("audio.PlayMusic %q: player: %v", music.Path, err)
		a.failed[music.Path] = true
		return
	}
	a.musicPlayer = p
	p.Play()
}

// StopMusic stops the current track and forgets it, so the next PlayMusic
// starts afresh even with the same Music.
func (a *AudioPlayer) StopMusic() {
	if a.music == nil && a.musicPlayer == nil {
		return
	}
	logf("audio.StopMusic (was %s)", sourceName(a.music))
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

// PlayVoice plays a recorded line once and reports how long it lasts, in
// seconds, so the line can stay up exactly that long; 0 means nothing
// played. A Voice already playing is stopped first: one line at a time.
func (a *AudioPlayer) PlayVoice(voice *Voice) float64 {
	a.StopVoice()
	if voice == nil || !a.playable(voice) {
		logf("audio.PlayVoice %s: no source", sourceName(voice))
		return 0
	}
	pcm, ok := a.pcm(voice)
	if !ok {
		return 0
	}
	p := audio.NewPlayerFromBytes(a.ensureContext(), pcm)
	p.Play()
	a.voicePlayer = p
	return float64(len(pcm)) / float64(audioSampleRate*audioFrameBytes)
}

// StopVoice cuts the recorded line short — a skipped line stops talking.
func (a *AudioPlayer) StopVoice() {
	if a.voicePlayer == nil {
		return
	}
	_ = a.voicePlayer.Close()
	a.voicePlayer = nil
}

// PlaySound plays a sound once, on top of anything else playing. The file
// is decoded on first use and kept; nil or a file that cannot be played
// does nothing.
func (a *AudioPlayer) PlaySound(sound *Sound) {
	if sound == nil || !a.playable(sound) {
		logf("audio.PlaySound %s: no source", sourceName(sound))
		return
	}
	pcm, ok := a.pcm(sound)
	if !ok {
		return
	}
	p := audio.NewPlayerFromBytes(a.ensureContext(), pcm)
	p.Play()
	a.sfxActive = append(a.sfxActive, p)
	a.pruneSfx()
}

// pcm decodes an Audible once into resampled 16-bit stereo PCM and caches
// it by Source.
func (a *AudioPlayer) pcm(audible Audible) ([]byte, bool) {
	source := audible.Source()
	if data, ok := a.sfxCache[source]; ok {
		return data, true
	}
	data, err := decodeAudioBytes(source, audioSampleRate)
	if err != nil {
		logf("audio.decode %q: %v", source, err)
		a.failed[source] = true
		return nil, false
	}
	a.sfxCache[source] = data
	return data, true
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
