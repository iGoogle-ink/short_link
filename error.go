package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

func errRsp(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s:", e.Status(), e.Error())
		errRspWithJson(w, e.Status(), err.Error())
	default:
		log.Println("default error:", err)
		errRspWithJson(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func errRspWithJson(w http.ResponseWriter, code int, payload interface{}) {
	rsp, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(rsp)
}
