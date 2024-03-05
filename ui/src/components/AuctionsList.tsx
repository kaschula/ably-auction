import {useEffect, useState} from "react";

export type AuctionResponse = {
    id: string
    item: string
    currentHighestBid: Bid
    bids: Bid[]
}

export type Bid = {
    userName: string
    price: number
}

const emptyAuctionList: AuctionResponse[] = []

type SignUpProps = {
    setAuctionId: (name: string) => void
    userName: string
}

export function AuctionsList({setAuctionId, userName}: SignUpProps) {
    const [auctions, setAuctions] = useState(emptyAuctionList)

    // api request for available auctions
    useEffect(() => {
        fetch("http://localhost:8080/auctions", {
            method: "GET", // *GET, POST, PUT, DELETE, etc.
            mode: "cors", // no-cors, *cors, same-origin
            headers: {"Content-Type": "application/json"},
        }).catch((err) => {
            console.error("fetch auctions error", err)
        }).then((response) => {
              return response?.json()
        }).then((jsonRes) => {
            setAuctions(jsonRes as AuctionResponse[])
        })
    },  []);


    const auctionItems = auctions.map((a: AuctionResponse) => {
        return (<div key={a.id} onClick={() => setAuctionId(a.id)} className="card" style={{width: "18rem"}}>
            <div className="card-body">
                <h5 className="card-title">{a.item}</h5>
            </div>
        </div>)
    })


    // show list of the items
    return <div className={"container"}>
        <div className={"row"}>
            <h1>Please select an Auction</h1>
            <p>User: {userName}</p>
            <div className="card-group">
                {auctionItems}
            </div>
        </div>
    </div>
}

