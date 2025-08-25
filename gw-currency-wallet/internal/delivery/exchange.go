package delivery

import (
	"context"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery/middleware"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/grpcclient"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
)

type exchangeRequest struct {
	ToCurrency   string  `json:"to_currency"`
	FromCurrency string  `json:"from_currency"`
	Amount       float64 `json:"amount"`
}

func (h *Handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	response, err := grpcclient.GetOnlyRates(h.ExchangeClient, ctx)
	if err != nil {
		h.Logger.Errorf("GetExchangeRates Ошибка получения курса валют: ", err)
		utils.InternalErrorResponse(w)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}

func (h *Handler) ExchangeFunds(w http.ResponseWriter, r *http.Request) {
	var receivedJson exchangeRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, utils.JSONEnveloper{"error": err}, nil)
		return
	}
	h.Logger.Debugf("ExchangeFunds: получен JSON: %v\n", receivedJson)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rate, err := grpcclient.GetOnlySpecificRate(h.ExchangeClient, ctx, receivedJson.FromCurrency, receivedJson.ToCurrency)
	rate = float32(math.Round(float64(rate * 100)))

	var user storages.User
	userID, ok := r.Context().Value("user").(middleware.ContextID)
	if !ok {
		h.Logger.Errorf("Ошибка получения id пользователя из контекста")
		utils.InternalErrorResponse(w)
		return
	}

	user.ID = uint(userID)
	err = h.Models.Users.FindUser(&user)
	if err != nil {
		if errors.Is(err, storages.ErrRecordNotFound) {
			h.Logger.Error("ID пользователя из контекста не соответствует ID в системе, ", err)
		} else {
			h.Logger.Error("Непредвиденная ошибка: ", err)
		}
		utils.InternalErrorResponse(w)
		return
	}

	amount := utils.Abs(int(math.Round(receivedJson.Amount * 100)))
	wallet, err := h.Models.Wallets.ExchangeBetweenCurrency(user.Wallet.ID, amount, int(rate), receivedJson.FromCurrency, receivedJson.ToCurrency)
	if err != nil {
		if errors.Is(err, storages.ErrLessThanZero) {
			h.Logger.Debug("Недостаточный баланс: ", err)
			utils.BadRequestResponse(w, "Insufficient funds or invalid currencies")
		} else {
			h.Logger.Error("Непредвиденная ошибка: ", err)
			utils.InternalErrorResponse(w)
		}
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{
		"message":          "Exchange successful",
		"exchanged_amount": receivedJson.Amount,
		"new_balance": map[string]float64{
			"RUB": float64(wallet.RubBalance),
			"USD": float64(wallet.UsdBalance),
			"EUR": float64(wallet.EurBalance),
		}}, nil)
	if err != nil {
		h.Logger.Error("Ошибка формирования json: ", err)
		utils.InternalErrorResponse(w)
	}
}
