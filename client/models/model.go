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

func ComposeDocument(opts map[string]string) *Document {
	dict := make(Document)
	if log_message, ok := opts["message"]; ok {
		dict["message"] = log_message
	} else {
		if opts["method"] != "DELETE" {
			dict["link"] = map[string]string{
				"rel":  "self",
				"href": "/api/cache/" + opts["key"]}
		}
		if opts["method"] == "GET" {
			dict["data"] = map[string]interface{}{
				opts["key"]:  opts["value"],
				"timeToLive": opts["timeToLive"]}
		} else {
			dict["message"] = opts["method"] + " complete."
		}
	}
	return &dict
}
