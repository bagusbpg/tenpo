package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	tenpoLog "github.com/bagusbpg/tenpo/kikai/log"
	"github.com/bagusbpg/tenpo/temochi"
)

func (ths handler) UpsertStocks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req temochi.UpsertStocksReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			err = fmt.Errorf("failed reading request body: %v", err)
			tenpoLog.Error(r.Context(), r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = ths.validator.StructCtx(r.Context(), req)
		if err != nil {
			err = fmt.Errorf("failed validating request body: %v", err)
			tenpoLog.Error(r.Context(), r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimPrefix(r.URL.Path, "/stocks/") != req.WarehouseID {
			err = errors.New("warehouse_id mismatch")
			tenpoLog.Error(r.Context(), r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		var res temochi.UpsertStocksRes
		err = ths.service.UpsertStocks(r.Context(), req, &res)
		if err != nil {
			err = fmt.Errorf("failed at service.UpsertStocks: %v", err)
			tenpoLog.Error(r.Context(), r.URL.Path, err)

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
		tenpoLog.Success(r.Context(), r.URL.Path)
	}
}
