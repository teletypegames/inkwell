// Package chatlog is a stand-alone inkbase widget: [ChatLog], the message
// log — player commands, in-world replies and system lines — drawn as a box
// with the newest line at the bottom, older lines scrolling up out of view.
//
// A [SpeakerStyle] says how a reply names who said it: inline in front of
// every reply, once as a header above a rule (which suits a panel given over
// to one voice at a time), or not at all. Deliveries keeps only the replies
// delivered one of the ways it lists, which is how a game with a radio, a
// tape or a narrator gives that voice a panel of its own without the widget
// knowing what the thing is; player commands and system lines are nobody's,
// so they are kept either way. Resting names who the header carries before
// anybody speaks.
//
// The host provides DrawText, Messages, SpeakerDelivery, SpeakerColor,
// SpeakerLabel and ChatLogColors() [Colors].
package chatlog
