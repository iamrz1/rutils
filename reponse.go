package rutils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	ISOLayout = "2006-01-02T15:04:05.000Z"
)

type Meta struct {
	Count    int64 `json:"count"`
	PageSize int64 `json:"page_size"`
}

// Response serializer util
type Response struct {
	Status    string      `json:"status,omitempty"`
	Message   string      `json:"message,omitempty"`
	Success   bool        `json:"success"`
	Meta      interface{} `json:"meta,omitempty"`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

func ServeJSONObject(w http.ResponseWriter, code int, message string, data interface{}, meta interface{}, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var resp interface{}
	type EmptyObject struct{}
	if data == nil {
		data = EmptyObject{}
	}

	resp = &Response{
		Status:    http.StatusText(code),
		Message:   message,
		Success:   success,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().Format(ISOLayout),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func ServeJSONList(w http.ResponseWriter, code int, message string, list interface{}, meta interface{}, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var resp interface{}
	type EmptyObject struct{}
	if list == nil {
		list = []EmptyObject{}
	}

	resp = &Response{
		Status:    http.StatusText(code),
		Message:   message,
		Success:   success,
		Data:      list,
		Meta:      meta,
		Timestamp: time.Now().Format(ISOLayout),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func HandleObjectError(w http.ResponseWriter, err error) {
	errMeta := map[string]string{}
	if os.Getenv("ENV") != "prod" {
		errMeta["error"] = err.Error()
	}

	switch v := err.(type) {
	case ValidationError:
		ServeJSONObject(w, http.StatusBadRequest, v.ErrorMessage(), nil, nil, false)
		return
	case GenericHttpError:
		ServeJSONObject(w, v.Code(), v.Error(), nil, nil, false)
		return
	default:
		ServeJSONObject(w, http.StatusInternalServerError, "Something went wrong", nil, errMeta, false)
		return
	}
}

func HandleListError(w http.ResponseWriter, err error) {
	log.Println("service error: ", err)

	errMeta := map[string]string{}
	if os.Getenv("ENV") != "prod" {
		errMeta["error"] = err.Error()
	}

	switch v := err.(type) {
	case ValidationError:
		ServeJSONList(w, http.StatusBadRequest, v.ErrorMessage(), nil, nil, false)
		return
	case GenericHttpError:
		ServeJSONList(w, v.Code(), v.Error(), nil, nil, false)
		return
	default:
		ServeJSONList(w, http.StatusInternalServerError, "Something went wrong", nil, errMeta, false)
		return
	}
}
