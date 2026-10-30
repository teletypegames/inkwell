// Package sprite holds the picture types characters and items are drawn
// with.
//
// A [Sprite] is a picture or a sequence of them: a sheet and the frames cut
// from it, in order, at FrameTime seconds each, looping or holding the last.
// No frames means the whole sheet is the one frame, so a still item and a
// walking character are the same type.
//
// A character has one Sprite per [Pose] the art covers, and the engine picks
// the pose: PoseWalk while it moves, PoseTalk while its line is up, PoseWait
// once it has stood still for [WaitAfter] seconds, PoseIdle otherwise — or
// the one SetPose holds until the character next walks or speaks. A pose the
// art does not cover falls back to PoseIdle.
//
// Direction is not a second sheet. A Sprite says which way its artwork
// [Facing] looks, the character remembers the way it last moved, and the
// frame is mirrored when the two differ, so one walk cycle serves both
// directions.
package sprite
