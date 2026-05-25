package pncdsl

// AudioPlayer is a stub for music/sfx playback. The action constructors
// (PlayMusic/PlaySound/StopMusic) call through here; on this milestone we
// just log the request so the game runs without an audio device.
type AudioPlayer struct {
	currentMusic string
}

func NewAudioPlayer() *AudioPlayer { return &AudioPlayer{} }

func (a *AudioPlayer) PlayMusic(name string) {
	if a.currentMusic == name {
		return
	}
	a.currentMusic = name
	logf("audio.PlayMusic %q", name)
}

func (a *AudioPlayer) StopMusic() {
	if a.currentMusic == "" {
		return
	}
	logf("audio.StopMusic (was %q)", a.currentMusic)
	a.currentMusic = ""
}

func (a *AudioPlayer) PlaySound(name string) {
	logf("audio.PlaySound %q", name)
}
