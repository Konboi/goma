package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/russross/blackfriday"
)

var (
	path = flag.String("path", "README.md", "Markdown file path")
	port = flag.Int("port", 8000, "Markdown file path")
)

func previewHandler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	html := blackfriday.MarkdownCommon(file)

	fmt.Fprintf(w, string(html))
}

func main() {
	flag.Parse()

	listen_port := ":" + strconv.Itoa(*port)
	http.HandleFunc("/", previewHandler)
	http.ListenAndServe(listen_port, nil)
}
