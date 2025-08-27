package delivery

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/delivery/middleware"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/grpcclient"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/validator"
	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/pkg/utils"
)

type exchangeRequest struct {
	ToCurrency   string  `json:"to_currency"`
	FromCurrency string  `json:"from_currency"`
	Amount       float64 `json:"amount"`
}

func (h *Handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	rates, ok := h.exchangeCache.Get("all_rates")
	h.Logger.Debugf("существование курса валют в кеше: %v, полученный курс из кеша: %v, ", ok, rates)

	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		h.Logger.Debug("запрос на удалённый сервер...")
		response, err := grpcclient.GetOnlyRates(h.ExchangeClient, ctx)
		if err != nil {
			h.Logger.Errorf("GetExchangeRates Ошибка получения курса валют: ", err)
			utils.InternalErrorResponse(w)
			return
		}
		h.Logger.Debug("полученный курс валют: ", response)

		// TODO поменяй там протобаф уже......
		for key, rate := range response {
			rates = append(rates, ExchangeCachedItem{
				FromCurrency: key[0:3],
				ToCurrency:   key[5:],
				Rate:         rate,
			})
		}

		h.exchangeCache.Set("all_rates", rates, time.Duration(time.Minute*3))
	}

	err := utils.WriteJSON(w, http.StatusOK, utils.JSONEnveloper{"rates": rates}, nil)
	if err != nil {
		h.Logger.Errorf("Ошибка формирования json: %v", err)
		utils.InternalErrorResponse(w)
	}
}

func (h *Handler) ExchangeFunds(w http.ResponseWriter, r *http.Request) {
	// ew what an ugly bastard
	// TODO вынести структуры в более подходящее место
	// TODO в принципе облагородить метод
	var receivedJson exchangeRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.UnprocessableEntityResponse(w, utils.JSONEnveloper{"error": err})
		return
	}
	h.Logger.Debugf("ExchangeFunds: получен JSON: %v\n", receivedJson)

	fromCurrency := strings.ToLower(receivedJson.FromCurrency)
	toCurrency := strings.ToLower(receivedJson.ToCurrency)

	v := validator.NewValidator()
	validator.ValidateExchangeRequest(v, fromCurrency, toCurrency, receivedJson.Amount)
	if !v.Valid() {
		utils.BadRequestResponse(w, v.Errors)
		return
	}

	var rate float32
	dir := fmt.Sprintf("%s->%s", fromCurrency, toCurrency)
	cache, ok := h.exchangeCache.Get(dir)
	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		rate, err := grpcclient.GetOnlySpecificRate(h.ExchangeClient, ctx, fromCurrency, toCurrency)
		if err != nil {
			h.Logger.Error("Ошибка получения данных от gRPC сервера: ", err)
			utils.InternalErrorResponse(w)
			return
		}

		rate = float32(math.Round(float64(rate * 100)))

		h.exchangeCache.Set(dir, []ExchangeCachedItem{{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         rate,
		}}, time.Minute*3)

	} else {
		rate = cache[0].Rate
	}

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
	h.Logger.Debug("wallet Id: ", user.Wallet.ID)
	amount := utils.Abs(int(math.Round(receivedJson.Amount * 100)))
	wallet, err := h.Models.Wallets.ExchangeBetweenCurrency(user.Wallet.ID, amount, int(rate), fromCurrency, toCurrency)
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

	type balance struct {
		RUB float64 `json:"RUB"`
		USD float64 `json:"USD"`
		EUR float64 `json:"EUR"`
	}
	type exchangeResponse struct {
		Message         string  `json:"message"`
		ExchangedAmount float64 `json:"exchanged_amount"`
		NewBalance      balance `json:"new_balance"`
	}
	err = utils.WriteJSON(w, http.StatusOK, exchangeResponse{
		Message:         "Exchange successful",
		ExchangedAmount: receivedJson.Amount * float64(rate/100),
		NewBalance: balance{
			USD: float64(wallet.UsdBalance),
			RUB: float64(wallet.RubBalance),
			EUR: float64(wallet.EurBalance),
		},
	}, nil)
	if err != nil {
		h.Logger.Error("Ошибка формирования json: ", err)
		utils.InternalErrorResponse(w)
	}
}
