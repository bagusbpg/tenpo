package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bagusbpg/tenpo/temochi"
)

func (ths client) DeleteStock(ctx context.Context, req temochi.DeleteStockReq, res *temochi.DeleteStockRes) error {
	u, err := url.ParseRequestURI(ths.url)
	if err != nil {
		return fmt.Errorf("failed to parse base url: %s", err.Error())
	}

	u.Path = temochi.BASE_PATH
	u.JoinPath(req.WarehouseID, req.SKU)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed creating request: %s", err.Error())
	}

	response, err := ths.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed sending request: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("got %d status code", response.StatusCode)
	}

	return nil
}
