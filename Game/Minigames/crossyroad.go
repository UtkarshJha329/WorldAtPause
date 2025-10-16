package Minigames

import (
	"WorldAtPause/Engine"
	"WorldAtPause/Game/Minigames/CrossyRoad"
	"log"
	"math"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CrossyRoadGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	ObjectMovementTimerSystem   Engine.TimerSystem
	objectMovementEveryDuration time.Duration

	Vector2InterpolationSystem    Engine.InterpolationSystem[Engine.Vector2]
	Vector2IntInterpolationSystem Engine.InterpolationSystem[Engine.Vector2Int]
	IntInterpolationSystem        Engine.InterpolationSystem[int]

	bufferedInputTimerSystem Engine.TimerSystem

	redTileEntityID    int
	blueTileEntityID   int
	greenTileEntityID  int
	yellowTileEntityID int

	testBlockTileEntityID int

	worldGridSize             Engine.Vector2Int
	worldGridTileSizeInPixels int

	totalNumWorldBoards int
	worldChunksPool     Engine.Pool[CrossyRoad.Chunk]

	totalNumSectionsInSectionPool int
	sectionsPool                  Engine.Pool[CrossyRoad.Section]

	worldPositionOffset Engine.Vector2Int

	curChunkIndex                int
	numChunksAliveAtOnceOnScreen int
	chunkSlideInterpolationTime  time.Duration
	chunkSlideInterpolating      bool

	personEntityID                     int
	characterInChunkIndex              int
	characterTilePosition              Engine.Vector2Int
	characterScreenPosition            Engine.Vector2
	characterIsInterpolating           bool
	characterMovementInterpolationTime time.Duration
	characterOnWalkableObject          bool

	characterNumTilesCrossed       int
	chunkSlideAfterNumTilesCrossed int

	inputBuffer                        []Engine.Vector2Int
	curInputInBufferWriteIndex         uint
	maxInputsInBuffer                  uint
	inputBufferTimeInSeconds           float64
	bufferInputClearResetTimerPoolItem *Engine.PoolItem[Engine.Timer]

	rowsToCrossForWin  int
	madeFirstMove      bool
	fellOutOfPlatforms bool
}

