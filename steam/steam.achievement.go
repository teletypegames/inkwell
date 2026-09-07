package steam

import (
	"log"

	inkwell "git.teletypegames.org/engines/inkwell"
	"github.com/bogdannbv/go-steamworks"
)

// Unlock records an achievement with Steam and stores the stats, which is
// what makes the client show its toast. Setting one the player already has
// is a no-op on Steam's side, so the engine reporting it again after a Load
// costs nothing.
func (steam *Steam) Unlock(achievement *inkwell.Achievement) {
	if !steam.running {
		return
	}
	stats := steamworks.SteamUserStats()
	if !stats.SetAchievement(achievement.Name) {
		log.Printf("steam: achievement %q not set", achievement.Name)
		return
	}
	if !stats.StoreStats() {
		log.Printf("steam: achievement %q set but stats not stored", achievement.Name)
	}
}
