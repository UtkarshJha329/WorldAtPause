package Engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var YesNoUITree *UITree
var yesButtonUIRectIndex int
var noButtonUIRectIndex int

func (world *World) MakeYesNoUITree() *UITree {

	YesNoUITree = &UITree{}

	backgroundUIRect := UIRect{
		Sprite:   world.UIRectSprites["UI Test Rect Sprite"],
		Position: Vector2{X: 128, Y: 112},
		Size:     Vector2{X: 64, Y: 16},
	}

	YesNoUITree.UIRects = append(YesNoUITree.UIRects, backgroundUIRect)

	yesUIRect := UIRect{
		Sprite:   world.UIRectSprites["Yes Button Sprite"],
		Position: Vector2{X: 0, Y: 0},
		Size:     Vector2{X: 32, Y: 16},
	}

	YesNoUITree.UIRects = append(YesNoUITree.UIRects, yesUIRect)
	yesButtonUIRectIndex = len(YesNoUITree.UIRects) - 1
	YesNoUITree.UIRects[0].ChildrenIndexInUITree = append(YesNoUITree.UIRects[0].ChildrenIndexInUITree, yesButtonUIRectIndex)

	noUIRect := UIRect{
		Sprite:   world.UIRectSprites["No Button Sprite"],
		Position: Vector2{X: 32, Y: 0},
		Size:     Vector2{X: 32, Y: 16},
	}

	YesNoUITree.UIRects = append(YesNoUITree.UIRects, noUIRect)
	noButtonUIRectIndex = len(YesNoUITree.UIRects) - 1
	YesNoUITree.UIRects[0].ChildrenIndexInUITree = append(YesNoUITree.UIRects[0].ChildrenIndexInUITree, noButtonUIRectIndex)

	return YesNoUITree
}

func ClickedYesInYesNoUITree() bool {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		mousePos := Vector2{X: float64(mouseX), Y: float64(mouseY)}

		return YesNoUITree.PointOnThisRectInUITree(0, yesButtonUIRectIndex, mousePos)
	}

	return false
}

func ClickedNoInYesNoUITree() bool {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		mouseX, mouseY := ebiten.CursorPosition()
		mousePos := Vector2{X: float64(mouseX), Y: float64(mouseY)}

		return YesNoUITree.PointOnThisRectInUITree(0, noButtonUIRectIndex, mousePos)
	}

	return false
}
