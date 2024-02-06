package test

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/jsonutils"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jon-kamis/klogger"
)

type MockPolygonApi struct {
	Handler MockPolygonHandler
	BaseUrl string
	Port    int
}

type MockPolygonHandler struct {
	JSONUtil jsonutils.JSONUtils
}

func (m *MockPolygonApi) Routes() http.Handler {
	// Create a router r
	r := chi.NewRouter()

	r.Get(fmt.Sprintf(constants.PolygonGetPrevCloseAPI, "{ticker}"), m.Handler.MockGetStockByTicker)
	r.Get(fmt.Sprintf(constants.PolygonGetDateRangeAPI, "{ticker}", "{startDt}", "{endDt}"), m.Handler.MockGetStockByTicker)
	return r
}

func (h *MockPolygonHandler) MockGetStockByTicker(w http.ResponseWriter, r *http.Request) {
	method := "polygonexternal.MockGetStockByTicker"
	klogger.Enter(method)

	ticker := chi.URLParam(r, "ticker")

	var data []restmodels.AggResponseItem

	i := restmodels.AggResponseItem{
		Ticker:              ticker,
		Low:                 1,
		High:                1,
		Open:                1,
		Close:               1,
		UnixTime:            int(time.Now().Unix()),
		TradeVolume:         110,
		VolumeWeightedPrice: 1,
	}

	data = append(data, i)

	pc := restmodels.AggResponse{
		Adjusted:     false,
		QueryCount:   1,
		RequestId:    "1",
		Results:      data,
		ResultsCount: 1,
		Status:       fmt.Sprintf("%d", http.StatusOK),
		Ticker:       ticker,
	}

	h.JSONUtil.WriteJSON(w, http.StatusOK, pc)
	klogger.Exit(method)
}
