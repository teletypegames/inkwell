package inkwell

import (
	ibaction "git.teletypegames.org/engines/inkwell/inkbase/action"
	"git.teletypegames.org/engines/inkwell/inkbase/chatlog"
	ibcond "git.teletypegames.org/engines/inkwell/inkbase/condition"
	"git.teletypegames.org/engines/inkwell/inkbase/endcard"
	"git.teletypegames.org/engines/inkwell/inkbase/keybinding"
	"git.teletypegames.org/engines/inkwell/inkbase/letterbox"
	"git.teletypegames.org/engines/inkwell/inkbase/menu"
	"git.teletypegames.org/engines/inkwell/inkbase/panel"
	"git.teletypegames.org/engines/inkwell/inkbase/statusline"
	"git.teletypegames.org/engines/inkwell/inkbase/textcard"
	"git.teletypegames.org/engines/inkwell/inkbase/topbar"
)

// host.go adds exported methods to *Game that satisfy the host interfaces
// act/ and cond/ actions type-assert to from ctx.Game. The host methods of
// the inkbase widgets live next to their re-exports in the widget.*.go
// files; the assertions below keep *Game satisfying every one of them, and
// the widgets that own part of the screen satisfying ClickBlocker.
var (
	_ ibaction.StageHost = (*Game)(nil)
	_ ibcond.StageReader = (*Game)(nil)
	_ chatlog.Host       = (*Game)(nil)
	_ endcard.Host       = (*Game)(nil)
	_ keybinding.Host    = (*Game)(nil)
	_ letterbox.Host     = (*Game)(nil)
	_ menu.Host          = (*Game)(nil)
	_ panel.Host         = (*Game)(nil)
	_ statusline.Host    = (*Game)(nil)
	_ textcard.Host      = (*Game)(nil)
	_ topbar.Host        = (*Game)(nil)
	_ ClickBlocker       = (*ChatLog)(nil)
	_ ClickBlocker       = (*Menu)(nil)
	_ ClickBlocker       = (*Panel)(nil)
	_ ClickBlocker       = (*TopBar)(nil)
)

// UseThemeByName satisfies inkbase/action.ThemeHost: it makes the theme
// registered under name the active one; an unknown name changes nothing.
func (g *Game) UseThemeByName(name string) {
	if g.ThemeManager.Has(name) {
		g.activeTheme = name
	}
}

// SetQuitting satisfies inkbase/action.Quitter.
func (g *Game) SetQuitting() { g.quitting = true }

// --- inventory host methods (for act/) ---

// InventoryAdd puts the named item in the player's Inventory; it satisfies
// [act.InventoryHost].
func (g *Game) InventoryAdd(name string) { g.Inventory.Add(name) }

// InventoryRemove takes the named item out of the player's Inventory; it
// satisfies [act.InventoryHost].
func (g *Game) InventoryRemove(name string) { g.Inventory.Remove(name) }

// InventoryHas reports whether the player carries the named item; it
// satisfies [act.InventoryHost].
func (g *Game) InventoryHas(name string) bool { return g.Inventory.Has(name) }

// InventorySelected satisfies cond.InventoryChecker.
func (g *Game) InventorySelected() string { return g.Inventory.Selected() }

// --- condition host methods (for cond/) ---

// CurrentStagePtr satisfies cond.StageChecker. It returns the current Stage
// as an any so that cond/ can compare by pointer identity without importing
// the Stage type.
func (g *Game) CurrentStagePtr() any { return g.currentStage }

// --- adventure host methods (for act/) ---

// NoteSpoken records that the voice behind a name has been heard.
func (g *Game) NoteSpoken(name string) { g.noteSpoken(name) }

// Talking is the name of the character or speaker whose line is showing.
func (g *Game) Talking() string { return g.talking }

// SetTalking sets the speaker whose line is up.
func (g *Game) SetTalking(name string) { g.talking = name }

// PlayVoice plays a voice recording and returns its duration in seconds.
func (g *Game) PlayVoice(v *Voice) float64 { return g.Audio.PlayVoice(v) }

// StopVoice stops the current voice playback.
func (g *Game) StopVoice() { g.Audio.StopVoice() }

// ConsumedClick reports (and consumes) a left click for skip-line.
func (g *Game) ConsumedClick() bool { return g.Input.ConsumedClick() }

// ChangeScene transitions to the named scene and queues the old scene's
// OnLeave and the new scene's OnEnter; an unknown name logs and does
// nothing. A script that has more to do in the new scene uses [GoTo], or
// [Game.EnterScene].
func (g *Game) ChangeScene(name string) { g.changeScene(name, nil) }

// EnterScene starts a transition to the named scene and returns the Runner
// that holds a script through it: the Runner waits for the swap, runs the
// old scene's OnLeave and the new scene's OnEnter in its own place, and is
// done when the new picture is fully uncovered, so what follows it in the
// script happens in the new scene. An unknown name logs and returns nil.
func (g *Game) EnterScene(name string) Runner {
	change := &sceneChange{}
	if !g.changeScene(name, change) {
		return nil
	}
	return &sceneChangeRunner{g: g, change: change}
}

// StartDialog opens a dialog tree.
func (g *Game) StartDialog(d *Dialog) { g.startDialog(d) }

// DialogActive reports whether a dialog is in progress.
func (g *Game) DialogActive() bool { return g.dialogActive() }

// EndDialog closes the active dialog.
func (g *Game) EndDialog() { g.endDialog() }

// GotoDialogNode moves to a node within the active dialog.
func (g *Game) GotoDialogNode(name string) { g.gotoDialogNode(name) }

// WalkCharacter sends a character toward a point through the walkboxes.
func (g *Game) WalkCharacter(name string, to Point) { g.walkCharacter(name, to) }

// CharacterMoving reports whether a character is en route.
func (g *Game) CharacterMoving(name string) bool { return g.characterMoving(name) }

// SetCharacterPose freezes a character's pose.
func (g *Game) SetCharacterPose(name string, pose Pose) {
	if rc, ok := g.runtimeChar(name); ok {
		rc.pose = pose
		rc.posed = true
		rc.poseElapsed = 0
	}
}

// SetVoiceEnabled toggles voice recording playback.
func (g *Game) SetVoiceEnabled(on bool) { g.settings.Voice = on }

// ShowEndCard puts up the end-of-game overlay.
func (g *Game) ShowEndCard(text string) { g.showEndCard(text) }
