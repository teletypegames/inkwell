package inkwell

type sayAction struct {
	speaker string
	text    Text
	answer  func(state *State) Text
}
type sayRunner struct {
	spec     *sayAction
	elapsed  float64
	duration float64
	started  bool
	bubble   bool
	line     string
}

// Say puts one line in a character's mouth. The line is a Text, so the same
// call carries every language it has been written in; which one reaches the
// player is decided when the runner starts, not when the content is
// authored.
func Say(speaker Talker, text Text) Action {
	return &sayAction{
		speaker: talkerName(speaker),
		text:    text,
	}
}

func (a *sayAction) Start() Runner { return &sayRunner{spec: a} }

func (r *sayRunner) Tick(ctx *Ctx) Status {
	if !r.started {
		r.started = true
		text := r.spec.text
		if r.spec.answer != nil {
			text = r.spec.answer(ctx.Game.State)
		}
		line := text.resolved()
		r.line = line.Text
		// duration scales with text length, with a 1.2s floor — unless the
		// line is recorded and voices are on, when the recording sets it
		r.duration = 1.2 + float64(len([]rune(r.line)))*0.05
		if line.Voice != nil && ctx.Game.VoiceOn() {
			if seconds := ctx.Game.Audio.PlayVoice(line.Voice); seconds > 0 {
				r.duration = seconds
			}
		}
		ctx.Game.talking = r.spec.speaker
		ctx.Game.noteSpoken(r.spec.speaker)
		r.bubble = ctx.Game.DeliveryOf(r.spec.speaker) == DeliveryBubble
		if r.bubble {
			ctx.Game.SetSpeech(r.spec.speaker, r.line)
		}
		ctx.Game.LogResponse(r.spec.speaker, r.line)
	}
	r.elapsed += ctx.DT
	// skip on click
	if ctx.Game.Input.consumedClick() {
		r.elapsed = r.duration
	}
	if r.elapsed >= r.duration {
		ctx.Game.Audio.StopVoice()
		if ctx.Game.talking == r.spec.speaker {
			ctx.Game.talking = ""
		}
		if r.bubble {
			ctx.Game.ClearSpeech()
		}
		return StatusDone
	}
	return StatusRunning
}
