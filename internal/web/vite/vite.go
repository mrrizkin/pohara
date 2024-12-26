package vite

import (
	"encoding/json"

	goviteparser "github.com/mrrizkin/go-vite-parser"
	"github.com/mrrizkin/pohara/internal/web/template"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/parser"
	"go.uber.org/fx"
)

type Vite struct {
	vite *goviteparser.ViteManifestInfo
}

type Result struct {
	fx.Out

	Vite *Vite
}

var Module = fx.Module("vite",
	fx.Provide(New),

	fx.Provide(
		template.AsControl(func(v *Vite) *exec.ControlStructureSet {
			return exec.NewControlStructureSet(map[string]parser.ControlStructureParser{
				"vite":         viteParser(v),
				"reactRefresh": reactRefreshParser(v),
			})
		}),
	),
)

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
