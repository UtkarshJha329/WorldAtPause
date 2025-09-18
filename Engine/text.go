package Engine

type TextAssetData struct {
	FontName   string `json:"Font Name"`
	TextName   string `json:"Text Name"`
	TextString string `json:"Text String"`
}

type Text struct {
	Font     *Font
	TextName string
}
