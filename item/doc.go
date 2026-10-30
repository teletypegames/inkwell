// Package item holds the [Item] type and its [Manager]: the things the
// player can carry.
//
// An item has a Name to register and look it up by, a Sprite for the
// inventory and the cursor, a Description the UI shows for it, and Lore, a
// paragraph for a dossier screen once the item is known. OnUseSelf answers
// the item used on its own; OnUseWith maps another *Item to what happens when
// the two are used together — with one in hand, clicking the other in the
// inventory runs whichever of the two pairs them. An item used on a hotspot
// is the hotspot's business (its OnUseWith and OnUseWithAny), with the
// game's UseWithFail as the answer nobody wrote.
//
// The inventory itself remembers items by name and lives on the Game.
package item
