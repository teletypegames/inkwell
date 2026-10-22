package inkwell

import (
	"fmt"
	"sort"
	"strings"
)

// core.solve.go plays a game with nobody at the controls and no window. A
// game's content is declared, not coded — hotspots, exits, triggers,
// dialog choices and the actions they run — so every move a player can
// make from a given moment can be listed, made on a copy of the game, and
// the copy compared with where it started.
//
// The unit of search is one step to a milestone: breadth first from a
// moment to the nearest flag raised that is not raised there, or to the
// goal. That step answers two questions:
//
//   - Situation, from the session as it stands: what is the nearest thing
//     the player can do that moves the story on, and how? This is what a
//     companion who gives hints needs to know.
//   - Solve, from the start of a new session: can the game be finished, and
//     how? It takes the nearest milestone, then the nearest from there, and
//     so on to the goal. A milestone from which no further one is in reach
//     is a dead end: Solve reports the first it meets, steps back and tries
//     the next-nearest milestone instead, depth first, never entering the
//     same moment twice.
//
// A breadth-first walk of every moment would count every dead end, but the
// number of moments in a game with several puzzles open at once grows with
// every order they can be done in, and that walk only ends in useful time
// for a small game — SolveOptions.Exhaustive asks for it.
//
// The copy runs the game's real actions through the real engine paths —
// InvokeVerb, the trigger loop, the scene change — with the clock run
// forward so that every line, wait, walk and fade finishes at once. Only
// what reaches the world is compared: scene, flags, vars, inventory and
// trigger bookkeeping, and the previous scene where an exit leads back to
// it. What the player has visited, come to know or said in a Once choice is
// carried along but not compared: a Once choice that matters raises a flag
// or hands over an item, and that is compared.

// Move is one thing a player can do: a verb on a hotspot, an item used on a
// hotspot or on another item, or a choice in an open dialog.
type Move struct {
	// Scene is where the move is made.
	Scene string
	// Hotspot is the Name of the hotspot the move targets — a painted one,
	// an exit ("exit:<scene>") or one on a HotspotPanel. Empty for an item
	// used on an item and for a dialog choice.
	Hotspot string
	// Label is what the target is called on screen.
	Label Text
	// Verb is the verb used. An item used on something is "use".
	Verb string
	// Item is the item held while the verb is used, nil when none is.
	Item *Item
	// With is the item Item is used on, for item on item.
	With *Item
	// Dialog and Choice name a choice in an open dialog, by its position
	// in the node's Choices; Text is the choice's line.
	Dialog *Dialog
	Choice int
	Text   Text
}

// String reads a move the way the log would: "alley: use hook on grate".
func (m Move) String() string {
	switch {
	case m.Dialog != nil:
		return fmt.Sprintf("%s: [%s] %q", m.Scene, m.Dialog.Name, m.Text.String())
	case m.With != nil:
		return fmt.Sprintf("%s: use %s on %s", m.Scene, m.Item.Name, m.With.Name)
	case m.Item != nil:
		return fmt.Sprintf("%s: use %s on %s", m.Scene, m.Item.Name, m.Hotspot)
	}
	return fmt.Sprintf("%s: %s %s", m.Scene, m.Verb, m.Hotspot)
}

// SolveOptions shapes a search.
type SolveOptions struct {
	// Goal is what counts as finished. Nil means the end card is up.
	Goal Condition
	// Limit caps the number of distinct moments explored, over the whole
	// of a Solve or a Situation. 0 means 1000000 for Solve and 20000 for
	// Situation.
	Limit int
	// Exhaustive makes Solve walk every reachable moment breadth first
	// instead of going from milestone to milestone, so that every dead end
	// is counted. Only a small game finishes that way.
	Exhaustive bool
}

