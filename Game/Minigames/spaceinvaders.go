package Minigames

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpaceInvadersGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	bulletEntityID int

	bulletMoveAmountPerFrame float64
	totalNumEnemyShipsAlive  int
}

func (spaceInvadersGameMode *SpaceInvadersGameMode) Init() {

	sceneIndex := spaceInvadersGameMode.World.SceneIndexByName[spaceInvadersGameMode.SceneRef.SceneName]
	Engine.ReloadSceneWithSceneData(spaceInvadersGameMode.SceneRef, &Engine.ScenesData[sceneIndex])

	spaceInvadersGameMode.bulletEntityID = spaceInvadersGameMode.SceneRef.EntityIDsByName["Space Invaders Bullet"]
	spaceInvadersGameMode.bulletMoveAmountPerFrame = 2.0

	spaceInvadersGameMode.totalNumEnemyShipsAlive = 4
}

func (spaceInvadersGameMode *SpaceInvadersGameMode) Update() {

	entityComponentsRef := spaceInvadersGameMode.SceneRef.EntityComponentsForScene
	curRoomIndex := Engine.Vector2Int{X: 0, Y: 0}

	spaceShipMoveAmountPerFrame := 4.0

	entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID].Y = 240 - 20

	inputX := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputX += 1.0
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputX -= 1.0
	}

	spaceShipCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
		CollideWithTiles:     true,
		CollideWithObstacles: false,
		SlideWhenCollide:     false,
		CollideWithPlayer:    false,
		MovementLock:         Engine.Vector2{X: 1.0, Y: 0.0},
	}

	spaceShipInputDirection := Engine.Vector2{X: inputX, Y: 0.0}
	spaceShipTotalMoveAmount := Engine.Multiply_Float_Vector2(spaceShipMoveAmountPerFrame, &spaceShipInputDirection)

	spaceInvadersGameMode.SceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, entityComponentsRef.PlayerEntityID, spaceShipInputDirection, spaceShipTotalMoveAmount, spaceShipMoveAmountPerFrame, spaceShipCollideAndMoveParameters)

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		bulletPosRef := &entityComponentsRef.Positions[spaceInvadersGameMode.bulletEntityID]
		spaceShipPosRef := entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]

		bulletPosRef.X = spaceShipPosRef.X
		bulletPosRef.Y = spaceShipPosRef.Y - 16

		entityComponentsRef.EntityDead[spaceInvadersGameMode.bulletEntityID] = false
	}

	bulletCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
		CollideWithTiles:     true,
		CollideWithObstacles: true,
		SlideWhenCollide:     false,
		CollideWithPlayer:    false,
		MovementLock:         Engine.Vector2{X: 0.0, Y: 1.0},
	}

	bulletMoveInputDirection := Engine.Vector2{X: 0.0, Y: -1.0}
	bulletTotalMoveAmount := Engine.Multiply_Float_Vector2(spaceInvadersGameMode.bulletMoveAmountPerFrame, &bulletMoveInputDirection)

	bulletMoveCollisionResult := spaceInvadersGameMode.SceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, spaceInvadersGameMode.bulletEntityID, bulletMoveInputDirection, bulletTotalMoveAmount, spaceInvadersGameMode.bulletMoveAmountPerFrame, bulletCollideAndMoveParameters)

	if bulletMoveCollisionResult.TilemapYMoveCollisionResult.Collided ||
		bulletMoveCollisionResult.ObstacleYMoveCollisionResult.CollidedWithObstacle {

		if bulletMoveCollisionResult.ObstacleYMoveCollisionResult.CollidedWithObstacle {
			enemySpaceShipEntityID := bulletMoveCollisionResult.ObstacleYMoveCollisionResult.CollidedWithObstacleEntityID
			entityComponentsRef.EntityDead[enemySpaceShipEntityID] = true
			spaceInvadersGameMode.totalNumEnemyShipsAlive--
		}

		entityComponentsRef.EntityDead[spaceInvadersGameMode.bulletEntityID] = true
	}

	if spaceInvadersGameMode.totalNumEnemyShipsAlive <= 0 {
		spaceInvadersGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
		spaceInvadersGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	} else {
		spaceInvadersGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS
	}

}

func (spaceInvadersGameMode *SpaceInvadersGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(spaceInvadersGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(spaceInvadersGameMode.SceneRef, screenRef)

}

func (spaceInvadersGameMode *SpaceInvadersGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

}
