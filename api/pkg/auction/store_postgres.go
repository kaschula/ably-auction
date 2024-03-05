package auction

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sort"
)


var (
	ErrNotFound = errors.New("auction not found")
)

type Record struct {
	ID string `json:"id"`
	Sort int `json:"sort"`
	Item string `json:"item"`
	Bids []byte `json:"bids"`
	CurrentHighestBid []byte `json:"currentHighestBid"`
	Running bool `json:"running"`
}

// live sync response
type AuctionSync struct {
	SequenceID int `json:"sequenceID"`
	Data *Auction `json:"data"`
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (store *PostgresStore) Create(ctx context.Context, auction *Auction) error {
	insertStatement := `INSERT INTO auctions (id, item, bids, current_highest_bid, running) VALUES ($1, $2, $3, $4, $5)`

	bidsJSON, err := marshallBids(auction.Bids)
	if err != nil {
		return err
	}

	bidJSON, err := marshallBid(auction.CurrentHighestBid)
	if err != nil {
		return err
	}

	_, err = store.db.ExecContext(ctx, insertStatement,
		auction.ID, auction.Item, bidsJSON, bidJSON, auction.Running)

	if err != nil {
		return fmt.Errorf("failed to create auction statement:%w", err)
	}


	return nil
}

func (store *PostgresStore) Get(ctx context.Context, id string) (*Auction, error) {
	sqlStatement := `SELECT * FROM auctions WHERE id = $1`
	var auction Record
	err := store.db.QueryRowContext(ctx, sqlStatement, id).
		Scan(&auction.ID, &auction.Sort, &auction.Item, &auction.Bids, &auction.CurrentHighestBid, &auction.Running)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
	}

	bids, err := unmarshallBids(auction.Bids)
	if err != nil {
		return nil, err
	}

	currentHighest, err := unmarshallBid(auction.CurrentHighestBid)
	if err != nil {
		return nil, err
	}

	return &Auction{
		ID:                auction.ID,
		Sort:              auction.Sort,
		Item:              auction.Item,
		Bids:              bids,
		CurrentHighestBid: currentHighest,
		Running:           auction.Running,
	}, nil
	//return &auction, nil
}

func (store *PostgresStore) GetSync(ctx context.Context, id string) (*AuctionSync, error) {
	var err error
	var tx *sql.Tx
	defer func() {
		if err != nil && tx != nil {
			if err := tx.Rollback(); err != nil {
				fmt.Printf("get sync rollback error %s", err.Error())
			}
		}
	}()

	tx, err = store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w start transaction error",err)
	}

	auction, err := getAuction(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	sequenceID, err := getOutboxSequenceID(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit error: %w",err)
	}



	return &AuctionSync{
		SequenceID: sequenceID,
		Data: auction,
	}, nil
}

func (store *PostgresStore) Update(ctx context.Context, original, toUpdate *Auction) error {
	var err error
	var tx *sql.Tx
	defer func() {
		if tx != nil && err != nil {
			err := tx.Rollback()
			if err != nil {
				fmt.Printf("transaction roll back error: %s", err.Error())
			}
		}
	}()

	// determine if bid is high
	newHighestBid := false
	if original.CurrentHighestBid.Price < toUpdate.CurrentHighestBid.Price {
		newHighestBid = true
	}

	tx, err = store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w ", err)
	}

	if err := updateAuction(ctx, tx, toUpdate, newHighestBid, uuid.NewString()); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit update failed transaction: %w", err)
	}
	

	return nil
}

func (store *PostgresStore) GetAll(ctx context.Context) ([]*Auction, error) {
	sqlStatement := `
		SELECT * FROM auctions
	`

	rows, err := store.db.QueryContext(ctx, sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error get all auctions: %w", err)
	}
	defer rows.Close()

	var auctions []*Auction
	// Iterate through the result set
	for rows.Next() {
		var record Record
		err := rows.Scan(&record.ID, &record.Sort, &record.Item, &record.Bids, &record.CurrentHighestBid, &record.Running)
		if err != nil {
			return nil,  fmt.Errorf("error scanning row: %w", err)
		}

		bids, err := unmarshallBids(record.Bids)
		if err != nil {
			fmt.Printf("marshall bids error: %v", err.Error())
			continue
		}


		bid, err := unmarshallBid(record.CurrentHighestBid)
		if err != nil {
			fmt.Printf("marshall bids error: %v", err.Error())
			continue
		}

		auctions = append(auctions, &Auction{
			ID:                record.ID,
			Sort:              record.Sort,
			Item:              record.Item,
			Bids:              bids,
			CurrentHighestBid: bid,
			Running:           record.Running,
		})
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating through result set:", err)
	}


	sort.Slice(auctions, func(i, j int) bool {
		return auctions[i].Sort > auctions[j].Sort
	})

	return auctions, nil
}


