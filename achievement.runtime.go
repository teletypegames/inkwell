package inkwell

// tickAchievements evaluates every registered Achievement against the world
// as it stands and hands the ones that have just come true to every
// Unlocker among the extensions. An achievement is reported once per
// process: unlocked remembers what has been handed over, so a condition
// that stays true does not fire every frame. A Load that finds a condition
// already met reports it again on the next run of the process, which is
// what lets an achievement earned offline reach the platform.
func (g *Game) tickAchievements(ctx *Ctx) {
	if g.AchievementManager.Len() == 0 {
		return
	}
	for _, a := range g.AchievementManager.All() {
		if g.unlocked[a.Name] || a.When == nil || !a.When.Eval(ctx) {
			continue
		}
		g.unlocked[a.Name] = true
		g.unlock(a)
	}
}
