package inkwell

import (
	"encoding/json"
	"errors"
	"fmt"
)

const settingsVersion = 1

// Settings is what a player sets once and expects to find again: the
// language, the theme, the window scale. Extra is the game's own — anything
// a game wants kept beside these without the engine having to know it —
// and it round-trips untouched.
type Settings struct {
	Version     int    `json:"version"`
	Lang        string `json:"lang,omitempty"`
	Theme       string `json:"theme,omitempty"`
	WindowScale int    `json:"window_scale,omitempty"`
	// Voice is whether recorded lines are played. On until switched off.
	Voice bool           `json:"voice"`
	Extra map[string]any `json:"extra,omitempty"`
}

const settingsKey = "settings"

// Settings is the document as it stands: what was loaded at boot, updated
// by every SaveSettings since. A game reads and writes its Extra here.
func (g *Game) Settings() *Settings { return &g.settings }

// VoiceOn reports whether recorded lines play. SetVoice is the action that
// flips it and saves the settings.
func (g *Game) VoiceOn() bool { return g.settings.Voice }

// loadSettings reads the stored document and applies it: the language when
// none has been chosen for this run, the theme when the stored one is
// registered, the window scale when the game named none. A first run has
// nothing stored, and that is not an error.
func (g *Game) loadSettings() error {
	data, err := g.store().Read(settingsKey)
	if errors.Is(err, ErrNotStored) {
		return nil
	}
	if err != nil {
		return err
	}
	stored := g.settings
	if err := json.Unmarshal(data, &stored); err != nil {
		return fmt.Errorf("inkwell: settings unmarshal: %w", err)
	}
	g.settings = stored
	if Lang == "" && stored.Lang != "" {
		Lang = stored.Lang
	}
	if stored.Theme != "" && g.ThemeManager.Has(stored.Theme) {
		g.activeTheme = stored.Theme
	}
	if g.WindowScale == 0 && stored.WindowScale > 0 {
		g.WindowScale = stored.WindowScale
	}
	return nil
}

// SaveSettings writes the document with what the game is doing now: the
// language every Text resolves into, the active theme, the window scale.
// UseLang and UseTheme call it, so a choice on a settings screen is kept
// the moment it is made.
func (g *Game) SaveSettings() error {
	g.settings.Version = settingsVersion
	g.settings.Lang = Lang
	g.settings.Theme = g.activeTheme
	g.settings.WindowScale = g.WindowScale
	data, err := json.MarshalIndent(g.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("inkwell: settings marshal: %w", err)
	}
	if err := g.store().Write(settingsKey, data); err != nil {
		logf("settings: %v", err)
		return err
	}
	return nil
}
