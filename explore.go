package inkwell

import (
	"fmt"
	"image/color"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

// ExplorerRoot is the page tree of the built-in entity explorer: the Game
// itself, field by field. A manager opens onto its entities, an entity onto
// its literal, an action onto the actions it is made of, a Text onto its
// translations — as deep as the game was written. It is data rather than a
// widget, like DebugMenuRoot, which carries it as its last page; Explore
// runs a game with nothing but this on screen.
//
// Every row is worked out by reflection when its page opens and rebuilt
// while the page is up, so a condition's row says what it evaluates to
// right now and a manager's count is the count as it stands. Unexported
// fields are read too — an action's speaker and line live in unexported
// fields of an unexported type, and a tool that could not see inside a Say
// would stop exactly where the writing starts.
var ExplorerRoot = MenuItem{
	Label:  explorerTitle,
	Expand: exploreGame,
}

func exploreGame(ctx *Ctx) []MenuItem {
	return explorePage(reflect.ValueOf(ctx.Game).Elem())
}

// explorePage is the rows a composite value opens onto: a struct's fields, a
// slice's elements, a map's entries by key, a manager's entities, a Text's
// translations, and for an action the action itself to run first.
func explorePage(value reflect.Value) []MenuItem {
	if all := exploreAll(value); all.IsValid() {
		return explorePage(all)
	}
	if text, ok := exploreAs[Text](value); ok {
		return exploreTranslations(text)
	}
	var rows []MenuItem
	if action, ok := exploreAs[Action](value); ok {
		rows = append(rows, MenuItem{
			Label: explorerRun,
			Do:    action,
		})
	}
	value = exploreUnwrap(value)
	if !value.IsValid() {
		return rows
	}
	switch value.Kind() {
	case reflect.Struct:
		value = exploreAddressable(value)
		for i := 0; i < value.NumField(); i++ {
			rows = append(rows, exploreRow(value.Type().Field(i).Name, exploreField(value, i)))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			element := value.Index(i)
			label := "[" + strconv.Itoa(i) + "]"
			if name := exploreName(element); name != "" {
				label = name
			}
			rows = append(rows, exploreRow(label, element))
		}
	case reflect.Map:
		keys := value.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return exploreSummary(nil, keys[i]) < exploreSummary(nil, keys[j])
		})
		for _, key := range keys {
			rows = append(rows, exploreRow(exploreSummary(nil, key), value.MapIndex(key)))
		}
	}
	return rows
}

func exploreTranslations(text Text) []MenuItem {
	rows := make([]MenuItem, 0, len(text))
	for _, translation := range text {
		row := MenuItem{
			Label: plain(translation.Lang),
			Value: FixedValue(plain(translation.Text)),
		}
		if translation.Voice != nil {
			voice := reflect.ValueOf(translation.Voice)
			row.Expand = func(*Ctx) []MenuItem { return explorePage(voice) }
		}
		rows = append(rows, row)
	}
	return rows
}

// exploreRow is one row for one value: its label, what it is worth now in
// the value column, and the page it opens onto if it is anything more than
// a scalar.
func exploreRow(label string, value reflect.Value) MenuItem {
	row := MenuItem{
		Label: plain(label),
		Value: func(ctx *Ctx) Text { return plain(exploreSummary(ctx, value)) },
	}
	if exploreOpens(value) {
		row.Expand = func(*Ctx) []MenuItem { return explorePage(value) }
	}
	return row
}

// exploreOpens reports whether a value has a page behind it. A nil, a
// scalar, a function and a colour do not; a struct, a non-empty slice or
// map, and a manager do.
func exploreOpens(value reflect.Value) bool {
	if exploreAll(value).IsValid() {
		return true
	}
	if _, ok := exploreAs[color.Color](value); ok {
		return false
	}
	value = exploreUnwrap(value)
	if !value.IsValid() {
		return false
	}
	switch value.Kind() {
	case reflect.Struct:
		return true
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() > 0
	}
	return false
}

