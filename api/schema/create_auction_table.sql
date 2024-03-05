CREATE TABLE IF NOT EXISTS auctions (
    id VARCHAR(255) PRIMARY KEY,
    sort SERIAL,
    item VARCHAR(255),
    bids JSONB[],
    current_highest_bid JSONB,
    running BOOLEAN
);