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
				http.Error(w, GetErrMsg("Record not found, please use a type which is in the database"), http.StatusNotFound)
				return
			}
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
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
		w.WriteHeader(http.StatusOK)
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
				http.Error(w, GetErrMsg("Record not found, please use a type which is in the database"), http.StatusNotFound)
				return
			}
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			http.Error(w, GetErrMsg("The id you provided is not an int"), http.StatusBadRequest)
			return
		}

		obj, err := operations.GetEntity(db, ett, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, GetErrMsg("Record not found, please use an id which is in the database"), http.StatusNotFound)
				return
			} else if errors.Is(err, operations.ErrIdDontMatchEntityType) {
				http.Error(w, GetErrMsg(err.Error()), http.StatusBadRequest)
			} else {
				http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
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
				http.Error(w, GetErrMsg("Record not found, please use a type which is in the database"), http.StatusNotFound)
				return
			}
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, GetErrMsg("The id you provided is not an int"), http.StatusBadRequest)
			return
		}
		obj, err := operations.GetEntity(db, ett, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, GetErrMsg("Record not found, please use an id which is in the database"), http.StatusNotFound)
				return
			} else if errors.Is(err, operations.ErrIdDontMatchEntityType) {
				http.Error(w, GetErrMsg(err.Error()), http.StatusBadRequest)
			} else {
				http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			}
			return
		}
		err = operations.DeleteEntity(db, obj)
		if err != nil {
			http.Error(w, GetErrMsg("Deletion failed"), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
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
				http.Error(w, GetErrMsg("Record not found, please use a type which is in the database"), http.StatusNotFound)
				return
			}
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, GetErrMsg("Can't open request body"), http.StatusBadRequest)
			return
		}
		var cr createReq
		json.Unmarshal(content, &cr)
		fmt.Println(cr)
		et, err := operations.CreateEntity(db, ett, cr.Attrs)
		if err != nil {
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Add("Location", fmt.Sprintf("/v1/objects/%s/%d", ett.Name, et.ID)) // HACK: we need a more efficient way to get the URL for an entity or an entitytype
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
				http.Error(w, GetErrMsg("Record not found, please use a type which is in the database"), http.StatusNotFound)
				return
			}
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}
		content, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, GetErrMsg("Can't open request body"), http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, GetErrMsg("The id you provided is not an int"), http.StatusBadRequest)
			return
		}
		obj, err := operations.GetEntity(db, ett, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, GetErrMsg("Record not found, please use an id which is in the database"), http.StatusNotFound)
				return
			} else if errors.Is(err, operations.ErrIdDontMatchEntityType) {
				http.Error(w, GetErrMsg(err.Error()), http.StatusBadRequest)
			} else {
				http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			}
			return
		}

		var mr modifyReq
		err = json.Unmarshal(content, &mr)
		if err != nil {
			http.Error(w, GetErrMsg(err.Error()), http.StatusInternalServerError)
			return
		}
		fmt.Println(mr.Attrs)

		operations.UpdateEntity(db, obj, mr.Attrs)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(obj.EncodeToJson())

	}
}

type modifyReq struct {
	Attrs map[string]interface{}
}
