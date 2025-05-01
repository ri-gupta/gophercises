package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type nextArc struct {
	Text string
	Arc  string
}

type storyArc struct {
	Title   string
	Story   []string
	Options []nextArc
}

func renderChapter(story map[string]storyArc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storyPath := r.URL.Path[1:] // extract intro from /intro

		if storyDetails, ok := story[storyPath]; ok {

			t, err := template.New("index.html").Parse(storyHTMLTemplate())
			if err != nil {
				fmt.Fprintln(w, "Hello reader, this is the wrong path, please check again!")
				return
			}

			err = t.Execute(w, storyDetails)
			if err != nil {
				fmt.Fprintln(w, "Hello reader, this is the wrong path, please check again!")
				return
			}

			return
		}

		http.Redirect(w, r, "/intro", http.StatusFound)
	}
}

func storyHTMLTemplate() string {
	return `
	<!DOCTYPE html>
	<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			{{range .Story}}
				<p>{{.}}</p>
			{{end}}

			<ul>
				{{range .Options}}
					<li> 
						<a href=/{{.Arc}}>{{.Text}}</a> 
					</li>
				{{end}}
			</ul>
		</body>
	</html>`
}

func main() {
	story := make(map[string]storyArc)

	storyData, err := os.ReadFile("gopher.json")
	if err != nil {
		fmt.Println("Error reading story json file: ", err)
	}

	err = json.Unmarshal(storyData, &story)
	if err != nil {
		fmt.Println("Error parsing json: ", err)
	}

	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", renderChapter(story))
}
