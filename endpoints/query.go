package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ditrit/sandbox_eav/eav/query"
	"gorm.io/gorm"
)

func Query(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, GetErrMsg("can't open request body"), http.StatusInternalServerError)
			return
		}
		var q query.Query
		err = json.Unmarshal(content, &q)
		if err != nil {
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}

		// fmt.Printf("Unmarshalled data: ")
		// debugS(&q)
		payload, err := q.Run(db)
		if err != nil {
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}

		// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(payload)
	}
}

// print json eq to stdout
func debugS(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
