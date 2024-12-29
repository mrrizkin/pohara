package template

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/web"
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/exec"
	"go.uber.org/fx"
)

type Template struct {
	fs     http.FileSystem
	config *config.App
	env    *exec.Environment

	templates map[string]*exec.Template
}

type Dependencies struct {
	fx.In

	Config *config.App
}

type Result struct {
	fx.Out

	Template *Template
}

var Module = fx.Module("template",
	fx.Provide(New),
	fx.Decorate(extendTemplate),
	fx.Invoke(loader),
)

func AsCtx(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"ctx"`),
	)
}

func AsFilter(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"filter"`),
	)
}

func AsTest(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"test"`),
	)
}

func AsControl(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"control"`),
	)
}

func New(deps Dependencies) Result {
	fs := http.FS(web.Views)
	env := &exec.Environment{
		Context: exec.EmptyContext().
			Update(builtins.GlobalFunctions).
			Update(builtins.GlobalVariables),
		Filters:           builtins.Filters,
		Tests:             builtins.Tests,
		ControlStructures: builtins.ControlStructures,
		Methods:           builtins.Methods,
	}

	return Result{
		Template: &Template{
			config: deps.Config,
			env:    env,

			fs: fs,
		},
	}
}

func (t *Template) Render(
	template string,
	data map[string]interface{},
) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, ok := t.templates[template]
	if !ok {
		return nil, fmt.Errorf("template %s not found", template)
	}

	err := tmpl.Execute(&buf, exec.NewContext(data))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type ExtendTemplateDependencies struct {
	fx.In

	Template         *Template
	Context          []*exec.Context             `group:"ctx"`
	Filter           []*exec.FilterSet           `group:"filter"`
	Test             []*exec.TestSet             `group:"test"`
	ControlStructure []*exec.ControlStructureSet `group:"control"`
}

func extendTemplate(deps ExtendTemplateDependencies) *Template {
	template := deps.Template
	for _, c := range deps.Context {
		if c != nil {
			template.env.Context.Update(c)
		}
	}

	for _, f := range deps.Filter {
		if f != nil {
			template.env.Filters.Update(f)
		}
	}

	for _, t := range deps.Test {
		if t != nil {
			template.env.Tests.Update(t)
		}
	}

	for _, cs := range deps.ControlStructure {
		if cs != nil {
			template.env.ControlStructures.Update(cs)
		}
	}

	return template
}

func loader(t *Template) error {
	templates := make(map[string]*exec.Template)
	loader, err := newHttpFileSystemLoader(t.fs, t.config.VIEW_DIRECTORY)
	if err != nil {
		return err
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil || info.IsDir() {
			return nil
		}

		if len(t.config.VIEW_EXTENSION) >= len(path) ||
			path[len(path)-len(t.config.VIEW_EXTENSION):] != t.config.VIEW_EXTENSION {
			return nil
		}

		rel, err := filepath.Rel(t.config.VIEW_DIRECTORY, path)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(filepath.ToSlash(rel), t.config.VIEW_EXTENSION)

		file, err := t.fs.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		buf, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		tmpl, err := fromBytes(buf, loader, t.env)
		if err != nil {
			return err
		}

		templates[name] = tmpl

		return err
	}

	info, err := stat(t.fs, t.config.VIEW_DIRECTORY)
	if err != nil {
		return err
	}
	if err := walkInternal(t.fs, t.config.VIEW_DIRECTORY, info, walkFn); err != nil {
		return err
	}

	t.templates = templates

	return nil
}
