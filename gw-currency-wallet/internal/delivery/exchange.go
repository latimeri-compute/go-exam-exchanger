package delivery

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/brocker"
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

// GetExchangeRates возвращает курс валют
//
//	@Summary	returns exchange rates
//	@Description
//	@Tags		exchange
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string								true	"JWT"	example("BEARER {JWT}")
//	@Success	200				{object}	string								"Returns exchange rates"
//	@Failure	401				{object}	delivery.errorUnauthorizedResponse	"Invalid credentials"
//	@Router		/exchange/rates [get]
func (h *Handler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	rates, ok := h.exchangeCache.Get("all_rates")
	h.Logger.Debugf("существование курса валют в кеше: %v, полученный курс из кеша: %v, ", ok, rates)

	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		h.Logger.Debug("запрос на удалённый сервер...")
		response, err := grpcclient.GetOnlyRates(h.ExchangeClient, ctx)
		if err != nil {
			h.Logger.Errorf("Ошибка получения курса валют: ", err)
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

type balance struct {
	RUB utils.Currency `json:"RUB"`
	USD utils.Currency `json:"USD"`
	EUR utils.Currency `json:"EUR"`
}
type exchangeResponse struct {
	Message         string         `json:"message"`
	ExchangedAmount utils.Currency `json:"exchanged_amount"`
	NewBalance      balance        `json:"new_balance"`
}

// ExchangeFunds
//
//	@Summary	exchange funds
//	@Description
//	@Tags		exchange
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string								true	"JWT"	example("BEARER {JWT}")
//	@Param		request			body		delivery.exchangeRequest			true	"Exchange funds request"
//	@Success	200				{object}	delivery.exchangeResponse			"Returns updated balance and exchanged amount"
//	@Failure	400				{object}	delivery.errorInsufficientFunds		"Insufficient funds or invalid currencies"
//	@Failure	401				{object}	delivery.errorUnauthorizedResponse	"Invalid credentials"	example(error:Unauthorized)
//	@Router		/exchange [post]
func (h *Handler) ExchangeFunds(w http.ResponseWriter, r *http.Request) {
	var receivedJson exchangeRequest
	err := utils.UnpackJSON(w, r, &receivedJson)
	if err != nil {
		h.Logger.Debugf("Ошибка распаковки json: %v\n", err)
		utils.UnprocessableEntityResponse(w, utils.JSONEnveloper{"error": err})
		return
	}
	h.Logger.Debugf("получен JSON: %v\n", receivedJson)

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

		rate = rate * 100

		h.Logger.Debug("полученный курс: ", rate)

		h.exchangeCache.Set(dir, []ExchangeCachedItem{{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         rate,
		}}, time.Minute*3)

	} else {
		rate = cache[0].Rate
		h.Logger.Debug("Курс из кеша: ", rate)
	}

	user, ok := r.Context().Value("user").(storages.User)
	if !ok {
		utils.InternalErrorResponse(w)
		return
	}

	exchangedAmount := utils.Currency(receivedJson.Amount * float64(rate/100))
	h.Logger.Debug("exchangedAmount: ", exchangedAmount)
	h.Logger.Debug("wallet Id: ", user.WalletID)
	amount := utils.Abs(int(math.Round(receivedJson.Amount * 100)))
	wallet, err := h.Models.Wallets.ExchangeBetweenCurrency(user.WalletID, amount, int(rate), fromCurrency, toCurrency)
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

	go func() {
		if exchangedAmount >= 30_000 || amount >= 30_000_00 {
			mes := brocker.TransactionMessage{
				WalletID:     wallet.ID,
				Type:         "exchange",
				FromCurrency: fromCurrency,
				ToCurrency:   toCurrency,
				AmountFrom:   int(receivedJson.Amount * 100),
				AmountTo:     int(exchangedAmount * 100),
				Timestamp:    wallet.UpdatedAt,
			}

			_, _, err := h.messenger.MessageTransaction(mes)
			if err != nil {
				h.Logger.DPanic(err)
				h.Logger.Error("Ошибка отправления сообщения: ", err)
			}
		}
	}()

	err = utils.WriteJSON(w, http.StatusOK, exchangeResponse{
		Message:         "Exchange successful",
		ExchangedAmount: exchangedAmount,
		NewBalance: balance{
			USD: utils.Currency(float64(wallet.UsdBalance) / 100),
			RUB: utils.Currency(float64(wallet.RubBalance) / 100),
			EUR: utils.Currency(float64(wallet.EurBalance) / 100),
		},
	}, nil)
	if err != nil {
		h.Logger.Error("Ошибка формирования json: ", err)
		utils.InternalErrorResponse(w)
	}
}
