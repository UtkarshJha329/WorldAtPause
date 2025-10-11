package Minigames

import (
	"WorldAtPause/Engine"
	"WorldAtPause/Game/Minigames/TowerDefence"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TowerDefenceGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	TimerSystem *Engine.TimerSystem

	redTileEntityID    int
	blueTileEntityID   int
	greenTileEntityID  int
	yellowTileEntityID int

	selected_tile_entity_id int

	towerDefenceEnemyEntityID int
	bulletEntityID            int

	currentSelectedTileCoords Engine.Vector2
	currentSelectedCellCoord  Engine.Vector2Int

	occupationBoardNumCols          int
	occupationBoardNumRows          int
	occupationBoardTileSizeInPixels int
	occupationBoard                 Engine.Board[int]

	enemyPathPoints []Engine.Vector2Int

	totalNumBulletsInBulletPool int
	bulletPool                  Engine.Pool[TowerDefence.Bullet]

	totalNumEnemiesInEnemyPool int
	enemyPool                  Engine.Pool[TowerDefence.Enemy]
	enemySpawner               *Engine.PoolItem[Engine.Timer]

	totalNumTowersInTowerPool int
	towerPool                 Engine.Pool[TowerDefence.PlacedTower]
	towerTiles                map[Engine.Vector2Int]*Engine.PoolItem[TowerDefence.PlacedTower]

	numEnemiesKilled       int
	numEnemiesToKillForWin int
}

func (towerDefenceGameMode *TowerDefenceGameMode) Init() {

	towerDefenceGameMode.TimerSystem = &Engine.TimerSystem{}
	towerDefenceGameMode.TimerSystem.InitWithTimers("Tower Defence Bullets Timer System", 200)

	towerDefenceGameMode.redTileEntityID = towerDefenceGameMode.SceneRef.EntityIDsByName["Tower Defence Red Grid Tile"]
	towerDefenceGameMode.blueTileEntityID = towerDefenceGameMode.SceneRef.EntityIDsByName["Tower Defence Blue Grid Tile"]
	towerDefenceGameMode.greenTileEntityID = towerDefenceGameMode.SceneRef.EntityIDsByName["Tower Defence Green Grid Tile"]
	towerDefenceGameMode.yellowTileEntityID = towerDefenceGameMode.SceneRef.EntityIDsByName["Tower Defence Yellow Grid Tile"]

	towerDefenceGameMode.selected_tile_entity_id = towerDefenceGameMode.SceneRef.EntityIDsByName["Tower Defence Selected Tile"]

	towerDefenceGameMode.bulletEntityID = towerDefenceGameMode.SceneRef.EntityIDsByName["Tower Defence Fire Ball"]
	towerDefenceGameMode.towerDefenceEnemyEntityID = towerDefenceGameMode.SceneRef.EntityIDsByName["Tower Defence Enemy"]

	towerDefenceGameMode.occupationBoardNumCols = 10
	towerDefenceGameMode.occupationBoardNumRows = 10
	towerDefenceGameMode.occupationBoardTileSizeInPixels = 16

	towerDefenceGameMode.occupationBoard.InitBoard(towerDefenceGameMode.occupationBoardNumCols, towerDefenceGameMode.occupationBoardNumRows, towerDefenceGameMode.occupationBoardTileSizeInPixels, true)

	for y := range towerDefenceGameMode.occupationBoardNumCols {
		for x := range towerDefenceGameMode.occupationBoardNumRows {
			towerDefenceGameMode.occupationBoard.BoardData[y][x] = 0
		}
	}

	towerDefenceGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}

	towerDefenceGameMode.enemyPathPoints = []Engine.Vector2Int{
		{X: 0, Y: 4},
		{X: 1, Y: 4},
		{X: 2, Y: 4},
		{X: 2, Y: 5},
		{X: 2, Y: 6},
		{X: 2, Y: 7},
		{X: 2, Y: 8},
		{X: 3, Y: 8},
		{X: 4, Y: 8},
		{X: 5, Y: 8},
		{X: 5, Y: 7},
		{X: 5, Y: 6},
		{X: 5, Y: 5},
		{X: 5, Y: 4},
		{X: 5, Y: 3},
		{X: 5, Y: 2},
		{X: 5, Y: 1},
		{X: 6, Y: 1},
		{X: 7, Y: 1},
		{X: 8, Y: 1},
		{X: 8, Y: 2},
		{X: 8, Y: 3},
		{X: 8, Y: 4},
		{X: 9, Y: 4},
	}

	for _, pathPoint := range towerDefenceGameMode.enemyPathPoints {
		towerDefenceGameMode.occupationBoard.BoardData[pathPoint.Y][pathPoint.X] = 1
	}

	towerDefenceGameMode.totalNumBulletsInBulletPool = 100
	towerDefenceGameMode.bulletPool.InitPool("Bullet Pool", towerDefenceGameMode.totalNumBulletsInBulletPool)

	towerDefenceGameMode.towerTiles = make(map[Engine.Vector2Int]*Engine.PoolItem[TowerDefence.PlacedTower])

	towerDefenceGameMode.totalNumEnemiesInEnemyPool = 100
	towerDefenceGameMode.enemyPool.InitPool("Enemy Pool", towerDefenceGameMode.totalNumEnemiesInEnemyPool)

	towerDefenceGameMode.totalNumTowersInTowerPool = 100
	towerDefenceGameMode.towerPool.InitPool("Tower Pool", towerDefenceGameMode.totalNumTowersInTowerPool)

	towerDefenceGameMode.enemySpawner = towerDefenceGameMode.TimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(time.Duration(float64(2.5)*float64(time.Second)), true, func() {
		curEnemyPoolItem := towerDefenceGameMode.enemyPool.GetAnUnusedItemFromPool()
		curEnemyPoolItem.Item.Health = 5
		curEnemyPoolItem.Item.Position = Engine.Vector2{X: 0.0, Y: 0.0}
		curEnemyPoolItem.Item.FollowingCurrentPathPoint = 0
		curEnemyPoolItem.Item.MoveSpeedPerFrame = 0.5
	})

	towerDefenceGameMode.numEnemiesKilled = 0
}

