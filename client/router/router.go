package router

import (
	"encoding/json"
	controller "labs/redis/controllers"
	model "labs/redis/models"
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

	/** @note
			Sample JSON representation of Data Transfer Object / Payload model recieved from client,
			and the server's subsequent response
		@request
			POST | PUT /api/cache
			{
				"cacheKey": "foo",
				"cacheValue": "car",
				"ttl": 83
			}
		@response
			{
				"link": {
						"href": "/api/cache/foo",
						"rel": "self"
				},
				"message": "POST complete."
			}
	*/
	router.POST(endPointsMap["createUpdateResource"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(r.URL.Path, r.Method, w)
		controller.CreateById(w, r, p, rdb)
	})

	/** @note
			Sample JSON representation of Document data model sent to client, given
			a request to the desired resource
		@request
			GET /api/cache/foo
		@response
			{
				"data": { "foo": "car",, "timeToLive": "1m17s" },
				"link": { "href": "/api/cache/foo", "rel": "self" }
				...
			}
	*/
	router.GET(endPointsMap["getResource"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(r.URL.Path, r.Method, w)
		resourceId := p.ByName("cacheKey")
		controller.ReadById(resourceId, r, w, rdb)
	})

	/** @note
		Sample JSON representation of Document data model sent to client, given
		a request to the desired resource
		@request
			DELETE /api/cache/foo
		@response
			{
				"message": "DELETEd foo."
				...
			}
	*/
	router.DELETE(endPointsMap["deleteResource"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(r.URL.Path, r.Method, w)
		controller.DeleteById(w, r, p, rdb)
	})
	router.GET(endPointsMap["root"], func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		debugRoute(endPointsMap["root"], r.Method, w)
	})
}

func debugRoute(uri string, method string, w http.ResponseWriter) {
	dict := make(model.Document)
	dict["api"] = map[string]string{"serviceName": "memCacheLabs", method: uri, "version": " 0.0.1"}
	document, _ := json.Marshal(dict)
	log.Println(string(document))
}