// Solution is what Solve found.
type Solution struct {
	// Solved says the goal was reached, and Path is the way it was reached:
	// the shortest way to each milestone in turn, which is not always the
	// shortest way to the end. Milestones is the flag each step raised.
	Solved     bool
	Path       []Move
	Milestones []string
	// DeadEnds counts the moments found from which the goal can no longer
	// be reached — every one of them when the search was Exhaustive, those
	// the milestone search stepped back from otherwise. DeadEnd is the way
	// into the first one, and DeadEndFrontier what still changed the world
	// there, which is where to look for the puzzle that does not open.
	DeadEnds        int
	DeadEnd         []Move
	DeadEndFrontier []Move
	// States is how many distinct moments were reached. Truncated says the
	// search hit SolveOptions.Limit before it could decide.
	States    int
	Truncated bool
	// Unvisited is every registered scene no moment the search reached
	// stands in.
	Unvisited []string
	// Runaway is every move whose actions did not settle — a script that
	// never finishes, a trigger that keeps re-firing — with the way to it.
	Runaway [][]Move
}

// Situation is where a session stands, as seen from the next milestone.
type Situation struct {
	// Found says a milestone is within reach of the search. Milestone is
	// the flag it raises, "" when it is the goal itself; Next is the
	// shortest way to it, so Next[0] is the next thing to do.
	Found     bool
	Milestone string
	Next      []Move
	// Frontier is every move from here that changes the world rather than
	// only where the player stands.
	Frontier []Move
	// States and Truncated are as in Solution.
	States    int
	Truncated bool
}

// simulating counts the searches in progress. A game's own content may ask
// for a Situation — a companion's hint is a line decided when it is spoken —
// and a search plays that content too, so a Situation or Solve asked for from
// inside a search answers empty instead of starting another search inside
// it, which would never end.
var simulating int

const (
	defaultSolveLimit     = 1000000
	defaultSituationLimit = 20000
)

// Solve searches a new session from its start: the start scene entered and
// the start script played, exactly as Run would begin the world. The game
// itself is left untouched, so Solve can run beside Run or instead of it.
func (g *Game) Solve(options SolveOptions) Solution {
	if simulating > 0 {
		return Solution{}
	}
	simulating++
	defer func() { simulating-- }()
	g.resolveDerived()
	if options.Limit <= 0 {
		options.Limit = defaultSolveLimit
	}
	sim := newSimulation(g)
	sim.game.startWorld()
	var solution Solution
	if !sim.settle() {
		solution.Runaway = append(solution.Runaway, nil)
	}
	stood := make(map[string]bool)
	start := origin{base: capture(sim.game)}
	if options.Exhaustive {
		solution.exhaustive(sim, start, options, stood)
	} else {
		solution.milestones(sim, start, options, stood)
	}
	for _, s := range g.SceneManager.All() {
		if !stood[s.Name] {
			solution.Unvisited = append(solution.Unvisited, s.Name)
		}
	}
	return solution
}

// exhaustive walks every reachable moment breadth first.
func (solution *Solution) exhaustive(sim *simulation, start origin, options SolveOptions, stood map[string]bool) {
	run := sim.search(start, options.Limit, func(*searchNode) visit { return expand }, options.Goal)
	solution.absorb(run, stood, nil)
	if goal := run.firstGoal(); goal >= 0 {
		solution.Solved = true
		solution.Path = run.path(goal)
	}
	if run.truncated {
		return
	}
	reaching := run.reachingGoal()
	for i := range run.nodes {
		if reaching[i] {
			continue
		}
		if solution.DeadEnds == 0 {
			solution.DeadEnd = run.path(i)
			solution.DeadEndFrontier = run.frontierOf(i)
		}
		solution.DeadEnds++
	}
}

// milestoneStep is one moment on the milestone path: the search from it,
// the milestones it found nearest first, and which one is being tried.
type milestoneStep struct {
	run        *search
	candidates []*searchNode
	next       int
	path       int
	milestones int
}

