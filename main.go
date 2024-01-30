package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func giveSVG(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-type", "image/svg+xml")
	queryValues := r.URL.Query()
	fmt.Fprint(w, GetSVG(queryValues.Get("url")))
}

func main() {
	router := httprouter.New()
	router.GET("/svg", giveSVG)

	log.Fatal(http.ListenAndServe(":80", router))
}
