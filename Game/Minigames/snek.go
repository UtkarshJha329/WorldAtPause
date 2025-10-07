package Minigames

import (
	"WorldAtPause/Engine"
	"math"
)

type Snake struct {
	segmentPositions []Engine.Vector2Int
}

func (snake *Snake) InitSnakeWithSegments(numSegments int, headPosition Engine.Vector2Int, segmentsPositionDelta Engine.Vector2Int) {
	snake.segmentPositions = make([]Engine.Vector2Int, numSegments)

	snake.segmentPositions[0] = headPosition
	for i := 1; i < numSegments; i++ {
		curSegmentDelta := Engine.Multiply_Float_Vector2Int(i, &segmentsPositionDelta)
		snake.segmentPositions[i] = Engine.Add_Vector2Int(&headPosition, &curSegmentDelta)
	}
}

func (snake *Snake) GetHeadPos() Engine.Vector2Int {
	return snake.segmentPositions[0]
}

func (snake *Snake) GetTailPos() Engine.Vector2Int {
	return snake.segmentPositions[len(snake.segmentPositions)-1]
}

func (snake *Snake) GetTailPosAndDirectionFromBody() (Engine.Vector2Int, Engine.Vector2Int) {
	ultimateSegmentPos := snake.segmentPositions[len(snake.segmentPositions)-1]
	penUltimteSegmentPos := snake.segmentPositions[len(snake.segmentPositions)-2]
	directionFromPenultimateSegment := Engine.Subtract_Vector2Int(&ultimateSegmentPos, &penUltimteSegmentPos)
	return ultimateSegmentPos, directionFromPenultimateSegment
}

func (snake *Snake) SnakeFormsClosedLoop() bool {

	snakeHeadPos := snake.GetHeadPos()
	snakeTailPos := snake.GetTailPos()

	deltaPosBetweenHeadAndTail := Engine.Subtract_Vector2Int(&snakeHeadPos, &snakeTailPos)

	return math.Abs(float64(deltaPosBetweenHeadAndTail.X)) <= 1 && math.Abs(float64(deltaPosBetweenHeadAndTail.Y)) <= 1
}

func (snake *Snake) MoveHeadOnBoardWithBody(deltaMoveAmount Engine.Vector2Int, growing bool) bool {

	newHeadPos := Engine.Add_Vector2Int(&snake.segmentPositions[0], &deltaMoveAmount)
	for i, snakeSegmentPos := range snake.segmentPositions {
		if i == 0 || (!growing && i == len(snake.segmentPositions)-1) {
			continue
		}
		if snakeSegmentPos.X == newHeadPos.X && snakeSegmentPos.Y == newHeadPos.Y {
			return false
		}
	}

	if growing {
		snake.segmentPositions = append([]Engine.Vector2Int{newHeadPos}, snake.segmentPositions...)
	} else {
		copy(snake.segmentPositions[1:], snake.segmentPositions[:len(snake.segmentPositions)-1])
		snake.segmentPositions[0] = newHeadPos
	}

	return true
}

func (snake *Snake) GetBoundingBox() (topLeft, bottomRight Engine.Vector2Int) {
	leftMost := snake.GetHeadPos()
	rightMost := snake.GetHeadPos()
	topMost := snake.GetHeadPos()
	bottomMost := snake.GetHeadPos()

	for _, segmentPosition := range snake.segmentPositions {
		if segmentPosition.X < leftMost.X {
			leftMost = segmentPosition
		}
		if segmentPosition.X > rightMost.X {
			rightMost = segmentPosition
		}
		if segmentPosition.Y < topMost.Y {
			topMost = segmentPosition
		}
		if segmentPosition.Y > bottomMost.Y {
			bottomMost = segmentPosition
		}
	}

	topLeft = Engine.Vector2Int{X: leftMost.X, Y: topMost.Y}
	bottomRight = Engine.Vector2Int{X: rightMost.X, Y: bottomMost.Y}

	return topLeft, bottomRight
}
