package inkwell

import "git.teletypegames.org/engines/inkwell/inkbase"

// Point is a position in screen pixels; see [inkbase.Point].
type Point = inkbase.Point

// Rectangle is an axis-aligned box at X, Y of W by H pixels, and a Shape;
// see [inkbase.Rectangle].
type Rectangle = inkbase.Rectangle

// Polygon is a closed outline through its Points, and a Shape — a hotspot
// area or a walkbox; see [inkbase.Polygon].
type Polygon = inkbase.Polygon

// Shape is an area a click can land in: a Rectangle or a Polygon; see
// [inkbase.Shape].
type Shape = inkbase.Shape

// Rect returns the Rectangle at x, y of w by h; see [inkbase.Rect].
var Rect = inkbase.Rect

// Poly returns the Polygon through the given points; see [inkbase.Poly].
var Poly = inkbase.Poly

// PolyRect returns the rectangle x, y, w, h as a four-cornered Polygon, for
// a walkbox that is a plain box; see [inkbase.PolyRect].
var PolyRect = inkbase.PolyRect