// milestones goes from milestone to milestone, depth first, stepping back
// from a dead end to the next-nearest milestone.
func (solution *Solution) milestones(sim *simulation, start origin, options SolveOptions, stood map[string]bool) {
	entered := make(map[string]bool)
	budget := options.Limit
	open := func(from origin) *milestoneStep {
		run, candidates := sim.candidates(from, budget, options.Goal)
		budget -= len(run.nodes)
		solution.absorb(run, stood, solution.Path)
		return &milestoneStep{
			run:        run,
			candidates: candidates,
			path:       len(solution.Path),
			milestones: len(solution.Milestones),
		}
	}
	stack := []*milestoneStep{open(start)}
	entered[stack[0].run.nodes[0].key] = true
	for len(stack) > 0 {
		if budget <= 0 {
			solution.Truncated = true
			return
		}
		top := stack[len(stack)-1]
		solution.Path = solution.Path[:top.path]
		solution.Milestones = solution.Milestones[:top.milestones]
		if top.next >= len(top.candidates) {
			if len(top.candidates) == 0 && !top.run.truncated {
				if solution.DeadEnds == 0 {
					solution.DeadEnd = append([]Move(nil), solution.Path...)
					solution.DeadEndFrontier = top.run.frontier()
				}
				solution.DeadEnds++
			}
			stack = stack[:len(stack)-1]
			continue
		}
		candidate := top.candidates[top.next]
		top.next++
		if entered[candidate.key] {
			continue
		}
		entered[candidate.key] = true
		solution.Path = append(solution.Path, top.run.path(candidate.index)...)
		solution.Milestones = append(solution.Milestones, candidate.milestone)
		if candidate.goal {
			solution.Solved = true
			return
		}
		stack = append(stack, open(origin{base: candidate.base, pending: candidate.pending}))
	}
}

// absorb adds what one search saw to the solution: its size, its runaways,
// prefixed with the path that led to where it started, and the scenes it
// stood in.
func (solution *Solution) absorb(run *search, stood map[string]bool, before []Move) {
	solution.States += len(run.nodes)
	solution.Truncated = solution.Truncated || run.truncated
	for _, path := range run.runaway {
		solution.Runaway = append(solution.Runaway, append(append([]Move(nil), before...), path...))
	}
	for _, n := range run.nodes {
		stood[n.scene] = true
	}
}

// Situation searches from the session as it stands. Whatever is running is
// not carried over — an open dialog, a script mid-way — only what it has
// already done to the world, so the answer is about the moment the player
// is next free to act. A game whose world has not begun is asked about the
// moment it will begin at, as Solve starts.
func (g *Game) Situation(options SolveOptions) Situation {
	if simulating > 0 {
		return Situation{}
	}
	simulating++
	defer func() { simulating-- }()
	if options.Limit <= 0 {
		options.Limit = defaultSituationLimit
	}
	sim := newSimulation(g)
	from := origin{base: capture(g)}
	if g.currentScene == "" {
		g.resolveDerived()
		sim.game.startWorld()
		sim.settle()
		from = origin{base: capture(sim.game)}
	}
	run, milestone := sim.milestone(from, options.Limit, options.Goal)
	situation := Situation{
		States:    len(run.nodes),
		Truncated: run.truncated,
		Frontier:  run.frontier(),
	}
	if milestone != nil {
		situation.Found = true
		situation.Milestone = milestone.milestone
		situation.Next = run.path(milestone.index)
	}
	return situation
}

// ----- the copy -----------------------------------------------------------

// simulation is a headless copy of a game: the same managers and content,
// its own runtime. ctx.Game is the copy, so every action runs against it.
type simulation struct {
	game *Game
	// suspended holds the runners a dialog choice interrupted: the
	// RunDialog that opened the conversation waits underneath the choice's
	// actions and resumes when they are done.
	suspended []Runner
}

// simulatedDT is the clock step of a simulated frame: long enough that any
// line, wait or fade is over by the next one.
const simulatedDT = 1e6

// settleLimit is how many simulated frames one move may take before it is
// reported as a runaway.
const settleLimit = 10000

func newSimulation(g *Game) *simulation {
	copied := *g
	copied.Audio = &AudioPlayer{silent: true}
	copied.Input = newInput()
	copied.Store = NewMemoryStore()
	copied.Extensions = nil
	copied.transition = &transition{}
	copied.chars = make(map[string]*runtimeChar)
	copied.exitHotspots = nil
	copied.triggerStates = nil
	copied.unlocked = make(map[string]bool)
	copied.queue = nil
	copied.scriptRunner = nil
	copied.scriptCtx = nil
	copied.dialog = nil
	copied.messages = nil
	copied.MaxLogLines = 1
	copied.pending = noChange
	copied.settings.Voice = false
	copied.State = NewState()
	copied.Inventory = NewInventory()
	copied.currentStage = worldStage(g)
	copied.started = true
	return &simulation{game: &copied}
}

