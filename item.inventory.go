package inkwell

// Inventory is what the player carries: item names in the order they were
// picked up, and the one selected for use on something, if any.
type Inventory struct {
	items    []string
	selected string
}

// NewInventory returns an empty Inventory with nothing selected.
func NewInventory() *Inventory { return &Inventory{} }

// Add puts the named item at the end of the inventory; an item already
// carried is not added twice.
func (i *Inventory) Add(name string) {
	if i.Has(name) {
		return
	}
	i.items = append(i.items, name)
}

// Remove takes the named item out of the inventory, and clears the
// selection if it was the selected one. An item not carried is ignored.
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

// Has reports whether the named item is carried.
func (i *Inventory) Has(name string) bool {
	for _, n := range i.items {
		if n == name {
			return true
		}
	}
	return false
}

// Select makes the named item the one in hand. An empty name clears the
// selection; an item not carried leaves it as it was.
func (i *Inventory) Select(name string) {
	if name != "" && !i.Has(name) {
		return
	}
	i.selected = name
}

// Selected returns the name of the item in hand, or "" for none.
func (i *Inventory) Selected() string { return i.selected }

// Items returns a copy of the carried item names, in pickup order.
func (i *Inventory) Items() []string { return append([]string(nil), i.items...) }
