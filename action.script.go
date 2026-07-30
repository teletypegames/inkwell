package inkwell

type Script struct {
	Name    string
	Actions Action
}

func (s Script) GetName() string   { return s.Name }
func (s Script) TypeLabel() string { return "script" }
