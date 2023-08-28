package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/bagusbpg/tenpo/temochi"
)

func (ths *client) GetStocks(ctx context.Context, req temochi.GetStocksReq, res *temochi.GetStocksRes) error {
	if res == nil {
		return fmt.Errorf("missing destination object")
	}

	u, err := url.ParseRequestURI(ths.url)
	if err != nil {
		return fmt.Errorf("failed to parse base url: %s", err.Error())
	}

	u.Path = temochi.BASE_PATH
	u.JoinPath(req.WarehouseID)

	params := make(url.Values)
	for _, sku := range req.SKUs {
		params.Add("skus", sku)
	}
	u.RawQuery = params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed reading response body: %s", err.Error())
	}

	err = json.Unmarshal(responseBody, res)
	if err != nil {
		return fmt.Errorf("failed parsing response body: %s", err.Error())
	}

	return nil
}
