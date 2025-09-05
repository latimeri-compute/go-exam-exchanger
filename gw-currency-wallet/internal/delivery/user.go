package delivery

import (
	"errors"
	"net/http"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery/middleware"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/validator"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type loginJSON struct {
	Email    string `json:"email" example:"test@test.com"`
	Password string `json:"password" example:"pa$$word"`
}

type registerJSON struct {
	Username string `json:"username" example:"admin"`
	Email    string `json:"email" example:"test@test.com"`
	Password string `json:"password" example:"pa$$word"`
}

// RegisterUser registering new users
//
//	@Summary	create new user
//	@Description
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		credentials	body		delivery.registerJSON		true	"Credentials"
//	@Success	201			{object}	delivery.messageResponse	"User created"
//	@Failure	400			{object}	delivery.errorResponse		"Username or email already exists"
//	@Failure	400			{object}	delivery.errorResponse		"JSON fields didn't pass validation"
//	@Failure	422			{object}	delivery.errorResponse		"Malformed json or invalid fields"
//	@Router		/register [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var receivedJson registerJSON
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}

	h.Logger.Debugf("Получена попытка регистрации в системе: %v", receivedJson)

	v := validator.NewValidator()
	validator.ValidateUser(v, receivedJson.Email, receivedJson.Password)
	if !v.Valid() {
		utils.BadRequestResponse(w, v.Errors)
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(receivedJson.Password), bcrypt.DefaultCost)
	h.Logger.Debug(bytes)
	if err != nil {
		utils.InternalErrorResponse(w)
		h.Logger.Panicf("Ошибка хеширования пароля: %v", err)
		return
	}

	newUser := &storages.User{
		Username:     receivedJson.Username,
		Email:        receivedJson.Email,
		PasswordHash: bytes,
	}

	h.Logger.Info("Создание пользователя в базе данных...")
	err = h.Models.Users.CreateUser(newUser)
	h.Logger.Info("База данных вернула ответ...")
	if err != nil {
		if errors.Is(err, storages.ErrRecordExists) {
			utils.BadRequestResponse(w, "Username or email already exists")
			return
		}
		h.Logger.Errorf("Ошибка создания пользователя: %v", err)
		utils.InternalErrorResponse(w)
		return
	}

	h.Logger.Debugf("Пользователь %s успешно зарегистрирован в системе", newUser.Email)
	utils.WriteJSON(w, http.StatusCreated, utils.JSONEnveloper{"message": "User registered successfully"}, nil)
}

// LoginUser log into system
//
//	@Summary	produces JWT token on successfull login
//	@Description
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Param		credentials	body		delivery.loginJSON			true	"Credentials"
//	@Success	200			{object}	delivery.messageResponse	"Successfully logged in"
//	@Failure	400			{object}	delivery.errorResponse		"JSON fields didn't pass validation"
//	@Failure	401			{object}	delivery.errorResponse		"Invalid credentials"
//	@Failure	422			{object}	delivery.errorResponse		"Malformed json or invalid fields"
//	@Router		/login [post]
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var receivedJson loginJSON
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err)
		return
	}

	h.Logger.Debugf("Получена попытка входа в систему: %v\n", receivedJson)

	v := validator.NewValidator()
	validator.ValidateUser(v, receivedJson.Email, receivedJson.Password)
	if !v.Valid() {
		utils.BadRequestResponse(w, v.Errors)
		return
	}

	user := &storages.User{
		Email: receivedJson.Email,
	}
	err = h.Models.Users.FindUser(user)
	if err != nil {
		if errors.Is(err, storages.ErrRecordNotFound) {
			utils.UnauthorizedResponse(w, "Invalid username or password")
		} else {
			h.Logger.Error("Ошибка получения пользователя из базы данных: ", err)
			utils.InternalErrorResponse(w)
		}
		return
	}

	h.Logger.Debug("Найден пользователь: ", user)
	h.Logger.Debug(user.PasswordHash)

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(receivedJson.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			utils.UnauthorizedResponse(w, "Invalid username or password")
		} else {
			h.Logger.Error("Ошибка сравнивания хешей паролей: ", err)
			utils.InternalErrorResponse(w)
		}
		return
	}

	jwtBytes, err := middleware.IssueNewJWT([]byte(h.JWTsource), int(user.ID))
	if err != nil {
		h.Logger.Error("Ошибка формирования JST: ", err)
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Выдан JST %s", jwtBytes)

	err = utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"authorization_token": string(jwtBytes)}, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}
