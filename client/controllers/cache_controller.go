package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	model "labs/redis/models"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

var ctx = context.Background()
var counter = 0

/** @note
 *   JSON formatter which reads the given incoming JSON Object. Returns an
 *   instance of the data transfer object (DTO) upon successful deserialization.
 */
func parseReqKey(r *http.Request) *model.Payload {
	/**
	Parse key from request body
	*/
	payload := model.Payload{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	parseErr := json.Unmarshal(body, &payload)
	if parseErr != nil {
		panic(parseErr)
	}
	return &payload
}

func ReadById(resourceId string, r *http.Request, w http.ResponseWriter, rdb *redis.Client) {
	data, err := rdb.Get(ctx, resourceId).Result()
	if err == redis.Nil { // if key does not exist
		sendResponse(map[string]string{"message": "Key `" + resourceId + "` does not exist."}, w)
	} else if err != nil {
		panic(err)
	} else {
		// else key does exist
		ttl, _ := rdb.TTL(ctx, resourceId).Result()
		stringify := ttl.String()
		sendResponse(map[string]string{"method": r.Method, "key": resourceId, "value": data, "timeToLive": stringify}, w)
	}
}

func CreateById(w http.ResponseWriter, r *http.Request, p httprouter.Params, rdb *redis.Client) {
	payload := parseReqKey(r)
	cacheKey := payload.CacheKey
	cacheValue := payload.CacheValue
	// if ttl value is sent in request, the cache value will persist for only x * s
	timeToLive := (time.Duration(payload.Ttl) * time.Second)
	err := rdb.Set(ctx, cacheKey, cacheValue, timeToLive).Err()
	if err != nil {
		panic(err)
	}
	// else key does exist
	sendResponse(map[string]string{"method": r.Method, "key": cacheKey}, w)
}

func DeleteById(w http.ResponseWriter, r *http.Request, p httprouter.Params, rdb *redis.Client) {
	cacheKey := p.ByName("cacheKey")
	err := rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		panic(err)
	}
	// else key does exist
	sendResponse(map[string]string{"message": r.Method + "d " + cacheKey}, w)
}

func sendResponse(opts map[string]string, w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	dict := make(model.Document)
	if log_message, ok := opts["message"]; ok {
		dict["message"] = log_message
	} else {
		dict["resourceLocation"] = "/api/cache/" + opts["key"]
		if len(opts["value"]) > 0 {
			dict["data"] = map[string]string{opts["key"]: opts["value"], "timeToLive": opts["timeToLive"]}
		} else {
			dict["message"] = opts["method"] + " complete."
		}
	}
	document, _ := json.Marshal(dict)
	log.Println(string(document))
	fmt.Fprintln(w, string(document))
}
