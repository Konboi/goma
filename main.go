package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"code.google.com/p/go.exp/fsnotify"
	"code.google.com/p/go.net/websocket"
	"github.com/russross/blackfriday"
)

type Markdown struct {
	Html template.HTML
}

var (
	path       = flag.String("path", "README.md", "Markdown file path")
	port       = 5858
	watcher, _ = fsnotify.NewWatcher()
)

func previewHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.New("goma")
	templ_file, err := Asset("index.html")

	if err != nil {
		fmt.Errorf("Load Asset Error: %s", err.Error())
	}

	t, err := templ.Parse(string(templ_file))

	if err != nil {
		fmt.Errorf("Load Template Error: %s", err.Error())
	}

	file, err := ioutil.ReadFile(*path)
	if err != nil {
		panic(err)
	}

	html := blackfriday.MarkdownCommon(file)
	md := Markdown{template.HTML(string(html))}

	t.Execute(w, md)
}

func reloadHandler(ws *websocket.Conn) {
	for {
		select {
		case ev := <-watcher.Event:
			if strings.Contains(ev.String(), "MODIFY") {
				websocket.Message.Send(ws, "update")
			}
		}
	}
}

func main() {
	flag.Parse()

	watcher.Watch(*path)
	http.HandleFunc("/", previewHandler)
	http.Handle("/ws", websocket.Handler(reloadHandler))

	http.ListenAndServe(":5858", nil)
}
