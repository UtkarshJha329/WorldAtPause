package main

import (
	"WorldAtPause/Engine"
	"WorldAtPause/Game/MainGame"
	"WorldAtPause/Game/Minigames"
	"embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed Assets/*
var EmbeddedAssetsFS embed.FS

type Game struct {
	world *Engine.World
}

func (g *Game) Update() error {

	g.world.UpdateCurrentSceneGameMode()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	// TODO : Sort drawings based on Layers order and sort order index and then draw all at once.
	g.world.DrawCurrentSceneGameMode(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	entityComponentsRef := g.world.Scenes[g.world.CurrentSceneIndex].EntityComponentsForScene
	return int(entityComponentsRef.CameraData.ScreenSize.X), int(entityComponentsRef.CameraData.ScreenSize.Y)
}

func SetScenesGameModes(g *Game) {
	for _, scene := range g.world.Scenes {
		if scene.SceneName == "Main Game" {
			scene.GameMode = &MainGame.MainGameMode{
				World:    g.world,
				SceneRef: scene,
			}
		}
		if scene.SceneName == "Breakout Game" {
			scene.GameMode = &Minigames.BreakoutGameMode{
				World:    g.world,
				SceneRef: scene,
			}
		}
		if scene.SceneName == "Magical Ride Game" {
			scene.GameMode = &Minigames.MagicalRideGameMode{
				World:    g.world,
				SceneRef: scene,
			}
		}
		if scene.SceneName == "Space Invaders Game" {
			scene.GameMode = &Minigames.SpaceInvadersGameMode{
				World:    g.world,
				SceneRef: scene,
			}
		}
		if scene.SceneName == "Match Three Game" {
			scene.GameMode = &Minigames.MatchThreeGameMode{
				World:    g.world,
				SceneRef: scene,
			}
		}
		if scene.SceneName == "Sokoban Game" {
			scene.GameMode = &Minigames.SokobanGameMode{
				World:    g.world,
				SceneRef: scene,
			}
		}
	}
}

func main() {

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ninja!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	Engine.EmbeddedAssetsFS = EmbeddedAssetsFS

	game := Game{
		world: Engine.CreateAndPopulateWorldScenesAndEntitiesAndComponentsFromGameData("Assets/AssetsData.json"),
	}

	SetScenesGameModes(&game)

	game.world.CurrentSceneIndex = 0
	game.world.InitCurrentSceneGameMode()

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
