package httproot

import "embed"

const (
	// AppDir directory contains app-specific files
	AppDir = "app"
)

//go:embed app/*
var FS embed.FS
