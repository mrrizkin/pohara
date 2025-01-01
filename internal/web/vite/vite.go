package vite

import (
	"encoding/json"
	stdTempl "html/template"

	goviteparser "github.com/mrrizkin/go-vite-parser"
	"github.com/mrrizkin/pohara/internal/web/inertia"
	"github.com/mrrizkin/pohara/internal/web/template"
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
		template.Extend(func(v *Vite) template.ExtendResult {
			return template.NewExtend(map[string]interface{}{
				"vite":         v.Entry,
				"reactRefresh": v.ReactRefresh,
			})
		}),
		inertia.Extend(func(v *Vite) inertia.ExtendResult {
			return inertia.NewExtend(map[string]interface{}{
				"vite": func(input string) stdTempl.HTML {
					return stdTempl.HTML(v.Entry(input))
				},
				"reactRefresh": func() stdTempl.HTML {
					return stdTempl.HTML(v.ReactRefresh())
				},
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
