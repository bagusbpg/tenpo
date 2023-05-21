package temochi

const (
	BASE_PATH                  string = "/stocks"
	PATH_UPSERT_STOCKS         string = "/stocks/:warehouse_id"
	PATH_GET_STOCKS            string = "/stocks/:warehouse_id"
	PATH_UPDATE_CHANNELS_STOCK string = "/stocks/:warehouse_id"
	PATH_DELETE_CHANNEL_STOCK  string = "/stocks/:warehouse_id"
	PATH_DELETE_STOCK          string = "/stocks/:warehouse_id/:sku"

	QUERY_SKUS string = "skus"
)
