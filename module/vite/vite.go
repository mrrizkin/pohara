package vite

import (
	"encoding/json"

	goviteparser "github.com/mrrizkin/go-vite-parser"
	"go.uber.org/fx"
)

type Vite struct {
	vite *goviteparser.ViteManifestInfo
}

type Result struct {
	fx.Out

	Vite *Vite
}

func New() (Result, error) {
	vite := goviteparser.Parse(goviteparser.Config{
		OutDir:       "/build/",
		ManifestPath: "public/build/manifest.json",
		HotFilePath:  "public/hot",
	})
	return Result{
		Vite: &Vite{
			vite: &vite,
		},
	}, nil
}

func (v *Vite) Content() ([]byte, error) {
	return json.Marshal(v.vite.Manifest)
}

func (v *Vite) Entry(entries ...string) string {
	if v.vite.IsDev() {
		return v.vite.RenderDevEntriesTag(entries...)
	}

	return v.vite.RenderEntriesTag(entries...)
}

func (v *Vite) ReactRefresh() string {
	if v.vite.IsDev() {
		return v.vite.RenderReactRefreshTag()
	}

	return ""
}
