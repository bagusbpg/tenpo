package handler

import (
	"net/http"
	"strings"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
	"github.com/bagusbpg/tenpo/temochi"
	"golang.org/x/exp/slog"
)

func (ths *handler) DeleteStock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)

		params := strings.Split(strings.TrimPrefix(r.URL.Path, "/stocks/"), "/")
		req := temochi.DeleteStockReq{
			WarehouseID: params[0],
			SKU:         params[1],
		}

		err := ths.service.DeleteStock(r.Context(), req, nil)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at DeleteStock",
				slog.String("causedBy", "failed at service.DeleteStock: "+err.Error()),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "failed at service.DeleteStock: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		slog.LogAttrs(
			r.Context(),
			slog.LevelInfo, "success processing request",
			slog.String("handler", "DeleteStock"),
			slog.String("path", r.URL.Path),
			slog.String("requestID", requestID),
		)
	}
}
