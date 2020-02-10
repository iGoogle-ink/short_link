package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type App struct {
	once     sync.Once
	validate *validator.Validate
	Router   *mux.Router
}

type ShortReq struct {
	URL           string `json:"url" validate:"required"`
	ExpireMinutes int    `json:"expire_minutes" validate:"min=0"`
}

type ShortLinkRsp struct {
	ShortLink string `json:"short_link"`
}

func (a *App) Initialize() {
	// set log format
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.BindRouter()
}

func (a *App) lazyInit() {
	a.once.Do(func() {
		a.validate = validator.New()
	})
}

func (a *App) ValidateStruct(obj interface{}) error {
	a.lazyInit()
	if err := a.validate.Struct(obj); err != nil {
		return err
	}
	return nil
}

func (a *App) BindRouter() {
	a.Router.HandleFunc("/api/shorten", a.createShortLink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortLink).Methods("GET")
	a.Router.HandleFunc("/{shortLink:[a-zA-Z0-9]{1-11}}", a.redirectShortLink).Methods("GET")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	req := new(ShortReq)
	err := json.NewDecoder(r.Body).Decode(req)
	defer r.Body.Close()
	if err != nil {
		errRsp(w, &StatusError{
			Code: http.StatusBadRequest,
			Err:  err,
		})
		return
	}
	if err = a.ValidateStruct(req); err != nil {
		errRsp(w, &StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("validate parameter error: %v", *req),
		})
		return
	}
	log.Println("ShortReq:", *req)
}

func (a *App) getShortLink(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	link := query.Get("shortLink")
	log.Println("Link:", link)
}

func (a *App) redirectShortLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("vars:", vars)
	log.Println("vars.shortLink:", vars["shortLink"])
}

func (a *App) Run(addr string) {
	http.ListenAndServe(addr, a.Router)
}
