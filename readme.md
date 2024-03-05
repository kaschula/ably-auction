# Ably Auction

The purpose of this application is to demo the ably livesync products. It uses the livesync models SDK and Database connector.

## Overview
The app consists of 3 parts. THe fronend UI (/ui), the backend API server (/api) and the Ably Database Connector (ADBC).

The frontend application uses React. Due the simplicity of the app there is no routing. The application flow is as follows:-

A user navigates to the site. Typically this is at `localhost:5173`
A user enter there username
The user then selects from one of the running auctions
Once inside the auction a should see the current auction price and who is the leader
user can place a bid

The front-end application uses React. Due to the simplicity of the app, there is no routing. The application flow is as follows:-

1. A user navigates to the site. Typically this is at `localhost:5173`
2. A user enters their username
3. The user then selects from one of the running auctions
4. Once inside the auction, a user should see the current auction price and who is the leader
user can place a bid


The backend API `/auctions/:id/bid` processes the Bid placement.

A new Bid is written to the Postgres `auctions` table and if the bid has a new high value then an `outbox` record is written to be processed by the ADBC.

The model SDK will update all the client's auction views when a new High bid is received.

## How to run

This app is designed to be run locally.

## 1 Create the DB layer

To run the backend DB and connector there are two config files that need to be created, using the example files `adbc.example.env`
`adbc.example.yaml`.

Create `adbc.env` and set the `ADBC_ABLY_API_KEY` environment variable to your ably API key.

Create a `adbc.yaml` file and set the `ably.apiKey` value to your ably API key.

Run `docker compose up --build -d` from the root of the project. This should create the ADBC and Postgres DB containers.
Once the postgres DB is running connect to the postgres container either with client or terminal and run the create auctions script found in `api/schema/create_auction_table.sql`

## 2 Run the Backend API

From the `/api` directory run:

```
export AUCTION_ABLY_API_KEY=your-ably-api-key && go run ./api/cmd/api/...
```
 
This will fail if no DB is running

## 3. Run the frontend

From the `/iu` directory run `npm run dev`. 

The Auction API handles the authentication of the Ably SDK.

Now visit it `http://localhost:5173/` 




## Improvements

There are no optimistic updates currently in this app. It would be good to add some.

One possibility is the Auction model currently holds a list of bids. The Bid object could be updated track the time at 
which a bid is created. When a new bid is added it could be optimistically added to this list. This would mean bids under the current total would be valid.