// worldStage is the Stage the copy plays on: the world, whichever Stage
// that is, so that triggers fire and world-only widgets answer.
func worldStage(g *Game) *Stage {
	if g.currentStage != nil && g.currentStage.World {
		return g.currentStage
	}
	if g.StartStage != nil && g.StartStage.World {
		return g.StartStage
	}
	for _, s := range g.StageManager.All() {
		if s.World {
			return s
		}
	}
	return nil
}

// moment is a game between moves, with nothing running: what a search
// node starts from.
type moment struct {
	flags    map[string]bool
	vars     map[string]any
	visited  map[string]int
	talked   map[string]int
	known    map[string]int
	items    []string
	scene    string
	previous string
	triggers map[string]triggerState
	endCard  string
}

func capture(g *Game) *moment {
	triggers := make(map[string]triggerState, len(g.triggerStates))
	for name, state := range g.triggerStates {
		triggers[name] = *state
	}
	return &moment{
		flags:    g.State.Flags(),
		vars:     g.State.Vars(),
		visited:  g.State.VisitedMap(),
		talked:   g.State.TalkedMap(),
		known:    g.State.KnownMap(),
		items:    g.Inventory.Items(),
		scene:    g.currentScene,
		previous: g.previousScene,
		triggers: triggers,
		endCard:  g.endCard,
	}
}

func (s *simulation) restore(m *moment) {
	g := s.game
	g.State = NewState()
	g.State.LoadFlags(m.flags)
	g.State.LoadVars(m.vars)
	g.State.LoadVisited(m.visited)
	g.State.LoadTalked(m.talked)
	g.State.LoadKnown(m.known)
	g.Inventory = NewInventory()
	for _, name := range m.items {
		g.Inventory.Add(name)
	}
	g.currentScene = m.scene
	g.previousScene = m.previous
	if m.scene != "" {
		g.placeActors(m.scene)
	}
	g.triggerStates = make(map[string]*triggerState, len(m.triggers))
	for name, state := range m.triggers {
		copied := state
		g.triggerStates[name] = &copied
	}
	g.endCard = m.endCard
	g.dialog = nil
	g.queue = nil
	g.scriptRunner = nil
	g.scriptCtx = nil
	g.pending = noChange
	*g.transition = transition{}
	s.suspended = nil
}

// settle runs the frames a move sets off until the player is free to act
// again — nothing running and no trigger left to fire — or a dialog waits
// for a choice, or the end card is up. It reports false when the frames
// never stop.
func (s *simulation) settle() bool {
	g := s.game
	for frame := 0; frame < settleLimit; frame++ {
		for g.transition.active {
			g.transition.update(simulatedDT)
		}
		g.pumpQueue()
		if g.scriptRunner != nil {
			ctx := g.makeCtx()
			ctx.DT = simulatedDT
			status := g.scriptRunner.Tick(ctx)
			s.arrive()
			g.pending = noChange
			if status != StatusRunning {
				g.scriptRunner = nil
				g.scriptCtx = nil
				if n := len(s.suspended); n > 0 {
					g.scriptRunner = s.suspended[n-1]
					g.scriptCtx = g.makeCtx()
					s.suspended = s.suspended[:n-1]
				}
				continue
			}
			if g.dialogActive() || g.endCard != "" {
				return true
			}
			continue
		}
		if len(g.queue) > 0 || g.transition.active {
			continue
		}
		if g.endCard != "" {
			return true
		}
		g.tickTriggers()
		if len(g.queue) == 0 {
			return true
		}
	}
	return false
}

// arrive puts every walking character where it was going: a simulated
// walk takes no time.
func (s *simulation) arrive() {
	for _, c := range s.game.chars {
		if c.moving {
			c.pos = c.target
			c.path = nil
			c.moving = false
		}
	}
}

// panelHotspotter is a widget that mounts hotspots on the HUD — a
// HotspotPanel, or anything that embeds one.
type panelHotspotter interface {
	panelHotspots(g *Game) []Hotspot
}

// targets is every hotspot the player can reach right now: the scene's own,
// its exits, and those on every visible HotspotPanel.
func (s *simulation) targets() []Hotspot {
	g := s.game
	var out []Hotspot
	if g.currentScene != "" {
		out = append(out, g.SceneHotspots(g.currentScene)...)
	}
	ctx := g.makeCtx()
	for _, w := range g.WidgetManager.All() {
		panel, ok := w.(panelHotspotter)
		if !ok || !WidgetVisible(ctx, w) {
			continue
		}
		out = append(out, panel.panelHotspots(g)...)
	}
	return out
}

