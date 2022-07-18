package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ditrit/sandbox_eav/eav"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ErrMsg(msg string, w http.ResponseWriter) {
	w.Write(
		[]byte(fmt.Sprintf(
			`{"error": %q}`,
			msg,
		)),
	)
}

func GetObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := eav.GetEntityTypeByName(db, entityType)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
				return
			}
		}
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			ErrMsg("The id you provided is not an int", w)
			return
		}

		obj, err := eav.GetEntity(db, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
				return
			}
		}
		if obj.EntityTypeId != ett.ID {
			http.NotFound(w, r)
			ErrMsg("This object doesn't belong to this type", w)
		}

		w.Write(obj.EncodeToJson())
	}
}

func DeleteObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := eav.GetEntityTypeByName(db, entityType)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
				return
			}
		}
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			ErrMsg("The id you provided is not an int", w)
			return
		}
		obj, err := eav.GetEntity(db, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
				return
			}
		}
		if obj.EntityTypeId != ett.ID {
			http.NotFound(w, r)
			ErrMsg("This object doesn't belong to this type", w)
		}
		db.Delete(obj)
	}
}

func CreateObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := eav.GetEntityTypeByName(db, entityType)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ErrMsg("can't find type", w)
				http.NotFound(w, r)
				return
			}
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			ErrMsg("can't open request body", w)
			http.NotFound(w, r)
		}
		var cr createReq
		json.Unmarshal(content, &cr)
		fmt.Println(cr)
		et, err := eav.CreateEntity(db, ett, cr.Attrs)
		if err != nil {
			http.NotFound(w, r)
			ErrMsg(err.Error(), w)
			return
		}
		w.Write(et.EncodeToJson())

	}
}

type createReq struct {
	Attrs map[string]interface{}
}
