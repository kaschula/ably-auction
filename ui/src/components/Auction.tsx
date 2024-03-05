import {useState} from "react";
import {AuctionWrapper} from "./AuctionItem.tsx";
import {SignIn} from "./SignIn.tsx";
import {AuctionsList} from "./AuctionsList.tsx";


export function Auction() {
    const [userName, setUserName] = useState("")
    const [auctionId, setAuctionId] = useState("")

    if (userName == "") {
        return <SignIn setUserName={setUserName}/>
    }

    // set auction
    if (auctionId === "") {
        return <AuctionsList setAuctionId={setAuctionId} userName={userName}/>
    }

    return <div>
        <AuctionWrapper userName={userName} auctionId={auctionId} backHandler={() => {
            setAuctionId("")
        }}/>
    </div>
}