package pncdsl

import (
	"hash/fnv"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type AssetManager = Manager[Asset]

// loadedAssets is the runtime cache for decoded image/audio resources.
// The Asset structs in the manager remain immutable spec; this is where
// the actual *ebiten.Image bytes live, lazily decoded on first access.
type loadedAssets struct {
	images map[string]*ebiten.Image
	defaultW, defaultH int
}

func newLoadedAssets(w, h int) *loadedAssets {
	return &loadedAssets{
		images:   make(map[string]*ebiten.Image),
		defaultW: w,
		defaultH: h,
	}
}

func (la *loadedAssets) image(am *AssetManager, name string) *ebiten.Image {
	if img, ok := la.images[name]; ok {
		return img
	}
	a, ok := am.Get(name)
	if !ok {
		img := placeholderImage(name, la.defaultW, la.defaultH)
		la.images[name] = img
		return img
	}
	img := loadImageFile(a.Path)
	if img == nil {
		img = placeholderImage(name, la.defaultW, la.defaultH)
	}
	la.images[name] = img
	return img
}

func loadImageFile(path string) *ebiten.Image {
	f, err := os.Open(path)
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
