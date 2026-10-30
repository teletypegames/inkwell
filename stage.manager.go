package inkwell

// StageManager is the registry of a game's Stages, keyed by name; see
// [Manager].
type StageManager = Manager[*Stage]
