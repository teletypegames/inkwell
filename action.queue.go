package inkwell

// Do queues an action. Queued actions run one at a time and in order —
// the engine starts the next one on the first frame the previous one is
// done, so a widget or a domain pump can hand work in at any moment.
func (g *Game) Do(a Action) {
	if a == nil {
		return
	}
	g.queue = append(g.queue, a)
}

func (g *Game) queueAction(a Action, label string) {
	if a == nil {
		return
	}
	logf("queueAction %s", label)
	g.Do(a)
}

// pumpQueue starts the next queued action when nothing else is running.
func (g *Game) pumpQueue() {
	if g.scriptRunner != nil || len(g.queue) == 0 {
		return
	}
	a := g.queue[0]
	g.queue = g.queue[1:]
	g.scriptRunner = a.Start()
	g.scriptCtx = g.makeCtx()
}
