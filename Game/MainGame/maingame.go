package MainGame

import (
	"WorldAtPause/Engine"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type MainGameMode struct {
	SceneRef *Engine.Scene

	waitingForGameSceneFromTriggerIndex int

	deactivateBreakoutTrigger    bool
	deactivateMagicalRideTrigger bool

	playerSpeed float64
}

func (mainGameMode *MainGameMode) Init() {
	mainGameMode.deactivateBreakoutTrigger = false
	mainGameMode.playerSpeed = 2.0
}

func (mainGameMode *MainGameMode) Update() {

	curSceneRef := mainGameMode.SceneRef
	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerMoveAmountPerFrame := mainGameMode.playerSpeed

	// collideAndSlide := true

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)
	curRoom, curRoomHasSomeData := curSceneRef.RoomsData[curRoomIndex]

	inputDirection := Engine.Vector2{X: 0.0, Y: 0.0}

	curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputDirection.X += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputDirection.X -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		inputDirection.Y += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		inputDirection.Y -= 1
	}

	if inputDirection.X != 0 || inputDirection.Y != 0 {

		normalisedInputDir := Engine.Normalise_Vector2(&inputDirection)
		totalMoveAmount := Engine.Multiply_Float_Vector2(playerMoveAmountPerFrame, &normalisedInputDir)

		playerCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
			CollideWithTiles:     true,
			CollideWithObstacles: true,
			SlideWhenCollide:     true,
			CollideWithPlayer:    false,
			MovementLock:         Engine.Vector2{X: 1.0, Y: 1.0},
		}

		curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, curSceneRef.EntityComponentsForScene.PlayerEntityID, inputDirection, totalMoveAmount, playerMoveAmountPerFrame, playerCollideAndMoveParameters)
	}

	if curRoomHasSomeData {
		playerCollisionShapeRef := entityComponentsRef.CollisionShapes[entityComponentsRef.PlayerEntityID]
		stoppageDistanceFromPlayer := 32.0
		skeleMoveAmountPerFrame := 1.0
		for _, enemyEntityID := range curRoom.EnemyEntityIDs {
			skelePosRef := &entityComponentsRef.Positions[enemyEntityID]

			if Engine.DistanceSquare_Vector2(skelePosRef, playerPosRef) > math.Pow(stoppageDistanceFromPlayer, 2) {
				directionToPlayer := Engine.Subtract_Vector2(playerPosRef, skelePosRef)
				directionToPlayerNormalised := Engine.Normalise_Vector2(&directionToPlayer)

				totalDisplacement := Engine.Multiply_Float_Vector2(skeleMoveAmountPerFrame, &directionToPlayerNormalised)

				skeletonCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
					CollideWithTiles:     true,
					CollideWithObstacles: true,
					SlideWhenCollide:     true,
					CollideWithPlayer:    true,
					MovementLock:         Engine.Vector2{X: 1.0, Y: 1.0},
				}

				curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, enemyEntityID, directionToPlayerNormalised, totalDisplacement, skeleMoveAmountPerFrame, skeletonCollideAndMoveParameters)
			}

			skeletonCollisionShapeRef := entityComponentsRef.CollisionShapes[enemyEntityID]
			if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(skelePosRef, skeletonCollisionShapeRef, playerPosRef, playerCollisionShapeRef); collided {
				// fmt.Println("Skeleton is colliding with player!")
			}
		}

		for _, itemEntityID := range curRoom.ItemEntityIDs {

			if curSceneRef.EntityComponentsForScene.IsEntityAlive(itemEntityID) {
				itemPosRef := &entityComponentsRef.Positions[itemEntityID]
				itemCollisionShapeRef := entityComponentsRef.CollisionShapes[itemEntityID]

				if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(itemPosRef, itemCollisionShapeRef, &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID], entityComponentsRef.CollisionShapes[entityComponentsRef.PlayerEntityID]); collided {
					// curRoom.ItemEntityIDs = append(curRoom.ItemEntityIDs[:index], curRoom.ItemEntityIDs[index+1:]...)
					curSceneRef.EntityComponentsForScene.EntityDead[itemEntityID] = true
					fmt.Printf("Picked up item entity ID : %d\n", itemEntityID)
				}
			}
		}

		if !mainGameMode.deactivateBreakoutTrigger {

			breakoutTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Breakout Ball Trigger"]
			breakoutTriggerEntityPos := entityComponentsRef.Positions[breakoutTriggerEntityID]
			breakoutTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[breakoutTriggerEntityID]

			if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &breakoutTriggerEntityPos, breakoutTriggerCollisionShapeRef); collided {

				curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WAITING_FOR_CHILD
				curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_CHILD
				curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeToIndex = 1
				mainGameMode.waitingForGameSceneFromTriggerIndex = breakoutTriggerEntityID
			}
		}
		if !mainGameMode.deactivateMagicalRideTrigger {

			magicalRideTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Magical Ride Trigger"]
			magicalRideTriggerEntityPos := entityComponentsRef.Positions[magicalRideTriggerEntityID]
			magicalRideTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[magicalRideTriggerEntityID]

			if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &magicalRideTriggerEntityPos, magicalRideTriggerCollisionShapeRef); collided {

				curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WAITING_FOR_CHILD
				curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_CHILD
				curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeToIndex = 2

				curSceneRef.SceneGameStateData.SceneChangeData.QuestIssuedDuringSceneChange = Engine.QuestData{
					QuestType:  Engine.QUEST_TYPE_TIME_TRIAL,
					QuestValue: 50.0,
				}

				mainGameMode.waitingForGameSceneFromTriggerIndex = magicalRideTriggerEntityID
			}
		}
	}
	currentLevelIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)
	playerInLevel := Engine.Vector2{X: float64(currentLevelIndex.X), Y: float64(currentLevelIndex.Y)}

	levelHalfSize := Engine.Multiply_Float_Vector2(0.5, &entityComponentsRef.Tilemap.WorldGridSize)
	currentLevelPos := Engine.Vector2{X: playerInLevel.X * entityComponentsRef.Tilemap.WorldGridSize.X, Y: playerInLevel.Y * entityComponentsRef.Tilemap.WorldGridSize.Y}
	currentLevelCentre := Engine.Add_Vector2(&currentLevelPos, &levelHalfSize)

	Engine.CameraFollowTarget(currentLevelCentre, entityComponentsRef)
}

