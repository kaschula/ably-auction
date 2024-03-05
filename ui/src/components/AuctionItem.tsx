import {useEffect, useState} from "react";
import {AuctionResponse} from "./AuctionsList.tsx";
import {ModelType, useAuctionModel} from "../hooks/use-auction-model.ts";
import {BidHistory} from "./BidHistory.tsx";
import {BidForm} from "./PlaceBidForm.tsx";

type ItemAuctionProps = {
    auctionModel: ModelType
    auctionId: string
    userName: string
    backHandler: () => void
}

type AuctionWrapperProps = {
    auctionId: string
    userName: string
    backHandler: () => void
}

export function AuctionWrapper({userName, auctionId, backHandler}: AuctionWrapperProps) {
    const auctionModel = useAuctionModel(auctionId)

    if (!auctionModel) {
        return <div>waiting for auction model</div>
    }

    return <AuctionItem auctionId={auctionId} userName={userName} backHandler={backHandler} auctionModel={auctionModel} />
}

export function AuctionItem({auctionId, userName, backHandler, auctionModel}: ItemAuctionProps) {
    const [auction, setAuction] = useState(auctionModel.data.confirmed)

    useEffect(() => {
        if (!auctionModel) {
            return
        }

        const onUpdate = (err: Error | null, auction?: AuctionResponse) => {
            if (err) {
                console.error("onUpdate() error")
                return
            }

            setAuction(auction!)
        }

        if (auctionModel.state !== "disposed") {
            auctionModel.subscribe(onUpdate)
            return
        }

        return () => {
            auctionModel.unsubscribe(onUpdate);
        };

    }, [auctionModel]);

    if (auction === undefined) {
        return <div>Loading auction</div>
    }

    return <div className={"container"}>
        <div className={"row"}>
            <h1>Item Auction: {auction.item}</h1>
            <p>User: {userName}</p>
        </div>
        <div className={"row"}>
            <div className="col-md-6">
                <p>Current Price: {auction?.currentHighestBid?.price || 0}</p>
                <p>Highest Bidder: {auction?.currentHighestBid?.userName || "No bidder"} </p>
                <BidForm auctionModel={auctionModel} userName={userName} auctionId={auctionId}/>
                <button className={"btn btn-link"} type={"submit"} onClick={backHandler}>Leave</button>
            </div>
            <div className="col-md-6">
                <BidHistory bids={auction.bids}/>
            </div>
        </div>
    </div>

}
