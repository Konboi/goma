package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/russross/blackfriday"
	"golang.org/x/net/websocket"
)

type Markdown struct {
	Html template.HTML
}

var (
	watcher, _ = fsnotify.NewWatcher()
)

func (h *Handler) previewHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.New("goma")
	templ_file, err := Asset("index.html")

	if err != nil {
		fmt.Fprintf(w, "Load Asset Error: "+err.Error())

	}

	t, err := templ.Parse(string(templ_file))

	if err != nil {
		fmt.Fprintf(w, "Load Template Error: "+err.Error())
	}

	file, err := os.ReadFile(h.path)
	if err != nil {
		panic(err)
	}

	html := blackfriday.MarkdownCommon(file)
	md := Markdown{template.HTML(string(html))}

	t.Execute(w, md)
}

func (h *Handler) reloadHandler(ws *websocket.Conn) {
	for {
		select {
		case ev := <-watcher.Events:
			if strings.Contains(ev.String(), "MODIFY") {
				websocket.Message.Send(ws, "update")
			}
		}
	}
}

type Handler struct {
	path string
}

func main() {
	var path string
	var port int
	flag.StringVar(&path, "path", "README.md", "Markdown file path")
	flag.IntVar(&port, "port", 5858, "Port number")
	flag.Parse()

	handler := &Handler{path: path}

	watcher.Add(path)
	http.HandleFunc("/", handler.previewHandler)
	http.Handle("/ws", websocket.Handler(handler.reloadHandler))

	fmt.Println(fmt.Sprintf("launched server at http://localhost:%d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
