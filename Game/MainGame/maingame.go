package MainGame

import (
	"WorldAtPause/Engine"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type MainGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	deactivateBreakoutTrigger    bool
	deactivateMagicalRideTrigger bool

	playerSpeed float64

	showCurrentUITree bool
	currentUITree     *Engine.UITree
}

func (mainGameMode *MainGameMode) Init() {
	mainGameMode.deactivateBreakoutTrigger = false
	mainGameMode.playerSpeed = 2.0

	mainGameMode.showCurrentUITree = false
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

		mainGameMode.IssueSceneTransitionQuests()
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

	drawOptions := ebiten.DrawImageOptions{}
	if mainGameMode.showCurrentUITree {
		mainGameMode.currentUITree.RenderUITree(0, Engine.Vector2{X: 0.0, Y: 0.0}, screenRef, &drawOptions)
	}
}

func (mainGameMode *MainGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {
	mainGameMode.SceneRef.SceneGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestCompleteLambda(previousGameStateData)
}

func (mainGameMode *MainGameMode) IssueSceneTransitionQuests() {

	curSceneRef := mainGameMode.SceneRef
	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	playerCollisionShapeRef := entityComponentsRef.CollisionShapes[entityComponentsRef.PlayerEntityID]

	shouldShowQuestAcceptUITree := false

	if !mainGameMode.deactivateBreakoutTrigger {

		breakoutTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Breakout Ball Trigger"]
		breakoutTriggerEntityPos := entityComponentsRef.Positions[breakoutTriggerEntityID]
		breakoutTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[breakoutTriggerEntityID]

		if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &breakoutTriggerEntityPos, breakoutTriggerCollisionShapeRef); collided {

			shouldShowQuestAcceptUITree = true

			if Engine.ClickedYesInYesNoUITree() {

				fmt.Println("Accepted Breakout Quest.")

				curSceneRef.SpawnChildSceneWithQuest(mainGameMode.World.SceneIndexByName["Breakout Game"],
					Engine.QuestData{
						QuestType:   Engine.QUEST_TYPE_LEVEL_COMPLETE,
						QuestValues: map[int]float64{},
						QuestCompleteLambda: func(previousGameStateData Engine.GameStateData) {
							if previousGameStateData.GameState == Engine.GAMEMODE_WON {
								fmt.Println("Won at Breakout! Speed increased!")
								mainGameMode.playerSpeed += 1.0
							} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
								fmt.Println("Lost at Breakout! Speed decreased.")
								mainGameMode.playerSpeed -= 1.0
							}
						},
					})
			}
		}
	}
	if !mainGameMode.deactivateMagicalRideTrigger {

		magicalRideTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Magical Ride Trigger"]
		magicalRideTriggerEntityPos := entityComponentsRef.Positions[magicalRideTriggerEntityID]
		magicalRideTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[magicalRideTriggerEntityID]

		if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &magicalRideTriggerEntityPos, magicalRideTriggerCollisionShapeRef); collided {

			shouldShowQuestAcceptUITree = true

			if Engine.ClickedYesInYesNoUITree() {

				fmt.Println("Accepted Magical Ride Quest.")

				curSceneRef.SpawnChildSceneWithQuest(mainGameMode.World.SceneIndexByName["Magical Ride Game"],
					Engine.QuestData{
						QuestType: Engine.QUEST_TYPE_TIME_TRIAL,
						QuestValues: map[int]float64{
							Engine.QUEST_TYPE_TIME_TRIAL: 20.0,
						},
						QuestCompleteLambda: func(previousGameStateData Engine.GameStateData) {
							if previousGameStateData.GameState == Engine.GAMEMODE_WON {
								fmt.Println("Won at Magical Ride! Speed increased!")
								mainGameMode.playerSpeed += 1.0
							} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
								fmt.Println("Lost at Magical Ride! Speed decreased.")
								mainGameMode.playerSpeed -= 1.0
							}
						},
					})
			}
		}
	}

	{
		spaceInvadersTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Space Invaders Trigger"]
		spaceInvadersTriggerEntityPos := entityComponentsRef.Positions[spaceInvadersTriggerEntityID]
		spaceInvadersTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[spaceInvadersTriggerEntityID]

		if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &spaceInvadersTriggerEntityPos, spaceInvadersTriggerCollisionShapeRef); collided {

			shouldShowQuestAcceptUITree = true

			if Engine.ClickedYesInYesNoUITree() {

				fmt.Println("Accepted Space Invaders Quest.")

				curSceneRef.SpawnChildSceneWithQuest(mainGameMode.World.SceneIndexByName["Space Invaders Game"],
					Engine.QuestData{
						QuestType:   Engine.QUEST_TYPE_LEVEL_COMPLETE,
						QuestValues: map[int]float64{},
						QuestCompleteLambda: func(previousGameStateData Engine.GameStateData) {
							if previousGameStateData.GameState == Engine.GAMEMODE_WON {
								fmt.Println("Won at Space Invaders! Speed increased!")
								mainGameMode.playerSpeed += 1.0
							} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
								fmt.Println("Lost at Space Invaders! Speed decreased.")
								mainGameMode.playerSpeed -= 1.0
							}
						},
					})
			}
		}
	}

	{
		matchThreeTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Match Three Trigger"]
		matchThreeTriggerEntityPos := entityComponentsRef.Positions[matchThreeTriggerEntityID]
		matchThreeTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[matchThreeTriggerEntityID]

		if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &matchThreeTriggerEntityPos, matchThreeTriggerCollisionShapeRef); collided {

			shouldShowQuestAcceptUITree = true

			if Engine.ClickedYesInYesNoUITree() {

				fmt.Println("Accepted Match Three Quest.")

				curSceneRef.SpawnChildSceneWithQuest(mainGameMode.World.SceneIndexByName["Match Three Game"],
					Engine.QuestData{
						QuestType: Engine.QUEST_TYPE_SCORE_LIMIT,
						QuestValues: map[int]float64{
							Engine.QUEST_TYPE_SCORE_LIMIT: 130.0,
						},
						QuestCompleteLambda: func(previousGameStateData Engine.GameStateData) {
							if previousGameStateData.GameState == Engine.GAMEMODE_WON {
								fmt.Println("Won at Match Three! Speed increased!")
								mainGameMode.playerSpeed += 1.0
							} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
								fmt.Println("Lost at Match Three! Speed decreased.")
								mainGameMode.playerSpeed -= 1.0
							}
						},
					})
			}
		}
	}

	{
		sokobanTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Sokoban Trigger"]
		sokobanTriggerEntityPos := entityComponentsRef.Positions[sokobanTriggerEntityID]
		sokobanTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[sokobanTriggerEntityID]

		if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &sokobanTriggerEntityPos, sokobanTriggerCollisionShapeRef); collided {

			shouldShowQuestAcceptUITree = true

			if Engine.ClickedYesInYesNoUITree() {

				fmt.Println("Accepted Sokoban Quest.")

				curSceneRef.SpawnChildSceneWithQuest(mainGameMode.World.SceneIndexByName["Sokoban Game"],
					Engine.QuestData{
						QuestType: Engine.QUEST_TYPE_SCORE_LIMIT,
						QuestValues: map[int]float64{
							Engine.QUEST_TYPE_SCORE_LIMIT: 130.0,
						},
						QuestCompleteLambda: func(previousGameStateData Engine.GameStateData) {
							if previousGameStateData.GameState == Engine.GAMEMODE_WON {
								fmt.Println("Won at Sokoban! Speed increased!")
								mainGameMode.playerSpeed += 1.0
							} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
								fmt.Println("Lost at Sokoban! Speed decreased.")
								mainGameMode.playerSpeed -= 1.0
							}
						},
					})
			}
		}
	}

	{
		snekTriggerEntityID := curSceneRef.EntityIDsByName["Main Game Snek Game Trigger"]
		snekTriggerEntityPos := entityComponentsRef.Positions[snekTriggerEntityID]
		snekTriggerCollisionShapeRef := entityComponentsRef.CollisionShapes[snekTriggerEntityID]

		if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(playerPosRef, playerCollisionShapeRef, &snekTriggerEntityPos, snekTriggerCollisionShapeRef); collided {

			shouldShowQuestAcceptUITree = true

			if Engine.ClickedYesInYesNoUITree() {

				fmt.Println("Accepted Snek Game Quest.")

				curSceneRef.SpawnChildSceneWithQuest(mainGameMode.World.SceneIndexByName["Snek Game"],
					Engine.QuestData{
						QuestType: Engine.QUEST_TYPE_SCORE_LIMIT,
						QuestValues: map[int]float64{
							Engine.QUEST_TYPE_SCORE_LIMIT: 9.0,
						},
						QuestCompleteLambda: func(previousGameStateData Engine.GameStateData) {
							if previousGameStateData.GameState == Engine.GAMEMODE_WON {
								fmt.Println("Won at Snek! Speed increased!")
								mainGameMode.playerSpeed += 1.0
							} else if previousGameStateData.GameState == Engine.GAMEMODE_LOST {
								fmt.Println("Lost at Snek! Speed decreased.")
								mainGameMode.playerSpeed -= 1.0
							}
						},
					})
			}
		}
	}

	mainGameMode.showCurrentUITree = shouldShowQuestAcceptUITree
	if shouldShowQuestAcceptUITree {
		mainGameMode.currentUITree = Engine.YesNoUITree
	}
}
