package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
	"github.com/bagusbpg/tenpo/temochi"
	"golang.org/x/exp/slog"
)

func (ths *handler) GetStocks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)

		// error from QueryUnescape is ignored since it will
		// cause ParseQuery returns error anyway
		unescapedQuery, _ := url.QueryUnescape(r.URL.RawQuery)

		query, err := url.ParseQuery(unescapedQuery)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at GetStocks",
				slog.String("causedBy", "failed parsing query: "+err.Error()),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "failed parsing query: "+err.Error(), http.StatusBadRequest)
			return
		}

		req := temochi.GetStocksReq{
			WarehouseID: strings.Split(strings.TrimPrefix(r.URL.Path, "/stocks/"), "?")[0],
			SKUs:        query["skus"],
		}

		var res temochi.GetStocksRes
		err = ths.service.GetStocks(r.Context(), req, &res)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at GetStocks",
				slog.String("causedBy", "failed at service.GetStocks: "+err.Error()),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "failed at service.GetStocks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": res,
		})
		slog.LogAttrs(
			r.Context(),
			slog.LevelInfo, "success processing request",
			slog.String("handler", "GetStocks"),
			slog.String("path", r.URL.Path),
			slog.String("requestID", requestID),
		)
	}
}
