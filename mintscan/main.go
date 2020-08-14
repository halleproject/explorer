package main

import (
	"context"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/client"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/config"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/db"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/handlers"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "Mintscan API ", log.Lshortfile|log.LstdFlags)

	cfg := config.ParseConfig()

	client := client.NewClient(
		cfg.Node,
		cfg.Market,
	)

	db := db.Connect(cfg.DB)
	err := db.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to ping database"))
	}
	db.CreateTables()

	r := mux.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:         []string{"*"},
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods:         []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:         []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:         []string{"Link"},
		AllowCredentials:       true,
		OptionsPassthrough:     true,
		MaxAge:                 3599, // Maximum value not ignored by any of major browsers
	})

	r.Use(cors.Handler)

	getR := r.Methods(http.MethodGet).PathPrefix("/v1").Subrouter()
	getR.HandleFunc("/account/{address}", handlers.NewAccount(l, client, db).GetAccount)
	getR.HandleFunc("/account/txs/{address}", handlers.NewAccount(l, client, db).GetAccountTxs)
	getR.HandleFunc("/asset", handlers.NewAsset(l, client, db).GetAsset)
	getR.HandleFunc("/assets", handlers.NewAsset(l, client, db).GetAssets)
	getR.HandleFunc("/assets/txs", handlers.NewAsset(l, client, db).GetAssetTxs)
	getR.HandleFunc("/asset-holders", handlers.NewAsset(l, client, db).GetAssetHolders)
	getR.HandleFunc("/assets-images", handlers.NewAsset(l, client, db).GetAssetsImages)
	getR.HandleFunc("/blocks", handlers.NewBlock(l, client, db).GetBlocks)
	getR.HandleFunc("/fees", handlers.NewFee(l, client, db).GetFees)
	getR.HandleFunc("/validators", handlers.NewValidator(l, client, db, cfg.Node.NetworkType).GetValidators)
	getR.HandleFunc("/validator/{address}", handlers.NewValidator(l, client, db, cfg.Node.NetworkType).GetValidator)
	getR.HandleFunc("/market", handlers.NewMarket(l, client, db).GetCoinMarketData)
	getR.HandleFunc("/market/chart", handlers.NewMarket(l, client, db).GetCoinMarketChartData)
	getR.HandleFunc("/orders/{id}", handlers.NewOrder(l, client, db).GetOrders)
	getR.HandleFunc("/stats/assets/chart", handlers.NewStatistic(l, client, db).GetAssetsChartHistory)
	getR.HandleFunc("/status", handlers.NewStatus(l, client, db).GetStatus)
	getR.HandleFunc("/tokens", handlers.NewToken(l, client, db).GetTokens)
	getR.HandleFunc("/txs", handlers.NewTransaction(l, client, db).GetTxs)
	getR.HandleFunc("/txs/{hash}", handlers.NewTransaction(l, client, db).GetTxByHash)
	getR.HandleFunc("/txs_address", handlers.NewTransaction(l, client, db).GetTxsByAddress)
	getR.HandleFunc("/two_auth/generate", handlers.NewTwoAuth(l, client, db).Generate)
	getR.HandleFunc("/two_auth/auth", handlers.NewTwoAuth(l, client, db).Auth)
	getR.HandleFunc("/app_version/set", handlers.NewAppVersion(l, client, db).SetVersion)
	getR.HandleFunc("/app_version/get", handlers.NewAppVersion(l, client, db).GetVersion)

	postR := r.Methods(http.MethodPost).PathPrefix("/v1").Subrouter()
	postR.HandleFunc("/txs", handlers.NewTransaction(l, client, db).GetTxsByType)

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // catch-all
		w.Write([]byte("No route is found matching the URL"))
	})

	// create a new server
	sm := &http.Server{
		Addr:         ":" + cfg.Web.Port,
		Handler:      r,
		ErrorLog:     l,
		ReadTimeout:  50 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Printf("Server is running on http://localhost:%s\n", cfg.Web.Port)

		err := sm.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sm.Shutdown(ctx)

	l.Println("Gracefully shutting down the server: ", sig)
}