func (s *simulation) target(name string) *Hotspot {
	targets := s.targets()
	for i := range targets {
		if targets[i].Name == name {
			return &targets[i]
		}
	}
	return nil
}

// moves lists what the player can do now. Only what the game itself
// answers is listed: a verb with no handler of the hotspot's own and an
// item paired with nothing fall to generic replies, which change nothing.
func (s *simulation) moves() []Move {
	g := s.game
	if g.endCard != "" {
		return nil
	}
	if g.dialog != nil {
		return s.choices()
	}
	var held []*Item
	for _, name := range g.Inventory.Items() {
		if item, ok := g.ItemManager.Get(name); ok {
			held = append(held, item)
		}
	}
	verbs := g.VerbManager.Names()
	var out []Move
	for _, h := range s.targets() {
		for _, verb := range verbs {
			if h.handler(verb) == nil {
				continue
			}
			out = append(out, Move{
				Scene:   g.currentScene,
				Hotspot: h.Name,
				Label:   hotspotLabel(&h),
				Verb:    verb,
			})
		}
		for _, item := range held {
			if h.OnUseWith[item] == nil && h.OnUseWithAny == nil {
				continue
			}
			out = append(out, Move{
				Scene:   g.currentScene,
				Hotspot: h.Name,
				Label:   hotspotLabel(&h),
				Verb:    "use",
				Item:    item,
			})
		}
	}
	for _, item := range held {
		for _, other := range held {
			if item == other || (item.OnUseWith[other] == nil && other.OnUseWith[item] == nil) {
				continue
			}
			out = append(out, Move{
				Scene: g.currentScene,
				Label: g.ItemLabel(other.Name),
				Verb:  "use",
				Item:  item,
				With:  other,
			})
		}
	}
	return out
}

// choices is the open dialog's choices the player can see, filtered the
// way DialogBox filters them.
func (s *simulation) choices() []Move {
	g := s.game
	node := g.dialog.Node
	ctx := g.makeCtx()
	var out []Move
	for i, c := range node.Choices {
		if c.Show != nil && !c.Show.Eval(ctx) {
			continue
		}
		if c.Once && g.State.Talked(choiceKey(g.dialog.Dialog, node, i)) > 0 {
			continue
		}
		out = append(out, Move{
			Scene:  g.currentScene,
			Dialog: g.dialog.Dialog,
			Choice: i,
			Text:   c.Text,
		})
	}
	return out
}

// apply makes a move, the way the engine's own input paths would: the hotspot
// through InvokeVerb, an item on an item the way InventoryBar pairs them, a
// dialog choice the way DialogBox picks it. A choice's actions run straight
// away, over the RunDialog that is waiting for it.
func (s *simulation) apply(m Move) bool {
	g := s.game
	switch {
	case m.Dialog != nil:
		if g.dialog == nil || m.Choice >= len(g.dialog.Node.Choices) {
			return false
		}
		choice := g.dialog.Node.Choices[m.Choice]
		if choice.Once {
			g.State.NoteTalked(choiceKey(g.dialog.Dialog, g.dialog.Node, m.Choice))
		}
		if len(choice.Actions) > 0 {
			if g.scriptRunner != nil {
				s.suspended = append(s.suspended, g.scriptRunner)
			}
			g.scriptRunner = Seq(choice.Actions...).Start()
			g.scriptCtx = g.makeCtx()
		}
	case m.With != nil:
		a := m.Item.OnUseWith[m.With]
		if a == nil {
			a = m.With.OnUseWith[m.Item]
		}
		g.queueAction(a, "simulated useWith items")
	default:
		h := s.target(m.Hotspot)
		if h == nil {
			return false
		}
		selected := ""
		if m.Item != nil {
			selected = m.Item.Name
		}
		g.Inventory.Select(selected)
		g.InvokeVerb(h, m.Verb)
	}
	return true
}

