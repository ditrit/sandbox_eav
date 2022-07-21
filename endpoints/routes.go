package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ditrit/sandbox_eav/eav/models"
	"github.com/ditrit/sandbox_eav/eav/operations"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ErrMsg(msg string, w http.ResponseWriter) {
	w.Write(
		[]byte(GetErrMsg(msg)),
	)
}

// return json formated string to be consumed by frontend or client
func GetErrMsg(msg string) string {
	return fmt.Sprintf(
		`{"error": %q}`,
		msg,
	)
}

// The handler reponsible for the retreival of entities (and filter it if needed)
func GetObjects(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := operations.GetEntityTypeByName(db, entityType)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
				return
			}
		}
		queryparams := r.URL.Query()
		var qp map[string]string = make(map[string]string)
		for k, v := range queryparams {
			qp[k] = v[0]
		}
		fmt.Println(qp)
		var collection []*models.Entity = operations.GetEntitiesWithParams(db, ett, qp)

		var b strings.Builder
		b.WriteString("[")
		var pairs []string
		for _, v := range collection {
			pairs = append(pairs, string(v.EncodeToJson()))
		}
		b.WriteString(strings.Join(pairs, ","))
		b.WriteString("]")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(b.String()))
	}
}

// The handler reponsible for the retreival of une entity
func GetObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := operations.GetEntityTypeByName(db, entityType)
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

		obj, err := operations.GetEntity(db, uint(id))
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

// The handler reponsible for the deletion of entities and their associated value
func DeleteObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := operations.GetEntityTypeByName(db, entityType)
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
		obj, err := operations.GetEntity(db, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				ErrMsg(err.Error(), w)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if obj.EntityTypeId != ett.ID {
			http.NotFound(w, r)
			ErrMsg("This object doesn't belong to this type", w)
		}
		err = operations.DeleteEntity(db, obj)
		if err != nil {
			ErrMsg("This object doesn't belong to this type", w)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}

}

// The handler reponsible for the creation of entities
func CreateObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := operations.GetEntityTypeByName(db, entityType)
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
		et, err := operations.CreateEntity(db, ett, cr.Attrs)
		if err != nil {
			http.NotFound(w, r)
			ErrMsg(err.Error(), w)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Add("Location", fmt.Sprintf("/v1/objects/%d", et.ID)) // HACK: we need a more efficient way to get the URL for an entity or an entitytype
		w.WriteHeader(http.StatusCreated)
		w.Write(et.EncodeToJson())

	}
}

type createReq struct {
	Attrs map[string]interface{}
}

// The handler reponsible for the updates of entities
func ModifyObject(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entityType := vars["type"]
		ett, err := operations.GetEntityTypeByName(db, entityType)
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
		obj, err := operations.GetEntity(db, uint(id))
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
		err = json.Unmarshal(content, &mr)
		if err != nil {
			http.Error(w, GetErrMsg(err.Error()), 500)
			return
		}
		fmt.Println(mr.Attrs)

		operations.UpdateEntity(db, obj, mr.Attrs)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(obj.EncodeToJson())

	}
}

type modifyReq struct {
	Attrs map[string]interface{}
}
