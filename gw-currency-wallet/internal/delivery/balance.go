package delivery

import (
	"errors"
	"math"
	"net/http"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery/middleware"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/validator"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
	"gorm.io/gorm"
)

type fundsRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"` // (USD, RUB, EUR)
}

type balanceResponse struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	RUB float64 `json:"RUB"`
}

// получение баланса
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debugf("Получен JST %v", r.Header.Get("Authentication"))
	userId, ok := r.Context().Value("user").(middleware.ContextID)
	if !ok {
		h.Logger.Errorf("Ошибка получения id пользователя из контекста")
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Пользователь из контекста: %v", userId)

	user := &storages.User{
		Model: gorm.Model{
			ID: uint(userId),
		},
	}
	err := h.Models.Users.FindUser(user)
	if err != nil {
		h.Logger.Errorf("GetBalance: oшибка получения пользователя: %v", err)
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Найден пользователь в системе: %v", user)

	wallet, err := h.Models.Wallets.GetBalance(user.ID)
	if err != nil {
		h.Logger.Errorf("GetBalance: ошибка получения баланса пользователя: %v", err)
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Баланс пользователя: %v", wallet)

	balanceResponse := balanceResponse{
		USD: float64(wallet.UsdBalance) / 100,
		EUR: float64(wallet.EurBalance) / 100,
		RUB: float64(wallet.RubBalance) / 100,
	}
	err = utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"balance": balanceResponse}, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}

// пополнение баланса
func (h *Handler) TopUpBalance(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debugf("Получен JST %v", r.Header.Get("Authentication"))

	var receivedJson fundsRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: ", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}
	h.Logger.Debugf("TopUpBalance: получен JSON: ", receivedJson)

	v := validator.NewValidator()
	v.CheckBalanceChange(receivedJson.Amount, receivedJson.Currency)
	if !v.Valid() {
		h.Logger.Debug("Json не прошёл проверку валидности: ", v.Errors)
		utils.BadRequestResponse(w, v.Errors)
		return
	}

	amount := int(math.Abs(receivedJson.Amount) * 100)
	h.ChangeBalance(w, r, amount, receivedJson.Currency)
}

// снятие с баланса
func (h *Handler) WithdrawFromBalance(w http.ResponseWriter, r *http.Request) {
	var receivedJson fundsRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}
	h.Logger.Debugf("WithdrawFromBalance: получен JSON: %v\n", receivedJson)

	v := validator.NewValidator()
	v.CheckBalanceChange(receivedJson.Amount, receivedJson.Currency)
	if !v.Valid() {
		h.Logger.Debug("Json не прошёл проверку валидности: ", v.Errors)
		utils.BadRequestResponse(w, v.Errors)
		return
	}

	amount := -int(math.Abs(receivedJson.Amount) * 100)
	h.ChangeBalance(w, r, amount, receivedJson.Currency)
}

// смена баланса, метод просто для уменьшения тавтологии
func (h *Handler) ChangeBalance(w http.ResponseWriter, r *http.Request, amount int, currency string) {
	userId, ok := r.Context().Value("user").(middleware.ContextID)
	if !ok {
		h.Logger.Errorf("Ошибка получения id пользователя из контекста")
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Пользователь из контекста: %v", userId)

	user := &storages.User{
		Model: gorm.Model{
			ID: uint(userId),
		},
	}
	err := h.Models.Users.FindUser(user)
	if err != nil {
		h.Logger.Errorf("ошибка получения пользователя: %v", err)
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Найден пользователь в системе: %v", user)

	wallet, err := h.Models.Wallets.ChangeBalance(user.WalletID, amount, currency)
	if err != nil {
		if errors.Is(err, storages.ErrLessThanZero) {
			utils.BadRequestResponse(w, "Insufficient funds or invalid amount")
			return
		}
		h.Logger.Error("WithdrawFromBalance: ошибка снятия с кошелька: ", err)
		utils.InternalErrorResponse(w)
		return
	}
	balanceResponse := balanceResponse{
		USD: float64(wallet.UsdBalance) / 10000,
		EUR: float64(wallet.EurBalance) / 10000,
		RUB: float64(wallet.RubBalance) / 10000,
	}
	err = utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"message": "Withdrawal successful", "new_balance": balanceResponse}, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}
