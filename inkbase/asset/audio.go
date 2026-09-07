package asset

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

	"git.teletypegames.org/engines/inkwell/inkbase"
)

const audioSampleRate = 48000
const audioFrameBytes = 4

// AudioPlayer plays music, recorded lines and one-shot sound effects
// through Ebiten's audio subsystem. Fail-soft: missing files degrade to
// no-ops.
type AudioPlayer struct {
	ctx *audio.Context

	music       *inkbase.Music
	musicPlayer *audio.Player

	voicePlayer *audio.Player

	sfxCache  map[string][]byte
	sfxActive []*audio.Player

	failed map[string]bool

	logf func(string, ...any)
}

func NewAudioPlayer(logf func(string, ...any)) *AudioPlayer {
	return &AudioPlayer{
		sfxCache: make(map[string][]byte),
		failed:   make(map[string]bool),
		logf:     logf,
	}
}

func (a *AudioPlayer) ensureContext() *audio.Context {
	if a.ctx != nil {
		return a.ctx
	}
	if c := audio.CurrentContext(); c != nil {
		a.ctx = c
		return c
	}
	a.ctx = audio.NewContext(audioSampleRate)
	return a.ctx
}

func (a *AudioPlayer) playable(audible inkbase.Audible) bool {
	if audible == nil {
		return false
	}
	source := audible.Source()
	return source != "" && !a.failed[source]
}

func sourceName(audible inkbase.Audible) string {
	if audible == nil {
		return "<nil>"
	}
	return audible.Source()
}

func (a *AudioPlayer) log(format string, args ...any) {
	if a.logf != nil {
		a.logf(format, args...)
	}
}

func (a *AudioPlayer) PlayMusic(music *inkbase.Music) {
	if a.music == music && a.musicPlayer != nil && a.musicPlayer.IsPlaying() {
		return
	}
	a.stopMusic()
	a.music = music

	if music == nil || !a.playable(music) {
		a.log("audio.PlayMusic %s: no source", sourceName(music))
		return
	}
	ctx := a.ensureContext()
	stream, length, err := decodeAudioStream(music.Path, audioSampleRate)
	if err != nil {
		a.log("audio.PlayMusic %q: %v", music.Path, err)
		a.failed[music.Path] = true
		return
	}
	var source io.Reader = stream
	if music.Loop {
		source = audio.NewInfiniteLoop(stream, length)
	}
	p, err := audio.NewPlayer(ctx, source)
	if err != nil {
		a.log("audio.PlayMusic %q: player: %v", music.Path, err)
		a.failed[music.Path] = true
		return
	}
	a.musicPlayer = p
	p.Play()
}

func (a *AudioPlayer) StopMusic() {
	if a.music == nil && a.musicPlayer == nil {
		return
	}
	a.log("audio.StopMusic (was %s)", sourceName(a.music))
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

// PlayVoice plays a recorded line once and reports its duration in seconds.
func (a *AudioPlayer) PlayVoice(voice *inkbase.Voice) float64 {
	a.StopVoice()
	if voice == nil || !a.playable(voice) {
		a.log("audio.PlayVoice %s: no source", sourceName(voice))
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

func (a *AudioPlayer) StopVoice() {
	if a.voicePlayer == nil {
		return
	}
	_ = a.voicePlayer.Close()
	a.voicePlayer = nil
}

func (a *AudioPlayer) PlaySound(sound *inkbase.Sound) {
	if sound == nil || !a.playable(sound) {
		a.log("audio.PlaySound %s: no source", sourceName(sound))
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

func (a *AudioPlayer) pcm(audible inkbase.Audible) ([]byte, bool) {
	source := audible.Source()
	if data, ok := a.sfxCache[source]; ok {
		return data, true
	}
	data, err := decodeAudioBytes(source, audioSampleRate)
	if err != nil {
		a.log("audio.decode %q: %v", source, err)
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

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

func decodeAudioStream(path string, sampleRate int) (audioStream, int64, error) {
	f, err := OpenAsset(path)
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

func decodeAudioBytes(path string, sampleRate int) ([]byte, error) {
	s, _, err := decodeAudioStream(path, sampleRate)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(s)
}
