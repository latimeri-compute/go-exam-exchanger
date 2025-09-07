package delivery

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/brocker"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/validator"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
)

type fundsRequest struct {
	Amount   float64 `json:"amount" example:"234.56"`
	Currency string  `json:"currency" example:"rub"` // (USD, RUB, EUR)
}

type balanceResponse struct {
	USD utils.Currency `json:"USD" example:"120.00"`
	EUR utils.Currency `json:"EUR" example:"10.00"`
	RUB utils.Currency `json:"RUB" example:"45.50"`
}

// GetBalance returns user's balance
//
//	@Summary	returns user's balance
//	@Description
//	@Tags		balance
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string								true	"JWT"	example("BEARER {JWT}")	example("BEARER %jwt%")
//	@Success	200				{object}	delivery.balanceResponse			"Returns balance"
//	@Failure	401				{object}	delivery.errorUnauthorizedResponse	"Invalid credentials"
//	@Router		/balance [get]
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(storages.User)
	if !ok {
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Пользователь из контекста: %v", user)

	wallet, err := h.Models.Wallets.GetBalance(user.WalletID)
	if err != nil {
		h.Logger.Error("ошибка получения баланса пользователя: ", err)
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debug("Баланс пользователя: ", wallet)

	balanceResponse := balanceResponse{
		USD: utils.Currency(wallet.UsdBalance) / 100,
		EUR: utils.Currency(wallet.EurBalance) / 100,
		RUB: utils.Currency(wallet.RubBalance) / 100,
	}
	err = utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"balance": balanceResponse}, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}

// TopUpBalance пополнение баланса
//
//	@Summary	top up user's balance
//	@Description
//	@Tags		balance
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string								true	"JWT"	example("BEARER {JWT}")
//	@Param		request			body		delivery.fundsRequest				true	"top up request"
//	@Success	200				{object}	string								"returns updated balance"
//	@Failure	401				{object}	delivery.errorUnauthorizedResponse	"Invalid credentials"
//	@Router		/deposit [post]
func (h *Handler) TopUpBalance(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debug("Получен JST ", r.Header.Get("Authorization"))

	var receivedJson fundsRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debug("Ошибка распаковки json: ", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}
	h.Logger.Debug("получен JSON: ", receivedJson)

	v := validator.NewValidator()
	v.CheckBalanceChange(receivedJson.Amount, receivedJson.Currency)
	if !v.Valid() {
		h.Logger.Debug("Json не прошёл проверку валидности: ", v.Errors)
		utils.BadRequestResponse(w, v.Errors)
		return
	}

	amount := int(math.Abs(receivedJson.Amount) * 100)
	h.ChangeBalance(w, r, amount, receivedJson.Currency, "deposit")
}

// WithdrawFromBalance снятие с баланса
//
//	@Summary	withdraw from user's balance
//	@Description
//	@Tags		balance
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string								true	"JWT"	example("BEARER {JWT}")	example("BEARER %jwt%")
//	@Param		request			body		delivery.fundsRequest				true	"withdrawal request"
//	@Success	200				{object}	string								"returns updated balance"
//	@Failure	400				{object}	delivery.errorInsufficientFunds		"Insufficient funds or invalid currencies"
//	@Failure	400				{object}	delivery.errorResponse				""
//	@Failure	401				{object}	delivery.errorUnauthorizedResponse	"Invalid credentials"
//	@Router		/withdraw [post]
func (h *Handler) WithdrawFromBalance(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debug("Получен JST ", r.Header.Get("Authorization"))

	var receivedJson fundsRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debug("Ошибка распаковки json: ", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}
	h.Logger.Debugf("получен JSON: %v\n", receivedJson)

	v := validator.NewValidator()
	v.CheckBalanceChange(receivedJson.Amount, receivedJson.Currency)
	if !v.Valid() {
		h.Logger.Debug("Json не прошёл проверку валидности: ", v.Errors)
		utils.BadRequestResponse(w, v.Errors)
		return
	}

	amount := -int(math.Abs(receivedJson.Amount) * 100)
	h.ChangeBalance(w, r, amount, receivedJson.Currency, "withdrawal")
}

// смена баланса, метод просто для уменьшения тавтологии
func (h *Handler) ChangeBalance(w http.ResponseWriter, r *http.Request, amount int, currency, method string) {
	user, ok := r.Context().Value("user").(storages.User)
	if !ok {
		utils.InternalErrorResponse(w)
		return
	}
	h.Logger.Debugf("Пользователь из контекста: %v", user)

	wallet, err := h.Models.Wallets.ChangeBalance(user.WalletID, amount, currency)
	if err != nil {
		if errors.Is(err, storages.ErrLessThanZero) {
			utils.BadRequestResponse(w, "Insufficient funds or invalid amount")
		} else {
			h.Logger.Error("ошибка изменения баланса кошелька: ", err)
			utils.InternalErrorResponse(w)
		}
		return
	}
	go func() {
		var t string
		if amount >= 0 {
			t = "deposit"
		} else {
			t = "withdraw"
		}
		if utils.Abs(amount) >= 30_000_00 {
			mes := brocker.TransactionMessage{
				WalletID:     wallet.ID,
				Type:         t,
				FromCurrency: currency,
				AmountFrom:   utils.Abs(amount),
				Timestamp:    wallet.UpdatedAt,
			}

			_, _, err := h.messenger.MessageTransaction(mes)
			if err != nil {
				h.Logger.DPanic(err)
				h.Logger.Error("Ошибка отправления сообщения: ", err)
			}
		}
	}()

	balanceResponse := balanceResponse{
		USD: utils.Currency(wallet.UsdBalance) / 100,
		EUR: utils.Currency(wallet.EurBalance) / 100,
		RUB: utils.Currency(wallet.RubBalance) / 100,
	}
	err = utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"message": fmt.Sprintf("%s successful", method), "new_balance": balanceResponse}, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}
