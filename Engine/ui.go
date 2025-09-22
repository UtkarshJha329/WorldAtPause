package Engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UIRect struct {
	Sprite   *Sprite
	Position Vector2
	Size     Vector2

	ChildrenIndexInUITree []int
}

func (uiRect *UIRect) IsPointInRect(point Vector2) bool {
	return point.X >= uiRect.Position.X &&
		point.Y >= uiRect.Position.Y &&
		point.X <= uiRect.Position.X+uiRect.Size.X &&
		point.Y <= uiRect.Position.Y+uiRect.Size.Y
}

type UITree struct {
	UIRects []UIRect
}

func (uiTree *UITree) RenderUITree(fromUIRectIndex int, positionOffset Vector2, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions) {

	curUIRectRef := uiTree.UIRects[fromUIRectIndex]
	curUIRectDrawPos := Add_Vector2(&positionOffset, &curUIRectRef.Position)

	curUIRectRef.Sprite.DrawSprite(screenRef, drawImgOptionsRef, &curUIRectDrawPos)

	for _, childIndex := range curUIRectRef.ChildrenIndexInUITree {
		uiTree.RenderUITree(childIndex, curUIRectDrawPos, screenRef, drawImgOptionsRef)
	}
}

func (uiTree *UITree) PointInUITree(fromUIRectIndex int, point Vector2) (int, bool) {

	curUIRectRef := uiTree.UIRects[fromUIRectIndex]

	if curUIRectRef.IsPointInRect(point) {

		currentPointInRect := fromUIRectIndex

		for _, childIndex := range curUIRectRef.ChildrenIndexInUITree {
			relativePosOfPoint := Subtract_Vector2(&point, &curUIRectRef.Position)
			pointInRectID, pointInRect := uiTree.PointInUITree(childIndex, relativePosOfPoint)

			if pointInRect {
				currentPointInRect = pointInRectID
				break
			}
		}

		return currentPointInRect, true
	}

	return -1, false
}

func (uiTree *UITree) PointOnThisRectInUITree(fromUIRectIndex int, uiRectIndexToCheck int, point Vector2) bool {
	rectIndex, inARect := uiTree.PointInUITree(fromUIRectIndex, point)
	return inARect && rectIndex == uiRectIndexToCheck
}
