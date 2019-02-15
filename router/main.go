package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	router := New(Root)
	router.GET("/hello", Hello)
	router.GET("/hello/:name", HelloUser)
	http.ListenAndServe(":8080", router)
}

// Root handler
func Root(w http.ResponseWriter, r *http.Request, params url.Values) {
	fmt.Println("Root!")
}

// Hello hello world
func Hello(w http.ResponseWriter, r *http.Request, params url.Values) {
	fmt.Println("Hello World")
}

// HelloUser greets the user
func HelloUser(w http.ResponseWriter, r *http.Request, params url.Values) {
	fmt.Println("Hello User", params["name"])
}
