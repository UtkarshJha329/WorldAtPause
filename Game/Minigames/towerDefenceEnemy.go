package Minigames

import "WorldAtPause/Engine"

type TowerDefenceEnemy struct {
	position                  Engine.Vector2
	followingCurrentPathPoint int
	health                    int
	moveSpeedPerFrame         float64
}

func (curEnemy *TowerDefenceEnemy) MoveEnemyToNextPathPoint(towerDefenceGameMode *TowerDefenceGameMode) bool {

	curFollowingPathPoint := towerDefenceGameMode.enemyPathPoints[curEnemy.followingCurrentPathPoint]
	tilePosOnScreen := towerDefenceGameMode.occupationBoard.GetPositionOfTileOnScreen(curFollowingPathPoint)
	curFollowingPathPointPos := Engine.Vector2{X: float64(tilePosOnScreen.X), Y: float64(tilePosOnScreen.Y)}

	distToCurPathPoint := Engine.Distance_Vector2(&curFollowingPathPointPos, &curEnemy.position)

	if distToCurPathPoint >= float64(towerDefenceGameMode.occupationBoardTileSizeInPixels*2) {
		curEnemy.position.X = curFollowingPathPointPos.X
		curEnemy.position.Y = curFollowingPathPointPos.Y
	} else if distToCurPathPoint <= float64(towerDefenceGameMode.occupationBoardTileSizeInPixels*2) && distToCurPathPoint >= 0.1 {
		dirToCurPathPoint := Engine.Subtract_Vector2(&curFollowingPathPointPos, &curEnemy.position)
		dirToCurPathPoint = Engine.Normalise_Vector2(&dirToCurPathPoint)

		totalMoveAmount := Engine.Multiply_Float_Vector2(curEnemy.moveSpeedPerFrame, &dirToCurPathPoint)
		newPosition := Engine.Add_Vector2(&curEnemy.position, &totalMoveAmount)

		curEnemy.position.X = newPosition.X
		curEnemy.position.Y = newPosition.Y
	} else {
		curEnemy.followingCurrentPathPoint++

		if curEnemy.followingCurrentPathPoint >= len(towerDefenceGameMode.enemyPathPoints) {
			curEnemy.followingCurrentPathPoint = len(towerDefenceGameMode.enemyPathPoints) - 1
			return true
		}
	}
	return false
}
