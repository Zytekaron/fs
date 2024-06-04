package server

import (
	"fs/src/cache"
	"fs/src/server/response"
	"log"
	"net/http"
)

func dumpCacheHandler(c *cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		size := c.Size()
		c.Clear()

		log.Println("dumping cache", size)

		response.WriteJSON(w, http.StatusOK, map[string]any{
			"size": size,
		})
	})
}
