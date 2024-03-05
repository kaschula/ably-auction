import Ably from "ably";
import { AblyProvider} from "ably/react";
import {Auction} from "./components/Auction.tsx";

function App() {
    const client = new Ably.Realtime.Promise({authUrl: "http://localhost:8080/auth/ably"})

    return (
        <AblyProvider client={client}>
            <Auction/>
        </AblyProvider>
    )
}

export default App
