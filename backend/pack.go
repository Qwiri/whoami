package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Card struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type CardPack struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Avatars     []Card `json:"avatars"`
}

var Packs []*CardPack

func init() {
	matches, err := filepath.Glob("packs/*.json")
	if err != nil {
		panic(err)
	}
	for _, match := range matches {
		var data []byte
		if data, err = os.ReadFile(match); err != nil {
			panic(err)
		}
		var pack CardPack
		if err = json.Unmarshal(data, &pack); err != nil {
			panic(err)
		}
		Packs = append(Packs, &pack)
	}
}
