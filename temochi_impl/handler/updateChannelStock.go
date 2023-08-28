package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *handler) UpdateChannelStocks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)
		if !ok {
			slog.LogAttrs(
				r.Context(),
				slog.LevelWarn, "request has no requestID",
			)
		}

		var req temochi.UpdateChannelStocksReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at UpdateChannelStocks",
				slog.String("causedBy", "failed reading request body: "+err.Error()),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "failed reading request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = ths.validator.StructCtx(r.Context(), req)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at UpdateChannelStocks",
				slog.String("causedBy", "failed validating request body: "+err.Error()),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "failed validating request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimPrefix(r.URL.Path, "/stocks/") != req.WarehouseID {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at UpdateChannelStocks",
				slog.String("causedBy", "warehouse_id mismatch"),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "warehouse_id mismatch", http.StatusForbidden)
			return
		}

		err = ths.service.UpdateChannelStocks(r.Context(), req, nil)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at UpdateChannelStocks",
				slog.String("causedBy", "failed at service.UpdateChannelStocks: "+err.Error()),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "failed at service.UpdateChannelStocks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		slog.LogAttrs(
			r.Context(),
			slog.LevelInfo, "success processing request",
			slog.String("handler", "UpdateChannelStocks"),
			slog.String("path", r.URL.Path),
			slog.String("requestID", requestID),
		)
	}
}