func (towerDefenceGameMode *TowerDefenceGameMode) Update() {

	towerDefenceGameMode.TimerSystem.UpdateAllTimerDeltasAndStates()

	x, y := ebiten.CursorPosition()
	mousePixelPos := Engine.Vector2{X: float64(x), Y: float64(y)}

	tileGridStartPos := Engine.Vector2{X: float64(towerDefenceGameMode.occupationBoard.TileGridTotalXOffsetInPixels), Y: float64(towerDefenceGameMode.occupationBoard.TileGridTotalYOffsetInPixels)}
	mousePosRelToTileGrid := Engine.Subtract_Vector2(&mousePixelPos, &tileGridStartPos)

	if mousePosRelToTileGrid.X < 0 ||
		mousePosRelToTileGrid.Y < 0 ||
		mousePosRelToTileGrid.X >= float64(towerDefenceGameMode.occupationBoard.TileGridWidthInPixels) ||
		mousePosRelToTileGrid.Y >= float64(towerDefenceGameMode.occupationBoard.TileGridHeightInPixels) {
		towerDefenceGameMode.currentSelectedTileCoords = Engine.Vector2{X: -1, Y: -1}
	} else {
		gridCellXExcess := int(mousePosRelToTileGrid.X) % towerDefenceGameMode.occupationBoard.TileGridTileSizeInPixels
		gridCellYExcess := int(mousePosRelToTileGrid.Y) % towerDefenceGameMode.occupationBoard.TileGridTileSizeInPixels

		towerDefenceGameMode.currentSelectedTileCoords = Engine.Vector2{X: (mousePosRelToTileGrid.X - float64(gridCellXExcess)), Y: (mousePosRelToTileGrid.Y - float64(gridCellYExcess))}

		// PLACE NEW TOWER
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

			if towerDefenceGameMode.currentSelectedTileCoords.X < 0 || towerDefenceGameMode.currentSelectedTileCoords.Y < 0 {
				towerDefenceGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
			} else {
				towerDefenceGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: int(towerDefenceGameMode.currentSelectedTileCoords.X) / towerDefenceGameMode.occupationBoard.TileGridTileSizeInPixels, Y: int(towerDefenceGameMode.currentSelectedTileCoords.Y) / towerDefenceGameMode.occupationBoard.TileGridTileSizeInPixels}
			}

			towerDefenceGameMode.occupationBoard.BoardData[towerDefenceGameMode.currentSelectedCellCoord.Y][towerDefenceGameMode.currentSelectedCellCoord.X]++
			if towerDefenceGameMode.occupationBoard.BoardData[towerDefenceGameMode.currentSelectedCellCoord.Y][towerDefenceGameMode.currentSelectedCellCoord.X] >= 4 {
				towerDefenceGameMode.occupationBoard.BoardData[towerDefenceGameMode.currentSelectedCellCoord.Y][towerDefenceGameMode.currentSelectedCellCoord.X] = 0
			}

			curTower := towerDefenceGameMode.towerPool.GetAnUnusedItemFromPool()
			towerDefenceGameMode.towerTiles[towerDefenceGameMode.currentSelectedCellCoord] = curTower

			curTower.Item.TowerCellCoords = towerDefenceGameMode.currentSelectedCellCoord
			curTowerPositionOnScreen := towerDefenceGameMode.occupationBoard.GetPositionOfTileOnScreen(curTower.Item.TowerCellCoords)
			curTower.Item.TowerCellPosition = Engine.Vector2{X: float64(curTowerPositionOnScreen.X), Y: float64(curTowerPositionOnScreen.Y)}

			// Set Bullet Spawner For Tower
			curTower.Item.TowerShootRatePerSecond = 0.25
			curTower.Item.BulletSpeed = 1.0
			curTower.Item.BulletHealthDamageAmount = 1

			curTower.Item.TowerBulletSpawnTimer = towerDefenceGameMode.TimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(time.Duration((1.0/curTower.Item.TowerShootRatePerSecond)*float64(time.Second)), true, func() {
				curTower.Item.SpawnBullet(&towerDefenceGameMode.bulletPool, towerDefenceGameMode.TimerSystem)
			})
			curTower.Item.TowerBulletSpawnTimer.Item.PauseTimer()
		}
	}

	// HELP TOWER FIND TARGET ENEMY
	towerDefenceGameMode.towerPool.PerformOperationOnAlivePoolItems(func(curPlacedTower *Engine.PoolItem[TowerDefence.PlacedTower]) {
		notFoundTarget := true
		closestDistToEnemyYet := 500.0
		for i := 0; i < towerDefenceGameMode.enemyPool.CurNumAliveItemsInPool; i++ {
			distToCurEnemy := Engine.Distance_Vector2(&towerDefenceGameMode.enemyPool.Items[i].Item.Position, &curPlacedTower.Item.TowerCellPosition)
			if distToCurEnemy < 5.0*float64(towerDefenceGameMode.occupationBoardTileSizeInPixels) && distToCurEnemy < closestDistToEnemyYet {
				curPlacedTower.Item.TowerBulletSpawnTimer.Item.UnPauseTimer()
				curPlacedTower.Item.CurTargetEnemy = towerDefenceGameMode.enemyPool.Items[i]
				closestDistToEnemyYet = distToCurEnemy
				notFoundTarget = false
			}
		}
		if notFoundTarget {
			curPlacedTower.Item.TowerBulletSpawnTimer.Item.PauseTimer()
		}
	})

	// MOVE BULLET TOWARDS TARGET
	towerDefenceGameMode.bulletPool.PerformOperationOnAlivePoolItems(func(curBullet *Engine.PoolItem[TowerDefence.Bullet]) {

		dirToTarget := Engine.Subtract_Vector2(&curBullet.Item.EnemyTargetPoolItem.Item.Position, &curBullet.Item.Position)
		dirToTarget = Engine.Normalise_Vector2(&dirToTarget)

		curBullet.Item.Position.X += dirToTarget.X * curBullet.Item.BulletMoveSpeed
		curBullet.Item.Position.Y += dirToTarget.Y * curBullet.Item.BulletMoveSpeed
	})

	// Kill Bullets That have reached the enemy.
	towerDefenceGameMode.bulletPool.PerformOperationOnAlivePoolItemsBackwards(func(curBullet *Engine.PoolItem[TowerDefence.Bullet]) {
		if Engine.DistanceSquare_Vector2(&curBullet.Item.EnemyTargetPoolItem.Item.Position, &curBullet.Item.Position) <= 1.0 {

			curBullet.Item.EnemyTargetPoolItem.Item.Health -= curBullet.Item.BulletHealthDamageAmount
			towerDefenceGameMode.TimerSystem.KillTimer(curBullet.Item.LifeSpanTimer)
			towerDefenceGameMode.bulletPool.KillItemInPool(curBullet)
		}
	})

	// KILL 0 HEALTH ENEMIES
	towerDefenceGameMode.enemyPool.PerformOperationOnAlivePoolItemsBackwards(func(curEnemyPoolItem *Engine.PoolItem[TowerDefence.Enemy]) {
		if curEnemyPoolItem.Item.Health <= 0 {
			towerDefenceGameMode.enemyPool.KillItemInPool(curEnemyPoolItem)
			towerDefenceGameMode.numEnemiesKilled++
		}
	})

	// MOVE ENEMY TOWARDS NEXT PATH POINT
	towerDefenceGameMode.enemyPool.PerformOperationOnAlivePoolItems(func(curEnemyPoolItem *Engine.PoolItem[TowerDefence.Enemy]) {
		if curEnemyPoolItem.Item.MoveEnemyToNextPathPoint(&towerDefenceGameMode.enemyPathPoints, &towerDefenceGameMode.occupationBoard) {
			towerDefenceGameMode.numEnemiesKilled--
			curEnemyPoolItem.Item.Health = 0
		}
	})

	// KILL ENEMIES THAT REACHED THE END PATH POINT
	towerDefenceGameMode.enemyPool.PerformOperationOnAlivePoolItemsBackwards(func(curEnemyPoolItem *Engine.PoolItem[TowerDefence.Enemy]) {
		if curEnemyPoolItem.Item.Health <= 0 {
			towerDefenceGameMode.enemyPool.KillItemInPool(curEnemyPoolItem)
		}
	})

	if towerDefenceGameMode.numEnemiesKilled >= 0 && towerDefenceGameMode.numEnemiesKilled < towerDefenceGameMode.numEnemiesToKillForWin {
		towerDefenceGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS
	} else if towerDefenceGameMode.numEnemiesKilled >= towerDefenceGameMode.numEnemiesToKillForWin {
		towerDefenceGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
		towerDefenceGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	} else if towerDefenceGameMode.numEnemiesKilled <= -3 {
		towerDefenceGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_LOST
		towerDefenceGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	}

}

