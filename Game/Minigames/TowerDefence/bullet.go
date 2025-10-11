package TowerDefence

import (
	"WorldAtPause/Engine"
)

type Bullet struct {
	EnemyTargetPoolItem *Engine.PoolItem[Enemy]

	Position Engine.Vector2

	BulletMoveSpeed          float64
	BulletHealthDamageAmount int

	LifeSpanTimer *Engine.PoolItem[Engine.Timer]
}
