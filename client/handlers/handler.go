package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

var ctx = context.Background()
var counter = 0

type Payload struct {
	CacheKey   string `json:"cacheKey"`
	CacheValue string `json:"cacheValue"`
	Ttl        int    `json:"ttl"`
}

func ReadKey(w http.ResponseWriter, r *http.Request, p httprouter.Params, rdb *redis.Client) {
	cacheKey := p.ByName("cacheKey")
	data, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil { // if key does not exist
		debugResponse("[KEY NOT FOUND] ", "", "", "", w)
	} else if err != nil {
		panic(err)
	} else {
		debugResponse("[GET] ", cacheKey, data, "", w) // else key does exist
	}
}

func CreateKey(w http.ResponseWriter, r *http.Request, p httprouter.Params, rdb *redis.Client) {
	payload := parseReqKey(r)
	cacheKey := payload.CacheKey
	cacheValue := payload.CacheValue
	oneMillisecond := (time.Second / time.Millisecond)
	timeToLive := (time.Duration(payload.Ttl) * oneMillisecond) // if ttl value is sent in request, the cache value will persist for only x * Ms
	err := rdb.Set(ctx, cacheKey, cacheValue, timeToLive).Err()
	if err != nil {
		panic(err)
	}
	var ttl = ", timeToLive: " + timeToLive.String()
	debugResponse("[PUT] ", cacheKey, cacheValue, ttl, w) // else key does exist
}

func DeleteKey(w http.ResponseWriter, r *http.Request, p httprouter.Params, rdb *redis.Client) {
	cacheKey := p.ByName("cacheKey")
	err := rdb.Del(ctx, cacheKey).Err()
	if err != nil {
		panic(err)
	}
	debugResponse("[DELETE] ", "", "", "", w) // else key does exist
}

func parseReqKey(r *http.Request) *Payload {
	/**
	Parse key from request body
	*/
	payload := Payload{}
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

func debugResponse(alert string, cacheKey string, cacheValue string, additional string, w http.ResponseWriter) {
	endPoint := alert + "{response: { " + cacheKey + ": " + cacheValue + additional + " } }"
	log.Println(endPoint)
	fmt.Fprintln(w, endPoint)
}
