package main

import (
	"github.com/googege/goo"
	"net/http"
)

func main() {
	http.ListenAndServeTLS(":445", "", "", nil)
	goo.Join(nil, "", "")
}
