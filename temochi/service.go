package temochi

//go: generate mockgen -source=./service.go -destination=./mock/mock.go -package=mock

import "context"

type Service interface {
	// GetStocks append the destination object res.Stocks with an array of Stock at
	// requested warehouseID and certain SKUs when exists, otherwise it appends
	// nothing.
	//
	// If SKUs are not provided, it will append nothing to res.Stocks.
	//
	// The destination object res is a pointer to non-empty struct. Thus, the caller
	// is prohibited to set res to nil
	GetStocks(ctx context.Context, req GetStocksReq, res *GetStocksRes) error

	// UpsertStocks stores new Stocks or updates them if they have already existed.
	// Valid Stocks must have following properties:
	// 1) Stock must be greater than or equal to zero
	// 2) Buffer Stock must be greater than or equal to zero
	// 3) Buffer Stock must be less than or equal to Stock
	// 4) Channel's Stock must be less than or equal to Stock minus Buffer Stock
	//
	// If a Stock has already existed with ChannelID A, B, C, then upon being
	// upserted with ChannelID A and B only, ChannelID C of this Stock will be
	// deleted, leaving channelID A and B after upsert.
	//
	// UpsertStocks is atomic in operation.
	//
	// The destination object res is a pointer to non-empty struct. Thus, the caller
	// is prohibited to set res to nil
	UpsertStocks(ctx context.Context, req UpsertStocksReq, res *UpsertStocksRes) error

	// UpdateChannelStocks updates existing Stocks. It is particularly useful in case
	// of stock alteration following a confirmed or failed order. It basically will
	// 1) Alter channel stock at a given amount of change (delta)
	// 2) Alter main stock (at inventory table) with the same amount of change
	// 3) Alter related channel stock of same warehouseID and SKUs but different
	//    gateID and channelID (This adheres to the principle of synchronized stocks)
	// 4) Alter buffer stock if necessary.
	//
	// The above operation strictly conform to the same properties of valid stocks as
	// defined in UpsertStock. Thus,
	// 1) When main stock after alteration is less than buffer stock but is greater than
	//    zero, then buffer stock will be altered with the same amount of change if possible.
	//    This is to maintain that buffer stock always less than or equal to main stock.
	// 2) When buffer stock after alteration is less than zero, then buffer stock will
	//    be set to zero. This is to maintain that buffer stock always greater than or
	//    equal to zero
	// 3) When main stock after alteration is less than zero, UpdateChannelStocs will return
	//    error.
	//
	// UpdateChannelStocks is atomic in operation.
	//
	// The destination object res is a pointer to empty struct. Thus, the caller may
	// safely set res to nil.
	UpdateChannelStocks(ctx context.Context, req UpdateChannelStocksReq, res *UpdateChannelStocksRes) error

	// DeleteChannelStock deletes all channel Stocks at certain warehouseID, gateID and
	// channelID. It is particularly useful when channelID is deleted, in which all channel
	// Stocks associated with it must be deleted.
	//
	// The destination object res is a pointer to empty struct. Thus, the caller may
	// safely set res to nil.
	DeleteChannelStock(ctx context.Context, req DeleteChannelStockReq, res *DeleteChannelStockRes) error

	// DeleteStock deletes particular Stock having certain warehouseID and SKU.
	//
	// The destination object res is a pointer to empty struct. Thus, the caller may
	// safely set res to nil.
	DeleteStock(ctx context.Context, req DeleteStockReq, res *DeleteStockRes) error
}

type GetStocksReq struct {
	WarehouseID string
	SKUs        []string
}

type GetStocksRes struct {
	Stocks []Stock `json:"stocks"`
}

type Stock struct {
	Inventory
	ChannelStocks []ChannelStock `json:"channelStocks"`
}

type UpsertStocksReq struct {
	ActorID          string            `json:"actorId"`
	ActorName        string            `json:"actorName"`
	WarehouseID      string            `json:"warehouseId" validate:"required"`
	UpsertStockSpecs []UpsertStockSpec `json:"upsertStockSpecs" validate:"required,dive"`
}

type UpsertStockSpec struct {
	SKU               string             `json:"sku" validate:"required"`
	Stock             uint32             `json:"stock" validate:"required,gte=0"`
	BufferStock       uint32             `json:"bufferStock"`
	ChannelStockSpecs []ChannelStockSpec `json:"channelStockSpecs"`
}

type ChannelStockSpec struct {
	GateID    string `json:"gateId" validate:"required"`
	ChannelID string `json:"channelId" validate:"required"`
	Stock     uint32 `json:"stock" validate:"required,gte=0"`
}

type UpsertStocksRes struct {
	FailedSpecs []FailedUpsertStockSpec `json:"failedItems"`
}

type FailedUpsertStockSpec struct {
	SKU     string `json:"sku"`
	Message string `json:"code"`
}

type UpdateChannelStocksReq struct {
	ActorID                 string                   `json:"actorId"`
	ActorName               string                   `json:"actorName"`
	WarehouseID             string                   `json:"warehouseId" validate:"required"`
	UpdateChannelStockSpecs []UpdateChannelStockSpec `json:"updateChannelStockSpecs" validate:"required,dive"`
}

type UpdateChannelStockSpec struct {
	SKU       string `json:"sku" validate:"required"`
	GateID    string `json:"gateId" validate:"required"`
	ChannelID string `json:"channelId" validate:"required"`
	Delta     int32  `json:"delta" validate:"required"`
}

type UpdateChannelStocksRes struct{}

type DeleteChannelStockReq struct {
	WarehouseID string `json:"warehouseId" validate:"required"`
	GateID      string `json:"gateId" validate:"required"`
	ChannelID   string `json:"channelId" validate:"required"`
}

type DeleteChannelStockRes struct{}

type DeleteStockReq struct {
	WarehouseID string
	SKU         string
}

type DeleteStockRes struct{}
