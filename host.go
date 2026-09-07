package inkwell

// host.go adds exported methods to *Game that satisfy the host interfaces
// act/ and cond/ actions type-assert to from ctx.Game.

// UseThemeByName satisfies inkbase/action.ThemeHost.
func (g *Game) UseThemeByName(name string) { g.activeTheme = name }

// SetQuitting satisfies inkbase/action.Quitter.
func (g *Game) SetQuitting() { g.quitting = true }

// --- inventory host methods (for act/) ---

// InventoryAdd satisfies act.InventoryHost.
func (g *Game) InventoryAdd(name string)    { g.Inventory.Add(name) }
func (g *Game) InventoryRemove(name string) { g.Inventory.Remove(name) }
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
func (g *Game) ConsumedClick() bool { return g.Input.consumedClick() }

// ChangeScene transitions to the named scene.
func (g *Game) ChangeScene(name string) { g.changeScene(name) }

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
