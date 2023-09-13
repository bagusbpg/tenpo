package handler

import (
	"fmt"
	"net/http"
	"strings"

	tenpoLog "github.com/bagusbpg/tenpo/kikai/log"
	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *handler) DeleteStock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := strings.Split(strings.TrimPrefix(r.URL.Path, "/stocks/"), "/")
		req := temochi.DeleteStockReq{
			WarehouseID: params[0],
			SKU:         params[1],
		}

		err := ths.service.DeleteStock(r.Context(), req, nil)
		if err != nil {
			err = fmt.Errorf("failed at service.DeleteStock: %v", err)
			tenpoLog.Error(r, "DeleteStock", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		tenpoLog.Success(r, "DeleteStock")
	}
}
