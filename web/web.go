package web

import "embed"

//go:embed views
var Views embed.FS

//go:embed inertia
var InertiaRoot embed.FS
