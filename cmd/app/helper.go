package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type envelope map[string]interface{}

func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect  JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be null")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) IsValidIIN(iin string) bool {
	//length should be 12
	len12Match, _ := regexp.MatchString("^\\d{12}$", iin)
	if !len12Match {
		return false
	}

	iinslice := []int{}

	for i := range iin {
		iinslice = append(iinslice, int(iin[i]-'0'))
	}

	// check year
	year := iinslice[0]*10 + iinslice[1]
	if year > 99 || year < 0 {
		return false
	}

	month := iinslice[2]*10 + iinslice[3]

	if month <= 0 || month > 12 {
		return false
	}

	day := iinslice[4]*10 + iinslice[5]

	if !(day >= 1 && day <= 31) {
		return false
	}

	if !(iinslice[6] >= 1 && iinslice[6] <= 6) {
		return false
	}

	a12 := (iinslice[0]*1 + iinslice[1]*2 + iinslice[2]*3 + iinslice[3]*4 + iinslice[4]*5 + iinslice[5]*6 + iinslice[6]*7 + iinslice[7]*8 + iinslice[8]*9 + iinslice[9]*10 + iinslice[10]*11) % 11

	if a12 == 10 {
		// check with other weights
		a12 = (iinslice[0]*3 + iinslice[1]*4 + iinslice[2]*5 + iinslice[3]*6 + iinslice[4]*7 + iinslice[5]*8 + iinslice[6]*9 + iinslice[7]*10 + iinslice[8]*11 + iinslice[9]*1 + iinslice[10]*2) % 11
		if a12 == 10 {
			return false
		}
	}

	if !(a12 >= 0 && a12 <= 9) {
		return false
	}

	return true
}
