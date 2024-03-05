package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github/kaschula/ably-auction/api/pkg/auction"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

type PostAuctionBid struct {
	MutationId string `json:"mutationId"`
	Bid auction.Bid `json:"bid"`
}

type HttpHandlerService struct {
	auctionStore auction.Store
}

func NewHttpHandlerService(auctionStore auction.Store) *HttpHandlerService {
	return &HttpHandlerService{auctionStore: auctionStore}
}
func (service *HttpHandlerService) HandlerListAuction(c echo.Context) error {
	c.Logger().Info("list_auction_handler")

	auctions, err := service.auctionStore.GetAll(c.Request().Context())
	if err != nil {
		c.Logger().Info("get all auctions error", auctions)
		return err
	}

	return c.JSON(http.StatusOK, auctions)
}

func (service *HttpHandlerService) HandlerGetAuction(c echo.Context) error {
	c.Logger().Info("auction_get_handler")

	auctionID := c.Param("id")
	if auctionID == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "no auction id provided"})
	}

	auctionItem, err := service.auctionStore.Get(c.Request().Context(), auctionID)
	if err != nil {
		c.Logger().Info("get auction error", auctionItem)
		return err
	}

	return c.JSON(http.StatusOK, auctionItem)
}


func (service *HttpHandlerService) HandlerGetSyncAuction(c echo.Context) error {
	c.Logger().Info("auction_get_sync_handler")

	auctionID := c.Param("id")
	if auctionID == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "no auction id provided"})
	}

	auctionItem, err := service.auctionStore.GetSync(c.Request().Context(), auctionID)
	if err != nil {
		c.Logger().Info("get auction sync error", auctionItem, err)
		return err
	}

	return c.JSON(http.StatusOK, auctionItem)
}



func (service *HttpHandlerService) HandlerPostAuctionBid(c echo.Context) error {
	c.Logger().Info("auction_post_new_bid_handler")

	auctionID := c.Param("id")
	if auctionID == "" {
		return c.JSON(http.StatusBadRequest, Error{Message: "no auction id provided"})
	}


	body := &PostAuctionBid{}
	if err := c.Bind(body); err != nil {
		return fmt.Errorf("body parse: %w", err)
	}

	if err := service.auctionStore.AppendBid(c.Request().Context(), auctionID, body.Bid, body.MutationId); err != nil {
		c.Logger().Info("append bid error", auctionID)
		return err
	}

	return c.JSON(http.StatusOK, struct{}{})
}