func (crossyRoadGameMode *CrossyRoadGameMode) Init() {

	sceneIndex := crossyRoadGameMode.World.SceneIndexByName[crossyRoadGameMode.SceneRef.SceneName]
	Engine.ReloadSceneWithSceneData(crossyRoadGameMode.SceneRef, &Engine.ScenesData[sceneIndex])

	crossyRoadGameMode.Vector2InterpolationSystem.InitWithInterpolationsInPool("Crossy Road Vector2 Interpolation System", 10)
	crossyRoadGameMode.Vector2IntInterpolationSystem.InitWithInterpolationsInPool("Crossy Road Vector2Int Interpolation System", 10)
	crossyRoadGameMode.IntInterpolationSystem.InitWithInterpolationsInPool("Crossy Road Int Interpolation System", 50)

	crossyRoadGameMode.ObjectMovementTimerSystem.InitWithTimers("Object Movement Timer System", 1)

	crossyRoadGameMode.redTileEntityID = crossyRoadGameMode.SceneRef.EntityIDsByName["Crossy Road Red Grid Tile"]
	crossyRoadGameMode.blueTileEntityID = crossyRoadGameMode.SceneRef.EntityIDsByName["Crossy Road Blue Grid Tile"]
	crossyRoadGameMode.greenTileEntityID = crossyRoadGameMode.SceneRef.EntityIDsByName["Crossy Road Green Grid Tile"]
	crossyRoadGameMode.yellowTileEntityID = crossyRoadGameMode.SceneRef.EntityIDsByName["Crossy Road Yellow Grid Tile"]

	crossyRoadGameMode.testBlockTileEntityID = crossyRoadGameMode.SceneRef.EntityIDsByName["Crossy Road Test Block"]

	crossyRoadGameMode.personEntityID = crossyRoadGameMode.SceneRef.EntityIDsByName["Crossy Road Player"]

	crossyRoadGameMode.totalNumWorldBoards = 3
	crossyRoadGameMode.worldGridSize = Engine.Vector2Int{X: 20, Y: 15}
	crossyRoadGameMode.worldGridTileSizeInPixels = 16

	crossyRoadGameMode.worldChunksPool.InitPool("Crossy Roads World Board Pool", crossyRoadGameMode.totalNumWorldBoards)

	crossyRoadGameMode.curChunkIndex = 0
	for crossyRoadGameMode.worldChunksPool.CurNumAliveItemsInPool < 2 {
		curWorldChunkPoolItem := crossyRoadGameMode.worldChunksPool.GetAnUnusedItemFromPool()
		curWorldChunkPoolItem.Item.ChunkIndex = crossyRoadGameMode.curChunkIndex

		crossyRoadGameMode.curChunkIndex++
	}

	crossyRoadGameMode.worldChunksPool.PerformOperationOnAlivePoolItems(func(curChunkPoolItem *Engine.PoolItem[CrossyRoad.Chunk]) {
		curChunkPoolItem.Item.ChunkBoard.InitBoard(crossyRoadGameMode.worldGridSize.X, crossyRoadGameMode.worldGridSize.Y, crossyRoadGameMode.worldGridTileSizeInPixels, true)
		// fmt.Println("Size of Board Data := ", len(curChunkPoolItem.Item.ChunkBoard.BoardData), len(curChunkPoolItem.Item.ChunkBoard.BoardData[0]))
		curChunkPoolItem.Item.RowsAndObjects = make([]CrossyRoad.Row, curChunkPoolItem.Item.ChunkBoard.Board_num_rows)

		for i := 0; i < curChunkPoolItem.Item.ChunkBoard.Board_num_rows; i++ {
			curChunkPoolItem.Item.RowsAndObjects[i].InRowObjectsPool.InitPool("Chunk "+strconv.Itoa(curChunkPoolItem.Item.ChunkIndex)+" Row "+strconv.Itoa(i)+" Objects Pool", curChunkPoolItem.Item.ChunkBoard.Board_num_cols)
		}
	})

	crossyRoadGameMode.totalNumSectionsInSectionPool = 15
	crossyRoadGameMode.sectionsPool.InitPool("Crossy Roads Sections Pool", crossyRoadGameMode.totalNumSectionsInSectionPool)

	// crossyRoadGameMode.FillTwoNextBoardsWithRandomSections()
	crossyRoadGameMode.objectMovementEveryDuration = time.Duration(float64(2.5) * float64(time.Second))
	crossyRoadGameMode.worldChunksPool.PerformOperationOnAlivePoolItems(func(curChunkPoolItem *Engine.PoolItem[CrossyRoad.Chunk]) {
		crossyRoadGameMode.FillChunkWithRandomSections(&curChunkPoolItem.Item)
		crossyRoadGameMode.FillChunkWithSectionAppropriateObjects(&curChunkPoolItem.Item)

		for i := 0; i < curChunkPoolItem.Item.ChunkBoard.Board_num_rows; i++ {
			curRowRef := &curChunkPoolItem.Item.RowsAndObjects[i]
			interpolateToPosition := 16
			curRowRef.ObjectsMovementInterpolationPoolItem = crossyRoadGameMode.IntInterpolationSystem.CreateNewInterpolation(0, &curRowRef.ObjectsMovementInterpolationOffset, &interpolateToPosition, crossyRoadGameMode.objectMovementEveryDuration, true, Engine.LinearInterpolationInt, func() {
				curRowRef.ObjectsMovementInterpolationOffset = 0
			})
		}
	})

	crossyRoadGameMode.worldPositionOffset = Engine.Vector2Int{X: 0, Y: 0}
	// NUM CHUNKS ON SCREEN AT ONCE = (SCREEN HEIGHT / CHUNKS HEIGHT) + 1
	crossyRoadGameMode.numChunksAliveAtOnceOnScreen = 2
	crossyRoadGameMode.chunkSlideInterpolationTime = time.Duration(float64(0.75) * float64(time.Second))
	crossyRoadGameMode.chunkSlideInterpolating = false

	crossyRoadGameMode.characterInChunkIndex = 0
	crossyRoadGameMode.characterTilePosition = Engine.Vector2Int{X: 9, Y: 0}
	crossyRoadGameMode.characterIsInterpolating = false
	crossyRoadGameMode.characterMovementInterpolationTime = time.Duration(float64(0.15) * float64(time.Second))
	crossyRoadGameMode.characterScreenPosition = *crossyRoadGameMode.GetCurrentCharacterScreenPos()
	crossyRoadGameMode.characterScreenPosition.Y += float64(crossyRoadGameMode.worldPositionOffset.Y)

	crossyRoadGameMode.characterNumTilesCrossed = 0
	crossyRoadGameMode.chunkSlideAfterNumTilesCrossed = 7

	crossyRoadGameMode.inputBufferTimeInSeconds = 0.25
	crossyRoadGameMode.curInputInBufferWriteIndex = 0
	crossyRoadGameMode.maxInputsInBuffer = uint(ebiten.DefaultTPS * crossyRoadGameMode.inputBufferTimeInSeconds)
	crossyRoadGameMode.inputBuffer = make([]Engine.Vector2Int, crossyRoadGameMode.maxInputsInBuffer)
	crossyRoadGameMode.bufferedInputTimerSystem.InitWithTimers("Buffered Inputs Timer System", 1)
	crossyRoadGameMode.bufferInputClearResetTimerPoolItem = crossyRoadGameMode.bufferedInputTimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(time.Duration(crossyRoadGameMode.inputBufferTimeInSeconds*float64(time.Second)), true, func() {
		crossyRoadGameMode.curInputInBufferWriteIndex = 0
	})

	crossyRoadGameMode.ObjectMovementTimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(crossyRoadGameMode.objectMovementEveryDuration, true, func() {
		crossyRoadGameMode.worldChunksPool.PerformOperationOnAlivePoolItems(func(curChunkPoolItem *Engine.PoolItem[CrossyRoad.Chunk]) {
			curChunkPoolItem.Item.SlideAllRowsInAppropriateDirectionsByOneTileOnBoard()
		})
	})

	crossyRoadGameMode.madeFirstMove = false
	crossyRoadGameMode.fellOutOfPlatforms = false
}

