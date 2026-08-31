package inkwell

type AssetKind int

const (
	AssetImage AssetKind = iota
	AssetAudio
)

type Asset struct {
	Name string
	Path string
	Kind AssetKind
}

func (a Asset) GetName() string   { return a.Name }
func (a Asset) TypeLabel() string { return "asset" }