// exploreSummary is the value column: a name for what has one, the current
// language's line for a Text, the answer for a Condition, the kind of action
// for an Action, a hex triplet for a colour, a count for a collection, and
// the value itself for a scalar. The context is nil where no Condition can
// be — a map key — and a Condition met there says only what it is.
func exploreSummary(ctx *Ctx, value reflect.Value) string {
	if !value.IsValid() {
		return "nil"
	}
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Func, reflect.Map, reflect.Slice:
		if value.IsNil() {
			return "nil"
		}
	}
	if all := exploreAll(value); all.IsValid() {
		return strconv.Itoa(all.Len())
	}
	if text, ok := exploreAs[Text](value); ok {
		return text.String()
	}
	if named, ok := exploreAs[Named](value); ok {
		return named.GetName()
	}
	if condition, ok := exploreAs[Condition](value); ok {
		if ctx == nil {
			return exploreKind(value)
		}
		return exploreKind(value) + " = " + strconv.FormatBool(condition.Eval(ctx))
	}
	if _, ok := exploreAs[Action](value); ok {
		return exploreKind(value)
	}
	if colour, ok := exploreAs[color.Color](value); ok {
		red, green, blue, alpha := colour.RGBA()
		return fmt.Sprintf("#%02X%02X%02X%02X", red>>8, green>>8, blue>>8, alpha>>8)
	}
	value = exploreUnwrap(value)
	switch value.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, 64)
	case reflect.String:
		return value.String()
	case reflect.Slice, reflect.Array, reflect.Map:
		return strconv.Itoa(value.Len())
	case reflect.Func:
		return "func"
	case reflect.Struct:
		if name := exploreName(value); name != "" {
			return name
		}
		return exploreKind(value)
	}
	return value.Type().String()
}

// exploreName is the name a value carries, if it carries one: what GetName
// answers, or a Name field.
func exploreName(value reflect.Value) string {
	if named, ok := exploreAs[Named](value); ok {
		return named.GetName()
	}
	value = exploreUnwrap(value)
	if value.IsValid() && value.Kind() == reflect.Struct {
		if field := value.FieldByName("Name"); field.IsValid() && field.Kind() == reflect.String {
			return field.String()
		}
	}
	return ""
}

// exploreKind is a value's type as an author knows it: "Say" for a
// sayAction, "Flag" for a flagCond, "Given" for a Given.
func exploreKind(value reflect.Value) string {
	value = exploreUnwrap(value)
	if !value.IsValid() {
		return "nil"
	}
	name := value.Type().String()
	if dot := strings.LastIndex(name, "."); dot >= 0 {
		name = name[dot+1:]
	}
	for _, suffix := range []string{"Action", "Cond"} {
		if trimmed := strings.TrimSuffix(name, suffix); trimmed != name && trimmed != "" {
			name = trimmed
			break
		}
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

// exploreAll is the entities of a manager as a slice, or an invalid value
// for anything that is not a manager.
func exploreAll(value reflect.Value) reflect.Value {
	if !value.IsValid() || value.Kind() != reflect.Ptr || value.IsNil() {
		return reflect.Value{}
	}
	if !strings.HasPrefix(value.Type().Elem().Name(), "Manager[") {
		return reflect.Value{}
	}
	all := value.MethodByName("All")
	if !all.IsValid() {
		return reflect.Value{}
	}
	return all.Call(nil)[0]
}

// exploreAs is the value as T when it is one, through pointers and
// interfaces, and whether it was.
func exploreAs[T any](value reflect.Value) (T, bool) {
	var zero T
	for value.IsValid() {
		if value.CanInterface() {
			if typed, ok := value.Interface().(T); ok {
				return typed, true
			}
		}
		if (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) && !value.IsNil() {
			value = value.Elem()
			continue
		}
		return zero, false
	}
	return zero, false
}

// exploreUnwrap follows pointers and interfaces to the value they hold. A
// nil ends as an invalid value.
func exploreUnwrap(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}

// exploreAddressable is the value where it stands, or a copy of it in fresh
// memory when it stands nowhere — an element read out of a map, or the
// struct behind an interface — so that its fields can be addressed.
func exploreAddressable(value reflect.Value) reflect.Value {
	if value.CanAddr() {
		return value
	}
	copied := reflect.New(value.Type()).Elem()
	copied.Set(value)
	return copied
}

// exploreField is a struct's field whether or not it is exported. An
// unexported field comes back read-only from reflect; the same memory read
// through a fresh pointer does not, which is the whole of what a tool that
// only ever reads needs.
func exploreField(structure reflect.Value, index int) reflect.Value {
	field := structure.Field(index)
	if field.CanInterface() {
		return field
	}
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}

// Explore runs a game with nothing on screen but the explorer: a Stage that
// is a Sticky Menu over the whole window, rooted at ExplorerRoot. The game
// is validated and wired exactly as Run wires it, so what the explorer shows
// is the game as it would have run.
func Explore(g *Game) error {
	explorer := &Stage{
		Name: "explorer",
	}
	g.StageManager.Register(explorer)
	g.WidgetManager.Register(&Menu{
		Name:   "explorer",
		When:   OnStage(explorer),
		Root:   ExplorerRoot,
		Bounds: Rect(0, 0, float64(g.Width), float64(g.Height)),
		Sticky: true,
	})
	g.StartStage = explorer
	return Run(g)
}

var (
	explorerTitle = Text{
		{Lang: "en", Text: "EXPLORER"},
		{Lang: "hu", Text: "EXPLORER"},
	}
	explorerRun = Text{
		{Lang: "en", Text: "run"},
		{Lang: "hu", Text: "futtasd"},
	}
)