func (crossyRoadGameMode *CrossyRoadGameMode) Update() {

	crossyRoadGameMode.Vector2InterpolationSystem.UpdateAllInterpolationDeltasAndStates()
	crossyRoadGameMode.Vector2IntInterpolationSystem.UpdateAllInterpolationDeltasAndStates()
	crossyRoadGameMode.IntInterpolationSystem.UpdateAllInterpolationDeltasAndStates()

	crossyRoadGameMode.bufferedInputTimerSystem.UpdateAllTimerDeltasAndStates()

	crossyRoadGameMode.ObjectMovementTimerSystem.UpdateAllTimerDeltasAndStates()

	if crossyRoadGameMode.madeFirstMove && !crossyRoadGameMode.characterOnWalkableObject {
		crossyRoadGameMode.fellOutOfPlatforms = true
	}

	if crossyRoadGameMode.characterOnWalkableObject {
		crossyRoadGameMode.madeFirstMove = true
		curInputDirection := Engine.Vector2Int{X: 0, Y: 0}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			curInputDirection.Y = -1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			curInputDirection.Y = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			curInputDirection.X = -1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			curInputDirection.X = 1
		}

		if curInputDirection.X != 0 || curInputDirection.Y != 0 {
			crossyRoadGameMode.inputBuffer[crossyRoadGameMode.curInputInBufferWriteIndex] = curInputDirection
			crossyRoadGameMode.curInputInBufferWriteIndex++

			if crossyRoadGameMode.curInputInBufferWriteIndex >= crossyRoadGameMode.maxInputsInBuffer {
				crossyRoadGameMode.curInputInBufferWriteIndex = 0
			}
		}

		// Annoying to have to wait for the chunk to slide before moving.
		if !crossyRoadGameMode.characterIsInterpolating && !crossyRoadGameMode.chunkSlideInterpolating {

			bufferedInput := Engine.Vector2Int{X: 0, Y: 0}
			crossyRoadGameMode.GetBufferedInput(&bufferedInput)

			crossedHalfChunk := false

			if bufferedInput.Y == -1 {
				crossyRoadGameMode.characterNumTilesCrossed++
			}
			if bufferedInput.Y == 1 {
				crossyRoadGameMode.characterNumTilesCrossed--
			}

			crossyRoadGameMode.characterTilePosition.Y += (-1 * bufferedInput.Y)
			crossyRoadGameMode.characterTilePosition.X += bufferedInput.X

			if crossyRoadGameMode.characterTilePosition.Y < 0 {
				crossyRoadGameMode.characterTilePosition.Y = crossyRoadGameMode.worldGridSize.Y - 1
				crossyRoadGameMode.characterInChunkIndex--
				foundPreviousChunk := false
				crossyRoadGameMode.worldChunksPool.PerformOperationOnAlivePoolItems(func(curChunkPoolItem *Engine.PoolItem[CrossyRoad.Chunk]) {
					if curChunkPoolItem.Item.ChunkIndex == crossyRoadGameMode.characterInChunkIndex {
						foundPreviousChunk = true
					}
				})
				if !foundPreviousChunk {
					crossyRoadGameMode.characterInChunkIndex++
					crossyRoadGameMode.characterTilePosition.Y = 0
					crossyRoadGameMode.characterNumTilesCrossed++
				}
			}
			if crossyRoadGameMode.characterTilePosition.Y >= crossyRoadGameMode.worldGridSize.Y {
				crossyRoadGameMode.characterTilePosition.Y = 0
				crossyRoadGameMode.characterInChunkIndex++
			}
			if crossyRoadGameMode.characterTilePosition.X < 0 {
				crossyRoadGameMode.characterTilePosition.X = 0
			}
			if crossyRoadGameMode.characterTilePosition.X >= crossyRoadGameMode.worldGridSize.X {
				crossyRoadGameMode.characterTilePosition.X = crossyRoadGameMode.worldGridSize.X - 1
			}

			if crossyRoadGameMode.characterNumTilesCrossed != 0 {
				if bufferedInput.Y == -1 && crossyRoadGameMode.characterNumTilesCrossed%crossyRoadGameMode.chunkSlideAfterNumTilesCrossed == 0 ||
					bufferedInput.Y == 1 && (crossyRoadGameMode.characterNumTilesCrossed+1)%crossyRoadGameMode.chunkSlideAfterNumTilesCrossed == 0 {
					crossedHalfChunk = true
				}
			}

			if crossedHalfChunk {
				targetWorldPositionOffsetY := crossyRoadGameMode.worldPositionOffset.Y + (-1 * bufferedInput.Y * (crossyRoadGameMode.chunkSlideAfterNumTilesCrossed * crossyRoadGameMode.worldGridTileSizeInPixels))

				targetWorldPositionOffset := Engine.Vector2Int{X: crossyRoadGameMode.worldPositionOffset.X, Y: targetWorldPositionOffsetY}
				if targetWorldPositionOffset.Y != crossyRoadGameMode.worldPositionOffset.Y {

					crossyRoadGameMode.Vector2IntInterpolationSystem.CreateNewInterpolation(crossyRoadGameMode.worldPositionOffset, &crossyRoadGameMode.worldPositionOffset, &targetWorldPositionOffset, crossyRoadGameMode.chunkSlideInterpolationTime, false, Engine.LinearInterpolationVector2Int, func() {
						crossyRoadGameMode.chunkSlideInterpolating = false
					})
					crossyRoadGameMode.chunkSlideInterpolating = true
				}
			}

			if bufferedInput.X != 0 || bufferedInput.Y != 0 {
				characterTargetScreenPosition := crossyRoadGameMode.GetCurrentCharacterScreenPos()
				crossyRoadGameMode.Vector2InterpolationSystem.CreateNewInterpolation(crossyRoadGameMode.characterScreenPosition, &crossyRoadGameMode.characterScreenPosition, characterTargetScreenPosition, crossyRoadGameMode.characterMovementInterpolationTime, false, Engine.LinearInterpolationVector2, func() {
					crossyRoadGameMode.characterIsInterpolating = false
				})
				crossyRoadGameMode.characterIsInterpolating = true
			}
		}
	}

	crossyRoadGameMode.worldChunksPool.PerformOperationOnAlivePoolItemsBackwards(func(curChunkPoolItem *Engine.PoolItem[CrossyRoad.Chunk]) {
		if curChunkPoolItem.Item.HasChunkBeenPassed(float64(crossyRoadGameMode.worldPositionOffset.Y)) {

			curChunkPoolItem.Item.EmptyAllRowsAndObjects()
			crossyRoadGameMode.worldChunksPool.KillItemInPool(curChunkPoolItem)

			curWorldChunkPoolItem := crossyRoadGameMode.worldChunksPool.GetAnUnusedItemFromPool()
			curWorldChunkPoolItem.Item.ChunkIndex = crossyRoadGameMode.curChunkIndex
			curWorldChunkPoolItem.Item.Filled = false

			crossyRoadGameMode.FillChunkWithRandomSections(&curWorldChunkPoolItem.Item)
			crossyRoadGameMode.FillChunkWithSectionAppropriateObjects(&curChunkPoolItem.Item)

			crossyRoadGameMode.curChunkIndex++
		}
	})

	crossyRoadGameMode.worldChunksPool.PerformOperationOnAlivePoolItems(func(curChunkPoolItem *Engine.PoolItem[CrossyRoad.Chunk]) {
		if curChunkPoolItem.Item.ChunkIndex == crossyRoadGameMode.characterInChunkIndex {
			crossyRoadGameMode.characterOnWalkableObject = curChunkPoolItem.Item.TileIsThreeFourthWayCoveredByObject(crossyRoadGameMode.characterTilePosition)
		}
	})

	if crossyRoadGameMode.characterNumTilesCrossed >= crossyRoadGameMode.rowsToCrossForWin {
		crossyRoadGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
		crossyRoadGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	} else if crossyRoadGameMode.fellOutOfPlatforms {
		crossyRoadGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_LOST
		crossyRoadGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	} else {
		crossyRoadGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS
	}

}

