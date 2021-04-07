package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func BuildErrorResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	var m string
	switch t := message.(type) {
	case string:
		m = t
	case error:
		m = t.Error()
	case fmt.Stringer:
		m = t.String()
	default:
		m = "Unknown Error"
	}

	fmt.Println(m)

	response := ErrorResponse{m}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func Build(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 65536))
	if err != nil {
		BuildErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if err = r.Body.Close(); err != nil {
		BuildErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	source := string(body)
	crc := GenerateCRC(source)

	err = nil
	var output []byte

	l := r.Header.Get("X-Language")
	if l == "" {
		l = "c"
	}

	v := r.Header.Get("X-Version")
	if v == "" {
		v = DefaultVersion()
	}

	b, _ := GetBinary(crc, v)
	if b != nil && b.Language == l && b.Results == "Success" {
		output = b.Fram
		err = nil
	} else {
		output, err = Compile(body, l, v)

		results := ""
		if err != nil {
			results = err.Error()
		} else {
			results = "Success"
		}

		b = &Binary{}
		b.New(crc, source, l, output, results, v)

		errStore := StoreBinary(b)

		if errStore != nil {
			BuildErrorResponse(w, http.StatusInternalServerError, errStore)
			return
		}
	}

	if err != nil {
		BuildErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	res := struct {
		Binary []byte `json:"binary"`
		Id     string `json:"id"`
	}{
		output,
		crc,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func Code(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	crc := vars["id"]

	b, err := GetBinary(crc, DefaultVersion())
	if err != nil {
		BuildErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	code := ""
	if b != nil {
		code = b.Source
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Language", b.Language)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, code)
}
