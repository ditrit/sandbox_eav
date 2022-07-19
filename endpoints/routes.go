package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ditrit/sandbox_eav/eav"
	"github.com/ditrit/sandbox_eav/eav/models"
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

func GetObjects(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
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
		var collection []*models.Entity
		db.Where("entity_type_id = ? ", ett.ID).Preload("Fields").Preload("Fields.Attribut").Preload("EntityType.Attributs").Preload("EntityType").Find(&collection)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var b strings.Builder
		b.WriteString("[")
		var pairs []string
		for _, v := range collection {
			pairs = append(pairs, string(v.EncodeToJson()))
		}
		b.WriteString(strings.Join(pairs, ","))
		b.WriteString("]")
		w.Write([]byte(b.String()))
	}
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
				ErrMsg(err.Error(), w)

				return
			}
		}
		if obj.EntityTypeId != ett.ID {
			http.NotFound(w, r)
			ErrMsg("This object doesn't belong to this type", w)
		}
		for _, v := range obj.Fields {
			db.Delete(v)
		}
		db.Delete(obj)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(et.EncodeToJson())

	}
}

type createReq struct {
	Attrs map[string]interface{}
}

func ModifyObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
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

		var mr modifyReq
		json.Unmarshal(content, &mr)
		fmt.Println(mr)
		et, err := eav.CreateEntity(db, ett, mr.Attrs)
		if err != nil {
			http.NotFound(w, r)
			ErrMsg(err.Error(), w)
			return
		}
		et.ID = obj.ID
		for _, f := range obj.Fields {
			db.Delete(f)
		}
		db.Delete(obj)

		db.Save(et)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(et.EncodeToJson())

	}
}

type modifyReq struct {
	id    int
	Attrs map[string]interface{}
}
