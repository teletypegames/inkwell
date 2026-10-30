// Package input provides consume-on-use input polling for inkbase hosts and
// widgets.
//
// The host calls [Input.Poll] once at the start of every frame, which
// snapshots the pointer position, the just-pressed state of both mouse
// buttons, the wheel's travel and every key that went down. Then the one
// click is authoritative: whoever consumes it first wins, and later checks
// see LeftClicked report false. A widget that takes the click calls
// ConsumeLeft or ConsumeRight; a click on a speech line only skips the line
// and does not also fire the hotspot underneath.
//
// Keys work the same way, and [Input.ClaimKey] — KeyPressed plus ConsumeKey —
// is the shape a widget binding half a dozen of them wants: a menu walking a
// list with the arrow keys claims them, so nothing else reads them on that
// frame. KeyPressed is edge-triggered like LeftClicked — holding a key
// answers true once — and a widget that wants a key held down reads
// ebiten.IsKeyPressed itself. The wheel is not consumed.
package input