// world is the moment without where the player stands: what a move that
// moves the story on changes.
func (s *simulation) world() string {
	g := s.game
	var b strings.Builder
	for _, name := range sortedKeys(g.State.Flags()) {
		if g.State.Flag(name) {
			b.WriteString(name)
			b.WriteString(",")
		}
	}
	b.WriteString("|")
	vars := g.State.Vars()
	for _, name := range sortedKeys(vars) {
		if vars[name] != nil {
			fmt.Fprintf(&b, "%s=%v,", name, vars[name])
		}
	}
	b.WriteString("|")
	items := g.Inventory.Items()
	sort.Strings(items)
	b.WriteString(strings.Join(items, ","))
	if g.endCard != "" {
		b.WriteString("|end")
	}
	return b.String()
}

// key is the moment as far as the search tells moments apart: the world,
// plus where the player stands and what the scene's triggers have seen.
func (s *simulation) key(world string) string {
	g := s.game
	var b strings.Builder
	b.WriteString(world)
	b.WriteString("|")
	b.WriteString(g.currentScene)
	b.WriteString("|")
	if s.backExit() {
		b.WriteString(g.previousScene)
	}
	b.WriteString("|")
	for _, name := range sortedKeys(g.triggerStates) {
		state := g.triggerStates[name]
		fmt.Fprintf(&b, "%s:%t:%t,", name, state.fired, state.lastTrue)
	}
	if g.dialog != nil {
		fmt.Fprintf(&b, "|dialog:%s/%s/%d", g.dialog.Dialog.Name, g.dialog.Node.Name, len(s.suspended))
	}
	return b.String()
}

