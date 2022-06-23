package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// this api's should be implemented
func (app *application) CreateRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var record *Record
	err := app.readJSON(w, r, &record)
	if err != nil {
		app.Logger.Printf("read record data err %v", err)
		app.badRequest(w, r, err.Error())
		return
	}

	newRec, err := app.Repository.CreateRecord(record)
	if err != nil {
		app.Logger.Printf("create record err: %v", err)
		app.serverError(w, r, nil)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"record": newRec}, nil)
	if err != nil {
		app.Logger.Printf("create record write err: %v", err)
		app.serverError(w, r, nil)
		return
	}
}
func (app *application) UpdateRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var record *Record
	err := app.readJSON(w, r, &record)

	if err != nil {
		app.Logger.Printf("read record data err %v", err)
		app.serverError(w, r, nil)
		return
	}

	updated, err := app.Repository.UpdateRecord(record)
	if err != nil {
		app.Logger.Printf("update record err: %v", err)
		if errors.Is(err, ErrorRecordNotExists) {
			app.badRequest(w, r, err.Error())
			return
		}
		app.serverError(w, r, nil)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"updated": updated}, nil)
	if err != nil {
		app.Logger.Printf("update record write err: %v", err)
		app.serverError(w, r, nil)
		return
	}
}
func (app *application) DeleteRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		app.badRequest(w, r, "invalid id")
		return
	}

	deleted, err := app.Repository.DeleteRecord(int64(id))
	if err != nil {
		app.Logger.Printf("delete record err: %v", err)
		app.serverError(w, r, nil)
		return
	}

	err = app.writeJson(w, http.StatusAccepted, envelope{"deleted_record": deleted}, nil)
	if err != nil {
		app.Logger.Printf("write delete record err: %v", err)
		app.serverError(w, r, nil)
		return
	}
}
func (app *application) GetAllRecords(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	records, err := app.Repository.GetAllRecords()

	if err != nil {
		app.Logger.Printf("get records list err: %v", err)
		app.serverError(w, r, nil)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"records": records}, nil)
	if err != nil {
		app.Logger.Printf("write records list err: %v", err)
		app.serverError(w, r, nil)
		return
	}
}

func (app *application) GetRecordById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		app.badRequest(w, r, "invalid id")
		return
	}

	record, err := app.Repository.GetRecordById(int64(id))

	if err != nil {
		app.Logger.Printf("get record by id: %v", err)
		if errors.Is(err, ErrorRecordNotExists) {
			app.badRequest(w, r, err.Error())
			return
		}
		app.serverError(w, r, nil)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"record": record}, nil)
	if err != nil {
		app.Logger.Printf("write record by id: %v", err)
		app.serverError(w, r, nil)
		return
	}
}
