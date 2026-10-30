package inkwell

// CharacterManager is the registry of a game's Characters, keyed by name;
// see [Manager].
type CharacterManager = Manager[*Character]
