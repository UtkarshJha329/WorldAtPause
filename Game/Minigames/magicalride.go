package Minigames

import (
	"WorldAtPause/Engine"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type MagicalRideGameMode struct {
	SceneRef *Engine.Scene

	playerMoveAmountPerFrame       float64
	playerCollideAndMoveParameters Engine.CollideAndMoveCollisionParameters

	enemyMoveAmountPerFrame float64
}

func (magicalRideGameMode *MagicalRideGameMode) Init() {

	magicalRideGameMode.playerMoveAmountPerFrame = 3.0
	magicalRideGameMode.playerCollideAndMoveParameters = Engine.CollideAndMoveCollisionParameters{
		CollideWithTiles:     true,
		CollideWithObstacles: false,
		CollideWithPlayer:    false,
		SlideWhenCollide:     false,
		MovementLock:         Engine.Vector2{X: 1.0, Y: 1.0},
	}

	magicalRideGameMode.enemyMoveAmountPerFrame = 2.0
}

func (magicalRideGameMode *MagicalRideGameMode) Update() {

	curSceneRef := magicalRideGameMode.SceneRef
	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)
	// curRoom, curRoomHasSomeData := curSceneRef.RoomsData[curRoomIndex]

	inputDirection := Engine.Vector2{X: 0.0, Y: 0.0}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		inputDirection.Y += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		inputDirection.Y -= 1
	}

	totalMoveAmount := Engine.Multiply_Float_Vector2(magicalRideGameMode.playerMoveAmountPerFrame, &inputDirection)
	curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, entityComponentsRef.PlayerEntityID, inputDirection, totalMoveAmount, magicalRideGameMode.playerMoveAmountPerFrame, magicalRideGameMode.playerCollideAndMoveParameters)

	magicalRideGameMode.HandleEnemiesMovement()

}

func (magicalRideGameMode *MagicalRideGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(magicalRideGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(magicalRideGameMode.SceneRef, screenRef)

}

func (magicalRideGameMode *MagicalRideGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

}

func (magicalRideGameMode *MagicalRideGameMode) HandleEnemiesMovement() {

	curSceneRef := magicalRideGameMode.SceneRef
	entityComponentsRef := magicalRideGameMode.SceneRef.EntityComponentsForScene

	playerPos := entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	playerCollisionShapeRef := entityComponentsRef.CollisionShapes[entityComponentsRef.PlayerEntityID]

	curRoomIndex := Engine.Vector2Int{X: 0, Y: 0}
	curRoom := curSceneRef.RoomsData[curRoomIndex]

	for _, enemyID := range curRoom.EnemyEntityIDs {

		curEnemyPosRef := &entityComponentsRef.Positions[enemyID]
		curEnemyCollisionShapeRef := entityComponentsRef.CollisionShapes[enemyID]

		if curEnemyPosRef.X < -16 {
			curEnemyPosRef.X = 340.0
			randYPosition := rand.IntN(208)
			curEnemyPosRef.Y = 16 + float64(randYPosition)
		} else {
			curEnemyPosRef.X -= magicalRideGameMode.enemyMoveAmountPerFrame
			if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(&playerPos, playerCollisionShapeRef, curEnemyPosRef, curEnemyCollisionShapeRef); collided {
				curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_LOST
				curSceneRef.SceneGameStateData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
			}
		}
	}
}
