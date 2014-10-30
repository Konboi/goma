package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/russross/blackfriday"
)

var (
	path = flag.String("path", "README.md", "Markdown file path")
	port = flag.Int("port", 8000, "Markdown file path")
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	file, err := GetFile(*path)
	if err != nil {
		panic(err)
	}
	html := blackfriday.MarkdownCommon(file)

	fmt.Fprintf(w, string(html))
}

func GetFile(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func main() {
	flag.Parse()

	http.HandleFunc("/", formHandler)
	http.ListenAndServe(":8000", nil)
}