func (store *PostgresStore) AppendBid(ctx context.Context, auctionID string, bid Bid, mutationID string) error {
	var tx *sql.Tx
	var err error
	defer func() {
		if tx != nil && err != nil {
			err := tx.Rollback()
			if err != nil {
				fmt.Printf("transaction roll back error: %s", err.Error())
			}
		}
	}()


	tx, err = store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w ", err)
	}


	// get the Auction
	auction , err := getAuction(ctx, tx, auctionID)
	if err != nil {
		return err
	}
	// compare the bids
	toUpdate := *auction
	toUpdate.Bids = append(toUpdate.Bids, bid)

	newHighBid := false
	if auction.CurrentHighestBid.Price < bid.Price {
		newHighBid = true
		toUpdate.CurrentHighestBid = bid
	}

	if err = updateAuction(ctx, tx, &toUpdate, newHighBid, mutationID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit update failed transaction: %w", err)
	}

	return nil
}


// todo rework this, its a bit of a mess
func updateAuction(ctx context.Context, tx *sql.Tx, toUpdate *Auction, newHighestBid bool, mutationID string) error {
	sqlStatement := `
		UPDATE auctions
		SET
			item = $1,
			bids = $2,
			current_highest_bid = $3,
			running = $4
		WHERE
			id = $5
	`


	var err error
	var latestBidJSON []byte
	if len(toUpdate.Bids) > 0 {
		latestBidJSON, err = marshallBid(toUpdate.Bids[len(toUpdate.Bids)-1]); if err != nil {
			return err
		}
	}

	bidsJSON, err := marshallBids(toUpdate.Bids)
	if err != nil {
		return err
	}

	currentHighestBidJSON, err := marshallBid(toUpdate.CurrentHighestBid)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, sqlStatement,
		toUpdate.Item,
		bidsJSON,
		currentHighestBidJSON,
		toUpdate.Running,
		toUpdate.ID,
	)
	if err != nil {
		return fmt.Errorf("error executing update statement: %w", err)
	}

	if len(latestBidJSON) > 0 {
		if err = insertOutbox(ctx, tx, toUpdate.ID, AUCTION_BID_PLACED, mutationID, latestBidJSON); err != nil {
			return err
		}
	}

	if newHighestBid {
		// this event does not optimistically update so set teh mutationID here
		if err = insertOutbox(ctx, tx, toUpdate.ID, AUCTION_NEW_HIGH_BID, uuid.NewString(), currentHighestBidJSON); err != nil {
			return err
		}
	}

	return nil
}

func getAuction(ctx context.Context, tx *sql.Tx, id string) (*Auction, error) {
	sqlStatement := `SELECT * FROM auctions WHERE id = $1`
	var auction Record
	err := tx.QueryRowContext(ctx, sqlStatement, id).
		Scan(&auction.ID, &auction.Sort, &auction.Item, &auction.Bids, &auction.CurrentHighestBid, &auction.Running)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
	}

	// parse bids
	bids, err := unmarshallBids(auction.Bids)
	if err != nil {
		return nil, err
	}

	currentHighest, err := unmarshallBid(auction.CurrentHighestBid)
	if err != nil {
		return nil, err
	}

	return &Auction{
		ID:                auction.ID,
		Sort:              auction.Sort,
		Item:              auction.Item,
		Bids:              bids,
		CurrentHighestBid: currentHighest,
		Running:           auction.Running,
	}, nil
}

func insertOutbox(ctx context.Context, tx *sql.Tx, channelName, messageName, id string, data []byte) error {
	var statement =  `INSERT INTO outbox (channel, name, mutation_id, data) VALUES ($1,$2,$3,$4);`

	if _, err := tx.ExecContext(ctx, statement, channelName, messageName, id, data); err != nil {
		return fmt.Errorf("outbox insert failed: %w", err)
	}

	return nil
}


func getOutboxSequenceID(ctx context.Context, tx *sql.Tx) (int, error) {
	var sequenceID int
	sequenceStatement := "SELECT COALESCE(MAX(sequence_id), 0) FROM outbox"
	err := tx.QueryRowContext(ctx, sequenceStatement).Scan(&sequenceID)
	if err != nil {
		return 0, fmt.Errorf("sequence query error: %w",err)
	}

	return sequenceID, nil
}

func unmarshallBids(data []byte) ([]Bid, error) {
	var bids []Bid
	err := json.Unmarshal(data, &bids)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall bids: %w", err)
	}

	return bids, nil
}

func unmarshallBid(data []byte) (Bid, error) {
	var bid Bid
	err := json.Unmarshal(data, &bid)
	if err != nil {
		return Bid{},  fmt.Errorf("failed to unmarshall bid: %w", err)
	}

	return bid, nil
}

func marshallBids(bids []Bid) ([]byte, error) {
	raw, err := json.Marshal(bids)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall bids: %w", err)
	}

	return raw, nil
}

func marshallBid(bid Bid) ([]byte, error) {
	raw, err := json.Marshal(bid)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall bid: %w", err)
	}

	return raw, nil
}





