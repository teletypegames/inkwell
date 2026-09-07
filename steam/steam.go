// Package steam is inkwell's Steamworks layer: one inkwell.Extension that
// initialises the Steam API before the window opens, pumps its callbacks
// every frame, reports achievements, keeps the save slots in Steam Cloud and
// offers the client's language as the game's SystemLang. It binds through
// github.com/bogdannbv/go-steamworks, which is purego rather than cgo, and
// runs on 64-bit Windows, Linux and macOS only.
//
// The root inkwell package does not import this one, so a game that never
// names it links none of it. A game switches Steam on by putting a *Steam in
// its Game literal:
//
//	Extensions: []inkwell.Extension{
//		&steam.Steam{
//			AppID: 480,
//			Languages: map[string]string{
//				"english":   "en",
//				"hungarian": "hu",
//			},
//			Cloud: []string{inkwell.SlotsKey},
//		},
//	},
//
// At runtime the Steam client must be running and the redistributable —
// steam_api64.dll, libsteam_api.so or libsteam_api.dylib — must sit in the
// working directory, which is where Steam launches a game from; the
// STEAMWORKS_LIB_PATH environment variable names another location. A run
// outside Steam needs steam_appid.txt beside the binary. When the library
// cannot be loaded or the API does not initialise, the extension logs the
// reason and stands down: the game runs without Steam rather than not at
// all.
package steam

import (
	"log"
	"os"

	inkwell "git.teletypegames.org/engines/inkwell"
	"github.com/bogdannbv/go-steamworks"
)

// Steam is the extension. Every field is a setting the game makes; the
// zero value is an app id of 0, no language mapping and no cloud, which is
// only good for finding out whether the client is there.
type Steam struct {
	// AppID is the game's Steam application id. RestartAppIfNecessary is
	// asked with it at Start, so a binary launched from outside Steam is
	// relaunched through it.
	AppID uint32
	// Languages maps the Steam API language name — "english", "hungarian",
	// see ISteamApps::GetCurrentGameLanguage — to the game's own language
	// code. The client's current language becomes Game.SystemLang when the
	// map names it; a language the map leaves out sets nothing.
	Languages map[string]string
	// Cloud lists the Store keys kept in Steam Cloud rather than on disk:
	// inkwell.SlotsKey for the saves, inkwell.SettingsKey for the
	// settings. Every other key stays with the platform's default Store.
	// Empty leaves the Store alone.
	Cloud []string

	running bool
}

// Start loads the library, relaunches through Steam when the binary was
// started outside it, initialises the API, and then hands the game what it
// learned: the client's language and, for the keys named in Cloud, a Store
// backed by Steam Cloud. It never returns an error: a missing client is
// logged and the game goes on without Steam.
func (steam *Steam) Start(game *inkwell.Game) error {
	if err := steamworks.Load(); err != nil {
		log.Printf("steam: %v; running without Steam", err)
		return nil
	}
	if steam.AppID != 0 && steamworks.RestartAppIfNecessary(steam.AppID) {
		os.Exit(0)
	}
	if err := steamworks.Init(); err != nil {
		log.Printf("steam: %v; running without Steam", err)
		return nil
	}
	steam.running = true

	if game.SystemLang == "" {
		game.SystemLang = steam.Languages[steamworks.SteamApps().GetCurrentGameLanguage()]
	}
	if game.Store == nil && len(steam.Cloud) > 0 {
		game.Store = newCloudStore(steam.Cloud, inkwell.NewFileStore(game.AppName()))
	}
	return nil
}

// Tick pumps the Steam callbacks. Steamworks wants this on the main thread
// every frame, which is what an Extension's Tick is.
func (steam *Steam) Tick(*inkwell.Ctx) {
	if steam.running {
		steamworks.RunCallbacks()
	}
}

// Stop shuts the API down once the game loop has returned.
func (steam *Steam) Stop() {
	if steam.running {
		steamworks.Shutdown()
		steam.running = false
	}
}

// Running reports whether the Steam API is up: the client was found and
// initialised at Start.
func (steam *Steam) Running() bool { return steam.running }