func (crossyRoadGameMode *CrossyRoadGameMode) Draw(screenRef *ebiten.Image) {

	redTileSpriteRef := &crossyRoadGameMode.SceneRef.EntityComponentsForScene.Sprites[crossyRoadGameMode.redTileEntityID]
	blueTileSpriteRef := &crossyRoadGameMode.SceneRef.EntityComponentsForScene.Sprites[crossyRoadGameMode.blueTileEntityID]
	greenTileSpriteRef := &crossyRoadGameMode.SceneRef.EntityComponentsForScene.Sprites[crossyRoadGameMode.greenTileEntityID]
	yellowTileSpriteRef := &crossyRoadGameMode.SceneRef.EntityComponentsForScene.Sprites[crossyRoadGameMode.yellowTileEntityID]

	testBlockTileSpriteRef := &crossyRoadGameMode.SceneRef.EntityComponentsForScene.Sprites[crossyRoadGameMode.testBlockTileEntityID]

	drawOptions := ebiten.DrawImageOptions{}

	crossyRoadGameMode.worldChunksPool.PerformOperationOnAlivePoolItems(func(curChunkPoolItem *Engine.PoolItem[CrossyRoad.Chunk]) {

		// fmt.Println("Drawing board := ", crossyRoadGameMode.worldChunksPool.Items[i].ItemIndex)

		if curChunkPoolItem.Item.IsChunkCurrentlyVisibleOnScreen(float64(crossyRoadGameMode.worldPositionOffset.Y)) {

			// fmt.Println("Cur chunk with ID : ", curChunkPoolItem.Item.ChunkIndex, " has y top := ", float64((0*curChunkPoolItem.Item.ChunkBoard.TileGridTileSizeInPixels)+curChunkPoolItem.Item.ChunkBoard.TileGridTotalYOffsetInPixels)+float64(crossyRoadGameMode.worldPositionOffset.Y)-float64(curChunkPoolItem.Item.ChunkIndex*crossyRoadGameMode.worldGridSize.Y*crossyRoadGameMode.worldGridTileSizeInPixels))

			numColsInBoard := curChunkPoolItem.Item.ChunkBoard.Board_num_cols
			numRowsInBoard := curChunkPoolItem.Item.ChunkBoard.Board_num_rows

			for y := numRowsInBoard - 1; y >= 0; y-- {
				for x := 0; x < numColsInBoard; x++ {
					drawPos := Engine.Vector2{X: float64((x * curChunkPoolItem.Item.ChunkBoard.TileGridTileSizeInPixels) + curChunkPoolItem.Item.ChunkBoard.TileGridTotalXOffsetInPixels), Y: float64((y*curChunkPoolItem.Item.ChunkBoard.TileGridTileSizeInPixels)+curChunkPoolItem.Item.ChunkBoard.TileGridTotalYOffsetInPixels) + float64(crossyRoadGameMode.worldPositionOffset.Y) - float64(curChunkPoolItem.Item.ChunkIndex*crossyRoadGameMode.worldGridSize.Y*crossyRoadGameMode.worldGridTileSizeInPixels)}
					switch curChunkPoolItem.Item.ChunkBoard.BoardData[numRowsInBoard-y-1][x] {
					case 0:
						redTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
					case 1:
						blueTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
					case 2:
						greenTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
					case 3:
						yellowTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
					}

					// if curChunkPoolItem.Item.ChunkBoard.BoardData[numRowsInBoard-y-1][x] >= 5 {
					// 	drawOptions.ColorScale.SetR(1.0)
					// 	drawOptions.ColorScale.SetG(0.0)
					// 	drawOptions.ColorScale.SetB(0.0)
					// 	testBlockTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)

					// 	drawOptions.ColorScale.Reset()
					// }

				}
			}

			curChunkPoolItem.Item.DrawChunkObjectsInPosition(screenRef, testBlockTileSpriteRef, &drawOptions, &crossyRoadGameMode.worldPositionOffset, &crossyRoadGameMode.worldGridSize)
		}
	})

	personTileSpriteRef := &crossyRoadGameMode.SceneRef.EntityComponentsForScene.Sprites[crossyRoadGameMode.personEntityID]
	// crossyRoadGameMode.characterScreenPosition.Y += float64(crossyRoadGameMode.worldPositionOffset.Y)
	// personTileSpriteRef.DrawSprite(screenRef, &drawOptions, &crossyRoadGameMode.characterScreenPosition)
	characterFinalScreenPosition := Engine.Vector2{X: crossyRoadGameMode.characterScreenPosition.X, Y: crossyRoadGameMode.characterScreenPosition.Y + float64(crossyRoadGameMode.worldPositionOffset.Y)}
	personTileSpriteRef.DrawSprite(screenRef, &drawOptions, &characterFinalScreenPosition)
}

