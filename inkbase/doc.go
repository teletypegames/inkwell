// Package inkbase is the generic foundation of inkwell: the registries, text
// in several languages, world state, the action and condition framework and
// the data types any Ebitengine game can build on without pulling in
// adventure-game machinery — no scenes, characters, dialogs or verbs. This
// package itself is pure Go with no dependency beyond the standard library;
// its sub-packages import Ebitengine only where they must.
//
// The root inkwell package builds on these types and re-exports them, so an
// adventure imports inkwell alone; a game that is not an adventure imports
// inkbase and the sub-packages it needs.
//
// # What is here
//
//   - [Manager] is the one registry shape: anything [Named] is registered by
//     name, Register panics on an empty or duplicate name, and every walk
//     is in insertion order. [Knowable] is a Named that also has a
//     TypeLabel, so a State can record it.
//   - [Text] is a string in several languages, one [Translation] per
//     language, resolved against [Lang] and then [Fallback] when it is
//     displayed; [Format] composes sentences in one language throughout and
//     [Plain] wraps a string that is not a translation of anything.
//   - [Point], [Rectangle], [Polygon] and [Shape] are the geometry, with
//     [Rect], [Poly] and [PolyRect] to build them.
//   - [State] is a run's mutable record: flags, vars, visit and talk counters,
//     and what the player has come to know.
//   - [Action], [Runner], [Status] and [Ctx] are the action framework, and
//     [Condition] the predicate every gate is made of.
//   - [Asset] (a picture), [SceneLayer], and [Audible] with [Voice], [Music]
//     and [Sound] are the media types.
//   - [Stage], [Script], [Speaker], [Talker] and [Delivery] are the generic
//     halves of what an inkwell game shows and says; [LogKind] and
//     [LogMessage] are the chat log's entries; [Align], [Size] and
//     [MouseButton] are small widget vocabulary.
//
// # Sub-packages
//
// The Ebitengine services are font (the glyph atlas and text drawing), input
// (consume-on-use input polling) and asset (the filesystem, the image cache
// and the audio player). The framework packages are action (the generic
// actions and combinators), condition (the generic conditions) and widget
// (the Widget contract, layers and ordering). The stand-alone widgets are
// menu, panel, textcard, keybinding, letterbox, endcard, statusline, topbar
// and chatlog.
//
// A stand-alone widget reaches its game through a narrow unexported host
// interface it type-asserts the UICtx's Game to, and an action or condition
// that needs more than the State does the same with an exported host
// interface. The host interfaces are listed in each package. inkwell's
// *Game satisfies the action and condition hosts; it carries its own
// versions of the widgets, wired to its theme, so the stand-alone widgets
// are for a game that brings its own host.
package inkbase