func (mainGameMode *MainGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(mainGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(mainGameMode.SceneRef, screenRef)

	breakoutTriggerEntityID := mainGameMode.SceneRef.EntityIDsByName["Main Game Breakout Ball Trigger"]
	breakoutTriggerEntityPosition := mainGameMode.SceneRef.EntityComponentsForScene.Positions[breakoutTriggerEntityID]
	breakoutTriggerSprite := mainGameMode.SceneRef.EntityComponentsForScene.Sprites[breakoutTriggerEntityID]

	drawOptions := ebiten.DrawImageOptions{}
	breakoutTriggerSprite.DrawSprite(screenRef, &drawOptions, &breakoutTriggerEntityPosition)

	magicalRideTriggerEntityID := mainGameMode.SceneRef.EntityIDsByName["Main Game Magical Ride Trigger"]
	magicalRideTriggerEntityPosition := mainGameMode.SceneRef.EntityComponentsForScene.Positions[magicalRideTriggerEntityID]
	magicalRideTriggerSprite := mainGameMode.SceneRef.EntityComponentsForScene.Sprites[magicalRideTriggerEntityID]

	magicalRideTriggerSprite.DrawSprite(screenRef, &drawOptions, &magicalRideTriggerEntityPosition)

}

func (mainGameMode *MainGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

	breakoutTriggerEntityID := mainGameMode.SceneRef.EntityIDsByName["Main Game Breakout Ball Trigger"]
	if mainGameMode.waitingForGameSceneFromTriggerIndex == breakoutTriggerEntityID {
		fmt.Println("Finished breakout.")
		mainGameMode.deactivateBreakoutTrigger = true

		if previousGameStateData.GameState == Engine.GAMEMODE_WON {
			fmt.Println("Won at breakout! Speed increased!")
			mainGameMode.playerSpeed += 1.0
		} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
			fmt.Println("Lost at breakout! Speed decreased.")
			mainGameMode.playerSpeed -= 1.0
		}
	}

	magicalRideTriggerEntityID := mainGameMode.SceneRef.EntityIDsByName["Main Game Magical Ride Trigger"]
	if mainGameMode.waitingForGameSceneFromTriggerIndex == magicalRideTriggerEntityID {
		fmt.Println("Finished Magical Ride.")
		mainGameMode.deactivateMagicalRideTrigger = true

		if previousGameStateData.GameState == Engine.GAMEMODE_WON {
			fmt.Println("Won at Magical Ride! Speed increased!")
			mainGameMode.playerSpeed += 1.0
		} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
			fmt.Println("Lost at Magical Ride! Speed decreased.")
			mainGameMode.playerSpeed -= 1.0
		}
	}

}
