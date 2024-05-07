package renderer

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin/render"

	"github.com/a-h/templ"
)

var Default = &HTMLTemplRenderer{}

type HTMLTemplRenderer struct {
	FallbackHtmlRenderer render.HTMLRender
}

type Renderer struct {
	Ctx       context.Context
	Status    int
	Component templ.Component
}

func (r *HTMLTemplRenderer) Instance(s string, d any) render.Render {
	templData, ok := d.(templ.Component)
	if !ok {
		if r.FallbackHtmlRenderer != nil {
			return r.FallbackHtmlRenderer.Instance(s, d)
		}
	}
	return &Renderer{
		Ctx:       context.Background(),
		Status:    -1,
		Component: templData,
	}
}

func New(ctx context.Context, status int, component templ.Component) *Renderer {
	return &Renderer{
		Ctx:       ctx,
		Status:    status,
		Component: component,
	}
}

func (r Renderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (r Renderer) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	if r.Status != -1 {
		w.WriteHeader(r.Status)
	}
	if r.Component != nil {
		return r.Component.Render(r.Ctx, w)
	}
	return nil
}
