import {Bid} from "./AuctionsList.tsx";
import {v4 as uuidv4} from "uuid";

type BidHistoryProps = {
    bids: Bid[]
}

export function BidHistory({bids}: BidHistoryProps) {
    const bidElms = bids.map((bid, i) => {
        if (i === bids.length-1) {
            return <li key={uuidv4()} className="list-group-item">{bid.userName} bid: £{bid.price}</li>
        }

        return <li key={uuidv4()} className="list-group-item">{bid.userName} bid: £{bid.price}</li>
    })

    return (
        <>
            <h2>Bid History</h2>
            <ul className="list-group">
                {bidElms}
            </ul>
        </>
    )
}