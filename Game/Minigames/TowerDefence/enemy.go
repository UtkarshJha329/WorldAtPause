package TowerDefence

import (
	"WorldAtPause/Engine"
)

type Enemy struct {
	Position                  Engine.Vector2
	FollowingCurrentPathPoint int
	Health                    int
	MoveSpeedPerFrame         float64
}

func (curEnemy *Enemy) MoveEnemyToNextPathPoint(enemyPathPoints *[]Engine.Vector2Int, occupationBoard *Engine.Board[int]) bool {

	curFollowingPathPoint := (*enemyPathPoints)[curEnemy.FollowingCurrentPathPoint]
	tilePosOnScreen := occupationBoard.GetPositionOfTileOnScreen(curFollowingPathPoint)
	curFollowingPathPointPos := Engine.Vector2{X: float64(tilePosOnScreen.X), Y: float64(tilePosOnScreen.Y)}

	distToCurPathPoint := Engine.Distance_Vector2(&curFollowingPathPointPos, &curEnemy.Position)

	if distToCurPathPoint >= float64(occupationBoard.TileGridTileSizeInPixels*2) {
		curEnemy.Position.X = curFollowingPathPointPos.X
		curEnemy.Position.Y = curFollowingPathPointPos.Y
	} else if distToCurPathPoint <= float64(occupationBoard.TileGridTileSizeInPixels*2) && distToCurPathPoint >= 0.1 {
		dirToCurPathPoint := Engine.Subtract_Vector2(&curFollowingPathPointPos, &curEnemy.Position)
		dirToCurPathPoint = Engine.Normalise_Vector2(&dirToCurPathPoint)

		totalMoveAmount := Engine.Multiply_Float_Vector2(curEnemy.MoveSpeedPerFrame, &dirToCurPathPoint)
		newPosition := Engine.Add_Vector2(&curEnemy.Position, &totalMoveAmount)

		curEnemy.Position.X = newPosition.X
		curEnemy.Position.Y = newPosition.Y
	} else {
		curEnemy.FollowingCurrentPathPoint++

		if curEnemy.FollowingCurrentPathPoint >= len((*enemyPathPoints)) {
			curEnemy.FollowingCurrentPathPoint = len((*enemyPathPoints)) - 1
			return true
		}
	}
	return false
}
