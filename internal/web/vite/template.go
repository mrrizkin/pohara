package vite

import (
	"fmt"
	"io"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type vite struct {
	entry string
	token *tokens.Token
	vite  *Vite
}

type reactRefresh struct {
	token *tokens.Token
	vite  *Vite
}

func viteParser(v *Vite) func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	return func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
		cs := &vite{
			token: p.Current(),
			vite:  v,
		}

		entry := args.Match(tokens.String)
		if entry == nil {
			return nil, args.Error("A vite name is required for vite controlStructure.", nil)
		}

		cs.entry = entry.Val
		if !args.Stream().End() {
			return nil, args.Error("Malformed vite controlStructure args.", nil)
		}

		return cs, nil
	}
}

func (cs *vite) Position() *tokens.Token {
	return cs.token
}

func (cs *vite) String() string {
	t := cs.Position()
	return fmt.Sprintf("ViteControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (cs *vite) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	_, err := io.WriteString(r.Output, cs.vite.Entry(cs.entry))
	return err
}

func reactRefreshParser(v *Vite) func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
	return func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
		cs := &reactRefresh{
			token: p.Current(),
			vite:  v,
		}

		if !args.Stream().End() {
			return nil, args.Error("Malformed vite react refresh controlStructure args.", nil)
		}

		return cs, nil
	}
}

func (cs *reactRefresh) Position() *tokens.Token {
	return cs.token
}
func (cs *reactRefresh) String() string {
	t := cs.Position()
	return fmt.Sprintf("ReactRefreshControlStructure(Line=%d Col=%d)", t.Line, t.Col)
}

func (cs *reactRefresh) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	_, err := io.WriteString(r.Output, cs.vite.ReactRefresh())
	return err
}
