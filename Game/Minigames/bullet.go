package Minigames

import (
	"WorldAtPause/Engine"
)

type Bullet struct {
	enemyTargetPoolItem *Engine.PoolItem[TowerDefenceEnemy]

	position Engine.Vector2

	bulletMoveSpeed          float64
	bulletHealthDamageAmount int

	lifeSpanTimer *Engine.PoolItem[Engine.Timer]
}
