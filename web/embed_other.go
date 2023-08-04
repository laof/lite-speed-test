//go:build !gomobile
// +build !gomobile

package web

import "embed"

//go:embed gui/dist/*
var guiStatic embed.FS

var FontBytes []byte
