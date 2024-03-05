import {ModelType} from "../hooks/use-auction-model.ts";

export const EVENT_BID_PLACED = "EVENT_BID_PLACED"

export async function placeBid(model: ModelType, auctionId: string, mutationId: string, userName: string, bidValue: number ) {
    const bid = {userName, price: bidValue}

    // do the optimistic update
    const [confirmation, cancel] = await model.optimistic({
        mutationID: mutationId,
        name: EVENT_BID_PLACED,
        data: bid,
    });

    try {
        // do the actual update
        const response = await fetch(`http://localhost:8080/auctions/${auctionId}/bid`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ mutationId, bid }),
        });

        // confirm optimistic update
        await confirmation
        if (!response.ok) {
            throw new Error(`PUT /auctions/:auctionId/bid: ${response.status} ${JSON.stringify(await response.json())}`);
        }

        return response.json();
    }catch (err) {
        cancel().catch((e) => {
            console.error("optimistic update cancel failed", e.message)
        });
    }
}