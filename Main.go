package main

import (
	"fmt"
	"net/http"
	"starapi/db"
	"starapi/handlers"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "Welcome")

}

func TestUrl(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "test")
	fmt.Fprintln(w, r.URL.Query().Get("test"))
}

var UserDb *db.Query

func main() {
	fmt.Println("Starwars api server start")

	UserDb = db.NewQuery()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/getuserdata", handlers.HandleGetUserData)

	http.ListenAndServe(":8000", router)
}
