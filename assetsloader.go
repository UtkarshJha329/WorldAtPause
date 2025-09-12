package main

type Components struct {
	Type             string
	Position         *Vector2
	Sprite           *Sprite
	tilemapAssetData *TilemapAssetData
}

type Prefab struct {
	Name       string       `json:"Name"`
	EntityType string       `json:"Entity Type"`
	Components []Components `json:"Components"`
}

type Root struct {
	Prefabs  []map[string][]Prefab `json:"Prefabs"`
	Entities []string              `json:"Entities"`
}
