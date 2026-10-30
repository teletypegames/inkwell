// Package keybinding is a stand-alone inkbase widget: [KeyBinding], a key and
// the action it queues — the whole widget. It draws nothing and occupies
// nothing; it exists so that "Escape opens the main menu" is a registration
// like everything else rather than a special case in the loop.
//
// The key is claimed, so nothing underneath sees it, and When is how a
// binding belongs to one Stage. DuringScript lets it answer while a script
// holds the frame, which is what a pause key needs.
//
// The host provides WidgetInput (an [input.Input]) and DoAction.
package keybinding
