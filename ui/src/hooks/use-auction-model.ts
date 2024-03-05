import {useEffect, useState} from "react";
import {useAbly} from "ably/react";
import ModelsClient, {
    backoffRetryStrategy,
    ConfirmedEvent,
    Model,
    OptimisticEvent,
    SyncReturnType
} from '@ably-labs/models';
import {AuctionResponse, Bid} from "../components/AuctionsList.tsx";
import {EVENT_BID_PLACED} from "../mutations/auction.ts";


export type ModelType = Model<() => SyncReturnType<AuctionResponse>>;

export function useAuctionModel(id: string): ModelType | undefined {
    const [model, setModel] = useState<ModelType|undefined>()
    const ably = useAbly()

    useEffect(() => {
        const modelsClient = new ModelsClient({
            ably,
            logLevel: "info",
            optimisticEventOptions: { timeout: 5000 },
            syncOptions: { retryStrategy: backoffRetryStrategy(2, 125,  1, 1000) }
        })

        async function initialise() {
            const model = modelsClient.models.get({
                name: "auctions:" + id, // this cause an error
                channelName: id,
                sync: async () => await getAuction(id),
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-expect-error
                merge
            })

            await model.sync()
            setModel(model)
            // sync model
        }

        if (!model) {
            initialise().catch((err) => console.error("useAuctionModel().initialise() error:", err))
        }

        // return clean up
        return () => {
            console.log("useAuctionModel() cancel")
            // test this dispose
            model?.dispose()
        }
    }, [id])

    return model
}

async function getAuction(id: string) {
    const response = await fetch(`http://localhost:8080/auctions/${id}/sync`)
    return response.json()
}

function merge(original: AuctionResponse, event: OptimisticEvent | ConfirmedEvent): AuctionResponse {
        const state = copyAuction(original)
        if (event.name === "NEW_HIGH_BID") {
            state.currentHighestBid = event.data as Bid
        }

        if (event.name === EVENT_BID_PLACED) {
            state.bids = [...state.bids, event.data as Bid]
        }

        return state
}

function copyAuction({id, item, currentHighestBid, bids }: AuctionResponse): AuctionResponse {
    return {
        id,
        item,
        currentHighestBid: {userName: currentHighestBid.userName, price: currentHighestBid.price},
        bids: [...bids]
    }
}