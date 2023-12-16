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

func (ths handler) DeleteChannelStock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req temochi.DeleteChannelStockReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			err = fmt.Errorf("failed to read request body: %v", err)
			tenpoLog.Error(r.Context(), err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = ths.validator.StructCtx(r.Context(), req)
		if err != nil {
			err = fmt.Errorf("failed to validate request body: %v", err)
			tenpoLog.Error(r.Context(), err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimPrefix(r.URL.Path, "/stocks/") != req.WarehouseID {
			err = errors.New("warehouse_id mismatch")
			tenpoLog.Error(r.Context(), err)
			http.Error(w, "warehouse_id mismatch", http.StatusForbidden)
			return
		}

		err = ths.service.DeleteChannelStock(r.Context(), req, nil)
		if err != nil {
			err = fmt.Errorf("failed at service.DeleteChannelStock: %v", err)
			tenpoLog.Error(r.Context(), err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		tenpoLog.Success(r.Context())
	}
}
