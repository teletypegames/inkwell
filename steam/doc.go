// Package steam is inkwell's Steamworks layer: one [inkwell.Extension] that
// initialises the Steam API before the window opens, pumps its callbacks
// every frame, reports achievements, keeps chosen documents in Steam Cloud
// and offers the client's language as the game's SystemLang. It binds
// through github.com/bogdannbv/go-steamworks, which is purego rather than
// cgo, and runs on 64-bit Windows, Linux and macOS only.
//
// The root inkwell package does not import this one, so a game that never
// names it links none of it, nor the binding underneath. A game switches
// Steam on by putting a [Steam] in its Game literal:
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
// # What it does
//
// [Steam.Start] loads the redistributable, relaunches the binary through
// Steam when it was started outside it (RestartAppIfNecessary with AppID),
// initialises the API, and then hands the game what it learned: the
// client's language through Languages, which maps Steam's names ("english",
// "hungarian") onto the game's codes, and, for the Store keys named in Cloud,
// a Store backed by Steam Cloud, one file per document. Every other key
// stays with the platform's default store.
//
// The cloud store degrades toward the disk. A cloud read that finds nothing
// falls back to the local document — which is how a save made before the
// game had Steam reaches the load screen and is uploaded on the next write —
// and a failed cloud write is written locally instead and logged.
//
// [Steam.Tick] runs the Steam callbacks, [Steam.Stop] shuts the API down,
// and [Steam.Unlock] sets an achievement by its Name and stores the stats;
// Steam answers a repeated unlock idempotently.
//
// When the library cannot be loaded or the API does not initialise, the
// extension logs the reason and stands down: the game runs without Steam
// rather than not at all. [Steam.Running] says which happened.
//
// # Runtime requirements
//
// The Steam client must be running, and the redistributable —
// steam_api64.dll, libsteam_api.so or libsteam_api.dylib — must sit in the
// working directory, which is where Steam launches a game from; the
// STEAMWORKS_LIB_PATH environment variable names another location. A run
// outside Steam needs a steam_appid.txt beside the binary; the STEAM_APPID
// environment variable makes the binding write one.
package steam
