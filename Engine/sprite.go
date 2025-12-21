package Engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type AnimationFrameData struct {
	StartFramePos Vector2Int `json:"StartFramePos"`
	EndFramePos   Vector2Int `json:"EndFramePos"`
}

type AnimationData struct {
	PerFrameTime float64 `json:"PerFrameTime"`

	AnimationFramesData []AnimationFrameData `json:"AnimationFramesData"`
}

type Animation struct {
	currentFrameCounter uint
	AnimationTimer      *PoolItem[Timer]
	animationFramesData []AnimationFrameData
}

type Sprite struct {
	Image           *ebiten.Image
	RenderRectStart Vector2Int
	RenderRectEnd   Vector2Int

	CurrentAnimationIndex uint
	Animations            []Animation
}

type SpriteAssetData struct {
	SpriteTextureLocation string          `json:"SpriteTextureLocation"`
	RenderRectStart       Vector2Int      `json:"RenderRectStart"`
	RenderRectEnd         Vector2Int      `json:"RenderRectEnd"`
	AnimationsData        []AnimationData `json:"AnimationsData"`
}

func (sprite *Sprite) ChangeSpriteAnimationIndexTo(index uint) {
	if sprite.CurrentAnimationIndex != index {
		sprite.Animations[sprite.CurrentAnimationIndex].AnimationTimer.Item.PauseTimer()
		sprite.CurrentAnimationIndex = index
		sprite.Animations[sprite.CurrentAnimationIndex].currentFrameCounter = 1
		sprite.Animations[sprite.CurrentAnimationIndex].AnimationTimer.Item.RestartTimer()
		sprite.RenderRectStart = sprite.Animations[sprite.CurrentAnimationIndex].animationFramesData[0].StartFramePos
		sprite.RenderRectEnd = sprite.Animations[sprite.CurrentAnimationIndex].animationFramesData[0].EndFramePos
	} else {
		sprite.Animations[sprite.CurrentAnimationIndex].AnimationTimer.Item.UnPauseTimer()
	}
}
