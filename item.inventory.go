package inkwell

type Inventory struct {
	items    []string
	selected string
}

func NewInventory() *Inventory { return &Inventory{} }

func (i *Inventory) Add(name string) {
	if i.Has(name) {
		return
	}
	i.items = append(i.items, name)
}

func (i *Inventory) Remove(name string) {
	for k, n := range i.items {
		if n == name {
			i.items = append(i.items[:k], i.items[k+1:]...)
			if i.selected == name {
				i.selected = ""
			}
			return
		}
	}
}

func (i *Inventory) Has(name string) bool {
	for _, n := range i.items {
		if n == name {
			return true
		}
	}
	return false
}

func (i *Inventory) Select(name string) {
	if name != "" && !i.Has(name) {
		return
	}
	i.selected = name
}

func (i *Inventory) Selected() string { return i.selected }
func (i *Inventory) Items() []string  { return append([]string(nil), i.items...) }
