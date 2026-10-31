// Package keybinding is an inkbase widget: [KeyBinding], a key and the action
// it queues — the whole widget. It draws nothing and occupies nothing; it
// exists so that "Escape opens the main menu" is a registration like
// everything else rather than a special case in the loop.
//
// The key is claimed, so nothing underneath sees it, and When is how a
// binding belongs to one Stage. DuringScript lets it answer while a script
// holds the frame, which is what a pause key needs.
//
// The widget reaches its game through [Host]: WidgetInput (an [input.Input])
// and Do. inkwell's *Game is a Host and the root package re-exports
// KeyBinding.
package keybinding
