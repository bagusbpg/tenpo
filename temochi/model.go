package temochi

type Inventory struct {
	WarehouseID string `json:"warehouseId"`
	SKU         string `json:"sku"`
	Stock       uint32 `json:"stock"`
	BufferStock uint32 `json:"bufferStock"`
	Version     uint64 `json:"version"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type ChannelStock struct {
	WarehouseID string `json:"warehouseId"`
	SKU         string `json:"sku"`
	GateID      string `json:"gateId"`
	ChannelID   string `json:"channelId"`
	Stock       uint32 `json:"stock"`
	Version     uint64 `json:"version"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}
