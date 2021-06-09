package ariang

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Handler{})
	httpcaddyfile.RegisterHandlerDirective("yacd", parseCaddyfile)
}

// content holds our static web server content.
// AriaNg 1.2.1
//go:embed www/*
var www embed.FS

// Handler implements an HTTP handler that ...
type Handler struct {
	Prefix string `json:"prefix,omitempty"`

	Handler http.Handler
}

// CaddyModule returns the Caddy module information.
func (Handler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.yacd",
		New: func() caddy.Module { return new(Handler) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Handler) Provision(ctx caddy.Context) error {
	m.Handler = http.FileServer(http.FS(FS(www)))
	if m.Prefix != "" {
		m.Handler = http.StripPrefix(m.Prefix, m.Handler)
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) (err error) {
	m.Handler.ServeHTTP(w, r)
	return nil
}

// UnmarshalCaddyfile unmarshals Caddyfile tokens into h.
func (h *Handler) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() {
		return d.ArgErr()
	}
	args := d.RemainingArgs()
	if len(args) > 0 {
		return d.ArgErr()
	}
	for nesting := d.Nesting(); d.NextBlock(nesting); {
		subdirective := d.Val()
		args := d.RemainingArgs()
		switch subdirective {
		case "prefix":
			if len(args) != 1 {
				return d.ArgErr()
			}
			if args[0] == "" {
				return errors.New("empty prefix")
			}
			if h.Prefix != "" {
				return errors.New("only one prefix is allowed")
			}
			h.Prefix = args[0]
		}
	}
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	m := &Handler{}
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Handler)(nil)
	_ caddyhttp.MiddlewareHandler = (*Handler)(nil)
	_ caddyfile.Unmarshaler       = (*Handler)(nil)
)

// FS is ...
type FS embed.FS

// Open is ...
func (fs FS) Open(name string) (fs.File, error) {
	return embed.FS(fs).Open(filepath.Join("www", filepath.FromSlash(path.Clean("/"+name))))
}
