package handlers

import (
	"fmt"
	"net/http"
	"starapi/db"

	"strconv"

	"github.com/julienschmidt/httprouter"
)

func HandleGetUserData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var userNo uint64

	userNo, err := strconv.ParseUint(r.URL.Query().Get("user_no"), 10, 64)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	fmt.Fprintf(w, "user_no=%d", userNo)

	dbId := 1
	query := fmt.Sprintf("select cash, gold, exp from tb_user where user_no=%d", userNo)
	results := db.UserDb[dbId-1].Query(query)
	if results != nil {
		for e := results.Front(); e != nil; e = e.Next() {
			row := e.Value.(map[string]interface{})
			fmt.Printf("cash=%d, gold=%d, exp=%d\n", row["cash"], row["gold"], row["exp"])
		}
	}

}
