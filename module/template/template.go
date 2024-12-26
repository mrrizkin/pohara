package template

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/module/cache"
	"github.com/mrrizkin/pohara/module/database/sql"
	"github.com/mrrizkin/pohara/module/logger"
	"github.com/mrrizkin/pohara/resources"
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/exec"
	"go.uber.org/fx"
)

type Template struct {
	fs     http.FileSystem
	config *config.App
	cache  *cache.Cache
	log    *logger.Logger
	env    *exec.Environment

	templates map[string]*exec.Template
}

type Dependencies struct {
	fx.In

	Config *config.App
	Cache  *cache.Cache
	Log    *logger.Logger
}

type Result struct {
	fx.Out

	Template *Template
}

var Module = fx.Module("template",
	fx.Provide(New),
	fx.Decorate(
		fx.Annotate(setupEnv, fx.ParamTags("",
			`group:"ctx"`,
			`group:"filter"`,
			`group:"test"`,
			`group:"control"`,
		)),
	),
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
	fs := http.FS(resources.Views)
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
			cache:  deps.Cache,
			log:    deps.Log,
			env:    env,

			fs: fs,
		},
	}
}

func (t *Template) cacheKey(template string, data map[string]interface{}) sql.StringNullable {
	if !t.config.VIEW_CACHE || !t.config.IsProduction() {
		return sql.StringNullable{
			Valid: false,
		}
	}

	var cacheKey string
	var hash [16]byte
	if data != nil {
		encodedData, err := json.Marshal(data)
		if err != nil {
			return sql.StringNullable{
				Valid: false,
			}
		}

		hash = md5.Sum(append([]byte(template), encodedData...))
	} else {
		hash = md5.Sum([]byte(template))
	}

	cacheKey = hex.EncodeToString(hash[:])

	return sql.StringNullable{
		Valid:  true,
		String: cacheKey,
	}
}

func (t *Template) Render(
	ctx *fiber.Ctx,
	template string,
	data map[string]interface{},
) error {
	cacheKey := t.cacheKey(template, data)
	if cacheKey.Valid {
		if value, ok := t.cache.Get(cacheKey.String); ok {
			if html, ok := value.([]byte); ok {
				return ctx.Type("html").Send(html)
			}

			t.log.Warn("cached template invalid type", "template", template)
		}
	}

	var buf bytes.Buffer
	tmpl, ok := t.templates[template]
	if !ok {
		return fmt.Errorf("template %s not found", template)
	}

	err := tmpl.Execute(&buf, exec.NewContext(data))
	if err != nil {
		return err
	}

	html := buf.Bytes()
	if t.config.VIEW_CACHE && t.config.IsProduction() && cacheKey.Valid {
		t.cache.Set(cacheKey.String, html)
	}

	return ctx.Type("html").Send(html)
}

func setupEnv(
	template *Template,
	ctx []*exec.Context,
	filter []*exec.FilterSet,
	test []*exec.TestSet,
	controlStructure []*exec.ControlStructureSet,
) *Template {
	for _, c := range ctx {
		if c != nil {
			template.env.Context.Update(c)
		}
	}

	for _, f := range filter {
		if f != nil {
			template.env.Filters.Update(f)
		}
	}

	for _, t := range test {
		if t != nil {
			template.env.Tests.Update(t)
		}
	}

	for _, cs := range controlStructure {
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
