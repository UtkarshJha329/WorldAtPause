package main

import (
	"encoding/json"
	"os"
)

type Components struct {
	Type     string
	Position *Vector2
	Sprite   *Sprite
}

func (c *Components) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, val := range raw {
		c.Type = key
		switch key {
		case "Position":
			var pos Vector2
			if err := json.Unmarshal(val, &pos); err != nil {
				return err
			}
			c.Position = &pos
		case "Sprite":
			var spr Sprite
			if err := json.Unmarshal(val, &spr); err != nil {
				return err
			}
			c.Sprite = &spr
		}
	}
	return nil
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

func LoadGameData(path string) (map[string]Prefab, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var root Root
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, nil, err
	}

	prefabMap := make(map[string]Prefab)

	for _, entry := range root.Prefabs {
		for _, prefabList := range entry {
			for _, p := range prefabList {
				prefabMap[p.Name] = p
			}
		}
	}

	return prefabMap, root.Entities, nil
}
