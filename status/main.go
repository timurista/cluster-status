package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// Site status
const (
	UP   = 1
	DOWN = 0
)

type StatusPageData struct {
	PageTitle string
	links     []string
}

func main() {

	tmpl := template.Must(template.ParseFiles("layout.html"))
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
		"http://twitter.com",
	}

	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c, tmpl)
	}

	for l := range c {
		go func(link string) {
			time.Sleep(time.Second * 5)
			checkLink(link, c, tmpl)
		}(l)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := StatusPageData{
			PageTitle: "Cluster Status Page",
			links:     links,
		}

		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":8081", nil)
}

func checkLink(link string, c chan string, tmpl *template.Template) {
	_, err := http.Get(link)
	// if err != nil {
	// 	return DOWN
	// }
	// return UP

	if err != nil {
		fmt.Println(link + " might be down")
		// tmpl.Execute(w, data)
		c <- link
	} else {
		fmt.Println(link + " is up")
		c <- link
	}
}
