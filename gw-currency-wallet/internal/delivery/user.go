package delivery

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery/middleware"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/validator"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var receivedJson loginJSON
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

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var receivedJson loginJSON
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.JSONEnveloper{"error": "Invalid username or password"}, nil)
			return
		}
	}

	h.Logger.Debugf("Найден пользователь: %v", user)

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(receivedJson.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.JSONEnveloper{"error": "Invalid username or password"}, nil)
			return
		} else {
			h.Logger.Panicf("Ошибка сравнивания хешей паролей: %v", err)
			utils.InternalErrorResponse(w)
			return
		}
	}

	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(int64(user.ID), 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 24))
	claims.Issuer = middleware.ModuleName
	claims.Audiences = []string{middleware.ModuleName}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(h.JWTsource))
	if err != nil {
		h.Logger.Errorf("Ошибка формирования JST: %v", err)
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Выдан JST %s", jwtBytes)

	err = utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"authentication_token": string(jwtBytes)}, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}
