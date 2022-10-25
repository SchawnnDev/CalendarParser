package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"time"
)

func startHttpServer(addr string) {
	println("Creating router...")
	r := chi.NewRouter()

	// Register middlewares
	registerMiddlewares(r)

	// Register routes
	registerRoutes(r)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("welcome"))
	})
	println("Starting HTTP server at " + addr)
	http.ListenAndServe(addr, r)
}

func registerMiddlewares(r *chi.Mux) {
	fmt.Printf("Registering middlewares...")
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	fmt.Printf(" done!\n")
}

func registerRoutes(r *chi.Mux) {
	fmt.Printf("Registering routes...")

	r.Route("/process", func(r chi.Router) {
		r.Post("/html", ImportHtml)
	})

	fmt.Printf(" done!\n")
}

func ImportHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Maximum 10 MB files.
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Could not parse Multipart Form", http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("data")
	defer file.Close()

	if err != nil {
		http.Error(w, "Could not process input file", http.StatusBadRequest)
		return
	}

	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, file); err != nil {
		http.Error(w, "Could not process input file", http.StatusBadRequest)
		return
	}

	data := parse(buf.Bytes())

	for _, d := range data {
		fmt.Printf("%s\n", d.Date)
	}

	result, _ := json.Marshal(data)

	w.Write(result)
}