// backExit reports whether the current scene has a way back the way the
// player came, which is the only place the previous scene matters.
func (s *simulation) backExit() bool {
	g := s.game
	if g.currentScene == "" {
		return false
	}
	for _, e := range g.SceneManager.MustGet(g.currentScene).Exits {
		if e.To == nil {
			return true
		}
	}
	return false
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ----- the search ---------------------------------------------------------

// origin is where a search starts: a quiet moment, and the dialog choices
// made since it, when the search starts inside a conversation.
type origin struct {
	base    *moment
	pending []Move
}

// searchNode is one moment the search reached. A moment inside an open
// dialog has a RunDialog waiting in it, which cannot be copied; such a node
// keeps the last quiet moment before it and the choices made since, and is
// reached again by replaying them.
type searchNode struct {
	index     int
	base      *moment
	pending   []Move
	key       string
	world     string
	scene     string
	flags     map[string]bool
	goal      bool
	milestone string
	parent    int
	move      Move
	edges     []int
}

type search struct {
	sim       *simulation
	goal      Condition
	nodes     []*searchNode
	truncated bool
	runaway   [][]Move
}

// visit is what a search does with a node it has just reached: look past
// it, stop there and look past the others, or stop altogether.
type visit int

const (
	expand visit = iota
	leaf
	halt
)

// raisedNew names the first flag n has raised that the origin had not, or
// marks n as the goal; it reports whether n is a milestone at all.
func raisedNew(n *searchNode, before map[string]bool) bool {
	if n.goal {
		return true
	}
	for _, name := range sortedKeys(n.flags) {
		if !before[name] {
			n.milestone = name
			return true
		}
	}
	return false
}

// milestone searches breadth first from an origin to the nearest milestone
// and returns the search and that node, nil when there is none within the
// limit.
func (s *simulation) milestone(from origin, limit int, goal Condition) (*search, *searchNode) {
	var found *searchNode
	var before map[string]bool
	run := s.search(from, limit, func(n *searchNode) visit {
		if before == nil {
			before = n.flags
			return expand
		}
		if raisedNew(n, before) {
			found = n
			return halt
		}
		return expand
	}, goal)
	return run, found
}

// candidates searches breadth first from an origin for every milestone that
// can be reached without passing another one first, nearest first, and
// returns the search and those nodes. Only the first node to raise a given
// flag is kept: the others raise the same flag the long way round.
func (s *simulation) candidates(from origin, limit int, goal Condition) (*search, []*searchNode) {
	var found []*searchNode
	var before map[string]bool
	kept := make(map[string]bool)
	run := s.search(from, limit, func(n *searchNode) visit {
		if before == nil {
			before = n.flags
			return expand
		}
		if !raisedNew(n, before) {
			return expand
		}
		if !kept[n.milestone] || n.goal {
			kept[n.milestone] = true
			found = append(found, n)
		}
		return leaf
	}, goal)
	return run, found
}

// search explores breadth first from an origin, asking visit about every
// node it reaches — the origin first — until the limit is reached or
// nothing is left.
func (s *simulation) search(from origin, limit int, see func(*searchNode) visit, goal Condition) *search {
	run := &search{sim: s, goal: goal}
	if limit <= 0 {
		limit = 1
	}
	root := &searchNode{
		base:    from.base,
		pending: from.pending,
		parent:  -1,
	}
	run.replay(root)
	run.describe(root)
	run.nodes = append(run.nodes, root)
	seen := map[string]int{root.key: 0}
	if see(root) == halt {
		return run
	}
	expandable := []bool{true}
	for next := 0; next < len(run.nodes); next++ {
		current := run.nodes[next]
		if current.goal || !expandable[next] {
			continue
		}
		run.replay(current)
		for _, m := range s.moves() {
			run.replay(current)
			if !s.apply(m) {
				continue
			}
			if !s.settle() {
				run.runaway = append(run.runaway, append(run.path(current.index), m))
				continue
			}
			child := &searchNode{
				parent: current.index,
				move:   m,
			}
			run.describe(child)
			if index, ok := seen[child.key]; ok {
				current.edges = append(current.edges, index)
				continue
			}
			if len(run.nodes) >= limit {
				run.truncated = true
				return run
			}
			if s.game.dialog != nil {
				child.base = current.base
				child.pending = append(append([]Move(nil), current.pending...), m)
			} else {
				child.base = capture(s.game)
			}
			child.index = len(run.nodes)
			seen[child.key] = child.index
			current.edges = append(current.edges, child.index)
			run.nodes = append(run.nodes, child)
			switch see(child) {
			case halt:
				return run
			case leaf:
				expandable = append(expandable, false)
			default:
				expandable = append(expandable, true)
			}
		}
	}
	return run
}

// describe fills in what a node is from the moment the simulation is at.
func (run *search) describe(n *searchNode) {
	s := run.sim
	g := s.game
	n.world = s.world()
	n.key = s.key(n.world)
	n.scene = g.currentScene
	n.flags = make(map[string]bool)
	for name, raised := range g.State.Flags() {
		if raised {
			n.flags[name] = true
		}
	}
	if run.goal != nil {
		n.goal = run.goal.Eval(g.makeCtx())
	} else {
		n.goal = g.endCard != ""
	}
}

// replay puts the simulation back at a node's moment.
func (run *search) replay(n *searchNode) {
	run.sim.restore(n.base)
	for _, m := range n.pending {
		run.sim.apply(m)
		run.sim.settle()
	}
}

func (run *search) path(index int) []Move {
	var out []Move
	for i := index; i > 0; i = run.nodes[i].parent {
		out = append(out, run.nodes[i].move)
	}
	for l, r := 0, len(out)-1; l < r; l, r = l+1, r-1 {
		out[l], out[r] = out[r], out[l]
	}
	return out
}

// frontier is every move from the origin that changed the world rather than
// only where the player stands.
func (run *search) frontier() []Move { return run.frontierOf(0) }

// firstGoal is the first goal node reached, which breadth first is the
// nearest; -1 when none was.
func (run *search) firstGoal() int {
	for i, n := range run.nodes {
		if n.goal {
			return i
		}
	}
	return -1
}

// reachingGoal marks every node from which some goal node can be reached.
func (run *search) reachingGoal() []bool {
	incoming := make([][]int, len(run.nodes))
	for i, n := range run.nodes {
		for _, to := range n.edges {
			incoming[to] = append(incoming[to], i)
		}
	}
	reaching := make([]bool, len(run.nodes))
	var stack []int
	for i, n := range run.nodes {
		if n.goal {
			reaching[i] = true
			stack = append(stack, i)
		}
	}
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, from := range incoming[top] {
			if !reaching[from] {
				reaching[from] = true
				stack = append(stack, from)
			}
		}
	}
	return reaching
}

// frontierOf is every move from a node that changed the world rather than
// only where the player stands.
func (run *search) frontierOf(index int) []Move {
	from := run.nodes[index]
	var out []Move
	for _, edge := range from.edges {
		if run.nodes[edge].world != from.world {
			out = append(out, run.nodes[edge].move)
		}
	}
	return out
}
