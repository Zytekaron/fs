package server

import (
	"fs/src/cache"
	"fs/src/config"
	"fs/src/middleware"
	"log"
	"net/http"
)

func Start(cfg *config.Config, cash *cache.Cache) {
	mux := http.NewServeMux()

	logger := middleware.Logger()
	reader := middleware.DataReader()

	getChain := middleware.Chain(
		logger,
		reader,
		middleware.RateLimit(cfg),
	)
	mux.Handle("GET /*", getChain(readFileHandler(cfg, cash)))

	adminChain := middleware.Chain(
		logger,
		reader,
		middleware.RequireToken([]string{cfg.Server.Tokens.Admin}),
	)
	mux.Handle("POST /admin/dump-cache", adminChain(dumpCacheHandler(cash)))

	log.Println("listening on", cfg.Server.Addr)
	err := http.ListenAndServe(cfg.Server.Addr, mux)
	if err != nil {
		log.Fatalln("listen:", err)
	}
}
