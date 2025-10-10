package Minigames

import (
	"WorldAtPause/Engine"
	"time"
)

type PlacedTower struct {
	towerCellCoords         Engine.Vector2Int
	towerCellPosition       Engine.Vector2
	towerShootRatePerSecond float64
	TowerBulletSpawnTimer   *Engine.PoolItem[Engine.Timer]

	bulletSpeed              float64
	bulletHealthDamageAmount int

	curTargetEnemy *Engine.PoolItem[TowerDefenceEnemy]
}

func (curPlacedTower *PlacedTower) SpawnBullet(towerDefenceGameMode *TowerDefenceGameMode) {
	curBulletPoolItem := towerDefenceGameMode.bulletPool.GetAnUnusedItemFromPool()
	curTilePosOnScreen := curPlacedTower.towerCellPosition
	curBulletPoolItem.Item.position = curTilePosOnScreen

	curBulletPoolItem.Item.bulletMoveSpeed = curPlacedTower.bulletSpeed
	curBulletPoolItem.Item.bulletHealthDamageAmount = curPlacedTower.bulletHealthDamageAmount

	curBulletPoolItem.Item.enemyTargetPoolItem = curPlacedTower.curTargetEnemy

	bulletExpiryTimeInSeconds := 2.0
	curBulletPoolItem.Item.lifeSpanTimer = towerDefenceGameMode.TimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(time.Duration(bulletExpiryTimeInSeconds*float64(time.Second)), false, func() {
		towerDefenceGameMode.bulletPool.KillItemInPool(curBulletPoolItem)
	})
}
