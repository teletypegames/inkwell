package inkwell

type flagCond struct{ name string }

// Flag holds while a named world-state flag is raised.
func Flag(name string) Condition       { return &flagCond{name: name} }
func (f *flagCond) Eval(ctx *Ctx) bool { return ctx.Game.State.Flag(f.name) }
