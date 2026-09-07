// Package asset provides the filesystem abstraction and image cache for
// inkbase. AssetFS must be set before use; images are loaded lazily and
// cached by Asset pointer.
package asset

import (
	"hash/fnv"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"git.teletypegames.org/engines/inkwell/inkbase"
)

// AssetFS, when non-nil, is the filesystem asset files are read from instead
// of the OS filesystem. Set it to an embed.FS before Run.
var AssetFS fs.FS

// OpenAsset opens an asset file honoring AssetFS.
func OpenAsset(path string) (fs.File, error) {
	if AssetFS != nil {
		return AssetFS.Open(path)
	}
	return os.Open(path)
}

// LoadedAssets is the runtime cache for decoded image resources.
type LoadedAssets struct {
	images             map[*inkbase.Asset]*ebiten.Image
	placeholders       map[*inkbase.Asset]bool
	missing            *ebiten.Image
	defaultW, defaultH int
}

func NewLoadedAssets(w, h int) *LoadedAssets {
	return &LoadedAssets{
		images:       make(map[*inkbase.Asset]*ebiten.Image),
		placeholders: make(map[*inkbase.Asset]bool),
		defaultW:     w,
		defaultH:     h,
	}
}

// Image loads an asset's picture once and hands back the same *ebiten.Image
// after. An asset that cannot be read gets a flat placeholder.
func (la *LoadedAssets) Image(a *inkbase.Asset) *ebiten.Image {
	if a == nil {
		if la.missing == nil {
			la.missing = PlaceholderImage("", la.defaultW, la.defaultH)
		}
		return la.missing
	}
	if img, ok := la.images[a]; ok {
		return img
	}
	img := LoadImageFile(a.Path)
	if img == nil {
		img = PlaceholderImage(a.Name+a.Path, la.defaultW, la.defaultH)
		la.placeholders[a] = true
	}
	la.images[a] = img
	return img
}

// IsPlaceholder reports whether the named image fell back to a generated
// rectangle.
func (la *LoadedAssets) IsPlaceholder(a *inkbase.Asset) bool {
	return a == nil || la.placeholders[a]
}

func LoadImageFile(path string) *ebiten.Image {
	f, err := OpenAsset(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(src)
}

// PlaceholderImage creates a deterministic colored rectangle for missing
// assets.
func PlaceholderImage(seed string, w, h int) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	h32 := fnv.New32a()
	_, _ = h32.Write([]byte(seed))
	sum := h32.Sum32()
	c := color.RGBA{
		R: 60 + uint8(sum&0x7F),
		G: 60 + uint8((sum>>8)&0x7F),
		B: 60 + uint8((sum>>16)&0x7F),
		A: 255,
	}
	img.Fill(c)
	return img
}
