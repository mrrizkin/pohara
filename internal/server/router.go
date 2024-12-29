package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
	"go.uber.org/fx"
)

type WebRouter struct {
	prefix string
	router func(*Router)
	name   []string
}

func NewWebRouter(prefix string, router func(*Router), names ...string) WebRouter {
	return WebRouter{
		prefix: prefix,
		router: router,
		name:   names,
	}
}

type WebRouters []WebRouter

func (wr WebRouters) Register(router *Router) {
	for _, r := range wr {
		router.Route(r.prefix, r.router, r.name...)
	}
}

func AsWebRouter(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"web_router"`),
	)
}

type ApiRouter struct {
	version string
	prefix  string
	router  func(*Router)
	name    []string
}

func NewApiRouter(version, prefix string, router func(*Router), names ...string) ApiRouter {
	return ApiRouter{
		prefix:  prefix,
		router:  router,
		name:    names,
		version: version,
	}
}

func AsApiRouter(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"api_router"`),
	)
}

type ApiRouters []ApiRouter

func (ar ApiRouters) Register(router *Router) {
	for _, r := range ar {
		prefix, err := url.JoinPath(r.version, r.prefix)
		if err != nil {
			continue
		}

		router.Route(prefix, r.router, r.name...)
	}
}

type ziggy struct {
	token *tokens.Token
	app   *fiber.App
}

func ziggyParser(
	v *fiber.App,
) func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	return func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
		cs := &ziggy{
			token: p.Current(),
			app:   v,
		}

		if !args.Stream().End() {
			return nil, args.Error("Malformed Ziggy controlStructure args.", nil)
		}

		return cs, nil
	}
}

func (cs *ziggy) Position() *tokens.Token {
	return cs.token
}
func (cs *ziggy) String() string {
	t := cs.Position()
	return fmt.Sprintf("ZiggyControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (cs *ziggy) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	_, err := io.WriteString(r.Output, cs.script())
	return err
}

func (cs *ziggy) script() string {
	routeMap := make(map[string]string)
	routesStack := cs.app.Stack()
	for _, routes := range routesStack {
		for _, route := range routes {
			if route.Name != "" {
				routeMap[route.Name] = route.Path
			}
		}
	}

	data, err := json.Marshal(routeMap)
	if err != nil {
		data = []byte("{}")
	}
	return fmt.Sprintf(`function $route(name) { return %s[name] }`, data)
}
