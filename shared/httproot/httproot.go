package httproot

import "embed"

const (
	// AppDir direcory contains app-specific files
	AppDir = "app"
	// VueUI directory contains common/shared files to create a Vue-SPA
	VueUI = "vue-ui"
)

//go:embed app/* vue-ui/*
var FS embed.FS