func (towerDefenceGameMode *TowerDefenceGameMode) Draw(screenRef *ebiten.Image) {

	// Engine.DrawActiveRoomInScene(towerDefenceGameMode.SceneRef, screenRef)

	Engine.DrawActiveTilemapInRoomInScene(towerDefenceGameMode.SceneRef, screenRef)

	drawOptions := ebiten.DrawImageOptions{}

	redTileSpriteRef := &towerDefenceGameMode.SceneRef.EntityComponentsForScene.Sprites[towerDefenceGameMode.redTileEntityID]
	blueTileSpriteRef := &towerDefenceGameMode.SceneRef.EntityComponentsForScene.Sprites[towerDefenceGameMode.blueTileEntityID]
	greenTileSpriteRef := &towerDefenceGameMode.SceneRef.EntityComponentsForScene.Sprites[towerDefenceGameMode.greenTileEntityID]
	yellowTileSpriteRef := &towerDefenceGameMode.SceneRef.EntityComponentsForScene.Sprites[towerDefenceGameMode.yellowTileEntityID]

	for y := range towerDefenceGameMode.occupationBoard.Board_num_rows {
		for x := range towerDefenceGameMode.occupationBoard.Board_num_cols {
			drawPos := Engine.Vector2{X: float64((x * towerDefenceGameMode.occupationBoard.TileGridTileSizeInPixels) + towerDefenceGameMode.occupationBoard.TileGridTotalXOffsetInPixels), Y: float64((y * towerDefenceGameMode.occupationBoard.TileGridTileSizeInPixels) + towerDefenceGameMode.occupationBoard.TileGridTotalYOffsetInPixels)}
			switch towerDefenceGameMode.occupationBoard.BoardData[y][x] {
			case 0:
				redTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			case 1:
				blueTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			case 2:
				greenTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			case 3:
				yellowTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			}
		}
	}

	if towerDefenceGameMode.currentSelectedTileCoords.X >= 0 && towerDefenceGameMode.currentSelectedTileCoords.Y >= 0 {
		selectionSprite := &towerDefenceGameMode.SceneRef.EntityComponentsForScene.Sprites[towerDefenceGameMode.selected_tile_entity_id]
		drawPos := Engine.Vector2{X: float64(towerDefenceGameMode.occupationBoard.TileGridTotalXOffsetInPixels) + towerDefenceGameMode.currentSelectedTileCoords.X, Y: float64(towerDefenceGameMode.occupationBoard.TileGridTotalYOffsetInPixels) + towerDefenceGameMode.currentSelectedTileCoords.Y}
		selectionSprite.DrawSprite(screenRef, &drawOptions, &drawPos)
	}

	bulletSpriteRef := &towerDefenceGameMode.SceneRef.EntityComponentsForScene.Sprites[towerDefenceGameMode.bulletEntityID]
	towerDefenceGameMode.bulletPool.PerformOperationOnAlivePoolItems(func(curBullet *Engine.PoolItem[TowerDefence.Bullet]) {
		bulletSpriteRef.DrawSprite(screenRef, &drawOptions, &curBullet.Item.Position)
	})

	enemySpriteRef := &towerDefenceGameMode.SceneRef.EntityComponentsForScene.Sprites[towerDefenceGameMode.towerDefenceEnemyEntityID]
	towerDefenceGameMode.enemyPool.PerformOperationOnAlivePoolItems(func(curEnemyPoolItem *Engine.PoolItem[TowerDefence.Enemy]) {
		enemySpriteRef.DrawSprite(screenRef, &drawOptions, &curEnemyPoolItem.Item.Position)
	})

	Engine.DrawActiveRoomObjectsInScene(towerDefenceGameMode.SceneRef, screenRef)
	// Engine.DrawPlayerInScene(towerDefenceGameMode.SceneRef, screenRef)

	Engine.DrawActiveRoomInSceneColliders(towerDefenceGameMode.SceneRef, screenRef)

}

func (towerDefenceGameMode *TowerDefenceGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {
	if previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestType == Engine.QUEST_TYPE_SCORE_LIMIT {
		towerDefenceGameMode.numEnemiesToKillForWin = int(previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestValues[Engine.QUEST_TYPE_SCORE_LIMIT])
	}
}
