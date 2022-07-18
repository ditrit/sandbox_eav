package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ditrit/sandbox_eav/eav"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ErrMsg(msg string, w http.ResponseWriter) {
	w.Write(
		[]byte(fmt.Sprintf(
			`{"msg": %q}`,
			msg,
		)),
	)
}

func GetObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			ErrMsg("The id you provided is not an int", w)
			return
		}

		obj := eav.GetEntity(db, uint(id))
		w.Write(obj.EncodeToJson())
	}
}

func CreateObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		ErrMsg("The id you provided is not an int", w)
		return
	}
	println(id)
}
