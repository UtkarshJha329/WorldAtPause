package TowerDefence

import (
	"WorldAtPause/Engine"
	"time"
)

type PlacedTower struct {
	TowerCellCoords         Engine.Vector2Int
	TowerCellPosition       Engine.Vector2
	TowerShootRatePerSecond float64
	TowerBulletSpawnTimer   *Engine.PoolItem[Engine.Timer]

	BulletSpeed              float64
	BulletHealthDamageAmount int

	CurTargetEnemy *Engine.PoolItem[Enemy]
}

func (curPlacedTower *PlacedTower) SpawnBullet(bulletPool *Engine.Pool[Bullet], TimerSystem *Engine.TimerSystem) {
	curBulletPoolItem := bulletPool.GetAnUnusedItemFromPool()
	curTilePosOnScreen := curPlacedTower.TowerCellPosition
	curBulletPoolItem.Item.Position = curTilePosOnScreen

	curBulletPoolItem.Item.BulletMoveSpeed = curPlacedTower.BulletSpeed
	curBulletPoolItem.Item.BulletHealthDamageAmount = curPlacedTower.BulletHealthDamageAmount

	curBulletPoolItem.Item.EnemyTargetPoolItem = curPlacedTower.CurTargetEnemy

	bulletExpiryTimeInSeconds := 2.0
	curBulletPoolItem.Item.LifeSpanTimer = TimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(time.Duration(bulletExpiryTimeInSeconds*float64(time.Second)), false, func() {
		bulletPool.KillItemInPool(curBulletPoolItem)
	})
}
