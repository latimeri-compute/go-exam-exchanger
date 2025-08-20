package delivery

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/validator"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	Models *storages.Models
	Logger *zap.SugaredLogger
}

type loginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var receivedJson loginJSON
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.BadRequestResponse(w, utils.JSONEnveloper{"error": err}, nil)
		return
	}

	h.Logger.Debugf("Получена регистрации в системе: %v\n", receivedJson)

	v := validator.NewValidator()
	validator.ValidateUser(v, receivedJson.Email, receivedJson.Password)
	if !v.Valid() {
		utils.BadRequestResponse(w, utils.JSONEnveloper{"error": v.Errors}, nil)
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(receivedJson.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Panicf("Ошибка хеширования пароля: %v", err)
		utils.InternalErrorResponse(w, nil)
		return
	}

	newUser := &storages.User{
		Email:        receivedJson.Email,
		PasswordHash: bytes,
	}

	err = h.Models.Users.CreateUser(newUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.BadRequestResponse(w, utils.JSONEnveloper{"error": "Username or email already exists"}, nil)
			return
		}
		h.Logger.Errorf("Ошибка создания пользователя: %v", err)
		utils.InternalErrorResponse(w, nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.JSONEnveloper{"message": "User registered successfully"}, nil)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var receivedJson loginJSON
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.BadRequestResponse(w, utils.JSONEnveloper{"error": err}, nil)
		return
	}

	h.Logger.Debugf("Получена попытка входа в систему: %v\n", receivedJson)

	v := validator.NewValidator()
	validator.ValidateUser(v, receivedJson.Email, receivedJson.Password)
	if !v.Valid() {
		utils.BadRequestResponse(w, utils.JSONEnveloper{"error": v.Errors}, nil)
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

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(receivedJson.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.JSONEnveloper{"error": "Invalid username or password"}, nil)
			return
		} else {
			h.Logger.Panicf("Ошибка сравнивания хешей паролей: %v", err)
			utils.InternalErrorResponse(w, nil)
			return
		}
	}

	//TODO
}

func GetBalance(w http.ResponseWriter, r *http.Request) {

}

func TopUpBalance(w http.ResponseWriter, r *http.Request) {

}

func WithdrawFromBalance(w http.ResponseWriter, r *http.Request) {

}

func GetExchangeRates(w http.ResponseWriter, r *http.Request) {

}

func ExchangeFunds(w http.ResponseWriter, r *http.Request) {

}
