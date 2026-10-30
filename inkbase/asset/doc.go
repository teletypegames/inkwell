// Package asset provides the filesystem abstraction, the image cache and the
// audio player for inkbase games.
//
// # Files
//
// Every file is opened through [OpenAsset], which reads from [AssetFS] — an
// embed.FS set before the game runs, to ship the assets inside the binary —
// or from the operating system's filesystem while AssetFS is nil.
//
// # Images
//
// An [inkbase.Asset] is only a spec, a Name and a Path. [LoadedAssets] opens
// the file the first time a picture is asked for, decodes it (PNG and JPEG)
// and caches it by Asset pointer. A missing or unreadable file falls back to
// a deterministic placeholder colour derived from the name
// ([PlaceholderImage]) and is marked as such, so the game runs without any
// art — useful for prototyping and CI — and [LoadedAssets.IsPlaceholder]
// lets a widget draw a stylised figure instead of a flat rectangle.
//
// # Audio
//
// [AudioPlayer] plays the three things a game is heard through: a Music
// track (streamed from disk, looping when asked, unchanged when the same
// track is asked for again), a recorded Voice line (whose duration it
// reports, so speech can last exactly as long) and a one-shot Sound (over
// everything else). The shared Ebitengine audio context is made lazily at
// 48 kHz, or taken over when a host process already made one. Voices and
// sounds are decoded once and cached. WAV, Ogg Vorbis and MP3 are read,
// picked by file extension. Every failure — a missing file, an unsupported
// codec, no audio device — degrades to a log line and silence, and a file
// that failed is not tried again.
package asset
