package model

//////////// data model ////////////

/** (formatted) data model
{"data":
	{cacheKey:cacheValue,"timeToLive":timeToLive}
	...
}
*/
type Document map[string]interface{}

////// data transfer object (DTO) //

type Payload struct {
	CacheKey   string `json:"cacheKey"`
	CacheValue string `json:"cacheValue"`
	Ttl        int    `json:"ttl"`
}

///////////////////////////////////