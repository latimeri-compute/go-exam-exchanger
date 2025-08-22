package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/grpcclient"
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
		return
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

}
