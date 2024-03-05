package auction

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"testing"
)

func TestNewPostgresStore_integration_update_writes_to_outbox(t *testing.T) {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=disable", "postgres", "postgres", "127.0.0.1:5432", "auction")
	ctx := context.Background()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatal("ping error : "+err.Error())
	}

	store := NewPostgresStore(db)

	auctionSync, err := store.GetSync(ctx, "auction-1")
	if err != nil {
		t.Fatalf("get auction %s", err.Error())
	}

	auction := auctionSync.Data
	bid := Bid{
		UserName: "user1",
		Price:    10,
	}



	toUpdate := *auction
	toUpdate.Bids = append(toUpdate.Bids, bid)
	toUpdate.CurrentHighestBid = bid

	if err := store.Update(ctx, auction, &toUpdate); err != nil {
		t.Fatalf("update error: %v", err.Error())
	}
}


func TestNewPostgresStore_integration_get_all_auction(t *testing.T) {

	connStr := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=disable", "postgres", "postgres", "127.0.0.1:5432", "auction")
	ctx := context.Background()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatal("ping error : "+err.Error())
	}

	store := NewPostgresStore(db)

	auctions, err := store.GetAll(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("%#v \n", auctions)
}


func TestNewPostgresStore_integration_create(t *testing.T) {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=disable", "postgres", "postgres", "127.0.0.1:5432", "auction")
	ctx := context.Background()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatal("ping error : "+err.Error())
	}

	store := NewPostgresStore(db)

	auction := &Auction{
		ID:                "auction-5",
		Item:              "Car",
		Bids:              []Bid{},
		CurrentHighestBid: Bid{},
		Running:           true,
	}

	if createErr := store.Create(ctx, auction); createErr != nil {
		t.Fatal("create error: ", createErr.Error())
	}
}