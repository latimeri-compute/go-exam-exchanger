package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/validator"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	pb "github.com/latimeri-compute/go-exam-exchanger/proto-exchange/exchange"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	Models         *storages.Models
	Logger         *zap.SugaredLogger
	ExchangeClient pb.ExchangeServiceClient
	JWT            string
}

type loginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(m *storages.Models, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		Models: m,
		Logger: logger,
	}
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
		utils.BadRequestResponse(w, utils.JSONEnveloper{"error": v.Errors}, nil)
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(receivedJson.Password), bcrypt.DefaultCost)
	h.Logger.Debug(bytes)
	if err != nil {
		utils.InternalErrorResponse(w, nil)
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
			utils.BadRequestResponse(w, utils.JSONEnveloper{"error": "Username or email already exists"}, nil)
			return
		}
		h.Logger.Errorf("Ошибка создания пользователя: %v", err)
		utils.InternalErrorResponse(w, nil)
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
	head := http.Header{}
	head.Set("Authorization", fmt.Sprintf("Bearer %s", h.JWT))

	utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"token": h.JWT}, head)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// TODO
}

type fundsRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"` // (USD, RUB, EUR)
}

func (h *Handler) TopUpBalance(w http.ResponseWriter, r *http.Request) {
	var receivedJson fundsRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}
	h.Logger.Debugf("TopUpBalance: получен JSON: %v\n", receivedJson)
	// TODO получение пользователя, пополнение валюты
}

func (h *Handler) WithdrawFromBalance(w http.ResponseWriter, r *http.Request) {

	var receivedJson fundsRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}
	h.Logger.Debugf("WithdrawFromBalance: получен JSON: %v\n", receivedJson)

	// TODO получение пользователя, снятие валюты
}

type ExchangeResponse struct {
	Currency string  `json:"currency"`
	Rate     float32 `json:"rate"`
}

func (h *Handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	var response []ExchangeResponse
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rates, err := h.ExchangeClient.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		h.Logger.Errorf("Ошибка получения курса валют: %v", err)
		utils.InternalErrorResponse(w, nil)
		return
	}
	for key, rate := range rates.Rates {
		response = append(response, ExchangeResponse{
			Currency: key,
			Rate:     rate,
		})
	}

	utils.WriteJSON(w, http.StatusOK, response, nil)
}

func (h *Handler) ExchangeFunds(w http.ResponseWriter, r *http.Request) {
	// TODO
}
