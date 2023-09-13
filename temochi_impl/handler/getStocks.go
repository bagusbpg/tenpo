package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	tenpoLog "github.com/bagusbpg/tenpo/kikai/log"
	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *handler) GetStocks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// error from QueryUnescape is ignored since it will
		// cause ParseQuery returns error anyway
		unescapedQuery, _ := url.QueryUnescape(r.URL.RawQuery)

		query, err := url.ParseQuery(unescapedQuery)
		if err != nil {
			err = fmt.Errorf("failed parsing query: %v", err)
			tenpoLog.Error(r, "GetStocks", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := temochi.GetStocksReq{
			WarehouseID: strings.Split(strings.TrimPrefix(r.URL.Path, "/stocks/"), "?")[0],
			SKUs:        query["skus"],
		}

		var res temochi.GetStocksRes
		err = ths.service.GetStocks(r.Context(), req, &res)
		if err != nil {
			err = fmt.Errorf("failed at service.GetStocks: %v", err)
			tenpoLog.Error(r, "GetStocks", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": res,
		})
		tenpoLog.Success(r, "GetStocks")
	}
}
