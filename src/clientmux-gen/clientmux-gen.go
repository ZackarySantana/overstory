// Package clientmuxgen turns a *ClientMux into strongly-typed
// client code.  Call Generate(mux, pkg, varBaseURL) to obtain ready-to-write
// Go source.
//
// The output looks like:
//
//	package api
//	type Client struct { base string; hc *http.Client }
//	func New(base string, hc *http.Client) *Client { … }
//	func (c *Client) CreateUser(ctx context.Context, in CreateUserReq)
//	     (CreateUserResp, error) { … }
//
// Each route gets one method whose signature matches the types recorded
// in RouteMeta.
//
// The generator keeps things simple on purpose: JSON body, http.Client,
// ctx-first signature, and error handling that bubbles up non-2xx as
// *StatusError (Status(), Error()).
package clientmux

import (
	"bytes"
	"go/format"
	"net/http"
	"text/template"

	"github.com/zackarysantana/overstory/src/clientmux"
)

// Generate builds client source for the supplied mux.
func Generate(mux *clientmux.ClientMux, pkgName, pkgImp string) ([]byte, error) {
	const fileTmpl = `
package {{.Pkg}}

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"{{.PkgImp}}"
)

// Code-gen: DO NOT EDIT.

type Client struct {
	base string
	hc   *http.Client
}

func New(baseURL string, hc *http.Client) *Client {
	if hc == nil { hc = http.DefaultClient }
	return &Client{base: baseURL, hc: hc}
}

type StatusError struct {
	Code int
	Body string
}

func (e *StatusError) Error() string { return fmt.Sprintf("http %d: %s", e.Code, e.Body) }
func (e *StatusError) Status() int   { return e.Code }

{{range .Routes}}
func (c *Client) {{.Func}}(ctx context.Context, in {{.Req}}) ({{.Resp}}, error) {
	var buf bytes.Buffer
	{{if .HasBody}}
	if err := json.NewEncoder(&buf).Encode(in); err != nil { return {{.ZeroResp}}, err }
	req, err := http.NewRequestWithContext(ctx, "{{.Method}}", c.base+"{{.Pattern}}", &buf)
	{{else}}
	req, err := http.NewRequestWithContext(ctx, "{{.Method}}", c.base+"{{.Pattern}}", nil)
	{{end}}
	if err != nil { return {{.ZeroResp}}, err }
	req.Header.Set("Content-Type", "application/json")

	res, err := c.hc.Do(req)
	if err != nil { return {{.ZeroResp}}, err }
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return {{.ZeroResp}}, &StatusError{Code: res.StatusCode, Body: string(body)}
	}

	var out {{.Resp}}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return {{.ZeroResp}}, err
	}
	return out, nil
}
{{end}}
`
	type routeCtx struct {
		Method, Pattern     string
		Func                string
		Req, Resp, ZeroResp string
		HasBody             bool
	}
	type ctx struct {
		Pkg    string
		Routes []routeCtx
		PkgImp string
	}

	c := ctx{Pkg: pkgName, PkgImp: pkgImp}
	for _, rt := range mux.Routes() {
		reqT, respT := rt.ReqType.String(), rt.RespType.String()
		fn := exportable(methodName(rt.Pattern))
		c.Routes = append(c.Routes, routeCtx{
			Method:   rt.Method,
			Pattern:  rt.Pattern,
			Func:     fn,
			Req:      reqT,
			Resp:     respT,
			ZeroResp: zero(respT),
			HasBody:  rt.Method != http.MethodGet && rt.Method != http.MethodHead,
		})
	}

	var buf bytes.Buffer
	if err := template.Must(template.New("").Parse(fileTmpl)).Execute(&buf, c); err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}

// Helpers -------------------------------------------------------------

func zero(typ string) string {
	switch typ {
	case "string":
		return `""`
	case "int", "int64", "float64":
		return "0"
	case "bool":
		return "false"
	default:
		return typ + "{}"
	}
}

func methodName(p string) string {
	// "/users/{id}/reset" → "UsersIdReset"
	var out []rune
	for i, r := range p {
		switch {
		case r == '/' || r == '-' || r == '_' || r == '{' || r == '}':
			continue
		default:
			if i == 0 || p[i-1] == '/' {
				r -= 32 // crude upper-case
			}
			out = append(out, r)
		}
	}
	if len(out) == 0 {
		return "Root"
	}
	return string(out)
}

func exportable(s string) string { return s } // placeholder for name-collision handling
