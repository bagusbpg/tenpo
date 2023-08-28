package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	tenpoHttp "github.com/bagusbpg/tenpo/kikai/http"
	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *handler) UpsertStocks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(tenpoHttp.REQUEST_ID_CONTEXT_KEY).(string)
		if !ok {
			slog.LogAttrs(
				r.Context(),
				slog.LevelWarn, "request has no requestID",
			)
		}

		var req temochi.UpsertStocksReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at UpsertStocks",
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
				slog.LevelError, "failed at UpsertStocks",
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
				slog.LevelError, "failed at UpsertStocks",
				slog.String("causedBy", "warehouse_id mismatch"),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)
			http.Error(w, "warehouse_id mismatch", http.StatusForbidden)
			return
		}

		var res temochi.UpsertStocksRes
		err = ths.service.UpsertStocks(r.Context(), req, &res)
		if err != nil {
			slog.LogAttrs(
				r.Context(),
				slog.LevelError, "failed at UpsertStocks",
				slog.String("causedBy", "failed at service.UpsertStocks: "+err.Error()),
				slog.String("path", r.URL.Path),
				slog.String("requestID", requestID),
			)

			if strings.Contains(err.Error(), "failed validating channel stock") {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": res,
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": res,
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": res,
		})
		slog.LogAttrs(
			r.Context(),
			slog.LevelInfo, "success processing request",
			slog.String("handler", "UpsertStocks"),
			slog.String("path", r.URL.Path),
			slog.String("requestID", requestID),
		)
	}
}
