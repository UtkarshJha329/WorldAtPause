package Engine

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FontAssetData struct {
	FontName     string `json:"Font Name"`
	FontLocation string `json:"Font Location"`
}

type Font struct {
	FontName string
	Font     *text.GoTextFaceSource
	Face     *text.GoTextFace
}

func (font *Font) LoadFontWithFontFromPath(fontPath string) {

	fontFileContents := LoadFileFromFileSystem(fontPath)

	var err error
	font.Font, err = text.NewGoTextFaceSource(bytes.NewReader(fontFileContents))
	if err != nil {
		log.Fatal(err)
	}
}
