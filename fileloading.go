package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed Assets/*
var EmbeddedAssetsFS embed.FS

func LoadImageFromFileSystem(imageLocation string) *ebiten.Image {
	spriteImage, _, err := ebitenutil.NewImageFromFileSystem(EmbeddedAssetsFS, imageLocation)
	if err != nil {
		log.Fatal("Failed to load sprite from" + imageLocation)
	}

	return spriteImage
}

func LoadFileFromFileSystem(pathOfFileToRead string) []byte {
	contents, err := fs.ReadFile(EmbeddedAssetsFS, pathOfFileToRead)
	if err != nil {
		fmt.Println(err, "Failed to read file: ", pathOfFileToRead)
	}
	return contents
}
