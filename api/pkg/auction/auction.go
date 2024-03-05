package auction

import (
	"context"
	"fmt"
	"github.com/ably/ably-go/ably"
)

type Auction struct {
	ID string `json:"id"`
	Sort int `json:"sort"`
	Item string `json:"item"`
	Bids []Bid `json:"bids"`
	CurrentHighestBid Bid `json:"currentHighestBid"`
	Running bool `json:"running"`
}



func NewAuction(ID string, item string, sort int) *Auction {
	return &Auction{ID: ID, Item: item, Bids: []Bid{}, CurrentHighestBid: Bid{}, Running: false, Sort: sort}
}


type Bid struct {
	UserName string `json:"userName"`
	Price float64 `json:"price"`
}

func NewBid(userName string, price float64) Bid {
	return Bid{UserName: userName, Price: price}
}


type Runner struct {
	auction *Auction
	ably  *ably.Realtime
	store Store
	ablyService MessageService
	done  chan bool
}

func NewRunner(auction *Auction, ably *ably.Realtime, store Store, ablyService MessageService, done chan bool) *Runner {
	return &Runner{auction: auction, ably: ably, store: store, ablyService: ablyService, done: done}
}

func (runner *Runner) Start(ctx context.Context)  {
	update := *runner.auction
	update.Running = true

	if err := runner.store.Update(ctx, runner.auction, &update); err != nil {
		panic(err.Error())
	}
	runner.auction = &update

	fmt.Printf("auction started %v(%v)\n", runner.auction.Item,  runner.auction.ID)

	<-runner.done
}
