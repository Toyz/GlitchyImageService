package engine

import (
	"image"
)

type Packet struct {
	ID   int
	To   Glitch
	From Glitch
}

type Glitch struct {
	Name        string
	Mime        string
	IsGif       bool
	Bounds      image.Rectangle
	Expressions []string
	File        []byte
}
