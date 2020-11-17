package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Site status
const (
	UP   = 1
	DOWN = 0
)

type StatusPageData struct {
	PageTitle         string
	Links             []string
	WebsocketLocation string
}

func main() {

	tmpl := template.Must(template.ParseFiles("pages/layout.html"))
	port := "8081"
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
		"http://twitter.com",
	}

	// c := make(chan string)

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			time.Sleep(time.Second * 5)

			c := make(chan string)

			for _, link := range links {
				go checkLink(link, c, conn)
			}

			for l := range c {
				go func(link string) {
					time.Sleep(time.Second * 5)
					checkLink(link, c, conn)
				}(l)
			}

		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("links")
		fmt.Println(links)
		// linkStrings := []string{}
		// for _, link := range links {
		// 	_, err := http.Get(link)
		// 	status := ""
		// 	if err != nil {
		// 		status = "X"
		// 	} else {
		// 		status = "UP"
		// 	}
		// 	linkStrings = append(linkStrings, status+" "+link)
		// }

		data := StatusPageData{
			PageTitle:         "Cluster Status Page",
			Links:             links,
			WebsocketLocation: "localhost:" + port,
		}

		tmpl.Execute(w, data)

	})
	http.ListenAndServe(":"+port, nil)
	fmt.Println("listening at " + port)
}

func checkLink(link string, c chan string, conn *websocket.Conn) {
	_, err := http.Get(link)
	// if err != nil {
	// 	return DOWN
	// }
	// return UP

	// msgType := websocket.TextMessage

	if err != nil {
		var v struct {
			link   string
			status string
		}
		// tmpl.Execute(w, data)
		v.link = link
		v.status = "down"
		// Partial JSON values.

		// data, _ := json.Marshal(v)

		msg := link + " might be down"
		if err := conn.WriteJSON(&v); err != nil {
			return
		}
		fmt.Println(msg)
		// tmpl.Execute(w, data)
		c <- link
	} else {
		type v struct {
			Link        string `json:"link"`
			Status      string `json:"status"`
			LastUpdated string `json:"lastUpdated"`
		}
		// v.link = link
		// v.status = "up"
		// Partial JSON values.
		m := v{
			Link:        link,
			Status:      "up",
			LastUpdated: time.Now().Local().String(),
		}
		// data, _ := json.Marshal(v)

		// msg := link + "is up"
		if err := conn.WriteJSON(m); err != nil {
			return
		}
		fmt.Println(m)
		c <- link
	}
}