func (crossyRoadGameMode *CrossyRoadGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {
	if previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestType == Engine.QUEST_TYPE_SCORE_LIMIT {
		crossyRoadGameMode.rowsToCrossForWin = int(previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestValues[Engine.QUEST_TYPE_SCORE_LIMIT])
	}
}

func (crossyRoadGameMode *CrossyRoadGameMode) FillChunkWithRandomSections(curChunk *CrossyRoad.Chunk) {

	curBoardHeightToFill := curChunk.ChunkBoard.Board_num_rows
	previousSectionType := CrossyRoad.CrossyRoads_SectionType_SPACE

	if crossyRoadGameMode.sectionsPool.CurNumAliveItemsInPool > 0 {
		if crossyRoadGameMode.sectionsPool.CurNumAliveItemsInPool > 1 {
			log.Fatal("FATAL ERROR : Chunk Index := ", curChunk.ChunkIndex, " Extra Unnecessary Sections are being generated. Num currently alive := ", crossyRoadGameMode.sectionsPool.CurNumAliveItemsInPool)
		} else {
			curBoardHeightToFill -= (crossyRoadGameMode.sectionsPool.Items[0].Item.SectionForNumRows - crossyRoadGameMode.sectionsPool.Items[0].Item.NumSectionRowsWrittenToBoard)
			previousSectionType = crossyRoadGameMode.sectionsPool.Items[0].Item.SectionType
		}
	}

	for curBoardHeightToFill > 0 {
		curSectionPoolItem := crossyRoadGameMode.sectionsPool.GetAnUnusedItemFromPool()
		curSectionPoolItem.Item.GenerateRandomSectionButNotSection(3, 3, previousSectionType)
		curBoardHeightToFill -= curSectionPoolItem.Item.SectionForNumRows
		previousSectionType = curSectionPoolItem.Item.SectionType
		curSectionPoolItem.Item.NumSectionRowsWrittenToBoard = 0
	}

	curSectionIndex := 0

	// fmt.Println("----------NEW BOARD-----------")
	// fmt.Println("Board Index := ", curBoardPoolItem.ItemIndex, " Num Rows In Board := ", curBoardPoolItem.Item.Board_num_rows)

	numRowsInBoard := curChunk.ChunkBoard.Board_num_rows
	for y := 0; y < numRowsInBoard; y++ {

		curSectionToUse := &crossyRoadGameMode.sectionsPool.Items[curSectionIndex].Item
		if crossyRoadGameMode.sectionsPool.CurNumAliveItemsInPool > 0 {

			numRowsInSection := curSectionToUse.SectionForNumRows

			// fmt.Println("Cur Section type := ", curSectionToUse.SectionType)

			curBoardRowIndexToWriteAt := 0
			for curSectionToUse.NumSectionRowsWrittenToBoard < numRowsInSection && y+curBoardRowIndexToWriteAt < numRowsInBoard {

				// fmt.Println(" cur section row to copy from : ", i, " cur board row to fill base : ", y, " Section value : ", crossyRoadGameMode.sectionsPool.Items[0].Item.SectionType)
				// fmt.Println("Cur board row to fill : ", y+curBoardRowIndexToWriteAt)

				for x := 0; x < curChunk.ChunkBoard.Board_num_cols; x++ {

					if curSectionToUse.SectionType == CrossyRoad.CrossyRoads_SectionType_SPACE {
						curChunk.ChunkBoard.BoardData[y+curBoardRowIndexToWriteAt][x] = 0 // red
					} else if curSectionToUse.SectionType == CrossyRoad.CrossyRoads_SectionType_ROAD {
						curChunk.ChunkBoard.BoardData[y+curBoardRowIndexToWriteAt][x] = 1 // blue
					} else if curSectionToUse.SectionType == CrossyRoad.CrossyRoads_SectionType_POOL {
						curChunk.ChunkBoard.BoardData[y+curBoardRowIndexToWriteAt][x] = 2 // green
					} else if curSectionToUse.SectionType == CrossyRoad.CrossyRoads_SectionType_PARKING {
						curChunk.ChunkBoard.BoardData[y+curBoardRowIndexToWriteAt][x] = 3 // yellow
					}
				}

				curBoardRowIndexToWriteAt++
				curSectionToUse.NumSectionRowsWrittenToBoard++
			}

			y = y + curBoardRowIndexToWriteAt - 1

			if curSectionToUse.NumSectionRowsWrittenToBoard >= numRowsInSection-1 {
				curSectionIndex++

				if y < numRowsInBoard-1 &&
					curSectionIndex >= crossyRoadGameMode.sectionsPool.CurNumAliveItemsInPool {

					log.Fatal("FATAL ERROR : Using sections not generated.", " Cur chunk no. := ", curChunk.ChunkIndex, " Cur Section no. := ", curSectionIndex, " Cur Num Section Alive In Pool := ", crossyRoadGameMode.sectionsPool.CurNumAliveItemsInPool)
				}
			}

		} else {
			log.Fatal("FATAL ERROR : Fewer Sections generated than accessed.")
			// Break out of loops since something has gone wrong
			// as there should've been more than enough sections in the pool to copy onto the boards.
		}
	}

	curChunk.Filled = true
	crossyRoadGameMode.sectionsPool.PerformOperationOnAlivePoolItemsBackwards(func(curSectionPoolItem *Engine.PoolItem[CrossyRoad.Section]) {
		if curSectionPoolItem.Item.NumSectionRowsWrittenToBoard >= curSectionPoolItem.Item.SectionForNumRows {
			crossyRoadGameMode.sectionsPool.KillItemInPool(curSectionPoolItem)
		}
	})
}

func (crossyRoadGameMode *CrossyRoadGameMode) FillChunkWithSectionAppropriateObjects(curChunk *CrossyRoad.Chunk) {

	previousMoveDir := -1
	for i := 0; i < len(curChunk.RowsAndObjects); i++ {

		// Decide section type
		// Get appropriate generation function based on section type
		// Generate objects using generation function
		// Update board with object data.

		curChunk.RowsAndObjects[i].ObjectsMoveDirectionX = previousMoveDir

		if previousMoveDir == -1 {
			previousMoveDir = 1
		} else {
			previousMoveDir = -1
		}

		sizeOfObjectsInRowInTiles := Engine.Vector2Int{X: 2 + rand.IntN(2), Y: 1}
		spacingBetweenObjectsInTiles := Engine.Vector2Int{X: 1 + rand.IntN(4), Y: 0}
		numObjectsInCurRow := (curChunk.ChunkBoard.Board_num_cols / (sizeOfObjectsInRowInTiles.X + spacingBetweenObjectsInTiles.X))

		rowObjectsStartingPosition := Engine.Vector2Int{X: 0, Y: i}

		for j := 0; j < numObjectsInCurRow; j++ {
			curObjectPoolItem := curChunk.RowsAndObjects[i].InRowObjectsPool.GetAnUnusedItemFromPool()
			curObjectPoolItem.Item.ObjectSizeInTiles = sizeOfObjectsInRowInTiles
			curObjectPoolItem.Item.HeadTilePositionIndex = rowObjectsStartingPosition

			curObjectTileIndex := rowObjectsStartingPosition
			for k := 0; k < sizeOfObjectsInRowInTiles.X; k++ {
				// fmt.Println(curObjectTileIndex)
				if curObjectTileIndex.X+k < curChunk.ChunkBoard.Board_num_cols {
					curChunk.ChunkBoard.BoardData[curObjectTileIndex.Y][curObjectTileIndex.X+k] = 5 + curObjectPoolItem.ItemIndex
				}
			}

			rowObjectsStartingPosition = Engine.Vector2Int{X: rowObjectsStartingPosition.X + sizeOfObjectsInRowInTiles.X, Y: rowObjectsStartingPosition.Y}
			rowObjectsStartingPosition = Engine.Add_Vector2Int(&rowObjectsStartingPosition, &spacingBetweenObjectsInTiles)
		}
	}
}

func (crossyRoadGameMode *CrossyRoadGameMode) GetCurrentCharacterScreenPos() *Engine.Vector2 {

	characterOffsetY := -((crossyRoadGameMode.characterInChunkIndex) * crossyRoadGameMode.worldGridSize.Y * crossyRoadGameMode.worldGridTileSizeInPixels)
	characterGenericPosY := (crossyRoadGameMode.worldGridSize.Y - crossyRoadGameMode.characterTilePosition.Y - 1) * crossyRoadGameMode.worldGridTileSizeInPixels
	characterPosY := characterGenericPosY + characterOffsetY
	characterPosX := crossyRoadGameMode.characterTilePosition.X * crossyRoadGameMode.worldGridTileSizeInPixels

	return &Engine.Vector2{X: float64(characterPosX), Y: float64(characterPosY)}
}

func (crossyRoadGameMode *CrossyRoadGameMode) GetBufferedInput(bufferedInput *Engine.Vector2Int) {

	bufferedInput.X = 0
	bufferedInput.Y = 0

	for i := range crossyRoadGameMode.curInputInBufferWriteIndex {
		bufferedInput.X += crossyRoadGameMode.inputBuffer[i].X
		bufferedInput.Y += crossyRoadGameMode.inputBuffer[i].Y
	}

	if bufferedInput.X != 0 {
		bufferedInput.X = int(math.Copysign(1, float64(bufferedInput.X)))
	}
	if bufferedInput.Y != 0 {
		bufferedInput.Y = int(math.Copysign(1, float64(bufferedInput.Y)))
	}

	crossyRoadGameMode.bufferInputClearResetTimerPoolItem.Item.ForceEndCurrentLoopOfTimerForNextUpdate()
}
