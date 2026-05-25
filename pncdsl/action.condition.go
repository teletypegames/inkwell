package pncdsl

type Condition interface {
	Eval(ctx *Ctx) bool
}

// ----- combinators ------------------------------------------------------

type notCond struct{ c Condition }

func Not(c Condition) Condition       { return &notCond{c: c} }
func (n *notCond) Eval(ctx *Ctx) bool { return !n.c.Eval(ctx) }

type andCond struct{ cs []Condition }

func And(cs ...Condition) Condition { return &andCond{cs: cs} }
func (a *andCond) Eval(ctx *Ctx) bool {
	for _, c := range a.cs {
		if !c.Eval(ctx) {
			return false
		}
	}
	return true
}

type orCond struct{ cs []Condition }

func Or(cs ...Condition) Condition { return &orCond{cs: cs} }
func (o *orCond) Eval(ctx *Ctx) bool {
	for _, c := range o.cs {
		if c.Eval(ctx) {
			return true
		}
	}
	return false
}

// ----- state predicates -------------------------------------------------

type flagCond struct{ name string }

func Flag(name string) Condition         { return &flagCond{name: name} }
func (f *flagCond) Eval(ctx *Ctx) bool   { return ctx.Game.State.Flag(f.name) }

type hasItemCond struct{ name string }

func HasItem(name string) Condition         { return &hasItemCond{name: name} }
func (h *hasItemCond) Eval(ctx *Ctx) bool   { return ctx.Game.Inventory.Has(h.name) }

type selectedItemCond struct{ name string }

func SelectedItem(name string) Condition       { return &selectedItemCond{name: name} }
func (s *selectedItemCond) Eval(ctx *Ctx) bool { return ctx.Game.Inventory.Selected() == s.name }

type inSceneCond struct{ name string }

func InScene(name string) Condition       { return &inSceneCond{name: name} }
func (i *inSceneCond) Eval(ctx *Ctx) bool { return ctx.Scene != nil && ctx.Scene.Name == i.name }

type varEqCond struct {
	name string
	v    any
}

func VarEq(name string, v any) Condition { return &varEqCond{name: name, v: v} }
func (e *varEqCond) Eval(ctx *Ctx) bool   { return ctx.Game.State.Var(e.name) == e.v }
