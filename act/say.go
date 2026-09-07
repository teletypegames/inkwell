package act

import "git.teletypegames.org/engines/inkwell/inkbase"

type sayAction struct {
	speaker string
	text    inkbase.Text
	answer  func(state *inkbase.State) inkbase.Text
}
type sayRunner struct {
	spec     *sayAction
	elapsed  float64
	duration float64
	started  bool
	bubble   bool
	line     string
}

// Say puts one line in a character's mouth.
func Say(speaker inkbase.Talker, text inkbase.Text) inkbase.Action {
	return &sayAction{
		speaker: speaker.TalkerName(),
		text:    text,
	}
}

func (a *sayAction) Start() inkbase.Runner { return &sayRunner{spec: a} }

func (r *sayRunner) Tick(ctx *inkbase.Ctx) inkbase.Status {
	host := ctx.Game.(SpeechHost)
	if !r.started {
		r.started = true
		text := r.spec.text
		if r.spec.answer != nil {
			text = r.spec.answer(ctx.State)
		}
		line := text.Resolved()
		r.line = line.Text
		r.duration = 1.2 + float64(len([]rune(r.line)))*0.05
		if line.Voice != nil && host.VoiceOn() {
			if seconds := host.PlayVoice(line.Voice); seconds > 0 {
				r.duration = seconds
			}
		}
		host.SetTalking(r.spec.speaker)
		host.NoteSpoken(r.spec.speaker)
		r.bubble = host.DeliveryOf(r.spec.speaker) == inkbase.DeliveryBubble
		if r.bubble {
			host.SetSpeech(r.spec.speaker, r.line)
		}
		host.LogResponse(r.spec.speaker, r.line)
	}
	r.elapsed += ctx.DT
	if host.ConsumedClick() {
		r.elapsed = r.duration
	}
	if r.elapsed >= r.duration {
		host.StopVoice()
		if host.Talking() == r.spec.speaker {
			host.SetTalking("")
		}
		if r.bubble {
			host.ClearSpeech()
		}
		return inkbase.StatusDone
	}
	return inkbase.StatusRunning
}
