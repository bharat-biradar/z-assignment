package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"task1/items_manager/pkg/models"
	"task1/items_manager/pkg/validator"
	"time"

	"github.com/julienschmidt/httprouter"
)

type userInfoForm struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	validator.Validator
}

type userLoginForm struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type itemData map[string]interface{}

type UserResponse struct {
	StatusCode int                    `json:"statusCode"`
	Status     string                 `json:"status"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

type UserSession struct {
	StatusCode   int    `json:"statusCode"`
	Status       string `json:"status"`
	SessionToken string `json:"token"`
}

func (app *Application) createUser(w http.ResponseWriter, r *http.Request) {
	var form userInfoForm
	defer r.Body.Close()
	err := decodeJsonBody(r, &form)

	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	app.validateSignupForm(&form)

	if !form.Valid() {
		app.badRequest(w, form.Validator, http.StatusBadRequest)
		return
	}

	user := models.User{
		Username:  form.UserName,
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Password:  form.Password,
	}

	err = app.DbClient.CreateUser(user)
	if err != nil {
		panic(err)
	}
	app.successResponse(w, nil)
}

func (app *Application) loginUser(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	defer r.Body.Close()
	err := decodeJsonBody(r, &form)

	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	user := models.User{
		Username: form.UserName,
		Password: form.Password,
	}

	err = app.DbClient.ValiatePassword(user)

	if errors.Is(err, models.ErrUserDoesNotExist) {
		var val validator.Validator
		val.FieldErrors = make(map[string]string)
		val.CheckField(false, "username", "invalid username")
		app.badRequest(w, val, http.StatusBadRequest)
		return
	}

	if errors.Is(err, models.ErrInvalidPassword) {
		var val validator.Validator
		val.FieldErrors = make(map[string]string)
		val.CheckField(false, "password", "wrong password")
		app.badRequest(w, val, http.StatusUnauthorized)
		return
	}

	data := make(map[string]interface{})
	if app.sessionManager.currentActiveSession(form.UserName) {
		data["session"] = "logged out of previous session"
		app.sessionManager.removeSession(app.sessionManager.userToKey[form.UserName])
	}

	sessionToken := app.getSessionToken(form.UserName)
	data["token"] = sessionToken
	app.successResponse(w, data)

}

func (app *Application) getItem(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	res, err := app.DbClient.GetItem(id)

	if err != nil {
		if errors.Is(err, models.ErrItemDoesNotExist) {
			app.notFound(w)
			return
		}

		app.serverError(w, err)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (app *Application) getAllItems(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value(contextUserKey)

	items, err := app.DbClient.GetAllUserItems(username.(string))

	if err != nil {
		app.serverError(w, err)
		return
	}

	var res UserResponse
	res.StatusCode = 200
	res.Status = http.StatusText(200)
	res.Data = make(map[string]interface{})
	res.Data["count"] = strconv.Itoa(len(items))
	res.Data["items"] = items
	json.NewEncoder(w).Encode(res)
}

func (app *Application) deleteItem(w http.ResponseWriter, r *http.Request) {
	// username := r.Context().Value(contextUserKey)
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	err := app.DbClient.DeleteItem(models.Item{Id: id})

	if err != nil {
		app.serverError(w, err)
		return
	}
	app.successResponse(w, nil)
}

func (app *Application) addItem(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(contextUserKey)
	var data itemData
	err := decodeJsonBody(r, &data)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	fmt.Println(username.(string))
	item := models.Item{
		Owner: username.(string),
		Data:  data,
	}
	res, err := app.DbClient.InsertItem(item)

	if err != nil {
		app.serverError(w, err)
		return
	}
	outData := make(map[string]interface{})
	outData["insertedId"] = res
	app.successResponse(w, outData)
}

func (app *Application) deletUser(w http.ResponseWriter, r *http.Request) {
	var form userInfoForm
	defer r.Body.Close()
	err := decodeJsonBody(r, &form)

	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	user := models.User{
		Username: form.UserName,
		Password: form.Password,
	}
	err = app.DbClient.DeleteUser(user)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.successResponse(w, nil)
}

func (app *Application) updateItem(w http.ResponseWriter, r *http.Request) {
	// username := r.Context().Value(contextUserKey)
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	var data itemData
	err := decodeJsonBody(r, &data)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	item := models.UpdateItem{
		Id:       id,
		Data:     data,
		Modified: time.Now(),
	}
	err = app.DbClient.UpdateItem(&item)

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.successResponse(w, nil)
}

func (app *Application) validateSignupForm(form *userInfoForm) {
	form.CheckField(validator.NotBlank(form.UserName), "username", "Username cannot be empty")
	form.CheckField(validator.NotBlank(form.FirstName), "firstName", "firstName cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	form.CheckField(validator.MaxChars(form.FirstName, 255), "firstName", "Name cannot exceed 255 characters")
	ok, err := app.DbClient.CheckUserNameExists(form.UserName)
	if errors.Is(err, models.ErrDatabaseOperation) {
		panic(err)
	}

	form.CheckField(!ok, "username", "username already taken")
}
