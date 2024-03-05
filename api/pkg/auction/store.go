package auction

import "context"

type Store interface{
	Get(ctx context.Context, id string) (*Auction, error)
	GetSync(ctx context.Context, id string) (*AuctionSync, error)
	Update(ctx context.Context, original, update *Auction) error
	GetAll(ctx context.Context) ([]*Auction, error)
	AppendBid(ctx context.Context, auctionID string, bid Bid, mutationID string) error
}


