package inkwell

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

// SceneLayer is one plane of a scene's picture: a full-screen image with
// transparency where the planes behind it show through, at a depth Z. The
// cast stands at 0, so a layer with a negative Z is behind every character
// and one with a positive Z is in front of every character; the further
// from 0, the further from the cast. The background is the layer at -1, a
// foreground the layer at 1. Two layers at one Z keep the order they were
// written in.
type SceneLayer struct {
	Asset *Asset
	Z     int
}

// sortedLayers is a scene's layers back to front.
func sortedLayers(s *Scene) []SceneLayer {
	out := slices.Clone(s.Layers)
	slices.SortStableFunc(out, func(a, b SceneLayer) int {
		return a.Z - b.Z
	})
	return out
}

// drawScene paints the picture and its cast: the layers behind the
// characters, then the characters back to front by Y, then the layers in
// front. The lowest layer is the picture, and a missing file there is the
// usual named placeholder; a missing layer above it is a missing piece of
// transparency, which is nothing, so it is skipped rather than painted as a
// plate over the scene.
func drawScene(dst *ebiten.Image, g *Game, s *Scene) {
	layers := sortedLayers(s)
	next := 0
	paint := func() {
		layer := layers[next]
		g.loaded.image(layer.Asset)
		if next == 0 || !g.loaded.isPlaceholder(layer.Asset) {
			drawFullScreen(dst, g, layer.Asset)
		}
		next++
	}
	for next < len(layers) && layers[next].Z <= 0 {
		paint()
	}
	for _, c := range sortedChars(g) {
		drawCharacter(dst, g, c)
	}
	for next < len(layers) {
		paint()
	}
}
