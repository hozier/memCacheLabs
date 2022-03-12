package model

//////////// data model ////////////

/** (formatted) data model
	{ "response": 
		{ cacheKey: cacheValue(, optional) } 
	}
*/

////// data transfer object (DTO) //

type Payload struct {
	CacheKey   string `json:"cacheKey"`
	CacheValue string `json:"cacheValue"`
	Ttl        int    `json:"ttl"`
}

///////////////////////////////////