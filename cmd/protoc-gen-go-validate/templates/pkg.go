// Package templates assembles the per-language template set.
//
// This is the Go-only subset of envoyproxy/protoc-gen-validate's
// templates/pkg.go. The upstream file also registers the cc, ccnop and java
// sets, which would pull their whole template trees in as dependencies for a
// plugin that only ever generates Go. Trimming to Go is why this file exists
// rather than importing the upstream package.
package templates

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	golang "github.com/onexstack/protobuf-go-plugins/cmd/protoc-gen-go-validate/templates/go"
	"github.com/onexstack/protobuf-go-plugins/cmd/protoc-gen-go-validate/templates/shared"
)

type (
	RegisterFn func(tpl *template.Template, params pgs.Parameters)
	FilePathFn func(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath
)

func makeTemplate(ext string, fn RegisterFn, params pgs.Parameters) *template.Template {
	tpl := template.New(ext)
	shared.RegisterFunctions(tpl, params)
	fn(tpl, params)
	return tpl
}

func Template(params pgs.Parameters) map[string][]*template.Template {
	return map[string][]*template.Template{
		"go": {makeTemplate("go", golang.Register, params)},
	}
}

// FilePathFor is the Go branch of the upstream switch: the output sits beside
// the source and takes a .validate.go extension.
func FilePathFor(tpl *template.Template) FilePathFn {
	return func(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath {
		out := ctx.OutputPath(f)
		out = out.SetExt(".validate." + tpl.Name())
		return &out
	}
}
