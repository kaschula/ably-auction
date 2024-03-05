import {ModelType} from "../hooks/use-auction-model.ts";
import {MouseEvent, useState} from "react";
import {placeBid} from "../mutations/auction.ts";
import {v4 as uuidv4} from "uuid";

type BidFormProps = {
    userName: string
    auctionId: string
    auctionModel: ModelType
}

export function BidForm({userName, auctionId, auctionModel}: BidFormProps) {
    const [bidValue, setBidValue] = useState("")

    const onClickHandler = (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        const bidValueNumber = parseFloat(bidValue)
        if (isNaN(bidValueNumber)) {
            console.error(`bid value: '${bidValueNumber}' NaN`)
            return
        }

        placeBid(auctionModel, auctionId, uuidv4(), userName, bidValueNumber)
            .then(() => setBidValue(""))

    }

    return (
        <form className="row g-2" onSubmit={(e) => e.preventDefault()}>
            <div className="col-auto">
                <label htmlFor="bidValue" className="visually-hidden">Your Bid:</label>
                <input name="bidValue" className="form-control" id="bidValue" placeholder="£"
                       value={bidValue} onChange={(e) => setBidValue(e.target.value)}/>
            </div>
            <div className="col-auto">
                <button type="button" className="btn btn-primary mb-3"
                        onClick={onClickHandler}>Place Bid
                </button>
            </div>
        </form>
    )
}