package inkwell

import (
	"hash/fnv"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type AssetManager = Manager[*Asset]

// loadedAssets is the runtime cache for decoded image/audio resources.
// The Asset structs in the manager remain immutable spec; this is where
// the actual *ebiten.Image bytes live, lazily decoded on first access.
type loadedAssets struct {
	images             map[*Asset]*ebiten.Image
	placeholders       map[*Asset]bool
	missing            *ebiten.Image
	defaultW, defaultH int
}

func newLoadedAssets(w, h int) *loadedAssets {
	return &loadedAssets{
		images:       make(map[*Asset]*ebiten.Image),
		placeholders: make(map[*Asset]bool),
		defaultW:     w,
		defaultH:     h,
	}
}

// image loads an asset's picture once and hands back the same *ebiten.Image
// after. An asset that cannot be read, or a nil one, gets a flat placeholder
// coloured from its name, so a missing file shows rather than crashes.
func (la *loadedAssets) image(a *Asset) *ebiten.Image {
	if a == nil {
		if la.missing == nil {
			la.missing = placeholderImage("", la.defaultW, la.defaultH)
		}
		return la.missing
	}
	if img, ok := la.images[a]; ok {
		return img
	}
	img := loadImageFile(a.Path)
	if img == nil {
		img = placeholderImage(a.Name+a.Path, la.defaultW, la.defaultH)
		la.placeholders[a] = true
	}
	la.images[a] = img
	return img
}

func (la *loadedAssets) isPlaceholder(a *Asset) bool {
	return a == nil || la.placeholders[a]
}

func loadImageFile(path string) *ebiten.Image {
	f, err := openAsset(path)
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

// placeholderImage creates a deterministic colored rectangle for missing
// assets so the game still runs without art on disk.
func placeholderImage(seed string, w, h int) *ebiten.Image {
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
