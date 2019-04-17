package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jorygeerts/perflog"
	"log"
	"net/http"
)

func main() {
	log.Printf(`Start serving`)

	store := perflog.NewStore()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HTML that does some Vue.js stuff?"))
	})

	r.Route(`/apps`, func(r chi.Router) {
		r.Get(`/`, func(writer http.ResponseWriter, request *http.Request) {
			rl := []render.Renderer{}

			for _, project := range store.GetProjects() {
				rl = append(rl, appRenderer{Id: project.Id, Name: project.Name})
			}

			render.RenderList(writer, request, rl)
		})

		r.Post(`/`, func(writer http.ResponseWriter, request *http.Request) {

			data := &newAppRequest{}
			if err := render.Decode(request, data); err != nil {
				render.Render(writer, request, ErrInvalidRequest(err))
				return
			}
			store.AddProject(data.Id, data.Name)
		})
	})

	http.ListenAndServe(`:1234`, r)
}

type appRenderer struct {
	Id   string
	Name string
}

func (ar appRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type newAppRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
