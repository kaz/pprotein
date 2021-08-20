package view

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var dist embed.FS

func FS() (fs.FS, error) {
	return fs.Sub(dist, "dist")
}
