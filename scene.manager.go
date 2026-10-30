package inkwell

// SceneManager is the registry of a game's Scenes, keyed by name; see
// [Manager].
type SceneManager = Manager[*Scene]
