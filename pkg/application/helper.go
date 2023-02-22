package application

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"task1/items_manager/pkg/validator"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(3, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) badRequest(w http.ResponseWriter, valid validator.Validator, code int) {
	var res UserResponse
	res.StatusCode = code
	res.Status = http.StatusText(code)
	m2 := make(map[string]interface{}, len(valid.FieldErrors))
	for k, v := range valid.FieldErrors {
		m2[k] = v
	}
	res.Data = m2
	json.NewEncoder(w).Encode(res)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func decodeJsonBody(r *http.Request, dst any) error {
	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, dst)

	if err != nil {
		return err
	}

	return nil
}

func (app *Application) successResponse(w http.ResponseWriter, data map[string]interface{}) {
	res := UserResponse{
		StatusCode: 200,
		Status:     "success",
	}
	if data != nil {
		res.Data = data
	}
	json.NewEncoder(w).Encode(res)
}
