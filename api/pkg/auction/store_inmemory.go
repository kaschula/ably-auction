package auction

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type InMemoryAuctionStore struct {
	data sync.Map
}

func NewInMemoryActionStore() *InMemoryAuctionStore {
	return &InMemoryAuctionStore{data: sync.Map{}}
}

func (store *InMemoryAuctionStore) Get(_ context.Context, id string) (*Auction, error) {
	item, found := store.data.Load(id)
	if !found {
		return nil, fmt.Errorf("acution not found")
	}

	auction, ok := item.(*Auction)
	if !ok {
		return nil, fmt.Errorf("can not cast item is type %t expected *Action", item)
	}

	return auction, nil
}

func (store *InMemoryAuctionStore) GetSync(_ context.Context, id string) (*AuctionSync, error) {
	item, found := store.data.Load(id)
	if !found {
		return nil, fmt.Errorf("acution not found")
	}

	auction, ok := item.(*Auction)
	if !ok {
		return nil, fmt.Errorf("can not cast item is type %t expected *Action", item)
	}

	return &AuctionSync{SequenceID: 0, Data: auction}, nil
}

func (store *InMemoryAuctionStore) Update(_ context.Context, _, toUpdate *Auction) error {
	store.data.Store(toUpdate.ID, toUpdate)

	return nil
}

func (store *InMemoryAuctionStore) GetAll(_ context.Context) ([]*Auction, error) {
	var as []*Auction
	store.data.Range(func(key, value any) bool {
		auction, ok := value.(*Auction)
		if !ok {
			fmt.Printf("get all error item not type Auction but %T\n", value)
			return true
		}

		as = append(as, auction)

		return true
	})

	sort.Slice(as, func(i, j int) bool {
		return as[i].Sort < as[j].Sort
	})

	return as, nil
}


func (store *InMemoryAuctionStore) AppendBid(ctx context.Context, auctionID string, bid Bid, mutatiionID string) error {
	return errors.New("not implemented")
}


