package action

import "git.teletypegames.org/engines/inkwell/inkbase"

// SettingsSaver can persist settings after a change.
type SettingsSaver interface {
	SaveSettings() error
}

type useLangAction struct{ lang string }

// UseLang switches the language every Text resolves into.
func UseLang(lang string) inkbase.Action { return &useLangAction{lang: lang} }

func (a *useLangAction) Start() inkbase.Runner { return a }

func (a *useLangAction) Tick(ctx *inkbase.Ctx) inkbase.Status {
	inkbase.Lang = a.lang
	if h, ok := ctx.Game.(SettingsSaver); ok {
		_ = h.SaveSettings()
	}
	return inkbase.StatusDone
}
