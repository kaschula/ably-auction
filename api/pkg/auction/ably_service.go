package auction

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ably/ably-go/ably"
)



type MessageService interface {
	UserBid(ctx context.Context, channel *ably.RealtimeChannel, auction *Auction, message *ably.Message)
}

type AblyDualWriteService struct {
	store Store
}

func NewAblyService(store Store) *AblyDualWriteService {
	return &AblyDualWriteService{store: store}
}

 func (service *AblyDualWriteService) UserBid(ctx context.Context, channel *ably.RealtimeChannel, auction *Auction,  message *ably.Message) {
	auctionID := auction.ID
	strData, ok := message.Data.(string)
	if !ok {
		fmt.Println("data is not string")
		return
	}

	bidMessage := &MessageBidReceived{}
	err := json.Unmarshal([]byte(strData), bidMessage)
	if err != nil {
		fmt.Printf("bid message marshall error %s\n", err.Error())
		return
	}

	bid := NewBid(bidMessage.UserName, bidMessage.Value)

	fmt.Printf("new bid for %v recieved %#v", auction.Item, bid)

	retrieved, err := service.store.Get(ctx, auctionID)
	if err != nil {
		fmt.Printf("error get auction %s\n", err.Error())
		return
	}

	toUpdate := *retrieved

	toUpdate.Bids = append(auction.Bids, bid)

	newHighBid := false
	if retrieved.CurrentHighestBid.Price < bid.Price {
		toUpdate.CurrentHighestBid = bid
		newHighBid = true
	}

	if err := service.store.Update(ctx, retrieved, &toUpdate); err != nil {
		fmt.Printf("store update %s", err.Error())
		return
	}
	if newHighBid {
		err = channel.Publish(ctx, AUCTION_NEW_HIGH_BID, bid)
		if err != nil {
			fmt.Printf("unable to publish %s", err.Error())
			return
		}
	}
}

type AblyLiveSyncService struct {
	store Store
}

func NewAblyLiveSyncService(store Store) *AblyLiveSyncService {
	return &AblyLiveSyncService{store: store}
}

func (service *AblyLiveSyncService) UserBid(ctx context.Context, _ *ably.RealtimeChannel, auction *Auction,  message *ably.Message) {
	auctionID := auction.ID
	strData, ok := message.Data.(string)
	if !ok {
		fmt.Println("data is not string")
		return
	}

	bidMessage := &MessageBidReceived{}
	err := json.Unmarshal([]byte(strData), bidMessage)
	if err != nil {
		fmt.Printf("bid message marshall error %s\n", err.Error())
		return
	}

	bid := NewBid(bidMessage.UserName, bidMessage.Value)

	fmt.Printf("new bid for %v recieved %#v", auction.Item, bid)

	retrieved, err := service.store.Get(ctx, auctionID)
	if err != nil {
		fmt.Printf("error get auction %s\n", err.Error())
		return
	}

	toUpdate := *retrieved
	toUpdate.Bids = append(auction.Bids, bid)

	if retrieved.CurrentHighestBid.Price < bid.Price {
		toUpdate.CurrentHighestBid = bid
	}

	if err := service.store.Update(ctx, retrieved, &toUpdate); err != nil {
		fmt.Printf("store update %s", err.Error())
		return
	}
	// new highest bid publish logic handled by store
}