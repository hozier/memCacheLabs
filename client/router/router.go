package router

import (
	"fmt"
	handler "labs/redis/handlers"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

func CreateCacheClient() *redis.Client {
	var redis_host = os.Getenv("REDIS_HOST")
	var redis_port = os.Getenv("REDIS_PORT")
	var no_password = os.Getenv("REDIS_PASS")
	return redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + redis_port,
		Password: no_password,
		// default db index ie. Db[0]
		DB: 0,
	})
}

func NewRouter() *httprouter.Router {
	rdb := CreateCacheClient()
	r := httprouter.New()
	CacheService(r, rdb)
	return r
}

/*
 HTTP method GET services all require param in the form /:cacheKey
 HTTP method POST services all require payload in the form { "cacheKey": "foo", "cacheValue": "bar" }
 HTTP method DELETE services all require param in the form /:cacheKey
*/
func CacheService(router *httprouter.Router, rdb *redis.Client) {
	endPointsMap := map[string]string{
		"createUpdateResource": "/api/cache", "getResource": "/api/cache/:cacheKey", "deleteResource": "/api/cache/:cacheKey", "root": "/",
	}
	router.POST(endPointsMap["createUpdateResource"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(endPointsMap["createUpdateResource"], r.Method, w)
		handler.CreateKey(w, r, p, rdb)
	})
	router.GET(endPointsMap["getResource"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(endPointsMap["getResource"], r.Method, w)
		handler.ReadKey(w, r, p, rdb)
	})
	router.DELETE(endPointsMap["deleteResource"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(endPointsMap["deleteResource"], r.Method, w)
		handler.DeleteKey(w, r, p, rdb)
	})
	router.GET(endPointsMap["root"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(endPointsMap["root"], r.Method, w)
	})
}

func debugRoute(ep string, method string, w http.ResponseWriter) {
	endPoint := "{api: { serviceName: memCacheLabs, endPoint: " + ep + ", method: " + method + ", version: 0.0.1 } }"
	log.Println(endPoint)
	fmt.Fprintln(w, endPoint)
}
