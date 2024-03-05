package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/ably/ably-go/ably"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github/kaschula/ably-auction/api/pkg/auction"
)

var API_KEY = ""

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Logger.SetLevel(log.DEBUG)
	err := resolveEnvVars()
	if err != nil {
		log.Fatal(err.Error())
	}

	done := make(chan bool)
	connStr := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=disable", "postgres", "postgres", "127.0.0.1:5432", "auction")
	ctx := context.Background()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	auctionStore := auction.NewPostgresStore(db)
	ablyService := auction.NewAblyService(auctionStore)
	realtimeClient, err := ably.NewRealtime(ably.WithKey(API_KEY))
	if err != nil {
		panic("ably connection error")
	}

	httpService := NewHttpHandlerService(auctionStore)

	startAuction(ctx, auctionStore, realtimeClient, ablyService, done)



	e.GET("/", HandlerIndex)
	e.GET("/auth/ably", HandlerAblyAuth)
	e.GET("/auctions", httpService.HandlerListAuction)
	e.GET("/auctions/:id", httpService.HandlerGetAuction)
	e.GET("/auctions/:id/sync", httpService.HandlerGetSyncAuction)
	e.POST("/auctions/:id/bid", httpService.HandlerPostAuctionBid)
	e.Logger.Fatal(e.Start(":8080"))
}

func startAuction(ctx context.Context, auctionStore *auction.PostgresStore, realtimeClient *ably.Realtime, ablyService auction.MessageService, done chan bool) {
	auctions, err := auctionStore.GetAll(ctx)
	if err != nil {
		panic("unable to get auctions: "+ err.Error())
	}

	for _, auctionToStart := range auctions {
		runner := auction.NewRunner(auctionToStart, realtimeClient, auctionStore, ablyService, done)
		go runner.Start(ctx)
	}
}


func resolveEnvVars() error {
	API_KEY = os.Getenv("AUCTION_ABLY_API_KEY")
	if API_KEY == "" {
		return errors.New("the env var AUCTION_ABLY_API_KEY is empty")
	}

	return nil
}