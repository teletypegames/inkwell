package inkwell

import iwdialog "git.teletypegames.org/engines/inkwell/dialog"

// Dialog is a conversation tree: named nodes and the node it opens on; see
// [iwdialog.Dialog].
type Dialog = iwdialog.Dialog

// DialogNode is one step of a Dialog: the lines spoken on reaching it and the
// choices offered after them; see [iwdialog.DialogNode].
type DialogNode = iwdialog.DialogNode

// DialogLine is one line of a DialogNode and the Talker who says it; see
// [iwdialog.DialogLine].
type DialogLine = iwdialog.DialogLine

// DialogChoice is an answer the player can pick in a DialogNode, with the
// Condition that offers it and the actions it runs; see
// [iwdialog.DialogChoice].
type DialogChoice = iwdialog.DialogChoice
