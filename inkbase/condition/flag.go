package condition

import "git.teletypegames.org/engines/inkwell/inkbase"

type flagCond struct{ name string }

// Flag holds while a named world-state flag is raised.
func Flag(name string) inkbase.Condition       { return &flagCond{name: name} }
func (f *flagCond) Eval(ctx *inkbase.Ctx) bool { return ctx.State.Flag(f.name) }
