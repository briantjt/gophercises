package main

import (
	"encoding/json"
	"fmt"
	Story "gophercises/cyoa/story"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

func storyHandler(stories map[string]Story.Story) http.HandlerFunc {
	t, err := template.New("story").Parse(Story.StoryTemplate)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		base := path.Base(r.URL.Path)
		err := t.Execute(w, stories[base])
		if err != nil {
			http.NotFound(w, r)
		}
	}
}
func main() {
	var stories map[string]Story.Story
	jsonContents, err := ioutil.ReadFile("./gopher.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(jsonContents, &stories)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/arcs/", storyHandler(stories))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/arcs/intro", http.StatusMovedPermanently)
	})
	fmt.Println("Starting the server on :3000")
	http.ListenAndServe(":3000", nil)
}
