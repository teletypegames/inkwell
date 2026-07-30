package inkwell

import (
	"io/fs"
	"os"
)

// AssetFS, when non-nil, is the filesystem asset files (images, audio)
// are read from instead of the OS filesystem. Set it to an embed.FS
// before Run to ship assets inside the binary — required for js/wasm
// builds, where there is no OS filesystem and every asset would fall
// back to a placeholder:
//
//	//go:embed assets
//	var assetsFS embed.FS
//
//	func init() { inkwell.AssetFS = assetsFS }
//
// Asset.Path values are used as-is, so they must be slash-separated
// and relative (e.g. "assets/bg/bedroom.png") — the same form that
// works with os.Open from the project root.
var AssetFS fs.FS

// openAsset opens an asset file honoring AssetFS. Every asset read in
// the package goes through here.
func openAsset(path string) (fs.File, error) {
	if AssetFS != nil {
		return AssetFS.Open(path)
	}
	return os.Open(path)
}
