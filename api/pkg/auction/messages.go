package auction

const (
	AUCTION_USER_JOINED = "USER_JOINED"
	AUCTION_BID_PLACED = "EVENT_BID_PLACED"
	AUCTION_NEW_HIGH_BID = "NEW_HIGH_BID"
)


type MessageBidReceived struct {
	UserName string `json:"userName"`
	Value float64 `json:"value"`
}