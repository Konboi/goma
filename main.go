package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"code.google.com/p/go.exp/fsnotify"
	"github.com/russross/blackfriday"
)

var (
	path       = flag.String("path", "README.md", "Markdown file path")
	port       = flag.Int("port", 8000, "Markdown file path")
	watcher, _ = fsnotify.NewWatcher()
)

func previewHandler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile(*path)

	if err != nil {
		panic(err)
	}
	html := blackfriday.MarkdownCommon(file)

	fmt.Fprintf(w, string(html))
}

func setWatcher() {
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if strings.Contains(ev.String(), "MODIFY") {
					fmt.Printf("hogehoeg \n")
				}
			}
		}
	}()
	watcher.Watch(*path)
}

func main() {
	flag.Parse()

	setWatcher()

	listen_port := ":" + strconv.Itoa(*port)
	http.HandleFunc("/", previewHandler)
	http.ListenAndServe(listen_port, nil)
}
