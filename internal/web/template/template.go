package template

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"github.com/mrrizkin/pohara/internal/ports"
	"github.com/mrrizkin/pohara/web"
	"github.com/nikolalohinski/gonja/v2/builtins"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/parser"
	"go.uber.org/fx"
)

type Template struct {
	fs     http.FileSystem
	config *config.App
	cache  ports.Cache
	log    ports.Logger
	env    *exec.Environment

	templates map[string]*exec.Template
}

type Dependencies struct {
	fx.In

	Config *config.App
	Cache  ports.Cache
	Log    ports.Logger
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

type ExtendResult struct {
	fx.Out

	Context          *exec.Context             `group:"ctx"`
	Filter           *exec.FilterSet           `group:"filter"`
	Test             *exec.TestSet             `group:"test"`
	ControlStructure *exec.ControlStructureSet `group:"control"`
}

func Extend(fn any) any {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("extend expects a function")
	}

	// Validate that the function returns exactly one value
	if fnType.NumOut() != 1 {
		panic("function must return 1 thing")
	}

	// Validate that the return value is a struct
	if fnType.Out(0).Kind() != reflect.Struct {
		panic("function must return a struct")
	}

	// Validate that the return type matches ExtendResult
	expectedType := reflect.TypeOf(ExtendResult{})
	if !fnType.Out(0).AssignableTo(expectedType) {
		panic("function must return ExtendResult")
	}

	return fn
}

func NewExtend(data map[string]interface{}) ExtendResult {
	ctx := make(map[string]interface{})
	control := make(map[string]parser.ControlStructureParser)
	filter := make(map[string]exec.FilterFunction)
	test := make(map[string]exec.TestFunction)

	for key, value := range data {
		// test if the value is a parser.ControlStructureParser
		if cs, ok := value.(parser.ControlStructureParser); ok {
			control[key] = cs
			continue
		}

		// test if the value is a exec.FilterFunction
		if f, ok := value.(exec.FilterFunction); ok {
			filter[key] = f
			continue
		}

		// test if the value is a exec.TestFunction
		if t, ok := value.(exec.TestFunction); ok {
			test[key] = t
			continue
		}

		ctx[key] = value
	}

	return ExtendResult{
		Context:          exec.NewContext(ctx),
		Filter:           exec.NewFilterSet(filter),
		Test:             exec.NewTestSet(test),
		ControlStructure: exec.NewControlStructureSet(control),
	}
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
			cache:  deps.Cache,
			log:    deps.Log,

			fs: fs,
		},
	}
}

func (t *Template) Render(
	c *fiber.Ctx,
	template string,
	data map[string]interface{},
) error {
	cacheKey := t.cacheKey(template, data)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if cacheKey.Valid {
		if value, ok := t.cache.Get(ctx, cacheKey.String); ok {
			if html, ok := value.([]byte); ok {
				return c.Type("html").Send(html)
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
	if t.config.IsCacheView() && cacheKey.Valid {
		t.cache.Set(ctx, cacheKey.String, html)
	}

	return c.Type("html").Send(html)
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

func (t *Template) cacheKey(template string, data map[string]interface{}) sql.StringNullable {
	if !t.config.IsCacheView() {
		return sql.StringNull()
	}

	var cacheKey string
	var hash [16]byte
	if data != nil {
		encodedData, err := json.Marshal(data)
		if err != nil {
			return sql.StringNull()
		}

		hash = md5.Sum(append([]byte(template), encodedData...))
	} else {
		hash = md5.Sum([]byte(template))
	}

	cacheKey = hex.EncodeToString(hash[:])

	return sql.String(cacheKey)
}
