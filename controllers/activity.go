package controllers

import (
	"casaattivita/lang"
	"casaattivita/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// BaseHandler will hold everything that controller needs
type BaseHandlerSqlx struct {
	db *sqlx.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandlerSqlx(db *sqlx.DB) *BaseHandlerSqlx {
	return &BaseHandlerSqlx{
		db: db,
	}
}

// swagger:model CommonError
type CommonError struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model CommonSuccess
type CommonSuccess struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model GetActivities
type GetActivities struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string             `json:"message"`
	Data    *models.Activities `json:"data"`
}

// swagger:model GetActivity
type GetActivity struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string `json:"message"`
	// Umidity value
	Data *models.Activity `json:"data"`
}

// ErrHandler returns error message response
func ErrHandler(errmessage string) *CommonError {
	errresponse := CommonError{}
	errresponse.Status = 0
	errresponse.Message = errmessage
	return &errresponse
}

// swagger:route GET /activities listActivity
// Get Activity list
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetActivities
func (h *BaseHandlerSqlx) GetActivitiesSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetActivities{}

	activities := models.GetActivitiesSqlx(h.db.DB)

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = activities

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route GET /lasthour lastHour
// Get list of last hour of activity values .... or the last value inserted
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetActivities
func (h *BaseHandlerSqlx) GetLastHourSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetActivities{}

	activities := models.GetLastHourSqlx(h.db.DB)

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = activities

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route POST /activity addActivities
// Create a new activity value
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetActivity
func (h *BaseHandlerSqlx) PostActivitySqlx(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := GetActivity{}

	decoder := json.NewDecoder(r.Body)
	var reqActivity *models.ReqAddActivity
	err := decoder.Decode(&reqActivity)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_request")))
		return
	}

	rain, errmessage := models.PostActivitySqlx(h.db.DB, reqActivity)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("insert_success")
	response.Data = rain
	json.NewEncoder(w).Encode(response)
}
