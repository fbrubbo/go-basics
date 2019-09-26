package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello</h1>")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>About</h1>")
	fmt.Fprintf(w, `
		<h6>You can even do ...</h6>
		<h5>multiple lines ...</h5>
		<h4>in one %s</h4>`, "formatted print")
}

type NewsAggPage struct {
	Title string
	News  string
}

func templ(w http.ResponseWriter, r *http.Request) {
	p := NewsAggPage{Title: "Amazing News Aggregator", News: "some news"}
	t, err := template.ParseFiles("basictemplating.html")
	fmt.Println(err)
	err = t.Execute(w, p)
	fmt.Println(err)
}

func main() {

	resp, err := http.Get("https://godoc.org/fmt")
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(bytes))
			fmt.Printf("Type: %T\n", bytes)
		}
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
	http.HandleFunc("/templ", templ)
	e := http.ListenAndServe(":3000", nil)
	if e != nil {
		fmt.Println(e)
	}

}
