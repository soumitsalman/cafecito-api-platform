# Unusual Whales API Inventory

| Field | Value |
|---|---|
| Status | **current** (external inventory) |
| Authority | Unusual Whales public OpenAPI as of retrieval date; **not** Cafecito contract |
| Audience | Product research |
| Last verified | 2026-08-25 |
| Owner role | Product research |
| Superseded by | n/a |

> Verified against the live Unusual Whales OpenAPI specification and operation documentation retrieved August 14, 2026.

- **Documented groups:** 35
- **Documented operations:** 208 total — 207 `GET`, 1 `POST`
- **Base URL:** `https://api.unusualwhales.com`
- **Authentication:** `Authorization: Bearer <API_TOKEN>`

## Executive summary

Unusual Whales provides a broad financial-market data platform centered on options activity and market microstructure. The current API reference covers options flow, option-contract analytics, dark-pool and lit-market trades, Greek exposure, volatility, stock fundamentals, institutional filings, congressional and politician trading, futures, crypto, FX, macro calendars, prediction markets, private markets, news, screeners, and real-time streams.

## Access and delivery

- REST API with JSON responses and Bearer-token authentication.
- WebSocket streaming for live option trades, GEX, flow, prices, news, trading halts, and related
  market feeds.
- Kafka enterprise streaming for high-throughput feeds such as option trades, insider trades,
  equities, and other market events.
- MCP server access to 50+ financial-data tools for AI assistants.

## Historical and downloadable data

The linked Data Shop advertises downloadable historical datasets for options flow, OHLC/price data, insider trades, GEX, dark-pool data, Market Tide, IV Rank, and 13F filings.

## How route entries are documented

Each route below is reconciled to the OpenAPI source. The entry includes the official operation summary and description excerpt, every documented path/query parameter, any documented JSON request body, the `200` response schema and top-level payload fields, and documented non-success responses. Payload field lists are schema metadata, not live response captures; nested objects are represented by their OpenAPI schema names.

## Primary sources

- [Public API overview](https://unusualwhales.com/public-api)
- [API documentation](https://api.unusualwhales.com/docs)
- [OpenAPI specification](https://api.unusualwhales.com/api/openapi)
- [WebSocket channel
  documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.channels)

## Exhaustive endpoint inventory

### Futures — 5

#### `GET /api/futures/{contract}/candles` — Futures Candles (OHLCV)

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.FuturesController.candles)

**What it does:**

OHLCV candles for a contract. Available on the Advanced API tier, or with the `futures` add-on —
contact oskar@unusualwhales.com for access.

**Parameters**

- **path `contract`** (required; type `string`) — Dated CME contract symbol, e.g. ESU6.
- **query `interval`** (optional; type `string`) — Candle interval: 1m (default), 5m, or 1d.
- **query `range`** (optional; type `string`) — History range: 1d (default), 5d, 1w, 1m, 3m, 1y.

**Response payload**

- `200`: `application/json` → `Futures Candles` — Futures Candles; fields: `data` (array<object>);
  data item `object` fields: `c` (string), `date` (string (date)), `end` (string (date-time)), `h`
  (string), `l` (string), `o` (string), `start` (string (date-time)), `v` (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/futures/contracts` — Futures Contracts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.FuturesController.contracts)

**What it does:**

Active CME futures contracts with trading activity over the last `days` (default 3, max 30), most
active first. Available on the Advanced API tier, or with the `futures` add-on — contact
oskar@unusualwhales.com for access.

**Parameters**

- **query `days`** (optional; type `integer`) — Lookback window in days (default 3, max 30).

**Response payload**

- `200`: `application/json` → `Futures Contracts` — Futures Contracts; fields: `data`
  (array<object>); data item `object` fields: `high` (string), `is_spread` (boolean), `last_trade`
  (string (date-time)), `low` (string), `name` (string), `product` (string), `security_id` (string),
  `trades` (integer), `volume` (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/futures/flow` — Futures Flow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.FuturesController.flow)

**What it does:**

Newest-first trade prints across ALL contracts. Optional server-side filters narrow the feed.
Available on the Advanced API tier, or with the `futures` add-on — contact oskar@unusualwhales.com
for access.

**Parameters**

- **query `limit`** (optional; type `integer`) — Max rows (default 100, max 500).
- **query `older_than`** (optional; type `string`) — Cursor: return trades executed before this
  ISO8601 timestamp.
- **query `newer_than`** (optional; type `string`) — Cursor: return trades executed after this
  ISO8601 timestamp.
- **query `products`** (optional; type `string`) — Comma-separated CME product codes to include,
  e.g. ES,NQ.
- **query `side`** (optional; type `string`) — Filter by aggressor side: buy or sell.
- **query `min_size`** (optional; type `number`) — Minimum trade size (contracts).
- **query `max_size`** (optional; type `number`) — Maximum trade size (contracts).
- **query `min_price`** (optional; type `number`) — Minimum trade price.
- **query `max_price`** (optional; type `number`) — Maximum trade price.
- **query `include_spreads`** (optional; type `string`) — Set to false to exclude calendar/spread
  contracts.

**Response payload**

- `200`: `application/json` → `Futures Flow` — Futures Flow; fields: `data` (array<object>); data
  item `object` fields: `crossed` (boolean), `executed_at` (string (date-time)), `nbbo_ask`
  (string), `nbbo_bid` (string), `price` (string), `product` (string), `side` (string), `size`
  (integer), `sym` (string), `trade_id` (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/futures/{contract}/stats` — Futures Session Stats

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.FuturesController.stats)

**What it does:**

Latest session statistics for a contract: official settlement (day-change reference), settlement
date, open interest and session volume. Available on the Advanced API tier, or with the `futures`
add-on — contact oskar@unusualwhales.com for access.

**Parameters**

- **path `contract`** (required; type `string`) — Dated CME contract symbol, e.g. ESU6.

**Response payload**

- `200`: `application/json` → `Futures Stats` — Futures Stats; fields: `data` (object)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/futures/{contract}/trades` — Futures Trades (Time & Sales)

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.FuturesController.trades)

**What it does:**

Newest-first trade prints for a contract, cursor-paginated via `older_than`/`newer_than`. Available
on the Advanced API tier, or with the `futures` add-on — contact oskar@unusualwhales.com for access.

**Parameters**

- **path `contract`** (required; type `string`) — Dated CME contract symbol, e.g. ESU6.
- **query `limit`** (optional; type `integer`) — Max rows (default 100, max 500).
- **query `older_than`** (optional; type `string`) — Cursor: return trades executed before this
  ISO8601 timestamp.
- **query `newer_than`** (optional; type `string`) — Cursor: return trades executed after this
  ISO8601 timestamp.

**Response payload**

- `200`: `application/json` → `Futures Trades` — Futures Trades; fields: `data` (array<object>);
  data item `object` fields: `crossed` (boolean), `executed_at` (string (date-time)), `nbbo_ask`
  (string), `nbbo_bid` (string), `price` (string), `side` (string), `size` (integer), `trade_id`
  (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### GEX / Greeks — 11

#### `GET /api/stock/{ticker}/gex-levels` — GEX Levels

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.gex_levels)

**What it does:**

The key gamma-exposure (GEX) price levels for a ticker on a given market date, derived from
per-strike net gamma exposure (`net = call_gex + put_gex`) evaluated relative to spot — the same
calculation used on the unusualwhales.com greek-exposure page:

- `call_wall` — strike above spot with the largest positive net gamma (resistance)
- `put_wall` — strike below spot with the largest positive net gamma (support)
- `gamma_magnet` — strike with the largest-magnitude net gamma (the strongest pin)
- `gamma_flip` — price where cumulative net gamma crosses zero, linearly interpolated between
  strikes and taken at the crossing nearest to spot (the zero-gamma level)

Any field may be `null` when there is no data for the date (or, for `gamma_flip`, no zero-crossing).
The spot-relative levels are also `null` when no spot price is available for the date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `GEX Levels` — GEX Levels; fields: `call_wall` (string), `gamma_flip`
  (string), `gamma_magnet` (string), `put_wall` (string)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/greek-exposure` — Greek Exposure

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.greek_exposure)

**What it does:**

Greek Exposure is the assumed greek exposure that market makers are exposed to. The most popular
greek exposure is gamma exposure (GEX). Investors and large funds lower risk and protect their money
by selling calls and buying puts. Market makers provide the liquidity to facilitate these trades.
GEX assumes that market makers are part of every transaction and that the bulk of their transactions
are buying calls and selling puts to investors hedging their portfolios. If a market maker has one
contract open with a gamma value of 0.05, then that market maker is exposed to 0.05 * [100 shares]
of gamma. The total market maker exposure is calculated by summing up the exposure values of all
open contracts determined by the daily open interest. Market makers profit from the bid-ask spreads
and as such, they constantly gamma hedge (they buy and sell shares to keep their positions delta
neutral). Long call positions are positive gamma - as the stock price increases and delta rises
(approaches 1), market makers hedge by selling shares, and they buy shares if the stock price
decreases and delta falls. Short put positions are negative gamma - as the stock price increases and
delta falls (approaches -1), market makers hedge by buying shares, and they sell shares if the stock
price decreases and delta rises. As such, in the event of large positive gamma, volatility is
suppressed as market makers will hedge by buying as the stock price decreases and selling as the
stock price increases … [see the official operation documentation for full notes]

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `timeframe`** (optional; type `Time frame`; default=1Y; example=2M) — The timeframe of the
  data to return. Can be one of the following formats:
  - YTD
  - 1D, 2D, etc.
  - 1W, 2W, etc.
  - 1M, 2M, etc.
  - 1Y, 2Y, etc.

**Response payload**

- `200`: `application/json` → `Greek Exposure` — Greek Exposure; fields: `call_charm` (Gex Call
  Charm), `call_delta` (Gex Call Delta), `call_gamma` (Gex Call Gamma), `call_vanna` (Gex Call
  Vanna), `date` (Market General Trading day), `put_charm` (Gex Put Charm), `put_delta` (Gex Put
  Delta), `put_gamma` (Gex Put Gamma), `put_vanna` (Gex Put Vanna)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/greek-exposure/expiry` — Greek Exposure By Expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.greek_exposure_by_expiry)

**What it does:**

The greek exposure of a ticker grouped by expiry dates across all contracts on a given market date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Greek Exposure By Expiry` — Greek Exposure By Expiry; fields:
  `call_charm` (Gex Call Charm), `call_delta` (Gex Call Delta), `call_gex` (Gex Call Gamma),
  `call_vanna` (Gex Call Vanna), `expiry` (Expiry), `put_charm` (Gex Put Charm), `put_delta` (Gex
  Put Delta), `put_gex` (Gex Put Gamma), `put_vanna` (Gex Put Vanna)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/greek-exposure/strike` — Greek Exposure By Strike

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.greek_exposure_by_strike)

**What it does:**

The greek exposure of a ticker grouped by strike price across all contracts on a given market date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Greek Exposure By Strike` — Greek Exposure By Strike; fields:
  `call_charm` (Gex Call Charm), `call_delta` (Gex Call Delta), `call_gex` (Gex Call Gamma),
  `call_vanna` (Gex Call Vanna), `put_charm` (Gex Put Charm), `put_delta` (Gex Put Delta), `put_gex`
  (Gex Put Gamma), `put_vanna` (Gex Put Vanna), `strike` (Strike)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/greek-exposure/strike-expiry` — Greek Exposure By Strike And Expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.greek_exposure_by_strike_expiry)

**What it does:**

The greek exposure of a ticker grouped by strike price for a specific expiry date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `expiry`** (required; type `Single expiry date`; example=2024-02-02) — A single expiry
  date in ISO date format.

**Response payload**

- `200`: `application/json` → `Greek Exposure By Strike And Expiry` — Greek Exposure By Strike And
  Expiry; fields: `call_charm` (Gex Call Charm), `call_delta` (Gex Call Delta), `call_gex` (Gex Call
  Gamma), `call_vanna` (Gex Call Vanna), `expiry` (Expiry), `put_charm` (Gex Put Charm), `put_delta`
  (Gex Put Delta), `put_gex` (Gex Put Gamma), `put_vanna` (Gex Put Vanna), `strike` (Strike)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/greek-flow` — Greek flow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.greek_flow)

**What it does:**

Returns the tickers greek flow (delta & vega flow) for the given market day broken down per minute.
Date must be the current or a past date. If no date is given, returns data for the current/last
market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Greek Flow` — Greek Flow; fields: `dir_delta_flow` (Dir Delta Flow),
  `dir_vega_flow` (Dir Vega Flow), `otm_dir_delta_flow` (OTM Dir Delta Flow), `otm_dir_vega_flow`
  (OTM Dir Vega Flow), `otm_total_delta_flow` (OTM Total Delta Flow), `otm_total_vega_flow` (OTM
  Total Vega Flow), `ticker` (Stock Ticker), `timestamp` (Timestamp), `total_delta_flow` (Total
  Delta Flow), `total_vega_flow` (Total Vega Flow), `transactions` (Transactions), `volume` (Volume)

#### `GET /api/stock/{ticker}/greek-flow/{expiry}` — Greek flow by expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.greek_flow_expiry)

**What it does:**

Returns the tickers greek flow (delta & vega flow) for the given market day broken down per minute &
expiry. Date must be the current or a past date. If no date is given, returns data for the
current/last market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **path `expiry`** (required; type `Single expiry date`; example=2024-02-02) — A single expiry date
  in ISO date format.

**Response payload**

- `200`: `application/json` → `Greek Flow Expiry` — Greek Flow Expiry; fields: `dir_delta_flow` (Dir
  Delta Flow), `dir_vega_flow` (Dir Vega Flow), `expiry` (Option Contract Expiry),
  `otm_dir_delta_flow` (OTM Dir Delta Flow), `otm_dir_vega_flow` (OTM Dir Vega Flow),
  `otm_total_delta_flow` (OTM Total Delta Flow), `otm_total_vega_flow` (OTM Total Vega Flow),
  `ticker` (Stock Ticker), `timestamp` (Timestamp), `total_delta_flow` (Total Delta Flow),
  `total_vega_flow` (Total Vega Flow), `transactions` (Transactions), `volume` (Volume)

#### `GET /api/stock/{ticker}/spot-exposures/strike` — Spot GEX exposures by strike

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.spot_exposures_by_strike)

**What it does:**

Returns the most recent spot GEX exposures across all strikes for the given ticker on a given date.
Calculated either with open interest or with volume. Spot GEX is the assumed $ value of the given
greek (ie. gamma) exposure that market makers need to hedge per 1% change of the underlying stock's
price movement. A positive value is long and a negative value is short. Investors and large funds
lower risk and protect their money by selling calls and buying puts. Market makers provide the
liquidity to facilitate these trades. GEX assumes that market makers are part of every transaction
and that the bulk of their transactions are buying calls and selling puts to investors hedging their
portfolios. If a market maker has one contract open with a gamma value of 0.05, then if the
underlying stock price moves by 1%, that market maker is exposed to $[0.05 * 100 shares * 0.01 *
stock price * underlying parameter of the greek variable (for gamma this variable is the stock
price)]. The total market maker spot exposure is calculated by summing up the spot exposure of all
open contracts determined by the daily open interest or by volume. Market makers profit from the
bid-ask spreads and as such, they constantly gamma hedge (they buy and sell shares to keep their
positions delta neutral). Long call positions are positive gamma - as the stock price increases and
delta rises (approaches 1), market makers hedge by selling shares, and they buy shares if the stock
price decreases and delta falls … [see the official operation documentation for full notes]

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `min_strike`** (optional; type `Min Strike`; minimum=0; example=120.5) — The minimum
  strike. Min: 0.
- **query `max_strike`** (optional; type `Max Strike`; minimum=0; example=1200) — The maximum
  strike. Min: 0.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Spot greek exposures by strike` — Spot greek exposures by strike;
  fields: `call_charm_ask` (Spot Call Charm Ask Side Exposure), `call_charm_bid` (Spot Call Charm
  Bid Side Exposure), `call_charm_oi` (Spot Call Charm Exposure), `call_charm_vol` (Spot Call Charm
  Exposure), `call_delta_ask` (Spot Call Delta Ask Side Exposure), `call_delta_bid` (Spot Call Delta
  Bid Side Exposure), `call_delta_oi` (Spot Call Delta Exposure), `call_delta_vol` (Spot Call Delta
  Exposure), `call_gamma_ask` (Spot Call Gamma Ask Side Exposure), `call_gamma_bid` (Spot Call Gamma
  Bid Side Exposure), `call_gamma_oi` (Spot Call Gamma Exposure), `call_gamma_vol` (Spot Call Gamma
  Exposure), `call_vanna_ask` (Spot Call Vanna Ask Side Exposure), `call_vanna_bid` (Spot Call Vanna
  Bid Side Exposure), `call_vanna_oi` (Spot Call Vanna Exposure), `call_vanna_vol` (Spot Call Vanna
  Exposure), `price` (Gex Underlying Price), `put_charm_ask` (Spot Put Charm Ask Side Exposure),
  `put_charm_bid` (Spot Put Charm Bid Side Exposure), `put_charm_oi` (Spot Put Charm Exposure),
  `put_charm_vol` (Spot Put Charm Exposure), `put_delta_ask` (Spot Put Delta Ask Side Exposure),
  `put_delta_bid` (Spot Put Delta Bid Side Exposure), `put_delta_oi` (Spot Put Delta Exposure),
  `put_delta_vol` (Spot Put Delta Exposure), `put_gamma_ask` (Spot Put Gamma Ask Side Exposure),
  `put_gamma_bid` (Spot Put Gamma Bid Side Exposure), `put_gamma_oi` (Spot Put Gamma Exposure),
  `put_gamma_vol` (Spot Put Gamma Exposure), `put_vanna_ask` (Spot Put Vanna Ask Side Exposure),
  `put_vanna_bid` (Spot Put Charm Bid Side Exposure), `put_vanna_oi` (Spot Put Vanna Exposure),
  `put_vanna_vol` (Spot Put Vanna Exposure), `strike` (Strike), `time` (Gex Calculation Time)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/spot-exposures/expiry-strike` — Spot GEX exposures by strike & expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.spot_exposures_by_strike_expiry_v2)

**What it does:**

Returns the most recent spot GEX exposures across all strikes for the given ticker & expiration on a
given date. Calculated either with open interest or with volume. Data is available since 2025-01-16.
Click here for the spot docs

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `expirations[]`** (required; type `Expiry dates`; example=[2024-02-02, 2024-01-26]) — An
  array of 1 or more expiry dates.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `min_strike`** (optional; type `Min Strike`; minimum=0; example=120.5) — The minimum
  strike. Min: 0.
- **query `max_strike`** (optional; type `Max Strike`; minimum=0; example=1200) — The maximum
  strike. Min: 0.
- **query `min_dte`** (optional; type `Min DTE`; minimum=0; example=1) — The minimum days to expiry.
  Min: 0.
- **query `max_dte`** (optional; type `Max DTE`; minimum=0; example=3) — The maximum days to expiry.
  Min: 0.

**Response payload**

- `200`: `application/json` → `Spot greek exposures by strike` — Spot greek exposures by strike;
  fields: `call_charm_ask` (Spot Call Charm Ask Side Exposure), `call_charm_bid` (Spot Call Charm
  Bid Side Exposure), `call_charm_oi` (Spot Call Charm Exposure), `call_charm_vol` (Spot Call Charm
  Exposure), `call_delta_ask` (Spot Call Delta Ask Side Exposure), `call_delta_bid` (Spot Call Delta
  Bid Side Exposure), `call_delta_oi` (Spot Call Delta Exposure), `call_delta_vol` (Spot Call Delta
  Exposure), `call_gamma_ask` (Spot Call Gamma Ask Side Exposure), `call_gamma_bid` (Spot Call Gamma
  Bid Side Exposure), `call_gamma_oi` (Spot Call Gamma Exposure), `call_gamma_vol` (Spot Call Gamma
  Exposure), `call_vanna_ask` (Spot Call Vanna Ask Side Exposure), `call_vanna_bid` (Spot Call Vanna
  Bid Side Exposure), `call_vanna_oi` (Spot Call Vanna Exposure), `call_vanna_vol` (Spot Call Vanna
  Exposure), `price` (Gex Underlying Price), `put_charm_ask` (Spot Put Charm Ask Side Exposure),
  `put_charm_bid` (Spot Put Charm Bid Side Exposure), `put_charm_oi` (Spot Put Charm Exposure),
  `put_charm_vol` (Spot Put Charm Exposure), `put_delta_ask` (Spot Put Delta Ask Side Exposure),
  `put_delta_bid` (Spot Put Delta Bid Side Exposure), `put_delta_oi` (Spot Put Delta Exposure),
  `put_delta_vol` (Spot Put Delta Exposure), `put_gamma_ask` (Spot Put Gamma Ask Side Exposure),
  `put_gamma_bid` (Spot Put Gamma Bid Side Exposure), `put_gamma_oi` (Spot Put Gamma Exposure),
  `put_gamma_vol` (Spot Put Gamma Exposure), `put_vanna_ask` (Spot Put Vanna Ask Side Exposure),
  `put_vanna_bid` (Spot Put Charm Bid Side Exposure), `put_vanna_oi` (Spot Put Vanna Exposure),
  `put_vanna_vol` (Spot Put Vanna Exposure), `strike` (Strike), `time` (Gex Calculation Time)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/spot-exposures/{expiry}/strike` — Spot GEX exposures by strike & expiry (Deprecated)

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.spot_exposures_by_strike_expiry)

**What it does:**

This endpoint has been deprecated and will be removed, please migrate to the new endpoint

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **path `expiry`** (required; type `Single expiry date`; example=2024-02-02) — A single expiry date
  in ISO date format.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `min_strike`** (optional; type `Min Strike`; minimum=0; example=120.5) — The minimum
  strike. Min: 0.
- **query `max_strike`** (optional; type `Max Strike`; minimum=0; example=1200) — The maximum
  strike. Min: 0.

**Response payload**

- `200`: `application/json` → `Spot greek exposures by strike` — Spot greek exposures by strike;
  fields: `call_charm_ask` (Spot Call Charm Ask Side Exposure), `call_charm_bid` (Spot Call Charm
  Bid Side Exposure), `call_charm_oi` (Spot Call Charm Exposure), `call_charm_vol` (Spot Call Charm
  Exposure), `call_delta_ask` (Spot Call Delta Ask Side Exposure), `call_delta_bid` (Spot Call Delta
  Bid Side Exposure), `call_delta_oi` (Spot Call Delta Exposure), `call_delta_vol` (Spot Call Delta
  Exposure), `call_gamma_ask` (Spot Call Gamma Ask Side Exposure), `call_gamma_bid` (Spot Call Gamma
  Bid Side Exposure), `call_gamma_oi` (Spot Call Gamma Exposure), `call_gamma_vol` (Spot Call Gamma
  Exposure), `call_vanna_ask` (Spot Call Vanna Ask Side Exposure), `call_vanna_bid` (Spot Call Vanna
  Bid Side Exposure), `call_vanna_oi` (Spot Call Vanna Exposure), `call_vanna_vol` (Spot Call Vanna
  Exposure), `price` (Gex Underlying Price), `put_charm_ask` (Spot Put Charm Ask Side Exposure),
  `put_charm_bid` (Spot Put Charm Bid Side Exposure), `put_charm_oi` (Spot Put Charm Exposure),
  `put_charm_vol` (Spot Put Charm Exposure), `put_delta_ask` (Spot Put Delta Ask Side Exposure),
  `put_delta_bid` (Spot Put Delta Bid Side Exposure), `put_delta_oi` (Spot Put Delta Exposure),
  `put_delta_vol` (Spot Put Delta Exposure), `put_gamma_ask` (Spot Put Gamma Ask Side Exposure),
  `put_gamma_bid` (Spot Put Gamma Bid Side Exposure), `put_gamma_oi` (Spot Put Gamma Exposure),
  `put_gamma_vol` (Spot Put Gamma Exposure), `put_vanna_ask` (Spot Put Vanna Ask Side Exposure),
  `put_vanna_bid` (Spot Put Charm Bid Side Exposure), `put_vanna_oi` (Spot Put Vanna Exposure),
  `put_vanna_vol` (Spot Put Vanna Exposure), `strike` (Strike), `time` (Gex Calculation Time)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/spot-exposures` — Spot GEX exposures per 1min

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.spot_exposures_one_minute)

**What it does:**

Returns the spot GEX exposures for the given ticker per minute. Spot GEX is the assumed $ value of
the given greek (ie. gamma) exposure that market makers need to hedge per 1% change of the
underlying stock's price movement. A positive value is long and a negative value is short. Investors
and large funds lower risk and protect their money by selling calls and buying puts. Market makers
provide the liquidity to facilitate these trades. GEX assumes that market makers are part of every
transaction and that the bulk of their transactions are buying calls and selling puts to investors
hedging their portfolios. If a market maker has one contract open with a gamma value of 0.05, then
if the underlying stock price moves by 1%, that market maker is exposed to $[0.05 * 100 shares *
0.01 * stock price * underlying parameter of the greek variable (for gamma this variable is the
stock price)]. The total market maker spot exposure is calculated by summing up the spot exposure of
all open contracts determined by the daily open interest or by volume. Market makers profit from the
bid-ask spreads and as such, they constantly gamma hedge (they buy and sell shares to keep their
positions delta neutral). Long call positions are positive gamma - as the stock price increases and
delta rises (approaches 1), market makers hedge by selling shares, and they buy shares if the stock
price decreases and delta falls … [see the official operation documentation for full notes]

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Spot GEX exposures per 1min` — Spot GEX exposures per 1min; fields:
  `charm_per_one_percent_move_dir` (CEX Per One Percent Move Directionalized Volume),
  `charm_per_one_percent_move_oi` (CEX Per One Percent Move OI), `charm_per_one_percent_move_vol`
  (CEX Per One Percent Move Volume), `gamma_per_one_percent_move_dir` (GEX Per One Percent Move
  Directionalized Volume), `gamma_per_one_percent_move_oi` (Gex Gamma Per One Percent Move OI),
  `gamma_per_one_percent_move_vol` (Gex Gamma Per One Percent Move Volume), `price` (Gex Underlying
  Price), `time` (Gex Calculation Time), `vanna_per_one_percent_move_dir` (VEX Per One Percent Move
  Directionalized Volume), `vanna_per_one_percent_move_oi` (VEX Per One Percent Move OI),
  `vanna_per_one_percent_move_vol` (VEX Per One Percent Move Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Alerts — 5

#### `GET /api/alerts/configuration` — Alert configurations

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AlertsController.configs)

**What it does:**

Returnst all alert configurations of the user. Users can create alerts for:

- Market tide
- Gamma exposure (GEX), Vanna exposure (VEX), Charm exposure (CEX)
- Interval Contract screeners (replicates and alerts on the Flow Feed)
- Analyst ratings, price targets, and actions
- Politician trades
- Insider trades
- OI changes for contract in premarket
- FDA
- Flow alerts
- Contract screener (replicates and alerts on the Hottest Chains)
- Stock screeners
- News
- Earnings
- Dividends
- Splits
- Trading stats (halts, unhalts)
- Economic release
- SEC filings

The alerts are the same alerts as the user created on https://unusualwhales.com/custom-alerts

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Alert config` — Alert config; fields: `config` (The configuration of
  an alert.), `created_at` (The creation timestamp), `id` (The alert id), `name` (The name of the
  alert), `noti_type` (The type of the alert), `status` (The status of the alert)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/alerts/filters` — Alert filters

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AlertsController.filters)

**What it does:**

Returns the available filters for creating alert configurations:

- `data`: per alert type (`noti_type`) the filter fields that can be used in the `config` object,
  each with its accepted values

- `access`: per alert type whether your account can create alerts of that type
- `noti_type_descriptions`: a description of each alert type
- `filter_descriptions`: per alert type a description of each filter field
- `rate_limits`: per alert type the maximum amount of alerts that trigger per day

Use this endpoint to determine which alert types the account can access and to build a structured
`config` object. Submit the resulting `noti_type` and `config` to `POST /api/alerts/configuration`
to create an alert. For more complex filters, fetch the Query language syntax and fields from `GET
/api/alerts/query/grammar`, then submit the Query expression to `POST /api/alerts/configuration`
using `input` instead of `config`. The `config` and `input` fields are mutually exclusive.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Alert filters` — Alert filters; fields: `access` (object), `data`
  (object), `filter_descriptions` (object), `noti_type_descriptions` (object), `rate_limits`
  (object) — The alert filters
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/alerts` — Alerts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AlertsController.alerts)

**What it does:**

Returns all the alerts that have been triggered for the user. Time filtering is available using the
`newer_than` and `older_than` parameters:

- The maximum lookback period is 14 days
- If no time range is specified, defaults to the last 14 days
- If only one time parameter is provided, the other is automatically calculated to maintain the
  14-day limit

- If both parameters are provided but exceed 14 days, the range is adjusted to 14 days from the
  `older_than` timestamp The alerts are the same alerts as the user created on
  https://unusualwhales.com/custom-alerts

**Parameters**

- **query `limit`** (optional; type `Default 50, Max 500 Min 1`; default=50; minimum=1; maximum=500;
  example=10) — How many items to return. Default: 50. Max: 500. Min: 1.
- **query `intraday_only`** (optional; type
  `DynLimitArqsjqhmnxasqdbwlhxpggrbxbzrywrqimvqctrexlomhtnfepj`; default=true; example=true) —
  Boolean flag whether to return only intraday alerts.
- **query `config_ids[]`** (optional; type `Configration IDs`;
  example=[ded2bee3-68fe-4aff-8316-f6ddbcf7ea67, c170e1cb-38c8-4cc9-91a0-e9b249ef7e00]) — A list of
  alert ids to filter by.
- **query `ticker_symbols`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list
  of tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `noti_types[]`** (optional; type `Noti Types`; example=[analyst_rating,
  option_contract_interval]) — A list of notification types.
- **query `newer_than`** (optional; type `string (datetime)`; format=datetime;
  example=2024-01-01T00:00:00+00:00) — Filter alerts created after this timestamp. When used alone,
  the time range is limited to 14 days forward from this timestamp. The format can either be a
  timestamp in unix seconds or milliseconds or a string in ISO format.
- **query `older_than`** (optional; type `string (datetime)`; format=datetime;
  example=2024-01-15T23:59:59+00:00) — Filter alerts created before this timestamp. When used alone,
  the time range is limited to 14 days backward from this timestamp. The format can either be a
  timestamp in unix seconds or milliseconds or a string in ISO format

**Response payload**

- `200`: `application/json` → `Alert` — Alert; fields: `created_at` (The creation timestamp), `id`
  (The alert id), `meta` (The raw data of the alert), `name` (The name of the alert), `noti_type`
  (The type of the alert), `symbol` (The stock ticker or option contract of the alert), `tape_time`
  (Alert Tape Time), `user_noti_config_id` (The alert id)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `POST /api/alerts/configuration` — Create alert configuration

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AlertsController.create_config)

**What it does:**

Creates a new alert configuration, or updates an existing one when `id` is given. The alert filters
can be given in one of two ways, but never both:

- `config`: a structured filter object for the given `noti_type`, the same configuration as returned
  by the alert configurations endpoint.

- `input`: a Query language expression. Use the Query language for complex filter expressions.
  Before creating an alert fetch the available filters:

- `GET /api/alerts/filters` lists per alert type (`noti_type`) the filter fields that can be used in
  `config`, their accepted values, and whether your account can create alerts of that type.

- `GET /api/alerts/query/grammar` returns the Query language syntax. `GET
  /api/alerts/query/grammar?target=option_trade` the full field reference for one target. Example
  with a structured `config`: Example with a Query language `input`: The created alerts are
  delivered like the alerts created on https://unusualwhales.com/custom-alerts and their triggers
  can be fetched from the alerts endpoint or streamed over the `custom_alerts` websocket channel.
  The alerts and usage of the websocket are the easiest way to setup and get notified on anything
  that matches your filter configs. The query language is a powerful tool to create complex
  configurations in a SQL like language.

**Parameters**

- None documented.

**Request body**

- `application/json` → `Create alert config request` — Create alert config request; fields: `config`
  (The configuration of an alert.), `id` (string (uuid)), `input` (string), `name` (The name of the
  alert), `noti_type` (The type of the alert), `status` (The status of the alert); sample keys:
  `input`, `name`, `noti_type`

**Response payload**

- `200`: `application/json` → `Created alert config` — Created alert config; fields: `data` (object)
  — The created or updated alert configuration
- `400`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Invalid alert configuration
- `404`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — User or alert configuration not found
- `500`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Internal Server Error

#### `GET /api/alerts/query/grammar` — Query grammar

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AlertsController.dsl_grammar)

**What it does:**

Returns the grammar of the alert Query language: available targets, syntax, fields, functions,
operators, scopes and example expressions. Without a `target` parameter an overview over all targets
is returned. Pass a `target` to get the full field reference for that alert type. Use the Query
language to create complex alert configurations via the create alert configuration endpoint. You can
compare fields `strike > spot`, add complex math operations `ln(strike / spot) / (iv * sqrt(dte /
365)) > 2` and write configs similar as to one writes SQL. If you are an Ai/Agent you should use
this endpoint to fetch the grammar and then post them to `/api/alerts/configuration`. See the
documentation at the
https://api.unusualwhales.com/docs/operations/PublicApi.AlertsController.create_config for more info
on the format and how to post an alert config.

**Parameters**

- **query `target`** (optional; type `string`; enum=[flow_alert, option_trade, option_contract,
  ticker_interval_flow, news, trading_state, multi_leg_trade]) — The alert target to return the
  detailed grammar for.

**Response payload**

- `200`: `application/json` → `Alert DSL grammar` — Alert DSL grammar; fields: Alert DSL grammar —
  The Query language grammar
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Commodities — 1

#### `GET /api/commodities/{name}` — Commodity Series

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CommoditiesController.show)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Returns a long-running
price series for a commodity. Available names: wti, brent, natural-gas, copper, aluminum, wheat,
corn, cotton, sugar, coffee, all-commodities.

**Parameters**

- **path `name`** (required; type `string`; enum=[wti, brent, natural-gas, copper, aluminum, wheat,
  corn, cotton, …])
- **query `interval`** (optional; type `string`; default=monthly; enum=[daily, weekly, monthly,
  quarterly, annual]) — Series cadence. Defaults to monthly.

**Response payload**

- `200`: `application/json` → `Macro Series Response` — Macro Series Response; fields: `data` (Macro
  Series) — Commodity Series
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Companies — 5

#### `GET /api/companies/{ticker}/dividends` — Company Dividends

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CompaniesController.dividends)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Historical dividend
events for a ticker (ex-date, declaration, record, payment, amount).

**Parameters**

- **path `ticker`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Dividends Response` — Dividends Response; fields: `data` (object) —
  Dividends
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/companies/{ticker}/profile` — Company Profile

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CompaniesController.profile)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Returns a normalized
company profile (sector, industry, description, market cap, P/E, dividend yield, etc.).

**Parameters**

- **path `ticker`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Company Profile Response` — Company Profile Response; fields: `data`
  (Company Profile) — Company Profile
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/companies/{ticker}/splits` — Company Stock Splits

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CompaniesController.splits)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Historical stock split
events for a ticker.

**Parameters**

- **path `ticker`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Splits Response` — Splits Response; fields: `data` (object) — Splits
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/companies/{ticker}/transcripts/{quarter}` — Earnings Call Transcript

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CompaniesController.transcript)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Full earnings-call
transcript for a ticker and quarter (e.g. 2024Q1). Returns speakers, sentiment, and statements.

**Parameters**

- **path `ticker`** (required; type `string`)
- **path `quarter`** (required; type `string`; pattern=^[0-9]{4}Q[1-4]$) — Year and quarter, e.g.
  "2024Q1".

**Response payload**

- `200`: `application/json` → `Transcript Response` — Transcript Response; fields: `data` (object) —
  Transcript
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/companies/{ticker}/earnings-estimates` — Forward Earnings Estimates

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CompaniesController.earnings_estimates)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Analyst-driven forward
earnings estimates by quarter/year.

**Parameters**

- **path `ticker`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Earnings Estimates Response` — Earnings Estimates Response; fields:
  `data` (object) — Earnings Estimates
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Congress — 4

#### `GET /api/congress/politicians` — List of Politicians with Trade Data

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CongressController.congress_politicians)

**What it does:**

Returns a distinct list of politicians for which trade data exists. Use `last_traded_within_months`
to filter to only recently-active politicians (e.g. 13 for last 13 months). Each entry includes
trade count, first/last trade date, party, chamber, and gender.

**Parameters**

- **query `last_traded_within_months`** (optional; type `integer`; minimum=1; maximum=240)

**Response payload**

- `200`: `application/json` → `Senate Stock` — Senate Stock; fields: `amounts` (Insider Trades
  Amount Range), `filed_at_date` (Insider Trades Filing Date), `is_active` (Is Active), `issuer`
  (Insider Trades Issuer), `member_type` (Insider Trades Member Type), `name` (Insider Trades
  Reporter's Standard Name), `notes` (Insider Trades Filing Notes), `politician_id` (Politician ID),
  `reporter` (Insider Trades Reporter), `ticker` (Stock Ticker), `transaction_date` (Insider Trades
  Transaction Date), `txn_type` (Insider Trades Transaction Type)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/congress/recent-trades` — Recent Congress Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CongressController.congress_recent_trades)

**What it does:**

Returns the latest transacted trades by congress members. If a date is given, will only return
reports, which's transaction date is <= the given input date.

**Parameters**

- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `ticker`** (optional; type `Optional Ticker`; example=IOVA) — Optional ticker symbol to
  filter results

**Response payload**

- `200`: `application/json` → `Senate Stock` — Senate Stock; fields: `amounts` (Insider Trades
  Amount Range), `filed_at_date` (Insider Trades Filing Date), `is_active` (Is Active), `issuer`
  (Insider Trades Issuer), `member_type` (Insider Trades Member Type), `name` (Insider Trades
  Reporter's Standard Name), `notes` (Insider Trades Filing Notes), `politician_id` (Politician ID),
  `reporter` (Insider Trades Reporter), `ticker` (Stock Ticker), `transaction_date` (Insider Trades
  Transaction Date), `txn_type` (Insider Trades Transaction Type)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/congress/late-reports` — Recent Late Reports

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CongressController.congress_late_reports)

**What it does:**

Returns the recent late reports by congress members. If a date is given, will only return recent
late reports, which's report date is <= the given input date.

**Parameters**

- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `ticker`** (optional; type `Optional Ticker`; example=IOVA) — Optional ticker symbol to
  filter results

**Response payload**

- `200`: `application/json` → `Senate Stock` — Senate Stock; fields: `amounts` (Insider Trades
  Amount Range), `filed_at_date` (Insider Trades Filing Date), `is_active` (Is Active), `issuer`
  (Insider Trades Issuer), `member_type` (Insider Trades Member Type), `name` (Insider Trades
  Reporter's Standard Name), `notes` (Insider Trades Filing Notes), `politician_id` (Politician ID),
  `reporter` (Insider Trades Reporter), `ticker` (Stock Ticker), `transaction_date` (Insider Trades
  Transaction Date), `txn_type` (Insider Trades Transaction Type)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/congress/congress-trader` — Recent Reports By Trader

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CongressController.congress_trader)

**What it does:**

Returns the recent reports by the given congress member. Supports pagination via `page` (1-indexed)
and `limit` parameters. Date range filtering: use `date` for upper bound and `date_from` for lower
bound. To query a single day, set both to the same date (e.g.
`?date_from=2025-06-15&date=2025-06-15`).

**Parameters**

- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `ticker`** (optional; type `Optional Ticker`; example=IOVA) — Optional ticker symbol to
  filter results
- **query `name`** (optional; type `Congress Member`; default=Nancy Pelosi; example=Adam Kinzinger)
  — The full name of a congress member. Cannot contain digits/numbers. Spaces and characters may
  need to be URI encoded, e.g. Adam Kinzinger -> Adam%20Kinzinger.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `date_from`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date
  in the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Senate Stock` — Senate Stock; fields: `amounts` (Insider Trades
  Amount Range), `filed_at_date` (Insider Trades Filing Date), `is_active` (Is Active), `issuer`
  (Insider Trades Issuer), `member_type` (Insider Trades Member Type), `name` (Insider Trades
  Reporter's Standard Name), `notes` (Insider Trades Filing Notes), `politician_id` (Politician ID),
  `reporter` (Insider Trades Reporter), `ticker` (Stock Ticker), `transaction_date` (Insider Trades
  Transaction Date), `txn_type` (Insider Trades Transaction Type)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Crypto — 4

#### `GET /api/crypto/{pair}/ohlc/{candle_size}` — Crypto OHLC Candles

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CryptoController.ohlc)

**What it does:**

Returns OHLC candle data for a crypto pair. Available candle sizes: 1m, 5m, 10m, 15m, 30m, 1h, 4h,
1d, 1w

**Parameters**

- **path `pair`** (required; type `string`)
- **path `candle_size`** (required; type `string`; enum=[1m, 5m, 10m, 15m, 30m, 1h, 4h, 1d, …])
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Crypto OHLC Response` — Crypto OHLC Response; fields: `data`
  (array<Crypto OHLC Candle>); data item `Crypto OHLC Candle` fields: `close` (string), `high`
  (string), `low` (string), `open` (string), `pair` (string), `start_time` (string (date-time)),
  `timestamp` (string (date-time)), `volume` (string) — Crypto OHLC Candles
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/crypto/{pair}/state` — Crypto Pair State

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CryptoController.state)

**What it does:**

Returns the current state for a crypto pair including 24h OHLCV data.

**Parameters**

- **path `pair`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Crypto Pair State Response` — Crypto Pair State Response; fields:
  `data` (one of Crypto Pair State, null) — Crypto Pair State
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/crypto/whale-transactions` — Crypto Whale Transactions

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CryptoController.whale_transactions)

**What it does:**

Returns recent whale transactions.

**Parameters**

- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.

**Response payload**

- `200`: `application/json` → `Crypto Whale Transactions Response` — Crypto Whale Transactions
  Response; fields: `data` (array<Crypto Whale Transaction>); data item `Crypto Whale Transaction`
  fields: `amount_formatted` (string), `blockchain` (string), `from_address` (string), `id`
  (integer), `timestamp` (string (date-time)), `to_address` (string), `token_name` (string),
  `token_symbol` (string), `transaction_hash` (string), `transaction_type` (string), `usd_value`
  (string), `whale_score` (string) — Crypto Whale Transactions
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/crypto/whales/recent` — Recent Crypto Whale Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.CryptoController.whales_recent)

**What it does:**

Returns recent large crypto trades (whale trades) across all pairs.

**Parameters**

- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.

**Response payload**

- `200`: `application/json` → `Crypto Recent Whales Response` — Crypto Recent Whales Response;
  fields: `data` (array<Crypto Whale Trade>); data item `Crypto Whale Trade` fields: `exchange`
  (string), `id` (string (uuid)), `pair` (string), `price` (string), `side` (string), `size`
  (string), `symbol` (string), `total` (string), `whaled_at` (string (date-time)) — Recent Crypto
  Whale Trades
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Dark pool — 3

#### `GET /api/darkpool/{ticker}/price-levels` — Darkpool Price Levels

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.DarkpoolController.darkpool_price_levels)

**What it does:**

Returns rounded darkpool and regular stock volume concentration by price level for one ticker. Each
price is grouped into a bucket to highlight areas where trading activity is concentrated. Date must
be the current or a past date. If no date is given, returns data for the current/last market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Darkpool Price Levels` — Darkpool Price Levels; fields: `data`
  (array<Darkpool Price Level>), `date` (Market General Trading day); data item `Darkpool Price
  Level` fields: `dark_pool_volume` (integer), `price` (Stock Price Level), `regular_volume`
  (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/darkpool/recent` — Recent Darkpool Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.DarkpoolController.darkpool_recent)

**What it does:**

Returns the latest darkpool trades. For real time streaming of darkpool trades, subscribe to the
`off_lit_trades` websocket channel, see
https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.off_lit_trades.

**Parameters**

- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `min_premium`** (optional; type `StockTradesMinPremium`; default=0; minimum=0;
  example=50000) — The minimum premium requested trades should have.
- **query `max_premium`** (optional; type `StockTradesMaxPremium`; example=150000) — The maximum
  premium requested trades should have.
- **query `min_size`** (optional; type `StockTradesMinSize`; default=0; minimum=0; example=50000) —
  The minimum size requested trades should have. Must be a positive integer.
- **query `max_size`** (optional; type `StockTradesMaxSize`; example=150000) — The maximum size
  requested trades should have. Must be a positive integer.
- **query `min_volume`** (optional; type `StockTradesMinVol`; default=0; minimum=0; example=50000) —
  The minimum consolidated volume requested trades should have. Must be a positive integer.
- **query `max_volume`** (optional; type `StockTradesMaxVol`; example=150000) — The maximum
  consolidated volume requested trades should have. Must be a positive integer.
- **query `order`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc]; example=asc) —
  Whether to sort descending or ascending. Descending by default.
- **query `order_by`** (optional; type `Darkpool order by field`; default=executed_at;
  enum=[executed_at, trf_executed_at, premium, size, volume]; example=premium) — The field to order
  the darkpool trades by. Defaults to executed_at.

**Response payload**

- `200`: `application/json` → `Darkpool Trade` — Darkpool Trade; fields: `canceled` (Single Trade Is
  Trade Cancelled), `executed_at` (Single Trade Execution Time), `ext_hour_sold_codes` (Single Trade
  External Hour Sold Code), `market_center` (Single Trade Market Center), `nbbo_ask` (ToBeDone),
  `nbbo_ask_quantity` (ToBeDone), `nbbo_bid` (ToBeDone), `nbbo_bid_quantity` (ToBeDone), `premium`
  (Option Contract Premium), `price` (Single Trade Price), `sale_cond_codes` (Single Trade Sale Cond
  Code), `size` (Single Trade Size), `ticker` (Stock Ticker), `tracking_id` (Single Trade Tracking
  ID), `trade_code` (Single Trade Trade Code), `trade_settlement` (Single Trade Settlement),
  `trf_executed_at` (Single Trade TRF Execution Time), `volume` (Stock Day Stock Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/darkpool/{ticker}` — Ticker Darkpool Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.DarkpoolController.darkpool_ticker)

**What it does:**

Returns the darkpool trades for the given ticker on a given day. Date must be the current or a past
date. If no date is given, returns data for the current/last market day. For real time streaming of
darkpool trades, subscribe to the `off_lit_trades` websocket channel, see
https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.off_lit_trades.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `newer_than`** (optional; type `NewerThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no older results will be returned. Can be used with `older_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `older_than`** (optional; type `OlderThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no newer results will be returned. Can be used with `newer_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `min_premium`** (optional; type `StockTradesMinPremium`; default=0; minimum=0;
  example=50000) — The minimum premium requested trades should have.
- **query `max_premium`** (optional; type `StockTradesMaxPremium`; example=150000) — The maximum
  premium requested trades should have.
- **query `min_size`** (optional; type `StockTradesMinSize`; default=0; minimum=0; example=50000) —
  The minimum size requested trades should have. Must be a positive integer.
- **query `max_size`** (optional; type `StockTradesMaxSize`; example=150000) — The maximum size
  requested trades should have. Must be a positive integer.
- **query `min_volume`** (optional; type `StockTradesMinVol`; default=0; minimum=0; example=50000) —
  The minimum consolidated volume requested trades should have. Must be a positive integer.
- **query `max_volume`** (optional; type `StockTradesMaxVol`; example=150000) — The maximum
  consolidated volume requested trades should have. Must be a positive integer.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `order`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc]; example=asc) —
  Whether to sort descending or ascending. Descending by default.
- **query `order_by`** (optional; type `Darkpool order by field`; default=executed_at;
  enum=[executed_at, trf_executed_at, premium, size, volume]; example=premium) — The field to order
  the darkpool trades by. Defaults to executed_at.

**Response payload**

- `200`: `application/json` → `Darkpool Trade` — Darkpool Trade; fields: `canceled` (Single Trade Is
  Trade Cancelled), `executed_at` (Single Trade Execution Time), `ext_hour_sold_codes` (Single Trade
  External Hour Sold Code), `market_center` (Single Trade Market Center), `nbbo_ask` (ToBeDone),
  `nbbo_ask_quantity` (ToBeDone), `nbbo_bid` (ToBeDone), `nbbo_bid_quantity` (ToBeDone), `premium`
  (Option Contract Premium), `price` (Single Trade Price), `sale_cond_codes` (Single Trade Sale Cond
  Code), `size` (Single Trade Size), `ticker` (Stock Ticker), `tracking_id` (Single Trade Tracking
  ID), `trade_code` (Single Trade Trade Code), `trade_settlement` (Single Trade Settlement),
  `trf_executed_at` (Single Trade TRF Execution Time), `volume` (Stock Day Stock Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Digital currencies — 2

#### `GET /api/digital-currencies/history` — Digital Currency Historical Series

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.DigitalCurrenciesController.history)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Daily, weekly, or monthly
OHLC bars for a digital asset.

**Parameters**

- **query `symbol`** (required; type `string`)
- **query `market`** (required; type `string`)
- **query `interval`** (optional; type `string`; default=daily; enum=[daily, weekly, monthly])

**Response payload**

- `200`: `application/json` → `Series Response` — Series Response; fields: `data` (Crypto Series) —
  Digital Currency History
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/digital-currencies/intraday` — Digital Currency Intraday Series

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.DigitalCurrenciesController.intraday)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Intraday OHLC bars for a
digital asset against a fiat market.

**Parameters**

- **query `symbol`** (required; type `string`; example=BTC)
- **query `market`** (required; type `string`; example=USD)
- **query `interval`** (optional; type `string`; default=5min; enum=[1min, 5min, 15min, 30min,
  60min])

**Response payload**

- `200`: `application/json` → `Series Response` — Series Response; fields: `data` (Crypto Series) —
  Digital Currency Intraday
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Earnings — 3

#### `GET /api/earnings/afterhours` — Afterhours

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EarningsController.afterhours)

**What it does:**

Returns the afterhours earnings for a given date. If you are looking to scan for extreme IV term
steepness before earnings use https://unusualwhales.com/skills/uw-earnings-vol-scan-skill.md

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `limit`** (optional; type `Default 50 Max 100 Min 1`; default=50; minimum=1; maximum=100;
  example=10) — How many items to return. Default: 50. Max: 100. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Earnings` — Earnings; fields: `actual_eps` (Actual EPS), `continent`
  (Continent), `country_code` (Country Code), `country_name` (Country Name), `ending_fiscal_quarter`
  (General ISO Date), `expected_move` (Expected Move), `expected_move_perc` (Expected Move Perc),
  `full_name` (Full Name), `has_options` (Has Options), `is_s_p_500` (Is Part of the S&P 500),
  `marketcap` (Marketcap), `post_earnings_close` (Stock Close Price), `post_earnings_date` (General
  ISO Date), `pre_earnings_close` (Stock Close Price), `pre_earnings_date` (General ISO Date),
  `reaction` (Reaction), `report_date` (General ISO Date), `report_time` (Report Time), `sector`
  (Sector), `source` (Source), `street_mean_est` (Street Mean Est), `symbol` (Ticker Symbol)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/earnings/{ticker}` — Historical Ticker Earnings

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EarningsController.ticker)

**What it does:**

Returns the historical earnings for the given ticker. The returned data includes information about
how the stock performed before and after the earnings event in its history. Furthermore, data about
how a long and short straddle would have performed in the past are included as well. If you are
looking to scan for extreme IV term steepness before earnings use
https://unusualwhales.com/skills/uw-earnings-vol-scan-skill.md

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Historical Ticker Earnings` — Historical Ticker Earnings; fields:
  `actual_eps` (Actual EPS), `ending_fiscal_quarter` (General ISO Date), `expected_move` (Expected
  Move), `expected_move_perc` (Expected Move Perc), `long_straddle_1d` (Long Straddle 1 Day),
  `long_straddle_1w` (Long Straddle 1 Week), `post_earnings_move_1d` (Post Earnings Move 1 Day),
  `post_earnings_move_1w` (Post Earnings Move 1 Week), `post_earnings_move_2w` (Post Earnings Move 2
  Week), `post_earnings_move_3d` (Post Earnings Move 13Day), `pre_earnings_move_1d` (Pre Earnings
  Move 1 Day), `pre_earnings_move_1w` (Pre Earnings Move 1 Week), `pre_earnings_move_2w` (Pre
  Earnings Move 2 Week), `pre_earnings_move_3d` (Pre Earnings Move 13Day), `report_date` (General
  ISO Date), `report_time` (Report Time), `short_straddle_1d` (Short Straddle 1 Day),
  `short_straddle_1w` (Short Straddle 1 Week), `source` (Source), `street_mean_est` (Street Mean
  Est)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/earnings/premarket` — Premarket

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EarningsController.premarket)

**What it does:**

Returns the premarket earnings for a given date. If you are looking to scan for extreme IV term
steepness before earnings use https://unusualwhales.com/skills/uw-earnings-vol-scan-skill.md

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `limit`** (optional; type `Default 50 Max 100 Min 1`; default=50; minimum=1; maximum=100;
  example=10) — How many items to return. Default: 50. Max: 100. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Earnings` — Earnings; fields: `actual_eps` (Actual EPS), `continent`
  (Continent), `country_code` (Country Code), `country_name` (Country Name), `ending_fiscal_quarter`
  (General ISO Date), `expected_move` (Expected Move), `expected_move_perc` (Expected Move Perc),
  `full_name` (Full Name), `has_options` (Has Options), `is_s_p_500` (Is Part of the S&P 500),
  `marketcap` (Marketcap), `post_earnings_close` (Stock Close Price), `post_earnings_date` (General
  ISO Date), `pre_earnings_close` (Stock Close Price), `pre_earnings_date` (General ISO Date),
  `reaction` (Reaction), `report_date` (General ISO Date), `report_time` (Report Time), `sector`
  (Sector), `source` (Source), `street_mean_est` (Street Mean Est), `symbol` (Ticker Symbol)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Economy — 1

#### `GET /api/economy/{indicator}` — Economic Indicator

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EconomyController.show)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Returns a long-running US
economic indicator series. Indicators: gdp, gdp-per-capita, treasury-yield, fed-funds, cpi,
inflation, retail-sales, durables, unemployment, payrolls. The optional `interval` and `maturity`
parameters apply only to indicators that support them (e.g. `treasury-yield` accepts a maturity like
`10year`).

**Parameters**

- **path `indicator`** (required; type `string`; enum=[gdp, gdp-per-capita, treasury-yield,
  fed-funds, cpi, inflation, retail-sales, durables, …])
- **query `interval`** (optional; type `string`; enum=[daily, weekly, monthly, quarterly, annual,
  semiannual])
- **query `maturity`** (optional; type `string`; enum=[3month, 2year, 5year, 7year, 10year, 30year])
  — Only for treasury-yield.

**Response payload**

- `200`: `application/json` → `Macro Series Response` — Macro Series Response; fields: `data` (Macro
  Series) — Economic Series
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### ETFs — 5

#### `GET /api/etfs/{ticker}/exposure` — Exposure

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EtfController.exposure)

**What it does:**

Returns all ETFs in which the given ticker is a holding

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Exposure` — Exposure; fields: `etf` (Etf Ticker), `full_name` (Etf
  Full Name), `last_price` (Candle Close), `prev_price` (Stock Prev Close Price), `shares` (Etf
  Shares), `weight` (Etf Weight)

#### `GET /api/etfs/{ticker}/holdings` — Holdings

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EtfController.holdings)

**What it does:**

Returns the holdings of the ETF

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Holdings` — Holdings; fields: `avg30_volume` (Stock Average 30 Day
  Volume), `bearish_premium` (Market General Bearish Premium), `bullish_premium` (Market General
  Bullish Premium), `call_premium` (Market General Call Premium), `call_volume` (Market General Call
  Volume), `close` (Candle Close), `has_options` (Stock Has Options), `high` (Candle High), `low`
  (Candle Low), `name` (Stock Company Name), `open` (Candle Open), `prev_price` (Stock Prev Close
  Price), `put_premium` (Market General Put Premium), `put_volume` (Market General Put Volume),
  `sector` (Market General Sector), `ticker` (Stock Ticker), `volume` (Stock Day Stock Volume),
  `week_52_high` (Stock Week 52 High), `week_52_low` (Stock Week 52 low)

#### `GET /api/etfs/{ticker}/in-outflow` — Inflow & Outflow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EtfController.in_outflow)

**What it does:**

Returns an ETF's inflow and outflow. Defaults to the last year of data. Maximum date range is 3
years.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `start_date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading
  date in the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `end_date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date
  in the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Inflow & Outflow` — Inflow & Outflow; fields: `change` (Etf Share
  Change), `change_prem` (Etf Premium Change), `close` (Stock Close Price), `date` (General ISO
  Date), `is_fomc` (Is FOMC Date), `volume` (Stock Day Stock Volume)

#### `GET /api/etfs/{ticker}/info` — Information

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EtfController.info)

**What it does:**

Returns the information about the given ETF ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Etf Info` — Etf Info; fields: `aum` (Etf AUM), `avg30_volume` (Stock
  Average 30 Day Volume), `call_vol` (Market General Call Volume), `description` (Etf Description),
  `domicile` (Etf Domicile), `etf_company` (Etf Company), `expense_ratio` (Etf Expense Ratio),
  `has_options` (Stock Has Options), `holdings_count` (Etf Holdings Count), `inception_date` (Etf
  Inception Date), `name` (Etf Full Name), `opt_vol` (Etf Total Volume), `put_vol` (Market General
  Put Volume), `stock_vol` (Stock Day Stock Volume), `website` (Etf Website Url)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/etfs/{ticker}/weights` — Sector & Country weights

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.EtfController.weights)

**What it does:**

Returns the sector & country weights for the given ETF ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Country & Sector exposure` — Country & Sector exposure; fields:
  `country` (Etf Countries), `sector` (Etf Sectors)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Forex — 3

#### `GET /api/forex/history` — FX Historical Series

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ForexController.history)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Daily, weekly, or monthly
OHLC bars for a currency pair.

**Parameters**

- **query `from`** (required; type `string`)
- **query `to`** (required; type `string`)
- **query `interval`** (optional; type `string`; default=daily; enum=[daily, weekly, monthly])

**Response payload**

- `200`: `application/json` → `Series Response` — Series Response; fields: `data` (Crypto Series) —
  FX History
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/forex/intraday` — FX Intraday Series

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ForexController.intraday)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Intraday OHLC bars for a
currency pair.

**Parameters**

- **query `from`** (required; type `string`)
- **query `to`** (required; type `string`)
- **query `interval`** (optional; type `string`; default=5min; enum=[1min, 5min, 15min, 30min,
  60min])

**Response payload**

- `200`: `application/json` → `Series Response` — Series Response; fields: `data` (Crypto Series) —
  FX Intraday
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/forex/rate` — FX Spot Rate

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ForexController.rate)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Realtime spot exchange
rate between two currencies.

**Parameters**

- **query `from`** (required; type `string`; example=USD) — From currency code (ISO 4217).
- **query `to`** (required; type `string`; example=EUR) — To currency code (ISO 4217).

**Response payload**

- `200`: `application/json` → `FX Rate Response` — FX Rate Response; fields: `data` (FX Spot Rate) —
  FX Spot Rate
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Group flow — 2

#### `GET /api/group-flow/{flow_group}/greek-flow` — Greek flow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.GroupFlowController.greek_flow)

**What it does:**

Returns the group flow's greek flow (delta & vega flow) for the given market day broken down per
minute. Date must be the current or a past date. If no date is given, returns data for the
current/last market day.

**Parameters**

- **path `flow_group`** (required; type `SingleFlowGroup`; enum=[airline, bank, basic materials,
  china, communication services, consumer cyclical, consumer defensive, crypto, …]; example=airline)
  — A flow group
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Group Flow` — Group Flow; fields: `dir_delta_flow` (Dir Delta Flow),
  `dir_vega_flow` (Dir Vega Flow), `flow_group` (Flow Group), `net_call_premium` (Market General Net
  Call Premium), `net_call_volume` (Market General Net Call Volume), `net_put_premium` (Market
  General Net Put Premium), `net_put_volume` (Market General Net Put Volume), `otm_dir_delta_flow`
  (OTM Dir Delta Flow), `otm_dir_vega_flow` (OTM Dir Vega Flow), `otm_total_delta_flow` (OTM Total
  Delta Flow), `otm_total_vega_flow` (OTM Total Vega Flow), `timestamp` (Timestamp),
  `total_delta_flow` (Total Delta Flow), `total_vega_flow` (Total Vega Flow), `transactions`
  (Transactions), `volume` (Volume)

#### `GET /api/group-flow/{flow_group}/greek-flow/{expiry}` — Greek flow by expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.GroupFlowController.greek_flow_expiry)

**What it does:**

Returns the group flow's greek flow (delta & vega flow) for the given market day broken down per
minute & expiry. Date must be the current or a past date. If no date is given, returns data for the
current/last market day.

**Parameters**

- **path `flow_group`** (required; type `SingleFlowGroup`; enum=[airline, bank, basic materials,
  china, communication services, consumer cyclical, consumer defensive, crypto, …]; example=airline)
  — A flow group
- **path `expiry`** (required; type `Single expiry date`; example=2024-02-02) — A single expiry date
  in ISO date format.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Group Flow Expiry` — Group Flow Expiry; fields: `dir_delta_flow` (Dir
  Delta Flow), `dir_vega_flow` (Dir Vega Flow), `expiry` (Option Contract Expiry), `flow_group`
  (Flow Group), `net_call_premium` (Market General Net Call Premium), `net_call_volume` (Market
  General Net Call Volume), `net_put_premium` (Market General Net Put Premium), `net_put_volume`
  (Market General Net Put Volume), `otm_dir_delta_flow` (OTM Dir Delta Flow), `otm_dir_vega_flow`
  (OTM Dir Vega Flow), `otm_total_delta_flow` (OTM Total Delta Flow), `otm_total_vega_flow` (OTM
  Total Vega Flow), `timestamp` (Timestamp), `total_delta_flow` (Total Delta Flow),
  `total_vega_flow` (Total Vega Flow), `transactions` (Transactions), `volume` (Volume)

### Insider trading — 4

#### `GET /api/insider/{ticker}` — Insiders

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InsiderController.insiders)

**What it does:**

Returns all insiders for the given ticker

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Insider` — Insider; fields: `cik` (Insider CIK), `id` (Insider ID),
  `is_person` (Insider Is Person), `logo_url` (Insider Logo URL), `name` (Insider Name), `name_slug`
  (Insider Name Slug), `social_links` (Insider Social Links), `ticker` (Stock Ticker)

#### `GET /api/insider/{sector}/sector-flow` — Sector Flow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InsiderController.sector_flow)

**What it does:**

Returns an aggregated view of the insider flow for the given sector. This can be used to quickly
examine the buy & sell insider flow for a given trading date

**Parameters**

- **path `sector`** (required; type `Sector`; enum=[Basic Materials, Communication Services,
  Consumer Cyclical, Consumer Defensive, Energy, Financial Services, Healthcare, Industrials, …];
  example=Technology) — A financial sector.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Insider Sector Flow` — Insider Sector Flow; fields: `avg_price`
  (Insider Average Price), `buy_sell` (Transaction Side), `date` (Market General Trading day),
  `has_more` (Has More), `premium` (Insider Premium), `transactions` (Insider Transactions),
  `uniq_insiders` (Insider Unique Insiders), `volume` (Insider Volume)

#### `GET /api/insider/{ticker}/ticker-flow` — Ticker Flow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InsiderController.ticker_flow)

**What it does:**

Returns an aggregated view of the insider flow for the given ticker. This can be used to quickly
examine the buy & sell insider flow for a given trading date

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Insider Ticker Flow` — Insider Ticker Flow; fields: `avg_price`
  (Insider Average Price), `buy_sell` (Transaction Side), `date` (Market General Trading day),
  `has_more` (Has More), `premium` (Insider Premium), `transactions` (Insider Transactions),
  `uniq_insiders` (Insider Unique Insiders), `volume` (Insider Volume)

#### `GET /api/insider/transactions` — Transactions

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InsiderController.transactions)

**What it does:**

Returns the latest insider transactions. By default all transacations that have been filled by the
same person on the same day with the same trade code are aggregated into a single row. Each of those
aggregated rows will contain a field `ids` which contains the ids of the single transactions that
were aggregated as well as the amount of transactions that were aggregated. If you want to disable
this behaviour you can set the `group` parameter to false to receive the single transacations as
they have been filled.

**Parameters**

- **query `ticker_symbol`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `min_value`** (optional; type `string`) — Minimum transaction value in dollars
- **query `max_value`** (optional; type `string`) — Maximum transaction value in dollars
- **query `min_price`** (optional; type `string`) — Minimum stock price at the time of transaction
- **query `max_price`** (optional; type `string`) — Maximum stock price at the time of transaction
- **query `owner_name`** (optional; type `string`) — Name of the insider who made the transaction
- **query `sectors[]`** (optional; type `Sectors`; example=[Consumer Cyclical, Technology,
  Utilities]) — Filter by company sector(s)
- **query `industries[]`** (optional; type `Industries`; example=[Semiconductors, Software -
  Infrastructure]) — Filter by company industry or industries
- **query `min_marketcap`** (optional; type `Min Marketcap`; minimum=0; example=1000000) — The
  minimum marketcap. Min: 0.
- **query `max_marketcap`** (optional; type `Max Marketcap`; minimum=0; example=250000000) — The
  maximum marketcap. Min: 0.
- **query `market_cap_size`** (optional; type `Single market cap size`; enum=[micro, small, mid,
  large, big]; example=large) — Size category of company market cap
- **query `min_earnings_dte`** (optional; type `Min Earnings DTE`; example=5) — The minimum days
  until the next earnings report.
- **query `max_earnings_dte`** (optional; type `Max Earnings DTE`; example=30) — The maximum days
  until the next earnings report.
- **query `min_amount`** (optional; type `string`) — Minimum number of shares in transaction
- **query `max_amount`** (optional; type `string`) — Maximum number of shares in transaction
- **query `is_director`** (optional; type `boolean`) — Filter transactions by company directors
- **query `is_officer`** (optional; type `boolean`) — Filter transactions by company officers
- **query `is_s_p_500`** (optional; type `boolean`) — Only include S&P 500 companies
- **query `is_ten_percent_owner`** (optional; type `boolean`) — Filter transactions by 10% owners
- **query `common_stock_only`** (optional; type `boolean`) — Only include common stock transactions
- **query `transaction_codes[]`** (optional; type `string`) — Filter by transaction codes
  (P=Purchase, S=Sale, etc.)
- **query `security_ad_codes`** (optional; type `string`) — Filter by security acquisition
  disposition codes
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `group`** (optional; type `boolean`) — Group insider transactions with same person and
  trade code on the same day into a single row
- **query `start_date`** (optional; type `string (date)`; format=date) — Filter for transaction
  dates occurring on or after this date (ISO format: YYYY-MM-DD)
- **query `end_date`** (optional; type `string (date)`; format=date) — Filter for transaction dates
  occurring on or before this date (ISO format: YYYY-MM-DD)

**Response payload**

- `200`: `application/json` → `Insider Trade Agg` — Insider Trade Agg; fields: `amount` (Insider
  Amount), `date_excercisable` (Insider Date Excercisable), `director_indirect` (Insider Director
  Indirect), `expiration_date` (Insider Expiration Date), `filing_date` (Insider Filing Date),
  `formtype` (Insider Form Type), `ids` (Insider IDs), `is_10b5_1` (Insider Is 10b5-1),
  `is_director` (Insider Is Director), `is_officer` (Insider Is Officer), `is_ten_percent_owner`
  (Insider Is 10 Owner), `natureofownership` (Insider Nature of Ownership), `officer_title` (Insider
  Officer Title), `owner_name` (Insider Owner Name), `price` (Insider Price), `price_excercisable`
  (Insider Price Excercisable), `security_ad_code` (Insider Security AD Code), `security_title`
  (Insider Security Title), `shares_owned_after` (Insider Shares Owned After), `shares_owned_before`
  (Insider Shares Owned Before), `ticker` (Stock Ticker), `transaction_code` (Insider Transaction
  Code), `transaction_date` (Insider Transaction Date), `transactions` (Insider Transactions)

### Institutions / 13F — 7

#### `GET /api/institution/{name}/activity/v2` — Institutional Activity

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InstitutionController.activity_v2)

**What it does:**

The trading activities for a given institution. If the result contains multiple report periods, it
will contain first the latest report period results ordered by the `order` and `order_direction`,
then the second latest report period ordered by the `order` and `order_direction` and so on, e.g.
first all results for report period 2025-09-30, then all results for report period 2025-06-30 ...
Per default the results per report period will be ordered by `units` descending. If you provide the
same date for `start_date` and `end_date` you will only get the results for filings with that date
as the report period end date e.g. if you set `start_date` and `end_date` to 2025-09-30 you will
only get the results for the report period ending on 2025-09-30. WARNING: Providing partial names
e.g. "VANGUARD" might lead to unexpected results. To ensure you get the expected result, it is
recommended to use the company CIK (e.g. 0000102909 for VANGUARD GROUP INC) as the name parameter.
You can use endpoint "api/institutions" with a partial name like "VANGUARD" to get all matching
institutions with their CIKs. If you are looking to build a database of insitutional positions use
https://unusualwhales.com/skills/institutional.md

**Parameters**

- **path `name`** (required; type `Institution`; example=VANGUARD GROUP INC or 0000102909) — A large
  entity that manages funds and investments for others. Queryable by name or cik.
- **query `ticker_symbol`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `start_date`** (optional; type `Institutional Reporting Start Date`; example=2023-01-01) —
  Start date to filter institutional reporting data in the format YYYY-MM-DD. Only data with dates
  on or after this date will be returned.
- **query `end_date`** (optional; type `Institutional Reporting End Date`; example=2023-03-31) — End
  date to filter institutional reporting data in the format YYYY-MM-DD. Only data with dates on or
  before this date will be returned.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `order`** (optional; type `Institutional Activity Order By`; default=units; enum=[units,
  units_change]; example=units) — Optional columns to order the result by per report period end
  date.
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.

**Response payload**

- `200`: `application/json` → `Institutional Activity` — Institutional Activity; fields: `avg_price`
  (Avg Price), `buy_price` (Buy Price), `close` (Stock Close Price), `filing_date` (General ISO
  Date), `price_on_filing` (Price On Filing), `price_on_report` (Price On Report), `put_call`
  (PutCall), `report_date` (General ISO Date), `security_type` (Security Type), `sell_price` (Sell
  Price), `shares_outstanding` (Shares Outstanding), `ticker` (Ticker Symbol), `units` (Units),
  `units_change` (Units Change)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/institution/{name}/activity` — Institutional Activity (Deprecated)

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InstitutionController.activity)

**What it does:**

The trading activities for a given institution. WARNING: Providing partial names e.g. "VANGUARD"
might lead to unexpected results. To ensure you get the expected result, it is recommended to use
the company CIK (e.g. 0000102909 for VANGUARD GROUP INC) as the name parameter. You can use endpoint
"api/institutions" with a partial name like "VANGUARD" to get all matching institutions with their
CIKs.

**Parameters**

- **path `name`** (required; type `Institution`; example=VANGUARD GROUP INC or 0000102909) — A large
  entity that manages funds and investments for others. Queryable by name or cik.
- **query `date`** (optional; type `Institutional Reporting End Date`; example=2023-03-31) — End
  date to filter institutional reporting data in the format YYYY-MM-DD. Only data with dates on or
  before this date will be returned.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `order`** (optional; type `Institutional Activity Order By`; default=units; enum=[units,
  units_change]; example=units) — Optional columns to order the result by per report period end
  date.
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.

**Response payload**

- `200`: `application/json` → `Institutional Activity` — Institutional Activity; fields: `avg_price`
  (Avg Price), `buy_price` (Buy Price), `close` (Stock Close Price), `filing_date` (General ISO
  Date), `price_on_filing` (Price On Filing), `price_on_report` (Price On Report), `put_call`
  (PutCall), `report_date` (General ISO Date), `security_type` (Security Type), `sell_price` (Sell
  Price), `shares_outstanding` (Shares Outstanding), `ticker` (Ticker Symbol), `units` (Units),
  `units_change` (Units Change)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/institution/{name}/holdings` — Institutional Holdings

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InstitutionController.holdings)

**What it does:**

Returns the holdings for a given institution. WARNING: Providing partial names e.g. "VANGUARD" might
lead to unexpected results. To ensure you get the expected result, it is recommended to use the
company CIK (e.g. 0000102909 for VANGUARD GROUP INC) as the name parameter. You can use endpoint
"api/institutions" with a partial name like "VANGUARD" to get all matching institutions with their
CIKs. If you are looking to build a database of insitutional positions use
https://unusualwhales.com/skills/institutional.md

**Parameters**

- **path `name`** (required; type `Institution`; example=VANGUARD GROUP INC or 0000102909) — A large
  entity that manages funds and investments for others. Queryable by name or cik.
- **query `ticker_symbol`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `start_date`** (optional; type `Institutional Report Start Date`; example=2023-01-01) — A
  date in the format of YYYY-MM-DD, only institutional holdings with `report_date` values on or
  after this date will be returned.
- **query `end_date`** (optional; type `Institutional Report End Date`; example=2023-03-31) — A date
  in the format of YYYY-MM-DD, only institutional holdings with `report_date` values on or before
  this date will be returned.
- **query `security_types`** (optional; type `Security Types`; example=[Share]) — An array of
  security types
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `order`** (optional; type `Institutional Holdings Order By`; enum=[date, ticker,
  security_type, put_call, first_buy, price_first_buy, units, units_change, …]; example=ticker) —
  Optional columns to order the result by
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.

**Response payload**

- `200`: `application/json` → `An Institution's Holdings` — An Institution's Holdings; fields:
  `avg_price` (Avg Price), `close` (Stock Close Price), `date` (General ISO Date), `first_buy`
  (General ISO Date), `full_name` (Full Name), `historical_units` (Historical Units),
  `perc_of_share_value` (Perc of Share Value), `perc_of_total` (Perc of Total), `price_first_buy`
  (Price First Buy), `put_call` (PutCall), `sector` (Sector), `security_type` (Security Type),
  `shares_outstanding` (Shares Outstanding), `ticker` (Ticker Symbol), `units` (Units),
  `units_change` (Units Change), `value` (Value)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/institution/{ticker}/ownership` — Institutional Ownership

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InstitutionController.ownership)

**What it does:**

The institutional ownership of a given ticker. If you are looking to build a database of
insitutional positions use https://unusualwhales.com/skills/institutional.md

**Parameters**

- **path `ticker`** (required; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `date`** (optional; type `Report Date`; format=date; example=2024-01-18) — The report date
  in the format of YYYY-MM-DD.
- **query `start_date`** (optional; type `Institutional Report Start Date`; example=2023-01-01) — A
  date in the format of YYYY-MM-DD, only institutional holdings with `report_date` values on or
  after this date will be returned.
- **query `end_date`** (optional; type `Institutional Report End Date`; example=2023-03-31) — A date
  in the format of YYYY-MM-DD, only institutional holdings with `report_date` values on or before
  this date will be returned.
- **query `tags[]`** (optional; type `Institution Tags`; example=[activist]) — An array of
  institution tags
- **query `order`** (optional; type `Institutional Ownership Order By`; enum=[name, short_name,
  filing_date, first_buy, units, units_change, units_changed, value, …]; example=name) — Optional
  columns to order the result by
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Institutional Ownership` — Institutional Ownership; fields:
  `avg_price` (Avg Price), `filing_date` (General ISO Date), `first_buy` (General ISO Date),
  `historical_units` (Historical Units), `inst_share_value` (Share Value), `inst_value` (Total
  Value), `name` (Name), `people` (People), `report_date` (General ISO Date), `shares_outstanding`
  (Shares Outstanding), `short_name` (Short Name), `tags` (Tags), `units` (Units), `units_change`
  (Units Change), `value` (Value)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/institutions/latest_filings` — Latest Filings

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InstitutionController.latest_filings)

**What it does:**

The latest institutional filings. If you are looking to build a database of insitutional positions
use https://unusualwhales.com/skills/institutional.md

**Parameters**

- **query `name`** (optional; type `Institution`; example=VANGUARD GROUP INC or 0000102909) — A
  large entity that manages funds and investments for others. Queryable by name or cik.
- **query `date`** (optional; type `string`) — Date in format YYYY-MM-DD
- **query `order`** (optional; type `Latest Institutional Filings Order By`; enum=[name, short_name,
  cik]; example=name) — Optional columns to order the result by. The filing date is always sorted in
  descending order.
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Latest Institutional Filings` — Latest Institutional Filings; fields:
  `cik` (CIK), `filing_date` (General ISO Date), `is_hedge_fund` (Is Hedge Fund), `name` (Name),
  `people` (People), `short_name` (Short Name), `tags` (Tags)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/institutions` — List of Institutions

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InstitutionController.list)

**What it does:**

Returns a list of institutions. If you are looking to build a database of insitutional positions use
https://unusualwhales.com/skills/institutional.md

**Parameters**

- **query `name`** (optional; type `Institution`; example=VANGUARD GROUP INC or 0000102909) — A
  large entity that manages funds and investments for others. Queryable by name or cik.
- **query `min_total_value`** (optional; type `MinValue`; example=0.5) — The min value of the given
  field.
- **query `max_total_value`** (optional; type `MaxValue`; example=10.0) — The max value of the given
  field.
- **query `min_share_value`** (optional; type `MinValue`; example=0.5) — The min value of the given
  field.
- **query `max_share_value`** (optional; type `MaxValue`; example=10.0) — The max value of the given
  field.
- **query `tags[]`** (optional; type `Institution Tags`; example=[activist]) — An array of
  institution tags
- **query `order`** (optional; type `Institutional List Order By`; enum=[name, call_value,
  put_value, share_value, call_holdings, put_holdings, share_holdings, total_value, …];
  example=name) — Optional columns to order the result by
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Institution Summary` — Institution Summary; fields: `buy_value` (Buy
  Value), `call_holdings` (Call Holding Units), `call_value` (Call Value), `cik` (CIK), `date`
  (Report Period End Date), `debt_holdings` (Debt Holding Units), `debt_value` (Debt Value),
  `description` (Description), `filing_date` (Filing Date), `founder_img_url` (Founder Image URL),
  `fund_holdings` (Fund Holding Units), `fund_value` (Fund Value), `is_hedge_fund` (Is Hedge Fund),
  `logo_url` (Logo URL), `name` (Name), `people` (People), `pfd_holdings` (Preferred Share Holding
  Units), `pfd_value` (Preferred Share Value), `put_holdings` (Put Holding Units), `put_value` (Put
  Value), `sell_value` (Sell Value), `share_value` (Share Value), `short_name` (Short Name), `tags`
  (Tags), `total_value` (Total Value), `warrant_holdings` (Warrant Holding Units), `warrant_value`
  (Warrant Value), `website` (Website)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/institution/{name}/sectors` — Sector Exposure

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.InstitutionController.sectors)

**What it does:**

The sector exposure for a given institution. WARNING: Providing partial names e.g. "VANGUARD" might
lead to unexpected results. To ensure you get the expected result, it is recommended to use the
company CIK (e.g. 0000102909 for VANGUARD GROUP INC) as the name parameter. You can use endpoint
"api/institutions" with a partial name like "VANGUARD" to get all matching institutions with their
CIKs. If you are looking to build a database of insitutional positions use
https://unusualwhales.com/skills/institutional.md

**Parameters**

- **path `name`** (required; type `Institution`; example=VANGUARD GROUP INC or 0000102909) — A large
  entity that manages funds and investments for others. Queryable by name or cik.
- **query `date`** (optional; type `Institutional Report End Date`; example=2023-03-31) — A date in
  the format of YYYY-MM-DD, only institutional holdings with `report_date` values on or before this
  date will be returned.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `An Institution's Sector Exposure` — An Institution's Sector Exposure;
  fields: `positions` (Positions), `positions_closed` (Positions Closed), `positions_decreased`
  (Positions Decreased), `positions_increased` (Positions Increased), `report_date` (General ISO
  Date), `sector` (Sector), `value` (Value)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Intelligence / reference data — 5

#### `GET /api/companies/listings` — Active or Delisted Securities

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.IntelController.listings)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** All US-traded securities,
optionally filtered to delisted as of a given date.

**Parameters**

- **query `status`** (optional; type `string`; default=active; enum=[active, delisted])
- **query `date`** (optional; type `string (date)`; format=date) — Only used when status=delisted.
  ISO date.

**Response payload**

- `200`: `application/json` → `Listings Response` — Listings Response; fields: `data` (object) —
  Listings
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/analytics/window` — Fixed-Window Analytics

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.IntelController.analytics_window)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Statistical analytics
over a fixed window across one or more tickers. Returns mean, stddev, correlation, etc.

**Parameters**

- **query `symbols`** (required; type `string`; example=AAPL,IBM,MSFT) — Comma-separated tickers.
- **query `range`** (required; type `string`; example=2month) — Either an ISO date for the start of
  the window (pair with the optional range_end) or a relative shorthand like '2month' or 'full'.
- **query `range_end`** (optional; type `string`; example=2024-08-31) — Optional end date when using
  ISO range_start.
- **query `interval`** (optional; type `string`; default=DAILY; enum=[1min, 5min, 15min, 30min,
  60min, DAILY, WEEKLY, MONTHLY])
- **query `ohlc`** (optional; type `string`; default=close; enum=[open, high, low, close])
- **query `calculations`** (required; type `string`; example=MEAN,STDDEV,CORRELATION) —
  Comma-separated. Supported: MIN, MAX, MEAN, MEDIAN, CUMULATIVE_RETURN, VARIANCE, STDDEV,
  MAX_DRAWDOWN, HISTOGRAM, AUTOCORRELATION, COVARIANCE, CORRELATION.

**Response payload**

- `200`: `application/json` → `Analytics Response` — Analytics Response; fields: `data` (object) —
  Analytics Window
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/calendar/ipo` — IPO Calendar

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.IntelController.ipo_calendar)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Upcoming IPOs in the next
3 months.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `IPO Calendar Response` — IPO Calendar Response; fields: `data`
  (object) — IPO Calendar
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/analytics/sliding` — Sliding-Window Analytics

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.IntelController.analytics_sliding)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Sliding window
statistical analytics across one or more tickers.

**Parameters**

- **query `symbols`** (required; type `string`)
- **query `range`** (required; type `string`)
- **query `range_end`** (optional; type `string`)
- **query `interval`** (optional; type `string`; default=DAILY)
- **query `ohlc`** (optional; type `string`; default=close)
- **query `window_size`** (optional; type `integer`; default=20; minimum=1)
- **query `calculations`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Analytics Response` — Analytics Response; fields: `data` (object) —
  Analytics Sliding
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/movers` — Top Movers

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.IntelController.movers)

**What it does:**

**Requires Advanced+ tier (Advanced, Enterprise, or Enterprise + Kafka).** Top gainers, top losers,
and most actively traded US tickers for the latest session.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Movers Response` — Movers Response; fields: `data` (object) — Movers
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Lit flow — 2

#### `GET /api/lit-flow/recent` — Recent Lit Flow Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.LitFlowController.lit_flow_recent)

**What it does:**

Returns the latest lit exchange trades.

**Parameters**

- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `min_premium`** (optional; type `StockTradesMinPremium`; default=0; minimum=0;
  example=50000) — The minimum premium requested trades should have.
- **query `max_premium`** (optional; type `StockTradesMaxPremium`; example=150000) — The maximum
  premium requested trades should have.
- **query `min_size`** (optional; type `StockTradesMinSize`; default=0; minimum=0; example=50000) —
  The minimum size requested trades should have. Must be a positive integer.
- **query `max_size`** (optional; type `StockTradesMaxSize`; example=150000) — The maximum size
  requested trades should have. Must be a positive integer.
- **query `min_volume`** (optional; type `StockTradesMinVol`; default=0; minimum=0; example=50000) —
  The minimum consolidated volume requested trades should have. Must be a positive integer.
- **query `max_volume`** (optional; type `StockTradesMaxVol`; example=150000) — The maximum
  consolidated volume requested trades should have. Must be a positive integer.

**Response payload**

- `200`: `application/json` → `Lit Trade` — Lit Trade; fields: `canceled` (Single Trade Is Trade
  Cancelled), `executed_at` (Single Trade Execution Time), `ext_hour_sold_codes` (Single Trade
  External Hour Sold Code), `market_center` (Single Trade Market Center), `nbbo_ask` (ToBeDone),
  `nbbo_ask_quantity` (ToBeDone), `nbbo_bid` (ToBeDone), `nbbo_bid_quantity` (ToBeDone), `premium`
  (Option Contract Premium), `price` (Single Trade Price), `sale_cond_codes` (Single Trade Sale Cond
  Code), `size` (Single Trade Size), `ticker` (Stock Ticker), `tracking_id` (Single Trade Tracking
  ID), `trade_code` (Single Trade Trade Code), `trade_settlement` (Single Trade Settlement),
  `volume` (Stock Day Stock Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/lit-flow/{ticker}` — Ticker Lit Flow Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.LitFlowController.lit_flow_ticker)

**What it does:**

Returns the lit exchange trades for the given ticker on a given day. Date must be the current or a
past date. If no date is given, returns data for the current/last market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `newer_than`** (optional; type `NewerThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no older results will be returned. Can be used with `older_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `older_than`** (optional; type `OlderThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no newer results will be returned. Can be used with `newer_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `min_premium`** (optional; type `StockTradesMinPremium`; default=0; minimum=0;
  example=50000) — The minimum premium requested trades should have.
- **query `max_premium`** (optional; type `StockTradesMaxPremium`; example=150000) — The maximum
  premium requested trades should have.
- **query `min_size`** (optional; type `StockTradesMinSize`; default=0; minimum=0; example=50000) —
  The minimum size requested trades should have. Must be a positive integer.
- **query `max_size`** (optional; type `StockTradesMaxSize`; example=150000) — The maximum size
  requested trades should have. Must be a positive integer.
- **query `min_volume`** (optional; type `StockTradesMinVol`; default=0; minimum=0; example=50000) —
  The minimum consolidated volume requested trades should have. Must be a positive integer.
- **query `max_volume`** (optional; type `StockTradesMaxVol`; example=150000) — The maximum
  consolidated volume requested trades should have. Must be a positive integer.
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.

**Response payload**

- `200`: `application/json` → `Lit Trade` — Lit Trade; fields: `canceled` (Single Trade Is Trade
  Cancelled), `executed_at` (Single Trade Execution Time), `ext_hour_sold_codes` (Single Trade
  External Hour Sold Code), `market_center` (Single Trade Market Center), `nbbo_ask` (ToBeDone),
  `nbbo_ask_quantity` (ToBeDone), `nbbo_bid` (ToBeDone), `nbbo_bid_quantity` (ToBeDone), `premium`
  (Option Contract Premium), `price` (Single Trade Price), `sale_cond_codes` (Single Trade Sale Cond
  Code), `size` (Single Trade Size), `ticker` (Stock Ticker), `tracking_id` (Single Trade Tracking
  ID), `trade_code` (Single Trade Trade Code), `trade_settlement` (Single Trade Settlement),
  `volume` (Stock Day Stock Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Market analytics — 12

#### `GET /api/market/correlations` — Correlations

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.correlations)

**What it does:**

Returns the correlations between a list of tickers. Date must be the current or a past date. If no
date is given, returns data for the current/last market day. You can filter the time period either
by: 1. Using the `interval` parameter (e.g. "1y", "6m", "3m", "1m") 2. Using `start_date` and
optionally `end_date` (if `end_date` is not provided, it defaults to the current date) If you send
`interval` alongside `start_date`, `interval` filter will take priority.

**Parameters**

- **query `tickers`** (required; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `interval`** (optional; type `Time frame`; default=1Y; example=2M) — The timeframe of the
  data to return. Can be one of the following formats:
  - YTD
  - 1D, 2D, etc.
  - 1W, 2W, etc.
  - 1M, 2M, etc.
  - 1Y, 2Y, etc.
- **query `start_date`** (optional; type `Correlation Start Date`; example=2023-01-01) — A date in
  the format of YYYY-MM-DD to start the correlation calculation from. This is optional and works
  alongside end_date to define a specific time range. Can be used as an alternative to the interval
  parameter.
- **query `end_date`** (optional; type `Correlation End Date`; example=2023-12-31) — A date in the
  format of YYYY-MM-DD to end the correlation calculation. This is optional and defaults to the
  current date when start_date is provided. Works alongside start_date to define a specific time
  range. Can be used as an alternative to the interval parameter.

**Response payload**

- `200`: `application/json` → `Correlation` — Correlation; fields: `correlation` (Correlation),
  `fst` (Stock Ticker), `max_date` (General ISO Date), `min_date` (General ISO Date), `rows` (Count
  of data points), `snd` (Stock Ticker)

#### `GET /api/market/{ticker}/etf-tide` — ETF Tide

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.etf_tide)

**What it does:**

The ETF tide is similar to the Market Tide. While the market tide is based on options activity of
the whole market the ETF tide is only based on the options activity of the holdings of the specified
ETF.

**Parameters**

- **path `ticker`** (required; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Daily Market Tide` — Daily Market Tide; fields: `net_call_premium`
  (Market General Net Call Premium), `net_put_premium` (Market General Net Put Premium),
  `net_volume` (Market General Net Volume), `timestamp` (General Tick time)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/economic-calendar` — Economic calendar

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.events)

**What it does:**

Returns the economic calendar.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Economic calendar` — Economic calendar; fields: `event` (Economic
  Event), `forecast` (Economic Forecast), `prev` (Economic Previous), `reported_period` (Economic
  Reported Period), `time` (Economic Time), `type` (Economic Type) — Economic Calendar response
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/fda-calendar` — FDA Calendar

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.fda_calendar)

**What it does:**

Returns FDA calendar data with filtering options. The FDA calendar contains information about:

- PDUFA (Prescription Drug User Fee Act) dates
- Advisory Committee Meetings
- FDA Decisions
- Clinical Trial Results
- New Drug Applications
- Biologics License Applications

##### Date format support

The target_date parameters support various FDA-specific date formats:

- Quarters: YYYY-Q[1-4] (e.g. 2024-Q1)
- Half years: YYYY-H[1-2] (e.g. 2024-H1)
- Mid-year: YYYY-MID (e.g. 2024-MID)
- Late-year: YYYY-LATE (e.g. 2024-LATE)
- Standard dates: YYYY-MM-DD

**Parameters**

- **query `announced_date_min`** (optional; type `Optional Market Date`; example=2024-01-18) —
  Minimum announced date (YYYY-MM-DD)
- **query `announced_date_max`** (optional; type `Optional Market Date`; example=2024-01-18) —
  Maximum announced date (YYYY-MM-DD)
- **query `target_date_min`** (optional; type `FDA Target Date`;
  pattern=^\d{4}-(Q[1-4]|H[1-2]|MID|LATE|\d{2}-\d{2})$; example=2024-Q1) — Minimum target date
  (supports Q1-Q4, H1-H2, MID, LATE formats)
- **query `target_date_max`** (optional; type `FDA Target Date`;
  pattern=^\d{4}-(Q[1-4]|H[1-2]|MID|LATE|\d{2}-\d{2})$; example=2024-Q1) — Maximum target date
  (supports Q1-Q4, H1-H2, MID, LATE formats)
- **query `drug`** (optional; type `Drug Name`; example=Keytruda) — Filter by drug name (partial
  match)
- **query `ticker`** (optional; type `Ticker`; example=AAPL,INTC) — Filter by ticker symbol
- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — Maximum number of results to return

**Response payload**

- `200`: `application/json` → `FDA Calendar` — FDA Calendar; fields: `catalyst` (FDA Calendar
  Catalyst), `description` (FDA Calendar Description), `drug` (FDA Calendar Drug), `end_date` (FDA
  Calendar End Date), `has_options` (Has Options), `indication` (FDA Calendar Indication),
  `marketcap` (Company Market Cap), `notes` (FDA Calendar Notes), `outcome` (FDA Calendar Outcome),
  `outcome_brief` (FDA Calendar Outcome Brief), `source_link` (FDA Calendar Source Link),
  `start_date` (FDA Calendar Start Date), `status` (FDA Calendar Status), `ticker` (FDA Calendar
  Ticker) — FDA Calendar response
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/market-tide` — Market Tide

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.market_tide)

**What it does:**

Market Tide is a proprietary tool that can be viewed from the Market Overview page. The Market Tide
chart provides real time data based on a proprietary formula that examines market wide options
activity and filters out 'noise'. Date must be the current or a past date. If no date is given,
returns data for the current/last market day. Per default data are returned in 1 minute intervals.
Use `interval_5m=true` to have this return data in 5 minute intervals instead. For example:

- $15,000 in calls transacted at the ask has the effect of increasing the daily net call premium by
  $15,000.

- $10,000 in calls transacted at the bid has the effect of decreasing the daily net call premium by
  $10,000.

The resulting net premium from both of these trades would be $5000 (+ $15,000 - $10,000).
Transactions taking place at the mid are not accounted for. In theory: The sentiment in the options
market becomes increasingly bullish if: 1. The aggregated CALL PREMIUM is increasing at a faster
rate. 2. The aggregated PUT PREMIUM is decreasing at a faster rate. The sentiment in the options
market becomes increasingly bearish if: 1. The aggregated CALL PREMIUM is decreasing at a faster
rate. 2. The aggregated PUT PREMIUM is increasing at a faster rate. If you are interested in which
tickers are influencing the market tide the most you can check
https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.top_net_impact. If you are
looking for a market tide but for specific sectors check out:
https://api.unusualwhales.com/docs/operations/PublicApi.MarketController … [see the official
operation documentation for full notes]

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `otm_only`** (optional; type `UseOtmOnly`; default=false; example=true) — Only include out
  of the money transactions.
- **query `interval_5m`** (optional; type `Use5mCandles`; default=true; example=false) — Return data
  in 5 minutes intervals.

**Response payload**

- `200`: `application/json` → `Daily Market Tide` — Daily Market Tide; fields: `net_call_premium`
  (Market General Net Call Premium), `net_put_premium` (Market General Net Put Premium),
  `net_volume` (Market General Net Volume), `timestamp` (General Tick time)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/net-flow/expiry` — Net Flow Expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.NetFlowController.expiry)

**What it does:**

Returns net premium flow by `tide_type` category, `moneyness` category, and `expiration` category,
allowing you to create chart variations like https://unusualwhales.com/zero-dte: About the query
parameters:

- **`tide_type`**: For example, setting `tide_type` to "equity_only" will filter out ETFs and
  indexes, leaving only net premium from single-name equities.

- **`moneyness`**: For example, setting `moneyness` to "otm" will filter out any contract that was
  not out of the money ("OTM") at the time of the transaction, leaving only net premium from
  contracts that were OTM at the time of the transaction.

- **`expiration`**: For example, setting `expiration` to "zero_dte" will filter out any contract not
  expiring this session, leaving only net premium from contracts expiring at 4PM eastern time today.

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — Market date to get
  data for (defaults to last market day)
- **query `moneyness`** (optional; type `array<string>`) — Moneyness filter (defaults to 'all')
- **query `tide_type`** (optional; type `array<string>`) — Tide type filter (defaults to 'all')
- **query `expiration`** (optional; type `array<string>`) — Expiration filter (defaults to
  ['weekly', 'zero_dte'])

**Response payload**

- `200`: `application/json` → `Expiry Tide Response` — Expiry Tide Response; fields: `data`
  (array<Expiry Tide Group>), `date` (string (date)), `expiration` (array<string>), `moneyness`
  (array<string>), `tide_type` (array<string>); data item `Expiry Tide Group` fields: `data`
  (array<Net Flow Data Point>), `moneyness` (string), `tide_type` (string) — Net Flow Expiry
  Response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/oi-change` — OI Change

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.oi_change)

**What it does:**

Returns the non-Index/non-ETF contracts and OI (open interest) change data with the highest OI
change (default: descending). The data updates once on trading days at around 6:45am EST in the
premarket. This endpoint is heavily used to confirm how much of volume has been translated into
actual open contracts. Date must be the current or a past date. If no date is given, returns data
for the current/last market day.

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.
- **query `order`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc]; example=asc) —
  Whether to sort descending or ascending. Descending by default.

**Response payload**

- `200`: `application/json` → `OI Change` — OI Change; fields: `avg_price` (ToBeDone), `curr_date`
  (General ISO Date), `curr_oi` (ToBeDone), `last_ask` (ToBeDone), `last_bid` (ToBeDone),
  `last_date` (General ISO Date), `last_fill` (ToBeDone), `last_oi` (ToBeDone), `oi_change`
  (ToBeDone), `oi_diff_plain` (ToBeDone), `option_symbol` (Option Contract Symbol),
  `percentage_of_total` (ToBeDone), `prev_ask_volume` (ToBeDone), `prev_bid_volume` (ToBeDone),
  `prev_mid_volume` (ToBeDone), `prev_multi_leg_volume` (ToBeDone), `prev_neutral_volume`
  (ToBeDone), `prev_stock_multi_leg_volume` (ToBeDone), `prev_total_premium` (ToBeDone), `rnk`
  (ToBeDone), `trades` (ToBeDone), `underlying_symbol` (ToBeDone), `volume` (ToBeDone)

#### `GET /api/market/sector-etfs` — Sector Etfs

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.sector_etfs)

**What it does:**

Returns the current trading days statistics for the SPDR sector etfs. If you are interested in
building some sort of market dashboard and want to quickly retrieve which sectors are up/down today
this endpoint returns all required data to build such dashboard/view. ---- This can be used to build
a market overview such as:

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Sector ETF` — Sector ETF; fields: `avg30_stock_volume` (Market
  General Avg 30 Day Stock Volume), `avg_30_day_call_volume` (Market General Avg 30 Day Call
  Volume), `avg_30_day_put_volume` (Market General Avg 30 Day Put Volume), `avg_7_day_call_volume`
  (Market General Avg 7 Day Call Volume), `avg_7_day_put_volume` (Market General Avg 7 Day Put
  Volume), `bearish_premium` (Market General Bearish Premium), `bullish_premium` (Market General
  Bullish Premium), `call_premium` (Market General Call Premium), `call_volume` (Market General Call
  Volume), `close` (Candle Close), `full_name` (Etf Sector Full Name), `high` (Candle High), `low`
  (Candle Low), `marketcap` (Stock Marketcap AUM), `open` (Candle Open), `prev_close` (Stock Prev
  Close Price), `prev_date` (Market General Previous Trading Day), `put_premium` (Market General Put
  Premium), `put_volume` (Market General Put Volume), `ticker` (Stock Ticker), `volume` (Stock Day
  Stock Volume), `week_52_high` (Stock Week 52 High), `week_52_low` (Stock Week 52 low)
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/{sector}/sector-tide` — Sector Tide

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.sec_indst)

**What it does:**

The Sector tide is similar to the Market Tide. While the market tide is based on options activity of
the whole market the sector tide is only based on the options activity of companies which are in
that specific sector

**Parameters**

- **path `sector`** (required; type `Single sector`; enum=[Basic Materials, Communication Services,
  Consumer Cyclical, Consumer Defensive, Energy, Financial Services, Healthcare, Industrials, …];
  example=Real Estate) — A singular sector.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Daily Market Tide` — Daily Market Tide; fields: `net_call_premium`
  (Market General Net Call Premium), `net_put_premium` (Market General Net Put Premium),
  `net_volume` (Market General Net Volume), `timestamp` (General Tick time)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/top-net-impact` — Top Net Impact

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.top_net_impact)

**What it does:**

Returns the top tickers by net premium (half bullish, half bearish). Defaults to last market day.
These tickers are representing the tickers which had the most influence in the market tide.

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `issue_types[]`** (optional; type `Issue types`; example=[Common Stock, Index]) — An array
  of 1 or more issue types.
- **query `limit`** (optional; type `Default 20 Max 100 Min 1`; default=20; minimum=1; maximum=100;
  example=20) — How many items to return. Default: 20. Max: 100. Min: 1.

**Response payload**

- `200`: `application/json` → `Top Net Impact` — Top Net Impact; fields: `net_premium` (number
  (float)), `ticker` (Stock Ticker)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/insider-buy-sells` — Total Insider Buy & Sells

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.insider_buy_sells)

**What it does:**

Returns the total amount of purchases & sells as well as notional values for insider transactions
across the market

**Parameters**

- **query `limit`** (optional; type `Min Limit 1`; minimum=1; example=10) — How many items to
  return. If no limit is given, returns all matching data. Min: 1.

**Response payload**

- `200`: `application/json` → `Insider statistics` — Insider statistics; fields: `filing_date`
  (Insider Trades Filing Date), `purchases` (Insider Trades Purchases), `purchases_notional`
  (Insider Trades Notional Purchases), `sells` (Insider Trades Sells), `sells_notional` (Insider
  Trades Notional Sells)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/market/total-options-volume` — Total Options Volume

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MarketController.total_options_volume)

**What it does:**

Returns the total options volume and premium for all trade executions that happened on a given
trading date. ---- This can be used to build a market options overview such as:

**Parameters**

- **query `limit`** (optional; type `Default 1 Max 500 Min 1`; default=1; minimum=1; maximum=500;
  example=10) — How many items to return. Default: 1. Max: 500. Min: 1.

**Response payload**

- `200`: `application/json` → `Market Options Volume` — Market Options Volume; fields:
  `call_premium` (Market General Call Premium), `call_volume` (Market General Call Volume), `date`
  (Market General Trading day), `put_premium` (Market General Put Premium), `put_volume` (Market
  General Put Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### News — 1

#### `GET /api/news/headlines` — News Headlines

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.NewsController.headlines)

**What it does:**

Returns the latest news headlines for financial markets. This endpoint provides access to news
headlines that may impact the markets, including company-specific news, sector news, and market-wide
events. Headlines can be filtered by source, content, ticker, and significance. Use the
`search_term` parameter to filter by headline content, or the `ticker` parameter to filter by
related ticker symbols. The data includes the headline text, source, related tickers, sentiment, and
whether it's considered a major news item.

**Parameters**

- **query `sources`** (optional; type `News Sources`; example=BusinessWire,MarketNews) — A
  comma-separated list of news sources to filter by (e.g., 'Reuters,Bloomberg').
- **query `search_term`** (optional; type `News Search Term`; example=earnings) — A search term to
  filter news headlines by content.
- **query `ticker`** (optional; type `News Ticker`; example=AAPL) — Filter news headlines to only
  those mentioning the specified ticker symbol.
- **query `major_only`** (optional; type `Major News Only`; default=false; example=true) — When set
  to true, only returns major/significant news.
- **query `limit`** (optional; type `Default 50 Max 100 Min 1`; default=50; minimum=1; maximum=100;
  example=10) — How many items to return. Default: 50. Max: 100. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Headline News` — Headline News; fields: `created_at` (string
  (date_time)), `headline` (string), `is_major` (boolean), `meta` (object), `sentiment` (string),
  `source` (string), `tags` (array<string>), `tickers` (array<string>)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Option contracts — 7

#### `GET /api/stock/{ticker}/expiry-breakdown` — Expiry Breakdown

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionContractController.expiry_breakdown)

**What it does:**

Returns all expirations for the given trading day for a ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Expiry breakdown` — Expiry breakdown; fields: `chains` (integer),
  `expiry` (Option Contract Expiry), `open_interest` (integer), `volume` (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-contract/{id}/flow` — Flow Data

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionContractController.flow)

**What it does:**

Returns the last 50 option trades for the given option chain. Optionally a min premium and a side
can be supplied in the query for further filtering. If no date is specified data for the last
trading day is being returned.

**Parameters**

- **path `id`** (required; type `OptionContract`; example=TSLA230526P00167500) — An option contract
  in the OSI format.
- **query `side`** (optional; type `Side`; default=ALL; enum=[ALL, ASK, BID, MID]; example=ASK) —
  The side of a stock trade. Must be one of ASK, BID, MID. If not set, will return all side's
  trades.
- **query `min_premium`** (optional; type `StockTradesMinPremium`; default=0; minimum=0;
  example=50000) — The minimum premium requested trades should have.
- **query `limit`** (optional; type `Min Limit 1`; minimum=1; example=10) — How many items to
  return. If no limit is given, returns all matching data. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Option Trade` — Option Trade; fields: `ask_vol` (Option Contract Ask
  Volume), `bid_vol` (Option Contract Bid Volume), `canceled` (Canceled), `delta` (Delta), `er_time`
  (Stock Earnings time), `ewma_nbbo_ask` (EWMA NBBO Ask), `ewma_nbbo_bid` (EWMA NBBO Bid),
  `exchange` (Exchange), `executed_at` (Executed At), `expiry` (Option Contract Expiry),
  `flow_alert_id` (Flow Alert ID), `full_name` (Stock Full Name), `gamma` (Gamma), `id` (Option
  Trade ID), `implied_volatility` (Implied Volatility), `industry_type` (Stock Industry Type),
  `is_agg` (boolean), `issue_type` (string), `marketcap` (Stock Marketcap AUM), `mid_vol` (Option
  Contract Mid Volume), `multi_vol` (Option Contract Multi Leg Volume), `nbbo_ask` (NBBO Ask),
  `nbbo_bid` (NBBO Bid), `next_earnings_date` (Stock Next Earnings Date), `no_side_vol` (Option
  Contract No Side Volume), `open_interest` (Option Contract Open interest), `option_chain_id`
  (Option Contract Symbol), `option_type` (Option Contract Option Type), `premium` (Premium),
  `price` (Fill Price), `report_flags` (Report Flags), `rho` (Rho), `rule_id` (Rule ID), `sector`
  (Market General Sector), `size` (Option Trade Size), `stock_multi_vol` (Option Contract Stock
  Multi Leg Volume), `strike` (Option Contract Strike), `tags` (Tags), `theo` (Theoretical Price),
  `theta` (Theta), `trade_ids` (array<string (uuid)>), `underlying_price` (Underlying Price),
  `underlying_symbol` (Option Contract Underlying Symbol), `upstream_condition_detail` (Upstream
  Condition Detail), `vega` (Vega), `volume` (Option Trade Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-contract/{id}/historic` — Historic Data

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionContractController.history)

**What it does:**

Returns for every trading day historic data for the given option contract. Data includes open, high,
low, close of the contract of fills. The percentage of the volume which was part of a multi leg
trade, stock multi leg trade, sweep, floor and cross. The high and low of the implied volatility is
included as well as the volume distributed per sides: Ask, bid, mid and neutral. Neutral is volume
that is either a cross trade or from trades that came in late. You can use this endpoint to retrieve
for a given chains historical details about how much volume has been traded in the past, when the OI
did start to build and much more.

**Parameters**

- **path `id`** (required; type `OptionContract`; example=TSLA230526P00167500) — An option contract
  in the OSI format.
- **query `limit`** (optional; type `Min Limit 1`; minimum=1; example=10) — How many items to
  return. If no limit is given, returns all matching data. Min: 1.

**Response payload**

- `200`: `application/json` → `Option contract ` — Option contract ; fields: Option contract
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-contract/{id}/intraday` — Intraday Data

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionContractController.intraday)

**What it does:**

Returns 1 minute interval intraday data for the given option contract. You can use this to build
volume profile charts for contracts to see when volume can in, and which side dominated the volume.
The volume and premium is segmented into bid, ask, mid and neutral. Date must be the current or a
past date. If no date is given, returns data for the current/last market day.

**Parameters**

- **path `id`** (required; type `OptionContract`; example=TSLA230526P00167500) — An option contract
  in the OSI format.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Option Contract Minute Ticks` — Option Contract Minute Ticks; fields:
  `avg_price` (Option Contract Avg Price), `close` (Option Contract Close), `expiry` (Option
  Contract Expiry), `high` (Option Contract High), `iv_high` (Option Contract IV High), `iv_low`
  (Option Contract IV Low), `low` (Option Contract Low), `open` (Option Contract Open),
  `option_symbol` (Option Contract Symbol), `premium_ask_side` (Option Contract Premium Ask Side),
  `premium_bid_side` (Option Contract Premium Bid Side), `premium_mid_side` (Option Contract Premium
  Mid Side), `premium_no_side` (Option Contract Premium No Side), `start_time` (General UTC
  Timestamp), `volume_ask_side` (Option Contract Ask Volume), `volume_bid_side` (Option Contract Bid
  Volume), `volume_mid_side` (Option Contract Mid Volume), `volume_multi` (Option Contract Multi Leg
  Volume), `volume_no_side` (Option Contract No Side Volume), `volume_stock_multi` (Option Contract
  Stock Multi Leg Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/option-stance` — Option Trade-Stance Ranking

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionContractController.option_stance)

**What it does:**

Ranks a ticker's option contracts by how well their mechanics fit a chosen trade **stance**, across
all recent expiries and strikes. Returns a 0–5 `fit_score` plus named 0–1 sub-scores (`iv_regime`,
`greeks_fit`, `dte_fit`, `liquidity`, `earnings_timing`) and a plain-language `explanation` for each
contract, alongside a `context` header (stock price, IV rank, IV percentile). This is
**descriptive** analysis of greeks / IV context / liquidity mechanics — not trade advice or a
recommendation to buy or sell. Every response carries a `disclaimer`. `stance` is one of:
`sell_premium`, `sell_vega`, `directional`, `leaps`, `cheapies`. Any stance threshold can be
overridden by passing the corresponding screener filter (e.g. `max_dte`, `min_premium`, `type`).
Pass `option_symbol` to score a **single specific contract** for the stance instead of ranking the
chain. The response's `data` then holds exactly that one scored contract (with the same `fit_score`,
sub-scores and `explanation`), even when it falls outside the stance's usual candidate profile.

**Parameters**

- **path `ticker`** (required; type `string`) — Ticker symbol, e.g. RBLX.
- **query `stance`** (required; type `string`; enum=[sell_premium, sell_vega, directional, leaps,
  cheapies]) — Trade stance to rank by.
- **query `limit`** (optional; type `integer`) — Max contracts to return. Default 25, max 100.
- **query `type`** (optional; type `string`; enum=[Calls, Puts]) — Restrict to Calls or Puts.
- **query `option_symbol`** (optional; type `string`) — Score a single specific contract (OCC option
  symbol, e.g. NVDA270115P00275000) instead of ranking the chain. `data` returns just that one
  scored contract.
- **query `date`** (optional; type `string`) — Optional historical date (YYYY-MM-DD).

**Response payload**

- `200`: `application/json` → `Option Trade-Stance Ranking response.` — Option Trade-Stance Ranking
  response.; fields: `as_of_date` (string (date)), `context` (object), `data` (array<object>),
  `disclaimer` (string), `stance` (string), `ticker` (string); data item `object` fields:
  `avg_price` (number), `components` (object), `delta` (number), `dte` (integer),
  `earnings_before_expiry` (boolean), `expires` (string (date)), `explanation` (array<string>),
  `fit_score` (number), `gamma` (number), `iv` (number), `next_earnings_date` (string (date)),
  `open_interest` (integer), `option_symbol` (string), `option_type` (string), `pop` (number),
  `premium` (number), `rank_metric` (object), `roc` (number), `strike` (number), `theta` (number),
  `vega` (number), `volume` (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/option-contracts` — Option contracts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionContractController.option_contracts)

**What it does:**

Returns all option contracts for the given ticker

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `expiry`** (optional; type `Single expiry date`; example=2024-02-02) — A single expiry
  date in ISO date format.
- **query `option_type`** (optional; type `OptionType`; enum=[call, Call, put, Put]) — The option
  type to filter by if specified.
- **query `vol_greater_oi`** (optional; type `boolean`) — Wether to only return chains where volume
  > open interest
- **query `exclude_zero_vol_chains`** (optional; type `boolean`) — Wether to only return chains
  where volume > 0
- **query `exclude_zero_dte`** (optional; type `boolean`) — Wether to only return chains which do
  not expire on the same day
- **query `exclude_zero_oi_chains`** (optional; type `boolean`) — Wether to only return chains where
  open interest > 0
- **query `maybe_otm_only`** (optional; type `boolean`) — Wether to only return chains which are out
  of the money
- **query `min_dte`** (optional; type `integer`) — Minimum days to expiration (expiry at least this
  many days from today).
- **query `max_dte`** (optional; type `integer`) — Maximum days to expiration (expiry at most this
  many days from today).
- **query `option_symbol[]`** (optional; type `array<string>`) — Options symbols to filter by
- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `Option contracts` — Option contracts; fields: `ask_volume` (Option
  Contract Ask Volume), `avg_price` (Option Contract Avg Price), `bid_volume` (Option Contract Bid
  Volume), `cross_volume` (Option Contract Cross Volume), `delta` (string), `floor_volume` (Option
  Contract Floor Volume), `gamma` (string), `high_price` (Option Contract High),
  `implied_volatility` (Option Contract Last Transaction IV), `last_price` (Option Contract Close),
  `last_tape_time` (string (date-time)), `low_price` (Option Contract Low), `mid_volume` (Option
  Contract Mid Volume), `multi_leg_volume` (Option Contract Multi Leg Volume), `nbbo_ask` (NBBO
  Ask), `nbbo_bid` (NBBO Bid), `no_side_volume` (Option Contract No Side Volume), `open_interest`
  (Option Contract Open interest), `option_symbol` (Option Contract Symbol), `prev_oi` (Option
  Contract Previous Open Interest), `rho` (string), `stock_multi_leg_volume` (Option Contract Stock
  Multi Leg Volume), `sweep_volume` (Option Contract Sweep Volume), `theta` (string),
  `total_premium` (Option Contract Premium), `vega` (string), `volume` (Option Contract Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-contract/{id}/volume-profile` — Volume Profile

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionContractController.volume_profile)

**What it does:**

Returns the volume profile (volume - sweep, floor, cross, ask, bid, etc. - per fill price) for an
option symbol on a given trading day. Date must be the current or a past date. If no date is given,
returns data for the current/last market day.

**Parameters**

- **path `id`** (required; type `OptionContract`; example=TSLA230526P00167500) — An option contract
  in the OSI format.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Option Contract Price Level Volume` — Option Contract Price Level
  Volume; fields: `ask_vol` (Option Contract Ask Volume), `bid_vol` (Option Contract Bid Volume),
  `cross_vol` (Option Contract Cross Volume), `date` (Market General Trading day), `floor_vol`
  (Option Contract Floor Volume), `mid_vol` (Option Contract Mid Volume), `multi_vol` (Option
  Contract Multi Leg Volume), `price` (Fill Price), `sweep_vol` (Option Contract Sweep Volume),
  `transactions` (Option Contract Total Trades Count), `volume` (Option Contract Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Option trades — 8

#### `GET /api/option-trades/exchange-breakdown/{date}` — Exchange & Trade Code Breakdown

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionTradeController.exchange_breakdown)

**What it does:**

Aggregates the option tape for one or more tickers on a single trading date, grouped by the
**options exchange** the prints executed on (e.g. `AMXO`, `ARCO`, `BATO`). For each ticker ×
exchange you get the trade count, total contracts, total premium (`sum(price × size × 100)`) and the
call/put trade split. Set `by_trade_code=true` to additionally break each row down by the upstream
trade condition code (`upstream_condition_detail`, e.g. `auto`, `mlat`, `slan`). Use `min_premium`
to drop smaller prints before aggregating. **Whole-universe mode:** omit `ticker[]` to get the
breakdown across the entire option universe. In that mode tickers are ranked by total contracts
traded that day (descending) and the response is paged by ticker via `limit` (tickers per page) and
`page` — so you receive the top names first. When `ticker[]` is supplied this paging is ignored and
exactly those tickers are returned. NOTICE: Access to this endpoint is only included in the Advanced
API subscription. Data is available for trading days after 2022-01-01.

**Parameters**

- **path `date`** (required; type `Market Date`; example=2024-01-18) — A trading date in the format
  of YYYY-MM-DD.
- **query `ticker[]`** (optional; type `array<string>`) — One or more underlying symbols to
  aggregate, e.g. `ticker[]=AAPL&ticker[]=NVDA`. Omit for whole-universe mode (top names by volume).
- **query `by_trade_code`** (optional; type `boolean`) — Whether to additionally group each row by
  the upstream trade condition code.
- **query `min_premium`** (optional; type `string`) — Only include individual prints whose premium
  (price × size × 100) is at least this value, before aggregating. Use to filter out smaller trades.
- **query `limit`** (optional; type `integer`) — Whole-universe mode only: tickers per page (ranked
  by volume). Default 100, max 500.
- **query `page`** (optional; type `integer`) — Whole-universe mode only: 1-based page of tickers.
  Default 1.
- **query `order`** (optional; type `string`) — Whole-universe mode only: how to rank tickers —
  `volume` (total contracts) or `premium` (total premium). Default `volume`.

**Response payload**

- `200`: `application/json` → `Option Trade Exchange Breakdown` — Option Trade Exchange Breakdown;
  fields: `call_trades` (integer), `contracts` (integer), `exchange` (Exchange), `premium` (number),
  `put_trades` (integer), `trade_code` (Upstream Condition Detail), `trade_count` (integer),
  `underlying_symbol` (Option Contract Underlying Symbol)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-trades/flow-alerts/{id}` — Flow Alert by ID

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionTradeController.flow_alert)

**What it does:**

Returns the trades that made up the specific alert. For multi leg flow alert trades it will return
all the related multi leg trades. For alerts such as RepeatedHits it will return all the
transactions that made up that alert.

**Parameters**

- **path `id`** (required; type `string (uuid)`; format=uuid) — The UUID of the flow alert
- **query `older_than`** (optional; type `OlderThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no newer results will be returned. Can be used with `newer_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).

**Response payload**

- `200`: `application/json` → `Flow Alert Detail` — Flow Alert Detail; fields: `alert` (object),
  `has_more` (boolean), `trades` (array<object>)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-trades/flow-alerts` — Flow Alerts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionTradeController.flow_alerts)

**What it does:**

Flow alerts are rule based aggregations on the full tape of option trades. While there are quite a
few different rules and alerts the most used one is the repeated hit family: RepeatedHits,
RepeatedHitsAscendingFill, RepeatedHitsDescendingFill Each of those represent an alert when there
have been multiple transactions on the same option contract within a few milliseconds. This can be
mean that a single order is being matched across multiple other orders and creating multiple
transactions. It can also just mean that there are multiple buyers/sellers at the same time. Trades
usually use the repeated hits with other data points to form a picture on whether there is some
urgency in entering/exiting a position in a contract/ticker. The full current options tape including
trades that do not form a RepeatedHits alert can be accessed through the Option Trades endpoint. By
setting `include_agg_trades` to true in the option trades endpoint you would also retrieve the
RepeatedHits from this endpoint. The difference between the 3 repeated hits alerts are:

- DescendingFill: Each transaction that comes after another in chronological order has either the
  same fill price as or a lower fill price than the previous transaction. The last transaction must
  be lower than the first transaction.

- AscendingFill: The opposite of DescendingFill. The fill prices increase instead of decreasing.
- RepeatedHits (neither ascending nor descending): When it does not fit into one of the first two
  categories. To express ascending and descending in a mathmatical notion … [see the official
  operation documentation for full notes]

**Parameters**

- **query `ticker_symbol`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `unusual`** (optional; type `boolean`) — Convenience preset for "unusual" flow, matching
  the live options flow default criteria: volume>OI, size>OI, all-opening, OTM, single-leg, DTE≤60,
  ask-side≥50%, premium≥$10k, size≥5, issue types ADR/Common Stock/ETF. Applied as defaults, so any
  of those filters you pass explicitly (e.g. `min_ask_perc=0.9`, `max_dte=40`) overrides the preset.
- **query `min_premium`** (optional; type
  `DynLimitAfqgdudluubpqrpugfwelhpltflyvanoarwazsdoafooqmrbidd`; minimum=0; example=12500.5) — The
  minimum premium on that alert. Min: 0.
- **query `max_premium`** (optional; type
  `DynLimitAvymamouaxziogtqbsvjvfcklkzcnjmyggksbqtonpvgdyshigg`; minimum=0; example=12500.5) — The
  maximum premium on that alert. Min: 0.
- **query `min_size`** (optional; type
  `DynLimitAxmuedlqljxucfnoacgaqzzmvargkkhfytyociqmhxxghjfhymg`; minimum=0; example=125) — The
  minimum size on that alert. Size is defined as the sum of the sizes of all transactions that make
  up the alert. Min: 0.
- **query `max_size`** (optional; type
  `DynLimitAcgxbajknkuejanqrwjzefbbafvhswjfxixwyanrttnzswffqth`; minimum=0; example=125) — The
  maximum size on that alert. Min: 0.
- **query `min_volume`** (optional; type
  `DynLimitAogjgquizxkwngreqdqyvzctimomotngnacvtujqoarvmjsmrdy`; minimum=0; example=125) — The
  minimum volume on that alert's contract at the time of the alert. Min: 0.
- **query `max_volume`** (optional; type
  `DynLimitAjcizbkysoxhjzexhobiqklvycebotgihsgwofdizljfnkodenj`; minimum=0; example=125) — The
  maximum volume on that alert's contract at the time of the alert. Min: 0.
- **query `min_open_interest`** (optional; type
  `DynLimitAtxxwwrjsnoapccwhgkibrcxsrckuridhfwpdjlmpoqtcqoveoo`; minimum=0; example=125) — The
  minimum open interest on that alert's contract at the time of the alert. Min: 0.
- **query `max_open_interest`** (optional; type
  `DynLimitAvsbeekweyimuakpcbpnvdcxdsqhwgaladbvtecyvjisexhunjy`; minimum=0; example=125) — The
  maximum open interest on that alert's contract at the time of the alert. Min: 0.
- **query `all_opening`** (optional; type
  `DynLimitAdjqwktdpdnzpcouxqzngzcjzedkvecopxdfnmprpijrognfbqu`; default=true; example=true) —
  Boolean flag whether all transactions are opening transactions based on OI, Size & Volume. Since
  Flow Alerts with rule_name values of RepeatedHits, RepeatedHitsAscendingFill, and
  RepeatedHitsDescendingFill are composed of many individual transactions, it is extremely unlikely
  that the all_opening value will be true, so if you are interested in these Flow Ale … [see the
  official operation documentation for full notes]
- **query `is_floor`** (optional; type
  `DynLimitAxcxofbssusivioivfebvadqrxwjgzgozxniphlsagedhncqzjb`; default=true; example=true) —
  Boolean flag whether a transaction is from the floor.
- **query `is_sweep`** (optional; type
  `DynLimitAjwtpopfarkirghubavsjbitibgeaowseailqertgoixxmsengh`; default=true; example=true) —
  Boolean flag whether a transaction is a intermarket sweep.
- **query `is_call`** (optional; type `DynLimitAmyrgupvcjpqzhjrjojzopogasvoxiniozgzktpboukwwxjtahv`;
  default=true; example=true) — Boolean flag whether a transaction is a call.
- **query `is_put`** (optional; type `DynLimitAqaxrpaopbbzymaqxniabhdwigqqhskegzxydooluxiiqvfwhtf`;
  default=true; example=true) — Boolean flag whether a transaction is a put.
- **query `is_ask_side`** (optional; type
  `DynLimitAmymfeccgyzrbmxtsbqnbiggzyyanmfdcmyndalnvkdneccuina`; default=true; example=true) —
  Boolean flag whether a transaction is ask side.
- **query `is_bid_side`** (optional; type
  `DynLimitAsxyauppacolqxkhhupanyjwmyrykbuwjmsvnppwvpjeykykoym`; default=true; example=true) —
  Boolean flag whether a transaction is bid side.
- **query `rule_name[]`** (optional; type `Rule Name`; example=[RepeatedHits,
  RepeatedHitsAscendingFill]) — An array of 1 or more rule name.
- **query `min_diff`** (optional; type `Min Contract Diff`; example=0.53) — The minimum OTM diff of
  a contract. Given a strike price of 120 and an underlying price of 98 the diff for a call option
  would equal to: (120 - 98) / 98 = 0.2245 The diff for a put option would equal to: -1 * (120 - 98)
  / 98 = -0.2245.
- **query `max_diff`** (optional; type `Min Contract Diff`; example=0.53) — The minimum OTM diff of
  a contract. Given a strike price of 120 and an underlying price of 98 the diff for a call option
  would equal to: (120 - 98) / 98 = 0.2245 The diff for a put option would equal to: -1 * (120 - 98)
  / 98 = -0.2245.
- **query `min_volume_oi_ratio`** (optional; type `Min Volume OI Ratio`; minimum=0; example=0.32) —
  The minimum ratio of contract volume to contract open interest. If the open interest of a contract
  is zero, then this ratio is evaluated as if the open interest of the contract was one (to avoid
  divide by zero errors). For example, if you set this ratio to 10, then a contract with zero open
  interest and 7 volume will NOT be included in your results.
- **query `max_volume_oi_ratio`** (optional; type `Max Volume OI Ratio`; minimum=0; example=1.58) —
  The maximum ratio of contract volume to contract open interest. If the open interest of a contract
  is zero, then this ratio is evaluated as if the open interest of the contract was one (to avoid
  divide by zero errors). For example, if you set this ratio to 50, then a contract with zero open
  interest and 75 volume will NOT be included in your results.
- **query `is_otm`** (optional; type `Is OTM Contract`; example=true) — Only include contracts which
  are currently out of the money.
- **query `issue_types[]`** (optional; type `Issue types`; example=[Common Stock, Index]) — An array
  of 1 or more issue types.
- **query `min_dte`** (optional; type `Min DTE`; minimum=0; example=1) — The minimum days to expiry.
  Min: 0.
- **query `max_dte`** (optional; type `Max DTE`; minimum=0; example=3) — The maximum days to expiry.
  Min: 0.
- **query `min_ask_perc`** (optional; type `Flow Alerts Min Ask Percentage`; minimum=0; maximum=1;
  example=0.25) — The minimum ask percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_ask_perc`** (optional; type `Flow Alerts Max Ask Percentage`; minimum=0; maximum=1;
  example=0.75) — The maximum ask percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_bid_perc`** (optional; type `Flow Alerts Min Bid Percentage`; minimum=0; maximum=1;
  example=0.25) — The minimum bid percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_bid_perc`** (optional; type `Flow Alerts Max Bid Percentage`; minimum=0; maximum=1;
  example=0.75) — The maximum bid percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_bull_perc`** (optional; type `Flow Alerts Min Bull Percentage`; minimum=0; maximum=1;
  example=0.5) — The minimum bull percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_bull_perc`** (optional; type `Flow Alerts Max Bull Percentage`; minimum=0; maximum=1;
  example=0.9) — The maximum bull percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_bear_perc`** (optional; type `Flow Alerts Min Bear Percentage`; minimum=0; maximum=1;
  example=0.5) — The minimum bear percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_bear_perc`** (optional; type `Flow Alerts Max Bear Percentage`; minimum=0; maximum=1;
  example=0.9) — The maximum bear percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_skew`** (optional; type `Flow Alerts Min Skew`; minimum=0; maximum=1; example=0.3) —
  The minimum skew. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_skew`** (optional; type `Flow Alerts Max Skew`; minimum=0; maximum=1; example=0.7) —
  The maximum skew. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_price`** (optional; type `Flow Alerts Min Price`; minimum=0; example=10.5) — The
  minimum price of the underlying asset. Min: 0.
- **query `max_price`** (optional; type `Flow Alerts Max Price`; minimum=0; example=500.75) — The
  maximum price of the underlying asset. Min: 0.
- **query `min_iv_change`** (optional; type `Flow Alerts Min IV Change`; example=0.01) — The minimum
  IV change. Unbounded decimal proxy for percentage (e.g., 0.01 for minimum +1% change).
- **query `max_iv_change`** (optional; type `Flow Alerts Max IV Change`; example=0.05) — The maximum
  IV change. Unbounded decimal proxy for percentage (e.g., 0.05 for maximum +5% change).
- **query `min_size_vol_ratio`** (optional; type `Flow Alerts Min Size Volume Ratio`; minimum=0;
  example=1.5) — The minimum size to volume ratio. Min: 0.
- **query `max_size_vol_ratio`** (optional; type `Flow Alerts Max Size Volume Ratio`; minimum=0;
  example=10.0) — The maximum size to volume ratio. Min: 0.
- **query `min_spread`** (optional; type `Flow Alerts Min Spread`; minimum=0; example=0.05) — The
  minimum spread. Min: 0.
- **query `max_spread`** (optional; type `Flow Alerts Max Spread`; minimum=0; example=5.0) — The
  maximum spread. Min: 0.
- **query `min_marketcap`** (optional; type `Min Marketcap`; minimum=0; example=1000000) — The
  minimum marketcap. Min: 0.
- **query `max_marketcap`** (optional; type `Max Marketcap`; minimum=0; example=250000000) — The
  maximum marketcap. Min: 0.
- **query `is_multi_leg`** (optional; type `Flow Alerts Is Multi Leg`; example=true) — Boolean flag
  whether the transaction is a multi-leg transaction.
- **query `size_greater_oi`** (optional; type `Flow Alerts Size Greater Than Open Interest`;
  example=true) — Only include alerts where the size is greater than the open interest.
- **query `vol_greater_oi`** (optional; type `Flow Alerts Volume Greater Than Open Interest`;
  example=true) — Only include alerts where the volume is greater than the open interest.
- **query `min_days_between_expiry_and_earnings`** (optional; type
  `MinDaysBetweenExpiryAndEarnings`; example=1) — Minimum value of (contract_expiry_date -
  underlying_next_earnings_date) in days. Negative = contract expires BEFORE earnings; zero = same
  day; positive = AFTER earnings. Use together with `max_days_between_expiry_and_earnings` to target
  a window around the next earnings announcement … [see the official operation documentation for
  full notes]
- **query `max_days_between_expiry_and_earnings`** (optional; type
  `MaxDaysBetweenExpiryAndEarnings`; example=6) — Maximum value of (contract_expiry_date -
  underlying_next_earnings_date) in days. Negative = contract expires BEFORE earnings; zero = same
  day; positive = AFTER earnings. Use together with `min_days_between_expiry_and_earnings` to target
  a window around the next earnings announcement … [see the official operation documentation for
  full notes]
- **query `newer_than`** (optional; type `NewerThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no older results will be returned. Can be used with `older_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `older_than`** (optional; type `OlderThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no newer results will be returned. Can be used with `newer_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.

**Response payload**

- `200`: `application/json` → `Flow Alert` — Flow Alert; fields: `alert_rule` (Alert Rule Name),
  `all_opening_trades` (Option Contract All Opening Trades), `created_at` (General UTC Timestamp),
  `expiry` (Option Contract Expiry), `expiry_count` (Option Contract Expiry Count), `has_floor`
  (Option Contract Has Floor), `has_multileg` (Single Trade Has Multileg), `has_singleleg` (Single
  Trade Is Single Leg), `has_sweep` (Single Trade Is Sweep), `issue_type` (Stock Issue Type),
  `open_interest` (ToBeDone), `option_chain` (Option Contract Symbol), `price` (ToBeDone), `strike`
  (Option Contract Strike), `ticker` (ToBeDone), `total_ask_side_prem` (ToBeDone),
  `total_bid_side_prem` (ToBeDone), `total_premium` (ToBeDone), `total_size` (ToBeDone),
  `trade_count` (ToBeDone), `type` (Option Contract Type), `underlying_price` (ToBeDone), `volume`
  (ToBeDone), `volume_oi_ratio` (ToBeDone)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-trades/full-tape/{date}` — Full Tape

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionTradeController.full_tape)

**What it does:**

Download all option transactions (the "full tape") for a given trading date. To filter option
transactions for the current trading day use the Option Trades endpoint. How far back you can
download is determined by your plan's historical lookback window with data available since
2022-01-01. If you are looking to build a data lake with the full tape use this skill with your
agent: https://unusualwhales.com/skills/uw-options-data-lake-skill.md You can download the data as a
zip file using wget. For example, to download data for Fri Jul 25th, 2025, if your API key is
"abc123":

**Parameters**

- **path `date`** (required; type `Market Date`; example=2024-01-18) — A trading date in the format
  of YYYY-MM-DD.

**Response payload**

- `200`: `application/zip` → no response schema documented

#### `GET /api/option-trades/multi-leg` — Multi-Leg Option Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MultiLegController.index)

**What it does:**

A live feed of detected multi-leg option strategies — vertical spreads, iron condors, butterflies,
calendars, diagonals and more — reconstructed from their individual legs, newest first. Each row
summarizes the strategy (fill/net price, net bid/ask, net premium, net greeks, strikes, DTE range,
breakevens, max profit/loss) across all legs. Use the per-strategy `id` with the `legs` endpoint to
fetch the individual contracts. Decimal values are returned as strings. Defaults to the current
trading day when no time bounds are given. Queries are limited to a 24-hour range. If both bounds
span more than 24 hours, `newer_than` is clamped to 24 hours before `older_than`.

**Parameters**

- **query `limit`** (optional; type `integer`) — Rows per page (default 50, max 500).
- **query `offset`** (optional; type `integer`) — Rows to skip for pagination (max 500).
- **query `ticker_symbol`** (optional; type `string`) — Restrict to a single underlying ticker.
- **query `newer_than`** (optional; type `string`) — Only strategies executed at/after this UTC
  timestamp (ISO-8601). Defaults to last market open. The query range is limited to 24 hours.
- **query `older_than`** (optional; type `string`) — Only strategies executed at/before this UTC
  timestamp (ISO-8601). The query range is limited to 24 hours.
- **query `strategy`** (optional; type `array<string>`) — Filter by detected strategy name(s).
  Repeatable array param (e.g. `strategy[]=iron_condor&strategy[]=call_vertical_spread`).
- **query `exclude_other`** (optional; type `boolean`) — Exclude strategies classified as `other`
  (unrecognized structures).
- **query `direction`** (optional; type `array<string>`) — Filter by direction (long/short).
  Repeatable array param.
- **query `net_side`** (optional; type `array<string>`) — Filter by the strategy's net aggressor
  side (bid/ask/mid). Repeatable array param.
- **query `min_size`** (optional; type `integer`) — Minimum total contracts across legs.
- **query `max_size`** (optional; type `integer`) — Maximum total contracts across legs.
- **query `min_premium`** (optional; type `string`) — Minimum net premium (supports abs()).
- **query `max_premium`** (optional; type `string`) — Maximum net premium (supports abs()).
- **query `min_dte`** (optional; type `integer`) — Minimum days-to-expiry.
- **query `max_dte`** (optional; type `integer`) — Maximum days-to-expiry.
- **query `min_leg_count`** (optional; type `integer`) — Minimum number of legs.
- **query `max_leg_count`** (optional; type `integer`) — Maximum number of legs.
- **query `all_otm`** (optional; type `boolean`) — Only strategies where every leg is
  out-of-the-money.
- **query `issue_types`** (optional; type `array<string>`) — Filter by underlying issue type(s),
  e.g. Common Stock, ETF. Repeatable array param.
- **query `sectors`** (optional; type `array<string>`) — Filter by underlying sector(s). Repeatable
  array param.

**Response payload**

- `200`: `application/json` → `Multi-Leg Trades Response` — Multi-Leg Trades Response; fields:
  `data` (array<Multi-Leg Trade>); data item `Multi-Leg Trade` fields: `all_opening_legs` (boolean),
  `all_otm` (boolean), `bid_ask_spread` (string), `breakevens` (array<string>), `code` (string),
  `diff_expirations` (boolean), `diff_strikes` (boolean), `diff_types` (boolean), `direction`
  (string), `executed_at` (string), `id` (string), `ivs` (array<string>), `leg_count` (integer),
  `max_dte` (integer), `max_loss` (string), `max_profit` (string), `max_strike` (string), `min_dte`
  (integer), `min_strike` (string), `net_ask` (string), `net_bid` (string), `net_delta` (string),
  `net_premium` (string), `net_price` (string), `net_side` (string), `net_theta` (string), `size`
  (integer), `strategy` (string), `strikes` (array<string>), `ticker` (string), `total_premium`
  (string), `txns` (integer), `underlying_price` (string), `uniq_exchanges` (array<string>) —
  Multi-Leg Trades
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-trades/multi-leg/{id}/legs` — Multi-Leg Trade Legs

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.MultiLegController.legs)

**What it does:**

Returns the individual option legs that make up a single multi-leg strategy, aggregated by contract
with size-weighted greeks and prices. Pass the strategy `id` returned by the Multi-Leg Option Trades
feed.

**Parameters**

- **path `id`** (required; type `string`) — The multi-leg strategy id (UUID) from the trades feed.

**Response payload**

- `200`: `application/json` → `Multi-Leg Legs Response` — Multi-Leg Legs Response; fields: `data`
  (array<Multi-Leg Leg>); data item `Multi-Leg Leg` fields: `delta` (string), `expiry` (string),
  `nbbo_ask` (string), `nbbo_bid` (string), `open_interest` (integer), `option_symbol` (string),
  `option_symbol_id` (integer), `premium` (string), `price` (string), `side` (string), `size`
  (integer), `strike` (string), `theta` (string), `type` (string), `volume` (integer) — Multi-Leg
  Legs
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-trades` — Option Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionTradeController.index)

**What it does:**

Fitler the full option trades tape. This endpoint returns the same data and supports the same filter
params as on the unusualwhales website https://unusualwhales.com/live-options-flow. This endpoint
only returns data for the latest trading day. To retrieve historical option trades, use the
`/api/option-trades/full-tape/:date` endpoint to download the full market file for a trading day.
List parameters may be supplied using repeated bracket notation, for example
`tags[]=ask_side&tags[]=bid_side`. Unix timestamps may be supplied in seconds or milliseconds.

**Parameters**

- **query `limit`** (optional; type `Default 50, Max 500 Min 1`; default=50; minimum=1; maximum=500;
  example=10) — How many items to return. Default: 50. Max: 500. Min: 1.
- **query `ticker_symbol`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `option_contracts[]`** (optional; type `array<OptionContract>`;
  example=[AAPL250117C00200000]) — Option contracts to include.
- **query `chain[]`** (optional; type `array<OptionContract>`; example=[AAPL250117C00200000]) —
  Alias for `option_contracts[]`.
- **query `strike`** (optional; type `Strike`; example=150.0) — The strike price of an option
  contract.
- **query `type`** (optional; type `OptionType`; enum=[call, Call, put, Put]) — The option type to
  filter by if specified.
- **query `newer_than`** (optional; type `NewerThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no older results will be returned. Can be used with `older_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `older_than`** (optional; type `OlderThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no newer results will be returned. Can be used with `newer_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `canceled`** (optional; type `Canceled`; example=false) — Whether the option trade was
  canceled.
- **query `is_multi_leg`** (optional; type `Flow Alerts Is Multi Leg`; example=true) — Boolean flag
  whether the transaction is a multi-leg transaction.
- **query `volume_greater_oi`** (optional; type `Volume Greater Than Open Interest Contract`;
  example=true) — Only include contracts where the volume is greater than the open interest.
- **query `exclude_deep_itm`** (optional; type `boolean`; example=true) — Exclude deep in-the-money
  contracts.
- **query `force_15_min_delay`** (optional; type `boolean`; example=true) — Only return trades that
  are at least 15 minutes old.
- **query `hide_expired`** (optional; type `boolean`; example=true) — Exclude expired option
  contracts.
- **query `include_agg_trades`** (optional; type `boolean`; example=true) — Whether to roll up
  related option transactions executed at the same time into a single transaction in the response.
  This allows filters to apply to their combined premium and size. For example, if one $25,000 order
  is reported as ten $2,500 transactions, it will only match `min_premium=20000` when
  `include_agg_trades=true`.
- **query `intraday_only`** (optional; type `boolean`; example=true) — Only return trades from the
  current trading day.
- **query `is_otm`** (optional; type `boolean`; example=true) — Filter out-of-the-money or
  in-the-money trades.
- **query `opening`** (optional; type `boolean`; example=true) — Filter opening or non-opening
  transactions.
- **query `opex_only`** (optional; type `boolean`; example=true) — Filter contracts by whether they
  expire on monthly OpEx Friday.
- **query `size_greater_oi`** (optional; type `boolean`; example=true) — Filter by whether trade
  size exceeds open interest.
- **query `exchanges[]`** (optional; type `array<string>`; example=[AMXO, MXOP]) — Options exchanges
  to include.
- **query `excluded_tags[]`** (optional; type `array<string>`; example=[bid_side]) — Exclude trades
  containing any of these tags.
- **query `expiry_dates[]`** (optional; type `Expiry dates`; example=[2024-02-02, 2024-01-26]) — An
  array of 1 or more expiry dates.
- **query `industries[]`** (optional; type `Industries`; example=[Semiconductors, Software -
  Infrastructure]) — An array of one or more industries.
- **query `issue_types[]`** (optional; type `Issue types`; example=[Common Stock, Index]) — An array
  of 1 or more issue types.
- **query `marketcap_size[]`** (optional; type `Market cap sizes`; example=[large, big]) — An array
  of one or more market capitalization size categories.
- **query `report_flag[]`** (optional; type `array<string>`; example=[intermarket_sweep]) — Trade
  report flags to include.
- **query `sectors[]`** (optional; type `Sectors`; example=[Consumer Cyclical, Technology,
  Utilities]) — An array of 1 or more sectors.
- **query `tags[]`** (optional; type `array<string>`; example=[ask_side]) — Include trades
  containing any of these tags.
- **query `trade_codes[]`** (optional; type `array<string>`; example=[auto, slan]) — OPRA trade code
  of the executed transaction.
- **query `min_ask_perc`** (optional; type `Flow Alerts Min Ask Percentage`; minimum=0; maximum=1;
  example=0.25) — The minimum ask percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_ask_perc`** (optional; type `Flow Alerts Max Ask Percentage`; minimum=0; maximum=1;
  example=0.75) — The maximum ask percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_bear_perc`** (optional; type `Flow Alerts Min Bear Percentage`; minimum=0; maximum=1;
  example=0.5) — The minimum bear percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_bear_perc`** (optional; type `Flow Alerts Max Bear Percentage`; minimum=0; maximum=1;
  example=0.9) — The maximum bear percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_bid_perc`** (optional; type `Flow Alerts Min Bid Percentage`; minimum=0; maximum=1;
  example=0.25) — The minimum bid percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_bid_perc`** (optional; type `Flow Alerts Max Bid Percentage`; minimum=0; maximum=1;
  example=0.75) — The maximum bid percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_bull_perc`** (optional; type `Flow Alerts Min Bull Percentage`; minimum=0; maximum=1;
  example=0.5) — The minimum bull percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_bull_perc`** (optional; type `Flow Alerts Max Bull Percentage`; minimum=0; maximum=1;
  example=0.9) — The maximum bull percentage. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_skew`** (optional; type `Flow Alerts Min Skew`; minimum=0; maximum=1; example=0.3) —
  The minimum skew. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `max_skew`** (optional; type `Flow Alerts Max Skew`; minimum=0; maximum=1; example=0.7) —
  The maximum skew. Decimal proxy for percentage (0 to 1). Min: 0. Max: 1.
- **query `min_days_between_expiry_and_earnings`** (optional; type
  `MinDaysBetweenExpiryAndEarnings`; example=1) — Minimum value of (contract_expiry_date -
  underlying_next_earnings_date) in days. Negative = contract expires BEFORE earnings; zero = same
  day; positive = AFTER earnings. Use together with `max_days_between_expiry_and_earnings` to target
  a window around the next earnings announcement … [see the official operation documentation for
  full notes]
- **query `max_days_between_expiry_and_earnings`** (optional; type
  `MaxDaysBetweenExpiryAndEarnings`; example=6) — Maximum value of (contract_expiry_date -
  underlying_next_earnings_date) in days. Negative = contract expires BEFORE earnings; zero = same
  day; positive = AFTER earnings. Use together with `min_days_between_expiry_and_earnings` to target
  a window around the next earnings announcement … [see the official operation documentation for
  full notes]
- **query `min_dte`** (optional; type `Min DTE`; minimum=0; example=1) — The minimum days to expiry.
  Min: 0.
- **query `max_dte`** (optional; type `Max DTE`; minimum=0; example=3) — The maximum days to expiry.
  Min: 0.
- **query `min_earnings_dte`** (optional; type `Min Earnings DTE`; example=5) — The minimum days
  until the next earnings report.
- **query `max_earnings_dte`** (optional; type `Max Earnings DTE`; example=30) — The maximum days
  until the next earnings report.
- **query `min_open_interest`** (optional; type `Min Open Interest`; minimum=0; example=10000) — The
  minimum open interest. Min: 0.
- **query `max_open_interest`** (optional; type `Max Open Interest`; minimum=0; example=35000) — The
  maximum open interest. Min: 0.
- **query `min_volume`** (optional; type `Min Contract Volume`; minimum=0; example=12300) — The
  minimum volume on the option contract. Min: 0.
- **query `max_volume`** (optional; type `Max Contract Volume`; minimum=0; example=55600) — The
  maximum volume on the option contract. Min: 0.
- **query `min_size`** (optional; type `integer (int64)`; format=int64; example=100) — Minimum trade
  size in contracts.
- **query `max_size`** (optional; type `integer (int64)`; format=int64; example=100) — Maximum trade
  size in contracts.
- **query `min_delta`** (optional; type `string`; example=abs(0.5)) — Minimum option delta.
- **query `max_delta`** (optional; type `string`; example=abs(0.5)) — Maximum option delta.
- **query `min_gamma`** (optional; type `string`; example=abs(0.05)) — Minimum option gamma.
- **query `max_gamma`** (optional; type `string`; example=abs(0.05)) — Maximum option gamma.
- **query `min_iv`** (optional; type `string`; example=0.5) — Minimum implied volatility as a
  decimal.
- **query `max_iv`** (optional; type `string`; example=0.5) — Maximum implied volatility as a
  decimal.
- **query `min_theta`** (optional; type `string`; example=abs(0.1)) — Minimum option theta.
- **query `max_theta`** (optional; type `string`; example=abs(0.1)) — Maximum option theta.
- **query `min_diff`** (optional; type `Min Contract Diff`; example=0.53) — The minimum OTM diff of
  a contract. Given a strike price of 120 and an underlying price of 98 the diff for a call option
  would equal to: (120 - 98) / 98 = 0.2245 The diff for a put option would equal to: -1 * (120 - 98)
  / 98 = -0.2245.
- **query `max_diff`** (optional; type `Max Contract Diff`; example=1.34) — The maximum OTM diff of
  a contract. Given a strike price of 120 and an underlying price of 98 the diff for a call option
  would equal to: (120 - 98) / 98 = 0.2245 The diff for a put option would equal to: -1 * (120 - 98)
  / 98 = -0.2245.
- **query `min_marketcap`** (optional; type `Min Marketcap`; minimum=0; example=1000000) — The
  minimum marketcap. Min: 0.
- **query `max_marketcap`** (optional; type `Max Marketcap`; minimum=0; example=250000000) — The
  maximum marketcap. Min: 0.
- **query `min_strike`** (optional; type `Min Strike`; minimum=0; example=120.5) — The minimum
  strike. Min: 0.
- **query `max_strike`** (optional; type `Max Strike`; minimum=0; example=1200) — The maximum
  strike. Min: 0.
- **query `min_vol_oi_ratio`** (optional; type `Min Volume OI Ratio`; minimum=0; example=0.32) — The
  minimum ratio of contract volume to contract open interest. If the open interest of a contract is
  zero, then this ratio is evaluated as if the open interest of the contract was one (to avoid
  divide by zero errors). For example, if you set this ratio to 10, then a contract with zero open
  interest and 7 volume will NOT be included in your results.
- **query `max_vol_oi_ratio`** (optional; type `Max Volume OI Ratio`; minimum=0; example=1.58) — The
  maximum ratio of contract volume to contract open interest. If the open interest of a contract is
  zero, then this ratio is evaluated as if the open interest of the contract was one (to avoid
  divide by zero errors). For example, if you set this ratio to 50, then a contract with zero open
  interest and 75 volume will NOT be included in your results.
- **query `min_premium`** (optional; type `string`; example=25000) — Minimum trade premium in
  dollars.
- **query `max_premium`** (optional; type `string`; example=25000) — Maximum trade premium in
  dollars.
- **query `min_price`** (optional; type `string`; example=5.25) — Minimum option trade price.
- **query `max_price`** (optional; type `string`; example=5.25) — Maximum option trade price.
- **query `min_spread`** (optional; type `string`; example=0.1) — Minimum bid-ask spread percentage.
- **query `max_spread`** (optional; type `string`; example=0.1) — Maximum bid-ask spread percentage.
- **query `min_underlying_price`** (optional; type `string`; example=195.50) — Minimum underlying
  price at execution.
- **query `max_underlying_price`** (optional; type `string`; example=195.50) — Maximum underlying
  price at execution.

**Response payload**

- `200`: `application/json` → `Option Trades Response` — Option Trades Response; fields: `data`
  (array<Option Trade>); data item `Option Trade` fields: `ask_vol` (Option Contract Ask Volume),
  `bid_vol` (Option Contract Bid Volume), `canceled` (Canceled), `delta` (Delta), `er_time` (Stock
  Earnings time), `ewma_nbbo_ask` (EWMA NBBO Ask), `ewma_nbbo_bid` (EWMA NBBO Bid), `exchange`
  (Exchange), `executed_at` (Executed At), `expiry` (Option Contract Expiry), `flow_alert_id` (Flow
  Alert ID), `full_name` (Stock Full Name), `gamma` (Gamma), `id` (Option Trade ID),
  `implied_volatility` (Implied Volatility), `industry_type` (Stock Industry Type), `is_agg`
  (boolean), `issue_type` (string), `marketcap` (Stock Marketcap AUM), `mid_vol` (Option Contract
  Mid Volume), `multi_vol` (Option Contract Multi Leg Volume), `nbbo_ask` (NBBO Ask), `nbbo_bid`
  (NBBO Bid), `next_earnings_date` (Stock Next Earnings Date), `no_side_vol` (Option Contract No
  Side Volume), `open_interest` (Option Contract Open interest), `option_chain_id` (Option Contract
  Symbol), `option_type` (Option Contract Option Type), `premium` (Premium), `price` (Fill Price),
  `report_flags` (Report Flags), `rho` (Rho), `rule_id` (Rule ID), `sector` (Market General Sector),
  `size` (Option Trade Size), `stock_multi_vol` (Option Contract Stock Multi Leg Volume), `strike`
  (Option Contract Strike), `tags` (Tags), `theo` (Theoretical Price), `theta` (Theta), `trade_ids`
  (array<string (uuid)>), `underlying_price` (Underlying Price), `underlying_symbol` (Option
  Contract Underlying Symbol), `upstream_condition_detail` (Upstream Condition Detail), `vega`
  (Vega), `volume` (Option Trade Volume)
- `400`: `application/json` → `Option Trades Error Response` — Option Trades Error Response; fields:
  `error` (object) — Invalid or unsupported option trade request
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error
- `502`: `application/json` → `Option Trades Error Response` — Option Trades Error Response; fields:
  `error` (object) — Option trades service unavailable

#### `GET /api/option-trades/optionable-tickers` — Optionable Tickers

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.OptionTradeController.optionable_tickers)

**What it does:**

Returns the current universe of underlying symbols that have listed options (sourced from the Nasdaq
option-root locators), sorted alphabetically. Query this daily to track new listings and delistings.
Pass `?ticker=SYMBOL` to check a single company instead — the response then reports whether that
symbol currently has listed options. NOTICE: Access to this endpoint is only included in the
Advanced API subscription.

**Parameters**

- **query `ticker`** (optional; type `Ticker`; example=AAPL,INTC) — Optional: check a single symbol
  instead of listing the whole universe.

**Response payload**

- `200`: `application/json` → `Optionable Tickers` — Optionable Tickers; fields: array of string
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Options Pulse — 4

#### `GET /api/options-pulse/total` — Market-wide Options Pulse

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.NasdaqPulseController.total)

**What it does:**

The market-wide Nasdaq Options Pulse gauge: latest snapshot + intraday series for a date.

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Nasdaq Options Pulse Result` — Nasdaq Options Pulse Result; fields:
  Nasdaq Options Pulse Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/options-pulse/top` — Options Pulse Scanner

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.NasdaqPulseController.top)

**What it does:**

Cross-symbol scanner: the latest intraday sentiment per ticker on a date, ranked by sentiment.
`direction` `bullish` ranks highest sentiment first, `bearish` lowest first. Filterable by `ticker`
(prefix), `min_score`/`max_score`, and `min_txn` (minimum opening-buy transactions).

**Parameters**

- **query `direction`** (optional; type `string`) — `bullish` (default) or `bearish`.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `ticker`** (optional; type `string`) — Restrict to tickers with this prefix.
- **query `min_score`** (optional; type `number`) — Minimum sentiment score.
- **query `max_score`** (optional; type `number`) — Maximum sentiment score.
- **query `min_txn`** (optional; type `integer`) — Minimum total opening-buy transactions (put +
  call).
- **query `limit`** (optional; type `integer`) — Max rows (default 50, max 500).

**Response payload**

- `200`: `application/json` → `Nasdaq Options Pulse Result` — Nasdaq Options Pulse Result; fields:
  Nasdaq Options Pulse Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/options-pulse/sectors` — Options Pulse by Sector/Industry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.NasdaqPulseController.sectors)

**What it does:**

Latest Nasdaq Options Pulse sentiment for each sector and industry on a date.

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Nasdaq Options Pulse Result` — Nasdaq Options Pulse Result; fields:
  Nasdaq Options Pulse Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/options-pulse` — Options Pulse for a ticker

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.NasdaqPulseController.ticker)

**What it does:**

The Nasdaq Options Pulse sentiment for a single ticker: the latest intraday snapshot plus the full
intraday series (`hr_min` buckets) for the trade date. `sntm_score` is the running daily sentiment,
`intvl_sntm_score` the per-bucket sentiment, and `put_txn`/`call_txn` the opening-buy transaction
counts.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Nasdaq Options Pulse Result` — Nasdaq Options Pulse Result; fields:
  Nasdaq Options Pulse Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Politician portfolios — 5

#### `GET /api/politician-portfolios/disclosures` — Annual Disclosures List

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PoliticianPortfoliosController.disclosures)

**What it does:**

Returns all annual disclosure file records for politicians. Can be filtered by politician_id and/or
limited to latest per politician. This is an enterprise only endpoint. Contact dev@unusualwhales.com
for details about accessing this data.

**Parameters**

- **query `politician_id`** (optional; type `string (uuid)`; format=uuid) — Filter by politician ID
- **query `latest_only`** (optional; type `boolean`) — If true, returns only the most recent
  disclosure per politician
- **query `year`** (optional; type `integer`) — Filter by disclosure year

**Response payload**

- `200`: `application/json` → `Annual Disclosures List` — Annual Disclosures List; fields: `data`
  (array<Annual Disclosure>); data item `Annual Disclosure` fields: `chamber` (string),
  `disclosure_year` (integer), `filing_date` (string (date)), `id` (string (uuid)), `name` (string),
  `politician_id` (string (uuid)), `politician_name` (string), `url` (string), `version` (integer) —
  Annual Disclosures
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/politician-portfolios/holders/{ticker}` — Politician Portfolio Holders by Ticker

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PoliticianPortfoliosController.holds_ticker)

**What it does:**

Returns all politician portfolio owner names, ID, and holdings for the specified ticker. This is an
enterprise only endpoint. Contact dev@unusualwhales.com for details about accessing this data.

**Parameters**

- **path `ticker`** (required; type `string`) — Stock ticker symbol (e.g., AAPL, TSLA)
- **query `aggregate_all_portfolios`** (optional; type `boolean`) — If true, aggregates all of a
  politicians portfolios into a single portfolio named 'aggregated'. Default is false. Note that
  this does not aggregate holdings within a portfolio, only across portfolios.

**Response payload**

- `200`: `application/json` → `Portfolio Holders` — Portfolio Holders; fields: `data`
  (array<Portfolio Holder>); data item `Portfolio Holder` fields: `full_name` (string), `id` (string
  (uuid)), `max_amount` (number), `mid_amount` (number), `min_amount` (number), `owner` (string) —
  Portfolio Holders
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/politician-portfolios/{politician_id}` — Politician Portfolios

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PoliticianPortfoliosController.portfolios)

**What it does:**

Returns all portfolios and holdings for a politician. This is an enterprise only endpoint. Contact
dev@unusualwhales.com for details about accessing this data.

**Parameters**

- **path `politician_id`** (required; type `string (uuid)`; format=uuid)
- **query `aggregate_all_portfolios`** (optional; type `boolean`) — If true, aggregates all
  portfolios into a single portfolio named 'aggregated'. Default is false. Note that this does not
  aggregate holdings within a portfolio, only across portfolios.

**Response payload**

- `200`: `application/json` → `Politician Portfolios` — Politician Portfolios; fields: `data`
  (array<Politician Portfolio>); data item `Politician Portfolio` fields: `crypto_holdings`
  (array<Crypto Holding>), `id` (string (uuid)), `last_annual_disclosure` (integer),
  `option_holdings` (array<Option Holding>), `owner` (string), `stock_holdings` (array<Stock
  Holding>) — Politician Portfolios
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/politician-portfolios/recent_trades` — Politician Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PoliticianPortfoliosController.recent_trades)

**What it does:**

Returns the latest transacted trades by congress members. If a date is given, will only return
reports, which's transaction date is <= the given input date. This is an enterprise only endpoint.
Contact dev@unusualwhales.com for details about accessing this data.

**Parameters**

- **query `limit`** (optional; type `Default 500 Max 500 Min 1`; default=500; minimum=1;
  maximum=500; example=10) — How many items to return. Default: 500. Max: 500. Min: 1.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `ticker`** (optional; type `Optional Ticker`; example=IOVA) — Optional ticker symbol to
  filter results
- **query `politician_id`** (optional; type `Politician ID`;
  example=18f9fc95-4661-444e-99f5-99d3778e0c31) — The Politician's ID that you'd like to filter by.
- **query `filter_late_reports`** (optional; type `Filter Late Reports Only`; default=false;
  example=true) — If true, will filter out any trades with a transaction date after the input date.
- **query `page`** (optional; type `Page`; example=1) — The page number for the data being requested
  (default is 1).
- **query `disclosure_newer_than`** (optional; type `NewerThan`; example=1715083417) — If provided,
  will only return trades with a disclosure date >= this date.
- **query `disclosure_older_than`** (optional; type `OlderThan`; example=1715083417) — If provided,
  will only return trades with a disclosure date <= this date.
- **query `transaction_newer_than`** (optional; type `NewerThan`; example=1715083417) — If provided,
  will only return trades with a transaction date >= this date.
- **query `transaction_older_than`** (optional; type `OlderThan`; example=1715083417) — If provided,
  will only return trades with a transaction date <= this date.

**Response payload**

- `200`: `application/json` → `Senate Stock` — Senate Stock; fields: `amounts` (Insider Trades
  Amount Range), `filed_at_date` (Insider Trades Filing Date), `is_active` (Is Active), `issuer`
  (Insider Trades Issuer), `member_type` (Insider Trades Member Type), `name` (Insider Trades
  Reporter's Standard Name), `notes` (Insider Trades Filing Notes), `politician_id` (Politician ID),
  `reporter` (Insider Trades Reporter), `ticker` (Stock Ticker), `transaction_date` (Insider Trades
  Transaction Date), `txn_type` (Insider Trades Transaction Type)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/politician-portfolios/people` — Politicians List

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PoliticianPortfoliosController.people)

**What it does:**

Returns all politician names and IDs. This is an enterprise only endpoint. Contact
dev@unusualwhales.com for details about accessing this data.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Politicians List` — Politicians List; fields: `data`
  (array<Politician>); data item `Politician` fields: `id` (string (uuid)), `name` (string) —
  Politicians List
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### President / POTUS — 2

#### `GET /api/potus/posts` — President's Posts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PotusController.posts)

**What it does:**

Returns the President's short-form social media posts (currently Truth Social), newest first. Use
`search_term` to filter posts by text content. Supports pagination via `limit` and `page` (page is
0-indexed).

**Parameters**

- **query `limit`** (optional; type `Default 200 Max 200 Min 1`; default=200; minimum=1;
  maximum=200; example=50) — How many items to return. Default: 200. Max: 200. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `search_term`** (optional; type `POTUS Post Search Term`; example=Walmart) — Optional
  search term to filter posts by their text content.

**Response payload**

- `200`: `application/json` → `POTUS Post` — POTUS Post; fields: `post` (string), `timestamp`
  (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/potus/schedule` — President's Schedule

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PotusController.schedule)

**What it does:**

Returns the President's public schedule for a given date, newest first by start time. If no `date`
is provided, returns entries for the most recent scheduled date. Supports pagination via `limit` and
`page` (page is 0-indexed).

**Parameters**

- **query `date`** (optional; type `Optional POTUS Schedule Date`; example=2026-07-09) — A date in
  the format of YYYY-MM-DD. This is optional and by default returns the most recent scheduled date.
- **query `limit`** (optional; type `Default 200 Max 200 Min 1`; default=200; minimum=1;
  maximum=200; example=50) — How many items to return. Default: 200. Max: 200. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.

**Response payload**

- `200`: `application/json` → `POTUS Schedule Entry` — POTUS Schedule Entry; fields: `coverage`
  (string), `date` (string), `details` (string), `location` (string), `time` (string), `type`
  (string), `url` (string), `video_url` (string)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Prediction markets — 9

#### `GET /api/predictions/market/{asset_id}` — Prediction Market Details

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.market)

**What it does:**

Returns prediction market details for a given asset ID.

**Parameters**

- **path `asset_id`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Prediction Market Detail Response` — Prediction Market Detail
  Response; fields: `data` (Prediction Market Detail) — Prediction Market Details
- `404`: `application/json` → `Error Message stating that the requested element was not found
  causing an empty result to be generated.` — Error Message stating that the requested element was
  not found causing an empty result to be generated.; fields: Error Message stating that the
  requested element was not found causing an empty result to be generated. — Not Found
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/insiders` — Prediction Market Insiders

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.insiders)

**What it does:**

Returns potential insider activity on prediction markets.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Prediction Market Insiders Response` — Prediction Market Insiders
  Response; fields: `data` (Prediction Market Insiders Data) — Prediction Market Insiders
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/market/{asset_id}/liquidity` — Prediction Market Liquidity

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.liquidity)

**What it does:**

Returns liquidity data for a given prediction market asset.

**Parameters**

- **path `asset_id`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Prediction Market Liquidity Response` — Prediction Market Liquidity
  Response; fields: `data` (Prediction Market Liquidity) — Prediction Market Liquidity
- `404`: `application/json` → `Error Message stating that the requested element was not found
  causing an empty result to be generated.` — Error Message stating that the requested element was
  not found causing an empty result to be generated.; fields: Error Message stating that the
  requested element was not found causing an empty result to be generated. — Not Found
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/market/{asset_id}/positions` — Prediction Market Positions

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.positions)

**What it does:**

Returns positions for a given prediction market asset.

**Parameters**

- **path `asset_id`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Prediction Market Positions Response` — Prediction Market Positions
  Response; fields: `data` (array<Prediction Market Position>); data item `Prediction Market
  Position` fields: `amount` (string), `avg_price` (string), `category_num_markets` (integer),
  `category_score` (string), `category_sum_invested` (string), `category_sum_pnl` (string),
  `category_win_rate` (string), `invested_usd` (string), `position_start_at` (string (date-time)),
  `realized_pnl` (string), `smart_score` (string), `tags` (array<string>), `total_bought` (string),
  `total_invested_usd` (string), `updated_at` (string (date-time)), `user_address` (string) —
  Prediction Market Positions
- `404`: `application/json` → `Error Message stating that the requested element was not found
  causing an empty result to be generated.` — Error Message stating that the requested element was
  not found causing an empty result to be generated.; fields: Error Message stating that the
  requested element was not found causing an empty result to be generated. — Not Found
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/user/{user_id}` — Prediction Market User

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.user)

**What it does:**

Returns a prediction market user profile by user/wallet ID.

**Parameters**

- **path `user_id`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Prediction Market User Response` — Prediction Market User Response;
  fields: `data` (Prediction Market User) — Prediction Market User
- `404`: `application/json` → `Error Message stating that the requested element was not found
  causing an empty result to be generated.` — Error Message stating that the requested element was
  not found causing an empty result to be generated.; fields: Error Message stating that the
  requested element was not found causing an empty result to be generated. — Not Found
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/whales` — Prediction Market Whales

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.whales)

**What it does:**

Returns large prediction market traders.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Prediction Market Whales Response` — Prediction Market Whales
  Response; fields: `data` (Prediction Market Whales Data) — Prediction Market Whales
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/smart-money` — Prediction Smart Money

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.smart_money)

**What it does:**

Returns profitable prediction market traders. Categories: Crypto, Culture, Finance, Politics,
Science, Sports, Technology.

**Parameters**

- **query `categories`** (optional; type `string`)
- **query `min_price`** (optional; type `number`)
- **query `max_price`** (optional; type `number`)

**Response payload**

- `200`: `application/json` → `Prediction Smart Money Response` — Prediction Smart Money Response;
  fields: `data` (Prediction Smart Money Data) — Prediction Smart Money
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/search-users` — Search Prediction Market Users

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.search_users)

**What it does:**

Search for prediction market users by query.

**Parameters**

- **query `q`** (required; type `string`)

**Response payload**

- `200`: `application/json` → `Search Prediction Market Users Response` — Search Prediction Market
  Users Response; fields: `data` (array<Prediction Market User Search Hit>); data item `Prediction
  Market User Search Hit` fields: `address` (string), `name` (string) — Search Prediction Market
  Users
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/predictions/unusual` — Unusual Prediction Markets

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PredictionController.unusual)

**What it does:**

Returns prediction markets with unusual activity. Categories: Crypto, Culture, Finance, Politics,
Science, Sports, Technology.

**Parameters**

- **query `categories`** (optional; type `string`)
- **query `limit`** (optional; type `integer`; default=50)
- **query `offset`** (optional; type `integer`; default=0)

**Response payload**

- `200`: `application/json` → `Unusual Prediction Markets Response` — Unusual Prediction Markets
  Response; fields: `data` (Unusual Prediction Markets Data) — Unusual Prediction Markets
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Private markets — 9

#### `GET /api/private-markets/companies` — List Private Markets Companies

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.companies)

**What it does:**

Returns Nasdaq Private Markets companies, optionally filtered by sector or name. Sorted
alphabetically. This is a premium endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **query `sector`** (optional; type `string`) — Exact-match sector filter (e.g. "Technology").
- **query `name`** (optional; type `string`) — Case-insensitive substring match against company
  name.
- **query `limit`** (optional; type `integer`; default=100; minimum=1; maximum=500)
- **query `offset`** (optional; type `integer`; default=0; minimum=0)

**Response payload**

- `200`: `application/json` → `Private Markets Companies` — Private Markets Companies; fields:
  `data` (array<Private Markets Company>); data item `Private Markets Company` fields: `description`
  (string), `doing_business_as` (string), `headquarters` (string), `id` (string (uuid)),
  `marketing_url` (string), `name` (string), `npm_ticker` (string), `sector` (string), `sub_sectors`
  (array<string>), `year_founded` (integer) — Private Markets Companies
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/companies/{npm_ticker}` — Private Markets Company Profile

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.company_profile)

**What it does:**

Returns a profile for a single private-markets company including the latest pricing tick, total
funding raised across known rounds, and the count of disclosed investors. This is a premium
endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **path `npm_ticker`** (required; type `string`) — Nasdaq Private Markets ticker (e.g. "OPENAI").

**Response payload**

- `200`: `application/json` → `Private Markets Company Profile Response` — Private Markets Company
  Profile Response; fields: `data` (Private Markets Company Profile) — Private Markets Company
  Profile
- `404`: `application/json` → `Error Message stating that the requested element was not found
  causing an empty result to be generated.` — Error Message stating that the requested element was
  not found causing an empty result to be generated.; fields: Error Message stating that the
  requested element was not found causing an empty result to be generated. — Not Found
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/companies/{npm_ticker}/funding` — Private Markets Funding Rounds

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.funding)

**What it does:**

Returns the funding round history for a single private-markets company, ordered most-recent first.
This is a premium endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **path `npm_ticker`** (required; type `string`)
- **query `limit`** (optional; type `integer`)
- **query `offset`** (optional; type `integer`)

**Response payload**

- `200`: `application/json` → `Private Markets Funding Rounds` — Private Markets Funding Rounds;
  fields: `data` (array<Private Markets Funding Round>); data item `Private Markets Funding Round`
  fields: `financing_round` (string), `id` (string (uuid)), `investment_amount` (number),
  `investment_date` (string (date)), `npm_ticker` (string) — Private Markets Funding Rounds
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/investors/{name}` — Private Markets Investor Profile

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.investor_profile)

**What it does:**

Returns the portfolio of companies for a specific investor (by name). This is a premium endpoint.
Contact dev@unusualwhales.com to request access.

**Parameters**

- **path `name`** (required; type `string`) — The investor's name (URL-encoded).

**Response payload**

- `200`: `application/json` → `Private Markets Investor Profile Response` — Private Markets Investor
  Profile Response; fields: `data` (Private Markets Investor Profile) — Private Markets Investor
  Profile
- `404`: `application/json` → `Error Message stating that the requested element was not found
  causing an empty result to be generated.` — Error Message stating that the requested element was
  not found causing an empty result to be generated.; fields: Error Message stating that the
  requested element was not found causing an empty result to be generated. — Not Found
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/companies/{npm_ticker}/investors` — Private Markets Investors for Company

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.investors)

**What it does:**

Returns the disclosed investors for a single private-markets company, ordered alphabetically. This
is a premium endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **path `npm_ticker`** (required; type `string`)
- **query `limit`** (optional; type `integer`)
- **query `offset`** (optional; type `integer`)

**Response payload**

- `200`: `application/json` → `Private Markets Investors` — Private Markets Investors; fields:
  `data` (array<Private Markets Investor>); data item `Private Markets Investor` fields: `id`
  (string (uuid)), `name` (string), `npm_ticker` (string) — Private Markets Investors
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/companies/{npm_ticker}/management` — Private Markets Management

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.management)

**What it does:**

Returns disclosed management/leadership for a single private-markets company. This is a premium
endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **path `npm_ticker`** (required; type `string`)
- **query `limit`** (optional; type `integer`)
- **query `offset`** (optional; type `integer`)

**Response payload**

- `200`: `application/json` → `Private Markets Management` — Private Markets Management; fields:
  `data` (array<Private Markets Management Member>); data item `Private Markets Management Member`
  fields: `first_name` (string), `id` (string (uuid)), `last_name` (string), `npm_ticker` (string),
  `title` (string) — Private Markets Management
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/companies/{npm_ticker}/pricing` — Private Markets Pricing History

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.pricing)

**What it does:**

Returns historical implied per-share pricing for a single private-markets company. This is a premium
endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **path `npm_ticker`** (required; type `string`)
- **query `start_date`** (optional; type `string (date)`; format=date) — Inclusive lower bound on
  pricing date (ISO 8601).
- **query `end_date`** (optional; type `string (date)`; format=date) — Inclusive upper bound on
  pricing date (ISO 8601).
- **query `limit`** (optional; type `integer`)
- **query `offset`** (optional; type `integer`)

**Response payload**

- `200`: `application/json` → `Private Markets Pricing History` — Private Markets Pricing History;
  fields: `data` (array<Private Markets Pricing Tick>); data item `Private Markets Pricing Tick`
  fields: `date` (string (date)), `id` (string (uuid)), `npm_ticker` (string), `price` (number) —
  Private Markets Pricing
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/search` — Search Private Markets

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.search)

**What it does:**

Substring-search across private-markets companies and investors. This is a premium endpoint. Contact
dev@unusualwhales.com to request access.

**Parameters**

- **query `query`** (required; type `string`) — Search string (case-insensitive substring match).
- **query `limit`** (optional; type `integer`; default=10; minimum=1; maximum=50)

**Response payload**

- `200`: `application/json` → `Private Markets Search` — Private Markets Search; fields: `data`
  (object) — Private Markets Search
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/private-markets/investors` — Top Private Markets Investors

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.PrivateMarketsController.top_investors)

**What it does:**

Returns the most prolific investors across the private-markets dataset, ordered by distinct company
count (descending). This is a premium endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **query `limit`** (optional; type `integer`)
- **query `offset`** (optional; type `integer`)

**Response payload**

- `200`: `application/json` → `Top Private Markets Investors` — Top Private Markets Investors;
  fields: `data` (array<Top Private Markets Investor>); data item `Top Private Markets Investor`
  fields: `company_count` (integer), `name` (string) — Top Private Markets Investors
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Screeners — 4

#### `GET /api/screener/analysts` — Analyst Rating

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ScreenerController.analyst_ratings)

**What it does:**

Returns the latest analyst rating for the given ticker.

**Parameters**

- **query `ticker`** (optional; type `SingleTicker`; example=AAPL) — A single ticker
- **query `limit`** (optional; type `integer`) — How many items to return. Default: 500, Max: 500,
  Min: 1
- **query `action`** (optional; type `Analyst Action`; enum=[initiated, reiterated, downgraded,
  upgraded, maintained]; example=maintained) — The action of the recommendation.
- **query `recommendation`** (optional; type `Analyst Recommendation`; enum=[buy, hold, sell];
  example=hold) — The recommendation the analyst gave out.

**Response payload**

- `200`: `application/json` → `Analyst Rating` — Analyst Rating; fields: `action` (Analyst Field
  Action), `analyst_name` (Analyst Field Name), `firm` (Analyst Field Firm), `recommendation`
  (Analyst Field Recommendation), `sector` (Market General Sector), `target` (Analyst Field Target
  Price), `ticker` (Analyst Sector), `timestamp` (Analyst Field Time) — Analyst rating response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/screener/option-contracts` — Hottest Chains

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ScreenerController.contract_screener)

**What it does:**

A contract screener endpoint to screen the market for contracts by a variety of filter options. For
an example of what can be build with this endpoint check out the Hottest Contracts on UnusualWhales.
For real time streaming of the same data, subscribe to the `contract_screener` websocket channel,
see https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.contract_screener.
NOTE: Contracts with a volume of less than 200 are not being returned

**Parameters**

- **query `ticker_symbol`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `sectors[]`** (optional; type `Sectors`; example=[Consumer Cyclical, Technology,
  Utilities]) — An array of 1 or more sectors.
- **query `unusual`** (optional; type `boolean`) — Convenience preset that returns only "unusual"
  contracts by applying the contract-expressible subset of the live options flow criteria:
  volume>OI, OTM, DTE≤60, ask-side≥50%, premium≥$10k, issue types ADR/Common Stock/ETF. These are
  applied as defaults, so any of those filters you pass explicitly (e.g. `max_dte=30`,
  `min_premium=25000`) overrides the preset.
- **query `min_underlying_price`** (optional; type `string`) — The minimum stock price.
- **query `max_underlying_price`** (optional; type `string`) — The maximum stock price.
- **query `is_otm`** (optional; type `boolean`) — Only include contracts which are currently out of
  the money.
- **query `exclude_ex_div_ticker`** (optional; type `boolean`) — When set to true, all tickers that
  trade ex-dividend today will be excluded. This is useful since on the day prior to the ex-dividend
  date, there will be above-average ITM call flow due to dividend arbitrage traders.
- **query `min_dte`** (optional; type `integer`) — The minimum days to expiry.
- **query `max_dte`** (optional; type `integer`) — The maximum days to expiry.
- **query `min_diff`** (optional; type `string`) — The minimum OTM diff of a contract.
- **query `max_diff`** (optional; type `string`) — The maximum OTM diff of a contract.
- **query `min_strike`** (optional; type `string`) — The minimum strike.
- **query `max_strike`** (optional; type `string`) — The maximum strike.
- **query `type`** (optional; type `OptionType`; enum=[call, Call, put, Put]) — The option type to
  filter by if specified.
- **query `expiry_dates[]`** (optional; type `Expiry dates`; example=[2024-02-02, 2024-01-26]) — An
  array of 1 or more expiry dates.
- **query `min_marketcap`** (optional; type `string`) — The minimum marketcap.
- **query `max_marketcap`** (optional; type `string`) — The maximum marketcap.
- **query `min_volume`** (optional; type `Min Contract Volume`; minimum=0; example=12300) — The
  minimum volume on the option contract. Min: 0.
- **query `max_volume`** (optional; type `Max Contract Volume`; minimum=0; example=55600) — The
  maximum volume on the option contract. Min: 0.
- **query `min_ticker_30_d_avg_volume`** (optional; type `integer`) — The minimum 30-day average
  stock volume for the underlying ticker.
- **query `max_ticker_30_d_avg_volume`** (optional; type `integer`) — The maximum 30-day average
  stock volume for the underlying ticker.
- **query `min_contract_30_d_avg_volume`** (optional; type `integer`) — The minimum 30-day average
  options contract volume for the underlying ticker.
- **query `max_contract_30_d_avg_volume`** (optional; type `integer`) — The maximum 30-day average
  options contract volume for the underlying ticker.
- **query `min_multileg_volume_ratio`** (optional; type `string`) — The minimum multi leg volume to
  contract volume ratio.
- **query `max_multileg_volume_ratio`** (optional; type `string`) — The maximum multi leg volume to
  contract volume ratio.
- **query `min_floor_volume_ratio`** (optional; type `string`) — The minimum floor volume to
  contract volume ratio.
- **query `max_floor_volume_ratio`** (optional; type `string`) — The maximum floor volume to
  contract volume ratio.
- **query `min_perc_change`** (optional; type `string`) — The minimum % price change of the contract
  to the previous day. Acceptable range: -1.00 to +inf.
- **query `max_perc_change`** (optional; type `string`) — The maximum % price change of the contract
  to the previous day. Acceptable range: -1.00 to +inf.
- **query `min_daily_perc_change`** (optional; type `string`) — The minimum intraday price change of
  the contract from open till now.
- **query `max_daily_perc_change`** (optional; type `string`) — The maximum intraday price change
  for the contract since market open.
- **query `min_premium`** (optional; type `string`) — The minimum premium on that contract.
- **query `max_premium`** (optional; type `string`) — The maximum premium on that contract.
- **query `min_avg_price`** (optional; type `string`) — The minimum average price of the contract.
- **query `max_avg_price`** (optional; type `string`) — The maximum average price of the contract.
- **query `min_volume_oi_ratio`** (optional; type `string`) — The minimum contract volume to open
  interest ratio.
- **query `max_volume_oi_ratio`** (optional; type `string`) — The maximum contract volume to open
  interest ratio.
- **query `min_open_interest`** (optional; type `integer`) — The minimum open interest on that
  contract.
- **query `max_open_interest`** (optional; type `integer`) — The maximum open interest on that
  contract.
- **query `min_floor_volume`** (optional; type `integer`) — The minimum floor volume on that
  contract.
- **query `max_floor_volume`** (optional; type `integer`) — The maximum floor volume on that
  contract.
- **query `vol_greater_oi`** (optional; type `boolean`) — Only include contracts where the volume is
  greater than the open interest.
- **query `issue_types[]`** (optional; type `Issue types`; example=[Common Stock, Index]) — An array
  of 1 or more issue types.
- **query `min_ask_perc`** (optional; type `string`) — The minimum ask percentage of volume that
  transacted on the ask.
- **query `max_ask_perc`** (optional; type `string`) — The maximum ask percentage of volume that
  transacted on the ask.
- **query `min_bid_perc`** (optional; type `string`) — The minimum bid percentage of volume that
  transacted on the bid.
- **query `max_bid_perc`** (optional; type `string`) — The maximum bid percentage of volume that
  transacted on the bid.
- **query `min_skew_perc`** (optional; type `string`) — The minimum skew percentage. Setting this to
  0.8 would return all contracts where either 80% of vol transacted on the ask or bid side
- **query `max_skew_perc`** (optional; type `string`) — The maximum skew percentage.Setting this to
  0.8 would return all contracts where max 80% of vol transacted on the ask or bid side
- **query `min_bull_perc`** (optional; type `string`) — The minimum bull percentage.
- **query `max_bull_perc`** (optional; type `string`) — The maximum bull percentage.
- **query `min_bear_perc`** (optional; type `string`) — The minimum bear percentage.
- **query `max_bear_perc`** (optional; type `string`) — The maximum bear percentage.
- **query `min_bid_side_perc_7_day`** (optional; type `string`) — The minimum percentage of days
  over the last 7 days where the contract traded primarily on the bid side
- **query `max_bid_side_perc_7_day`** (optional; type `string`) — The maximum percentage of days
  over the last 7 days where the contract traded primarily on the bid side
- **query `min_ask_side_perc_7_day`** (optional; type `string`) — The minimum percentage of days
  over the last 7 days where the contract traded primarily on the ask side
- **query `max_ask_side_perc_7_day`** (optional; type `string`) — The maximum percentage of days
  over the last 7 days where the contract traded primarily on the ask side
- **query `min_days_of_oi_increases`** (optional; type `integer`) — The minimum days of consecutive
  trading days where the open interest increased
- **query `max_days_of_oi_increases`** (optional; type `integer`) — The maximum days of consecutive
  trading days where the open interest increased
- **query `min_days_of_vol_greater_than_oi`** (optional; type `integer`) — The minimum days of
  consecutive days where volume was greater than open interest.
- **query `max_days_of_vol_greater_than_oi`** (optional; type `integer`) — The maximum days of
  consecutive days where volume was greater than open interest.
- **query `min_iv_perc`** (optional; type `string`) — The minimum implied volatility percentage.
- **query `max_iv_perc`** (optional; type `string`) — The maximum implied volatility percentage.
- **query `min_delta`** (optional; type `string`) — The minimum delta. Acceptable range: -1.00 to
  +1.00.
- **query `max_delta`** (optional; type `string`) — The maximum delta. Acceptable range: -1.00 to
  +1.00.
- **query `min_gamma`** (optional; type `string`) — The minimum gamma. Acceptable range: 0.00 to
  +inf.
- **query `max_gamma`** (optional; type `string`) — The maximum gamma. Acceptable range: 0.00 to
  +inf.
- **query `min_theta`** (optional; type `string`) — The minimum theta. Acceptable range: -inf to
  0.00.
- **query `max_theta`** (optional; type `string`) — The maximum theta. Acceptable range: -inf to
  0.00.
- **query `min_vega`** (optional; type `string`) — The minimum vega. Acceptable range: 0.00 to +inf.
- **query `max_vega`** (optional; type `string`) — The maximum vega. Acceptable range: 0.00 to +inf.
- **query `min_return_on_capital_perc`** (optional; type `string`) — The minimum return on capital
  percentage (ROC).
- **query `max_return_on_capital_perc`** (optional; type `string`) — The maximum return on capital
  percentage (ROC).
- **query `min_oi_change_perc`** (optional; type `string`) — The minimum open interest change
  percentage. Acceptable range: -1.00 to +inf.
- **query `max_oi_change_perc`** (optional; type `string`) — The maximum open interest change
  percentage. Acceptable range: -1.00 to +inf.
- **query `min_oi_change`** (optional; type `integer`) — The minimum open interest change as an
  absolute change.
- **query `max_oi_change`** (optional; type `integer`) — The maximum open interest change as an
  absolute change.
- **query `min_volume_ticker_vol_ratio`** (optional; type `string`) — The minimum ratio of contract
  volume to total option volume of the underlying. Acceptable range: 0.00 to 1.00.
- **query `max_volume_ticker_vol_ratio`** (optional; type `string`) — The maximum ratio of contract
  volume to total option volume of the underlying. Acceptable range: 0.00 to 1.00.
- **query `min_sweep_volume_ratio`** (optional; type `string`) — The minimum sweep volume ratio.
  Acceptable range: 0.00 to 1.00.
- **query `max_sweep_volume_ratio`** (optional; type `string`) — The maximum sweep volume ratio.
  Acceptable range: 0.00 to 1.00.
- **query `min_from_low_perc`** (optional; type `string`) — The minimum percentage change of the
  current price from todays low. Acceptable range: -1.00 to +inf.
- **query `max_from_low_perc`** (optional; type `string`) — The maximum percentage change of the
  current price from todays low. Acceptable range: -1.00 to +inf.
- **query `min_from_high_perc`** (optional; type `string`) — The minimum percentage change of the
  current price from todays high. Acceptable range: -1.00 to +inf.
- **query `max_from_high_perc`** (optional; type `string`) — The maximum percentage change of the
  current price from todays high. Acceptable range: -1.00 to +inf.
- **query `min_earnings_dte`** (optional; type `Min Earnings DTE`; example=5) — The minimum days
  until the next earnings report.
- **query `max_earnings_dte`** (optional; type `Max Earnings DTE`; example=30) — The maximum days
  until the next earnings report.
- **query `min_days_between_expiry_and_earnings`** (optional; type
  `MinDaysBetweenExpiryAndEarnings`; example=1) — Minimum value of (contract_expiry_date -
  underlying_next_earnings_date) in days. Negative = contract expires BEFORE earnings; zero = same
  day; positive = AFTER earnings. Use together with `max_days_between_expiry_and_earnings` to target
  a window around the next earnings announcement … [see the official operation documentation for
  full notes]
- **query `max_days_between_expiry_and_earnings`** (optional; type
  `MaxDaysBetweenExpiryAndEarnings`; example=6) — Maximum value of (contract_expiry_date -
  underlying_next_earnings_date) in days. Negative = contract expires BEFORE earnings; zero = same
  day; positive = AFTER earnings. Use together with `min_days_between_expiry_and_earnings` to target
  a window around the next earnings announcement … [see the official operation documentation for
  full notes]
- **query `min_transactions`** (optional; type `integer`) — The minimum number of transactions.
- **query `max_transactions`** (optional; type `integer`) — The maximum number of transactions.
- **query `min_close`** (optional; type `string`) — The minimum contract price (not underlying
  price).
- **query `max_close`** (optional; type `string`) — The maximum contract price (not underlying
  price).
- **query `order`** (optional; type `Screener contract order by field`; enum=[bid_ask_vol,
  bull_bear_vol, contract_pricing, daily_perc_change, diff, dte, earnings, expires, …];
  example=volume) — The field to order by.
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.
- **query `limit`** (optional; type `Default 50, Max 250 Min 1`; default=1; minimum=1; maximum=250;
  example=10) — How many items to return. Default: 50. Max: 250. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `is_new`** (optional; type `boolean`) — Return only new option contracts
- **query `opex_only`** (optional; type `boolean`) — Return only monthly option expirations

**Response payload**

- `200`: `application/json` → `Option Contract Screener response.` — Option Contract Screener
  response.; fields: `ask_side_volume` (Option Contract Ask Volume), `avg_price` (Option Contract
  Avg Price), `bid_side_volume` (Option Contract Bid Volume), `chain_prev_close` (Option Contract
  Previous Close Price), `close` (Option Contract Close), `cross_volume` (Option Contract Cross
  Volume), `er_time` (Stock Earnings time), `floor_volume` (Option Contract Floor Volume), `high`
  (Option Contract High), `last_fill` (Option Contract Last Transaction Time), `low` (Option
  Contract Low), `mid_volume` (Option Contract Mid Volume), `multileg_volume` (Option Contract Multi
  Leg Volume), `next_earnings_date` (Stock Next Earnings Date), `no_side_volume` (Option Contract No
  Side Volume), `open` (Option Contract Open), `open_interest` (Option Contract Open interest),
  `option_symbol` (Option Contract Symbol), `premium` (Option Contract Premium), `sector` (Market
  General Sector), `stock_multi_leg_volume` (Option Contract Stock Multi Leg Volume), `stock_price`
  (Stock Close Price), `sweep_volume` (Option Contract Sweep Volume), `ticker_vol` (Stock Total
  Volume), `total_ask_changes` (Option Contract Total Ask Changes), `total_bid_changes` (Option
  Contract Total Bid Changes), `trades` (Option Contract Total Trades Count), `volume` (Option
  Contract Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/screener/stocks` — Stock Screener

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ScreenerController.stock_screener)

**What it does:**

A stock screener endpoint to screen the market for stocks by a variety of filter options. For an
example of what can be build with this endpoint check out the Stock Screener on UnusualWhales.

**Parameters**

- **query `ticker`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `issue_types[]`** (optional; type `Issue types`; example=[Common Stock, Index]) — An array
  of 1 or more issue types.
- **query `min_change`** (optional; type `string`) — The minimum % change to the previous trading
  day.
- **query `max_change`** (optional; type `string`) — The maximum % change to the previous trading
  day.
- **query `min_underlying_price`** (optional; type `string`) — The minimum stock price.
- **query `max_underlying_price`** (optional; type `string`) — The maximum stock price.
- **query `is_s_p_500`** (optional; type `boolean`) — Boolean whether to only include stocks which
  are part of the S&P 500. Setting this to false has no effect.
- **query `has_dividends`** (optional; type `boolean`) — Boolean wheter to only include stocks which
  pay dividends. Setting this to false has no effect.
- **query `sectors[]`** (optional; type `Sectors`; example=[Consumer Cyclical, Technology,
  Utilities]) — An array of 1 or more sectors.
- **query `min_marketcap`** (optional; type `string`) — The minimum marketcap.
- **query `max_marketcap`** (optional; type `string`) — The maximum marketcap.
- **query `min_perc_3_day_total`** (optional; type `string`) — The minimum ratio of options volume
  vs 3 day avg options volume.
- **query `max_perc_3_day_total`** (optional; type `string`) — The maximum ratio of options volume
  vs 3 day avg options volume.
- **query `min_perc_3_day_call`** (optional; type `string`) — The minimum ratio of call options
  volume vs 3 day avg call options volume.
- **query `max_perc_3_day_call`** (optional; type `string`) — The maximum ratio of call options
  volume vs 3 day avg call options volume.
- **query `min_perc_3_day_put`** (optional; type `string`) — The minimum ratio of put options volume
  vs 3 day avg put options volume.
- **query `max_perc_3_day_put`** (optional; type `string`) — The maximum ratio of put options volume
  vs 3 day avg put options volume.
- **query `min_perc_30_day_total`** (optional; type `string`) — The minimum ratio of options volume
  vs 30 day avg options volume.
- **query `max_perc_30_day_total`** (optional; type `string`) — The maximum ratio of options volume
  vs 30 day avg options volume.
- **query `min_perc_30_day_call`** (optional; type `string`) — The minimum ratio of call options
  volume vs 30 day avg call options volume.
- **query `max_perc_30_day_call`** (optional; type `string`) — The maximum ratio of call options
  volume vs 30 day avg call options volume.
- **query `min_perc_30_day_put`** (optional; type `string`) — The minimum ratio of put options
  volume vs 30 day avg put options volume.
- **query `max_perc_30_day_put`** (optional; type `string`) — The maximum ratio of put options
  volume vs 30 day avg put options volume.
- **query `min_total_oi_change_perc`** (optional; type `string`) — The minimum open interest change
  compared to the previous day.
- **query `max_total_oi_change_perc`** (optional; type `string`) — The maximum open interest change
  compared to the previous day.
- **query `min_call_oi_change_perc`** (optional; type `string`) — The minimum open interest change
  of call contracts compared to the previous day.
- **query `max_call_oi_change_perc`** (optional; type `string`) — The maximum open interest change
  of call contracts compared to the previous day.
- **query `min_put_oi_change_perc`** (optional; type `string`) — The minimum open interest change of
  put contracts compared to the previous day.
- **query `max_put_oi_change_perc`** (optional; type `string`) — The maximum open interest change of
  put contracts compared to the previous day.
- **query `min_implied_move`** (optional; type `string`) — The minimum implied move.
- **query `max_implied_move`** (optional; type `string`) — The maximum implied move.
- **query `min_implied_move_perc`** (optional; type `string`) — The minimum implied move perc.
- **query `max_implied_move_perc`** (optional; type `string`) — The maximum implied move perc.
- **query `min_volatility`** (optional; type `string`) — The minimum volatility.
- **query `max_volatility`** (optional; type `string`) — The maximum volatility.
- **query `min_iv_rank`** (optional; type `string`) — The minimum iv rank.
- **query `max_iv_rank`** (optional; type `string`) — The maximum iv rank.
- **query `min_volume`** (optional; type `integer`) — The minimum options volume.
- **query `max_volume`** (optional; type `integer`) — The maximum options volume.
- **query `min_call_volume`** (optional; type `integer`) — The minimum call options volume.
- **query `max_call_volume`** (optional; type `integer`) — The maximum call options volume.
- **query `min_put_volume`** (optional; type `integer`) — The minimum put options volume.
- **query `max_put_volume`** (optional; type `integer`) — The maximum put options volume.
- **query `min_premium`** (optional; type `string`) — The minimum options premium.
- **query `max_premium`** (optional; type `string`) — The minimum options premium.
- **query `min_call_premium`** (optional; type `string`) — The minimum call options premium.
- **query `max_call_premium`** (optional; type `string`) — The minimum call options premium.
- **query `min_put_premium`** (optional; type `string`) — The minimum put options premium.
- **query `max_put_premium`** (optional; type `string`) — The minimum put options premium.
- **query `min_net_premium`** (optional; type `string`) — The minimum net options premium.
- **query `max_net_premium`** (optional; type `string`) — The minimum net options premium.
- **query `min_net_call_premium`** (optional; type `string`) — The minimum net call options premium.
- **query `max_net_call_premium`** (optional; type `string`) — The maximum net call options premium.
- **query `min_net_put_premium`** (optional; type `string`) — The minimum net put options premium.
- **query `max_net_put_premium`** (optional; type `string`) — The maximum net put options premium.
- **query `min_oi`** (optional; type `integer`) — The minimum open interest.
- **query `max_oi`** (optional; type `integer`) — The maximum open interest.
- **query `min_oi_vs_vol`** (optional; type `string`) — The minimum open interest vs options volume
  ratio.
- **query `max_oi_vs_vol`** (optional; type `string`) — The maximum open interest vs options volume
  ratio.
- **query `min_put_call_ratio`** (optional; type `string`) — The minimum put to call ratio.
- **query `max_put_call_ratio`** (optional; type `string`) — The maximum put to call ratio.
- **query `order`** (optional; type `Screener order by field`; enum=[avg_30_day_call_oi,
  avg_30_day_call_volume, avg_30_day_put_oi, avg_30_day_put_volume, avg_3_day_call_volume,
  avg_3_day_put_volume, avg_7_day_call_volume, avg_7_day_put_volume, …]; example=premium) — The
  field to order by.
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.
- **query `min_stock_volume_vs_avg30_volume`** (optional; type `string`) — The minimum stock volume
  vs average 30 day volume.
- **query `max_avg30_volume`** (optional; type `string`) — The maximum stock volume vs average 30
  day volume.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Stock Screener response` — Stock Screener response; fields:
  `avg_30_day_call_volume` (Market General Avg 30 Day Call Volume), `avg_30_day_put_volume` (Market
  General Avg 30 Day Put Volume), `avg_3_day_call_volume` (Market General Avg 3 Day Call Volume),
  `avg_3_day_put_volume` (Market General Avg 3 Day Put Volume), `avg_7_day_call_volume` (Market
  General Avg 7 Day Call Volume), `avg_7_day_put_volume` (Market General Avg 7 Day Put Volume),
  `bearish_premium` (Market General Bearish Premium), `bullish_premium` (Market General Bullish
  Premium), `call_open_interest` (Market General Call Open Interest), `call_premium` (Market General
  Call Premium), `call_volume` (Market General Call Volume), `call_volume_ask_side` (Market General
  Call Volume Ask Side), `call_volume_bid_side` (Market General Call Volume Bid Side), `close`
  (Stock Close Price), `er_time` (Stock Earnings time), `implied_move` (Stock Implied Move),
  `implied_move_perc` (Stock Implied Move Perc), `is_index` (Stock Is Index Ticker), `issue_type`
  (Stock Issue Type), `iv30d` (Stock IV 30d), `iv30d_1d` (Stock IV 30d 1D), `iv30d_1m` (Stock IV 30d
  1M), `iv30d_1w` (Stock IV 30d 1W), `iv_rank` (Stock IV 30d 1M), `marketcap` (Stock Marketcap AUM),
  `net_call_premium` (Market General Net Call Premium), `net_put_premium` (Market General Net Put
  Premium), `next_dividend_date` (Stock Next Dividend Date), `next_earnings_date` (Stock Next
  Earnings Date), `prev_call_oi` (Market General Previous Call Open interest), `prev_close` (Stock
  Prev Close Price), `prev_put_oi` (Market General Previous Put Open interest), `put_call_ratio`
  (Market General Put Call Ratio), `put_open_interest` (Market General Put Open Interest),
  `put_premium` (Market General Put Premium), `put_volume` (Market General Put Volume),
  `put_volume_ask_side` (Market General Put Volume Ask Side), `put_volume_bid_side` (Market General
  Put Volume Bid Side), `relative_volume` (Stock Relative Volume), `sector` (Market General Sector),
  `ticker` (Stock Ticker), `total_open_interest` (Stock Total Open Interest), `volatility` (Stock
  Volatility), `week_52_high` (Stock Week 52 High), `week_52_low` (Stock Week 52 low)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/option-activity/unusual` — Unusual Options Activity

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ScreenerController.unusual_activity)

**What it does:**

Returns option contracts flagged as unusual — the API equivalent of the live options flow "unusual"
view. This is the Hottest Chains screener with the `unusual` preset applied by default: the
contract-expressible subset of the live-options-flow criteria — volume > OI, out-of-the-money, DTE ≤
60, ask-side ≥ 50%, premium ≥ $10k, issue types ADR/Common Stock/ETF. Every Hottest Chains filter is
accepted here and overrides the preset (e.g. `?min_premium=25000&max_dte=30`); pass `unusual=false`
to return the full screen. For the complete list of filters and the response shape see the Hottest
Chains endpoint. For real time streaming of the same data, subscribe to the `contract_screener`
websocket channel, see
https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.contract_screener. NOTE:
Contracts with a volume of less than 200 are not being returned

**Parameters**

- **query `ticker_symbol`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `sectors[]`** (optional; type `Sectors`; example=[Consumer Cyclical, Technology,
  Utilities]) — An array of 1 or more sectors.
- **query `unusual`** (optional; type `boolean`) — Whether to apply the unusual preset (volume>OI,
  OTM, DTE≤60, ask-side≥50%, premium≥$10k). Defaults to true for this endpoint; pass false to return
  the full screen.
- **query `min_premium`** (optional; type `integer`) — The minimum total premium on the contract.
- **query `max_dte`** (optional; type `integer`) — The maximum days to expiry.
- **query `issue_types[]`** (optional; type `Issue types`; example=[Common Stock, Index]) — An array
  of 1 or more issue types.
- **query `order`** (optional; type `string`) — The field to order the results by.
- **query `order_direction`** (optional; type `string`) — Whether to sort ascending or descending.
  Defaults to descending.
- **query `limit`** (optional; type `integer`) — How many items to return. Defaults to 100. Max 200.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Option Contract Screener response.` — Option Contract Screener
  response.; fields: `ask_side_volume` (Option Contract Ask Volume), `avg_price` (Option Contract
  Avg Price), `bid_side_volume` (Option Contract Bid Volume), `chain_prev_close` (Option Contract
  Previous Close Price), `close` (Option Contract Close), `cross_volume` (Option Contract Cross
  Volume), `er_time` (Stock Earnings time), `floor_volume` (Option Contract Floor Volume), `high`
  (Option Contract High), `last_fill` (Option Contract Last Transaction Time), `low` (Option
  Contract Low), `mid_volume` (Option Contract Mid Volume), `multileg_volume` (Option Contract Multi
  Leg Volume), `next_earnings_date` (Stock Next Earnings Date), `no_side_volume` (Option Contract No
  Side Volume), `open` (Option Contract Open), `open_interest` (Option Contract Open interest),
  `option_symbol` (Option Contract Symbol), `premium` (Option Contract Premium), `sector` (Market
  General Sector), `stock_multi_leg_volume` (Option Contract Stock Multi Leg Volume), `stock_price`
  (Stock Close Price), `sweep_volume` (Option Contract Sweep Volume), `ticker_vol` (Stock Total
  Volume), `total_ask_changes` (Option Contract Total Ask Changes), `total_bid_changes` (Option
  Contract Total Bid Changes), `trades` (Option Contract Total Trades Count), `volume` (Option
  Contract Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Seasonality — 4

#### `GET /api/seasonality/{ticker}/monthly` — Average return per month

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SeasonalityController.monthly)

**What it does:**

Returns the average return by month for the given ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Seasonality Monthly` — Seasonality Monthly; fields: `avg_change`
  (Seasonality Average Change), `max_change` (Seasonality Max Change), `median_change` (Seasonality
  Median Change), `min_change` (Seasonality Min Change), `month` (Seasonality Month Number),
  `positive_closes` (Seasonality Positive Closes Count), `positive_months_perc` (Seasonality
  Positive Months Percent), `years` (Seasonality Year Count)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/seasonality/market` — Market Seasonality

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SeasonalityController.market_seasonality)

**What it does:**

Returns the average return by month for the tickers SPY, QQQ, IWM, XLE, XLC, XLK, XLV, XLP, XLY,
XLRE, XLF, XLI, XLB .

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Seasonality Market` — Seasonality Market; fields: `avg_change`
  (Seasonality Average Change), `max_change` (Seasonality Max Change), `median_change` (Seasonality
  Median Change), `min_change` (Seasonality Min Change), `month` (Seasonality Month Number),
  `positive_closes` (Seasonality Positive Closes Count), `positive_months_perc` (Seasonality
  Positive Months Percent), `ticker` (Stock Ticker), `years` (Seasonality Year Count)
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/seasonality/{month}/performers` — Month Performers

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SeasonalityController.month_performers)

**What it does:**

Returns the tickers with the highest performance in terms of price change in the month over the
years. Per default the result is ordered by 'positive_months_perc' descending, then 'median_change'
descending, then 'marketcap' descending.

**Parameters**

- **path `month`** (required; type `SingleMonthNumber`; enum=[1, 2, 3, 4, 5, 6, 7, 8, …]; example=3)
  — A month number indicating the month, e.g. 1 -> January, 2 -> February, ...
- **query `min_years`** (optional; type `Min Years`; default=10; minimum=1; example=3) — The minimum
  amount of years data for the month need to be available for the ticker. Default: 10. Min: 1.
- **query `ticker_for_sector`** (optional; type `SingleTickerForSector`; example=AAPL) — A single
  ticker. The result will only contain tickers in the same sector as the given ticker, e.g. 'MSFT'
  will only yield result tickers in sector 'Technology'.
- **query `s_p_500_nasdaq_only`** (optional; type `Nasdaq Only`)
- **query `min_oi`** (optional; type `Min Open Interest`; minimum=0; example=10000) — The minimum
  open interest. Min: 0.
- **query `limit`** (optional; type `Seasonality Performers Limit`; default=50; minimum=1;
  example=10) — How many items to return. Default: 50. Min: 1.
- **query `order`** (optional; type `Seasonality Performance Order By`; enum=[month,
  positive_closes, years, positive_months_perc, median_change, avg_change, max_change, min_change];
  example=ticker) — Optional columns to order the result by
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.

**Response payload**

- `200`: `application/json` → `Seasonality Performers` — Seasonality Performers; fields:
  `avg_change` (Seasonality Average Change), `marketcap` (Stock Marketcap AUM), `max_change`
  (Seasonality Max Change), `median_change` (Seasonality Median Change), `min_change` (Seasonality
  Min Change), `month` (Seasonality Month Number), `positive_closes` (Seasonality Positive Closes
  Count), `positive_months_perc` (Seasonality Positive Months Percent), `sector` (Market General
  Sector), `ticker` (Stock Ticker), `years` (Seasonality Year Count)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/seasonality/{ticker}/year-month` — Price change per month per year

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SeasonalityController.year_month)

**What it does:**

Returns the relative price change for all past months over multiple years.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Seasonality Year Month` — Seasonality Year Month; fields: `change`
  (Seasonality Change), `close` (Seasonality Close Price), `month` (Seasonality Month Number),
  `open` (Seasonality Open Price), `year` (Seasonality Year)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Short interest — 7

#### `GET /api/shorts/{ticker}/ftds` — Failures to Deliver

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ShortController.failures_to_deliver)

**What it does:**

Returns the short failures to deliver per day for the given ticker starting from the given date. If
no date is given, returns the data for the current/last market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Failures to Deliver` — Failures to Deliver; fields: `date` (date),
  `price` (Stock Close Price), `quantity` (string)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/shorts/{ticker}/data` — Short Data

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ShortController.short_data)

**What it does:**

Returns short data including rebate rate and short shares available for a ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `newer_than`** (optional; type `NewerThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no older results will be returned. Can be used with `older_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `older_than`** (optional; type `OlderThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no newer results will be returned. Can be used with `newer_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).

**Response payload**

- `200`: `application/json` → `Short Data` — Short Data; fields: `currency` (string), `fee_rate`
  (string), `name` (string), `rebate_rate` (string), `short_shares_available` (integer), `symbol`
  (string), `timestamp` (string (date_time))
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/short_screener` — Short Screener

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ShortController.short_screener)

**What it does:**

Returns short interest and float data for percentage calculations based off search params. This
endpoint provides information about the percentage of float that is shorted, the float size, and the
days to cover metric.

**Parameters**

- **query `tickers`** (optional; type `Ticker`; example=AAPL,INTC) — A comma separated list of
  tickers. To exclude certain tickers prefix the first ticker with a `-`.
- **query `limit`** (optional; type `DynLimitAfczisxafgezlfuwzepbhiwqbukvmdtgddaknczwtmanpqjuaaa`;
  default=100; minimum=0; maximum=500; example=125) — The limit of results. Default: 100. Max: 500.
  Min: 0.
- **query `offset`** (optional; type `DynLimitAmgvwdedprykkkwubboqoucydelwfunvmsgjwjjjzoktwvuggur`;
  default=0; minimum=0; example=100) — The offset of results. Default: 0. Min: 0.
- **query `min_short_interest`** (optional; type
  `DynLimitAwvnhnanzzxgxlnzynfeguzfrgjiiyczvdojvmbhcwxcajcrnid`; example=5000) — The min short
  interest.
- **query `max_short_interest`** (optional; type
  `DynLimitAaiarpkmqeltxpeikbejqwzovkgqrujryjtyfdorrxiqrlvwhsq`; example=100000) — The max short
  interest.
- **query `min_days_to_cover`** (optional; type
  `DynLimitAmxmdyrsybvtrzmcmrakrncwtrghafdvtxplwnvfxotlemwclyu`; example=1.1) — The min days to
  cover.
- **query `max_days_to_cover`** (optional; type
  `DynLimitAxpqibvjazjdihecinkbfzmdoygmtcdoppivmncwxlmauoodfri`; example=5.3) — The max days to
  cover.
- **query `min_si_float`** (optional; type
  `DynLimitAmiwzvcvlafbptywihmivtogxztamfplnstoilspyiqonfmgxis`; example=0.1) — The min short
  interest float ratio.
- **query `max_si_float`** (optional; type
  `DynLimitAnbnwbvghyjjamoaeliijkuqyrcuxtbuhunosqyizueyzjfwpet`; example=0.05) — The max short
  interest float ratio.
- **query `min_si_float_with_synth_long_pct_of_total_shares`** (optional; type
  `DynLimitArpfwrvvgapdircxpncyosfvlwcxiiwymysekvyxfcgkugfgsuv`; example=0.05) — The min
  si_float_with_synth_long_pct_of_total_shares.
- **query `max_si_float_with_synth_long_pct_of_total_shares`** (optional; type
  `DynLimitAiqpvlkbsmwtkdsuqytestyjpjsjgxgdzhocfxyomeutzvqtyax`; example=0.05) — The max
  si_float_with_synth_long_pct_of_total_shares.
- **query `min_total_float`** (optional; type
  `DynLimitAyelzaduwvxwoezfgredydrnihuyxezmypesolmohfvsexmnink`; example=10000000) — The min total
  float.
- **query `max_total_float`** (optional; type
  `DynLimitAxjheprcsljmvplbcjldheerzymjvlsybquzhatrknvmtdlzowe`; example=100000000) — The max total
  float.
- **query `order_by`** (optional; type `ShortScreenerOrderBy`; example=market_date) — ordering
  options.
- **query `order_direction`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc];
  example=asc) — Whether to sort descending or ascending. Descending by default.
- **query `min_market_date`** (optional; type `Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD.
- **query `max_market_date`** (optional; type `Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD.
- **query `min_fee_rate`** (optional; type
  `DynLimitApkrxcwicmouljbhmfmrltarukskjwfyuwwrelhciynqydpizqa`; example=0.01) — The min fee rate.
- **query `max_fee_rate`** (optional; type
  `DynLimitAdqyxibqjfhjmzzebfpxvixnofkvffvlemjqutahaihzuarylmn`; example=0.45) — The max fee rate.
- **query `min_rebate_rate`** (optional; type
  `DynLimitApdnwejavjdeghjpywbfjnbuehgzwmislusryzphyunaxnoiglr`; example=0.11) — The min rebate
  rate.
- **query `max_rebate_rate`** (optional; type
  `DynLimitAopcmddgvncuyuibnwdxqqvwzgorwokcmttygcreolepxwinlnk`; example=5.34) — The max rebate
  rate.
- **query `min_short_shares_available`** (optional; type
  `DynLimitAwjdnnnqekdxstvfbyqrdsstoysnoqnykdginwxdwamkskkivkj`; example=0) — The min short shares
  available.
- **query `max_short_shares_available`** (optional; type
  `DynLimitAgrozoonsqobedhiajcjfotmdzknpovygauiiwzttzhsqtlwemf`; example=10000) — The max short
  shares available.

**Response payload**

- `200`: `application/json` → `Short Screener` — Short Screener; fields: `days_to_cover` (string),
  `fee_rate` (date), `market_date` (date), `rebate_rate` (date), `short_interest` (integer),
  `short_shares_available` (date), `si_float` (string),
  `si_float_with_synth_long_pct_of_total_shares` (string), `symbol` (string), `total_float`
  (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/shorts/{ticker}/volumes-by-exchange` — Short Volume By Exchange

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ShortController.short_volume_by_exchange)

**What it does:**

Returns short volume data broken down by exchange for a ticker. If no date range is given the
default is the last 2 years. The max lookback range is 5 years; requests for a wider range are
clamped to the most recent 5 years ending at `older_than`. The `newer_than` and `older_than` query
params accept Unix-format timestamps (milliseconds or seconds) as well as ISO-format dates
(`YYYY-MM-DD`).

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `newer_than`** (optional; type `NewerThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no older results will be returned. Can be used with `older_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).
- **query `older_than`** (optional; type `OlderThan`; example=1715083417) — The unix time in
  milliseconds or seconds at which no newer results will be returned. Can be used with `newer_than`
  to paginate by time. Also accepts an ISO date or RFC 3339 datetime (example: 2024-01-25).

**Response payload**

- `200`: `application/json` → `Short Volume By Exchange` — Short Volume By Exchange; fields: `date`
  (date), `exchange_name` (string), `market_center` (string), `short_volume` (integer),
  `total_volume` (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/shorts/{ticker}/volume-and-ratio` — Short Volume and Ratio

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ShortController.short_volume_and_ratio)

**What it does:**

Returns short volume and short ratio data for a ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Short Volume` — Short Volume; fields: `close_price` (Stock Close
  Price), `market_date` (date), `short_volume` (string), `short_volume_ratio` (string),
  `total_volume` (string)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/shorts/{ticker}/interest-float` — V1 Short Interest and Float (Deprecated)

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ShortController.short_interest_and_float)

**What it does:**

(This endpoint has been deprecated, use V2 now) Returns short interest and float data for percentage
calculations for a ticker. This endpoint provides information about the percentage of float that is
shorted, the float size, and the days to cover metric.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `V1 Short Interest and Float (Deprecated)` — V1 Short Interest and
  Float (Deprecated); fields: `created_at` (string (date_time)), `days_to_cover_returned` (string),
  `market_date` (date), `percent_returned` (string), `si_float_returned` (integer), `symbol`
  (string), `total_float_returned` (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/shorts/{ticker}/interest-float/v2` — V2 Short Interest and Float

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.ShortController.short_interest_and_float_v2)

**What it does:**

Returns short interest, float size, days-to-cover, and related percentage calculations for a given
ticker. Broker-dealers report short interest to FINRA twice per month then FINRA publishes the
aggregated results approximately one week after each submission window closes. Because updates
arrive on this fixed schedule, data in the response from this endpoint can lag the current date by
several weeks. This is a FINRA reporting constraint, not a data delivery issue or bug. This endpoint
updates twice a month and pulls data from FINRA. The update schedule can be viewed on the FINRA
website here

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `V2 Short Interest and Float` — V2 Short Interest and Float; fields:
  `days_to_cover` (string), `fee_rate` (date), `market_date` (date), `rebate_rate` (date),
  `short_interest` (integer), `short_shares_available` (date), `si_float` (string),
  `si_float_with_synth_long_pct_of_total_shares` (string), `symbol` (string), `total_float`
  (integer)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Stock data — 38

#### `GET /api/stock/{ticker}/atm-chains` — ATM Chains

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.atm_chains)

**What it does:**

Returns the ATM chains for the given expirations

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `expirations[]`** (required; type `Expiry dates`; example=[2024-02-02, 2024-01-26]) — An
  array of 1 or more expiry dates.

**Response payload**

- `200`: `application/json` → `Option Contract Screener response.` — Option Contract Screener
  response.; fields: `ask_side_volume` (Option Contract Ask Volume), `avg_price` (Option Contract
  Avg Price), `bid_side_volume` (Option Contract Bid Volume), `chain_prev_close` (Option Contract
  Previous Close Price), `close` (Option Contract Close), `cross_volume` (Option Contract Cross
  Volume), `er_time` (Stock Earnings time), `floor_volume` (Option Contract Floor Volume), `high`
  (Option Contract High), `last_fill` (Option Contract Last Transaction Time), `low` (Option
  Contract Low), `mid_volume` (Option Contract Mid Volume), `multileg_volume` (Option Contract Multi
  Leg Volume), `next_earnings_date` (Stock Next Earnings Date), `no_side_volume` (Option Contract No
  Side Volume), `open` (Option Contract Open), `open_interest` (Option Contract Open interest),
  `option_symbol` (Option Contract Symbol), `premium` (Option Contract Premium), `sector` (Market
  General Sector), `stock_multi_leg_volume` (Option Contract Stock Multi Leg Volume), `stock_price`
  (Stock Close Price), `sweep_volume` (Option Contract Sweep Volume), `ticker_vol` (Stock Total
  Volume), `total_ask_changes` (Option Contract Total Ask Changes), `total_bid_changes` (Option
  Contract Total Bid Changes), `trades` (Option Contract Total Trades Count), `volume` (Option
  Contract Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/balance-sheets` — Balance Sheets

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AvFundamentalController.balance_sheets)

**What it does:**

Returns balance sheet data including assets, liabilities, equity, debt structure, intangibles,
goodwill, receivables, inventory, and shares outstanding.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Balance Sheets Response` — Balance Sheets Response; fields: `data`
  (array<Balance Sheet>); data item `Balance Sheet` fields:
  `accumulated_depreciation_amortization_ppe` (string), `capital_lease_obligations` (string),
  `cash_and_cash_equivalents` (string), `cash_and_short_term_investments` (string), `common_stock`
  (string), `common_stock_shares_outstanding` (string), `current_accounts_payable` (string),
  `current_debt` (string), `current_long_term_debt` (string), `current_net_receivables` (string),
  `deferred_revenue` (string), `fiscal_date_ending` (string (date)), `goodwill` (string),
  `inserted_at` (string (date-time)), `intangible_assets` (string),
  `intangible_assets_excluding_goodwill` (string), `inventory` (string), `investments` (string),
  `long_term_debt` (string), `long_term_debt_noncurrent` (string), `long_term_investments` (string),
  `other_current_assets` (string), `other_current_liabilities` (string), `other_non_current_assets`
  (string), `other_non_current_liabilities` (string), `property_plant_equipment` (string),
  `report_type` (string), `reported_currency` (string), `retained_earnings` (string),
  `short_long_term_debt_total` (string), `short_term_debt` (string), `short_term_investments`
  (string), `ticker` (string), `total_assets` (string), `total_current_assets` (string),
  `total_current_liabilities` (string), `total_liabilities` (string), `total_non_current_assets`
  (string), `total_non_current_liabilities` (string), `total_shareholder_equity` (string),
  `treasury_stock` (string), `updated_at` (string (date-time)) — Balance Sheets
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/net-prem-ticks` — Call/Put Net/Vol Ticks

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.net_prem_ticks)

**What it does:**

Returns the net premium ticks for a given ticker which can be used to build the following chart:
---- Each tick is resembling the data for a single minute tick. To build a daily chart you would
have to add the previous data to the current tick:

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Net Prem Tick response.` — Net Prem Tick response.; fields:
  `call_volume` (Market General Call Volume), `call_volume_ask_side` (Market General Call Volume Ask
  Side), `call_volume_bid_side` (Market General Call Volume Bid Side), `date` (Market General
  Trading day), `net_call_premium` (Market General Net Call Premium), `net_call_volume` (Market
  General Net Call Volume), `net_delta` (string), `net_put_premium` (Market General Net Put
  Premium), `net_put_volume` (Market General Net Put Volume), `put_volume` (Market General Put
  Volume), `put_volume_ask_side` (Market General Put Volume Ask Side), `put_volume_bid_side` (Market
  General Put Volume Bid Side), `tape_time` (General Tick time)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/cash-flows` — Cash Flow Statements

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AvFundamentalController.cash_flows)

**What it does:**

Returns cash flow data including operating, investing, and financing cashflows, capital
expenditures, dividend payouts, stock buybacks, and stock-based compensation.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Cash Flows Response` — Cash Flows Response; fields: `data`
  (array<Cash Flow Statement>); data item `Cash Flow Statement` fields: `capital_expenditures`
  (string), `cashflow_from_financing` (string), `cashflow_from_investment` (string),
  `change_in_cash_and_cash_equivalents` (string), `change_in_exchange_rate` (string),
  `change_in_inventory` (string), `change_in_operating_assets` (string),
  `change_in_operating_liabilities` (string), `change_in_receivables` (string),
  `depreciation_depletion_and_amortization` (string), `dividend_payout` (string),
  `dividend_payout_common_stock` (string), `dividend_payout_preferred_stock` (string),
  `fiscal_date_ending` (string (date)), `inserted_at` (string (date-time)), `net_income` (string),
  `operating_cashflow` (string), `payments_for_operating_activities` (string),
  `payments_for_repurchase_of_common_stock` (string), `payments_for_repurchase_of_equity` (string),
  `payments_for_repurchase_of_preferred_stock` (string), `proceeds_from_issuance_of_common_stock`
  (string), `proceeds_from_issuance_of_long_term_debt` (string),
  `proceeds_from_issuance_of_preferred_stock` (string), `proceeds_from_operating_activities`
  (string), `proceeds_from_repayments_of_short_term_debt` (string),
  `proceeds_from_repurchase_of_equity` (string), `proceeds_from_sale_of_treasury_stock` (string),
  `profit_loss` (string), `report_type` (string), `reported_currency` (string),
  `stock_based_compensation` (string), `ticker` (string), `updated_at` (string (date-time)) — Cash
  Flow Statements
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{sector}/tickers` — Companies in Sector

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.companies_in_sector)

**What it does:**

Returns a list of tickers which are in the given sector.

**Parameters**

- **path `sector`** (required; type `Sector`; enum=[Basic Materials, Communication Services,
  Consumer Cyclical, Consumer Defensive, Energy, Financial Services, Healthcare, Industrials, …];
  example=Technology) — A financial sector.

**Response payload**

- `200`: `application/json` → `Market Sector Tickers` — Market Sector Tickers; fields: array of
  unspecified
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/earnings` — Earnings History

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AvFundamentalController.earnings)

**What it does:**

Returns earnings data including reported EPS, estimated EPS, surprise amount, surprise percentage,
and report timing (pre/post market).

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Earnings Response` — Earnings Response; fields: `data`
  (array<Earnings Report>); data item `Earnings Report` fields: `estimated_eps` (string),
  `fiscal_date_ending` (string (date)), `inserted_at` (string (date-time)), `report_date` (string
  (date)), `report_time` (string), `report_type` (string), `reported_eps` (string), `surprise`
  (string), `surprise_percentage` (string), `ticker` (string), `updated_at` (string (date-time)) —
  Earnings History
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/flow-alerts` — Flow Alerts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.flow_alerts)

**What it does:**

This endpoint has been deprecated and will be removed. Please migrate to this Flow Alerts endpoint,
which provides a more detailed response:
https://api.unusualwhales.com/docs#/operations/PublicApi.OptionTradeController.flow_alerts

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `limit`** (optional; type `Default 100 Max 200 Min 1`; default=100; minimum=1;
  maximum=200; example=10) — How many items to return. Default: 100. Max: 200. Min: 1.
- **query `is_ask_side`** (optional; type
  `DynLimitAbfzzygkzgahgkmwimqsjlmzmfuwqmuqrzsmnjgltforgqucjva`; default=true; example=true) —
  Boolean flag whether a transaction is ask side.
- **query `is_bid_side`** (optional; type
  `DynLimitAlofwxyfeqzdkubghajqfovcqmodyvslygzpzeqbpptoodjtrsh`; default=true; example=true) —
  Boolean flag whether a transaction is bid side.

**Response payload**

- `200`: `application/json` → `Flow Alert` — Flow Alert; fields: `alert_rule` (Alert Rule Name),
  `all_opening_trades` (Option Contract All Opening Trades), `created_at` (General UTC Timestamp),
  `expiry` (Option Contract Expiry), `expiry_count` (Option Contract Expiry Count), `has_floor`
  (Option Contract Has Floor), `has_multileg` (Single Trade Has Multileg), `has_singleleg` (Single
  Trade Is Single Leg), `has_sweep` (Single Trade Is Sweep), `issue_type` (Stock Issue Type),
  `open_interest` (ToBeDone), `option_chain` (Option Contract Symbol), `price` (ToBeDone), `strike`
  (Option Contract Strike), `ticker` (ToBeDone), `total_ask_side_prem` (ToBeDone),
  `total_bid_side_prem` (ToBeDone), `total_premium` (ToBeDone), `total_size` (ToBeDone),
  `trade_count` (ToBeDone), `type` (Option Contract Type), `underlying_price` (ToBeDone), `volume`
  (ToBeDone), `volume_oi_ratio` (ToBeDone)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/flow-per-expiry` — Flow per expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.flow_per_expiry)

**What it does:**

Returns the option flow per expiry for the last trading day

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Flow per expiry` — Flow per expiry; fields: `call_otm_premium`
  (Market General Call OTM Premium), `call_otm_trades` (Market General Call OTM Trades),
  `call_otm_volume` (Market General Call OTM Volume), `call_premium` (Market General Call Premium),
  `call_premium_ask_side` (Market General Call Premium Ask Side), `call_premium_bid_side` (Market
  General Call Premium Bid Side), `call_trades` (Market General Call Trades), `call_volume` (Market
  General Call Volume), `call_volume_ask_side` (Market General Call Volume Ask Side),
  `call_volume_bid_side` (Market General Call Volume Bid Side), `date` (Market General Trading day),
  `expiry` (Option Contract Expiry), `put_otm_premium` (Market General Put OTM Premium),
  `put_otm_trades` (Market General Put OTM Trades), `put_otm_volume` (Market General Put OTM
  Volume), `put_premium` (Market General Put Premium), `put_premium_ask_side` (Market General Put
  Premium Ask Side), `put_premium_bid_side` (Market General Put Premium Bid Side), `put_trades`
  (Market General Put Trades), `put_volume` (Market General Put Volume), `put_volume_ask_side`
  (Market General Put Volume Ask Side), `put_volume_bid_side` (Market General Put Volume Bid Side),
  `ticker` (Stock Ticker)

#### `GET /api/stock/{ticker}/flow-per-strike` — Flow per strike

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.flow_per_strike)

**What it does:**

Returns the option flow per strike for a given trading day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Flow per strike` — Flow per strike; fields: `call_premium` (Market
  General Call Premium), `call_premium_ask_side` (Market General Call Premium Ask Side),
  `call_premium_bid_side` (Market General Call Premium Bid Side), `call_trades` (Market General Call
  Trades), `call_volume` (Market General Call Volume), `call_volume_ask_side` (Market General Call
  Volume Ask Side), `call_volume_bid_side` (Market General Call Volume Bid Side), `date` (Market
  General Trading day), `put_premium` (Market General Put Premium), `put_premium_ask_side` (Market
  General Put Premium Ask Side), `put_premium_bid_side` (Market General Put Premium Bid Side),
  `put_trades` (Market General Put Trades), `put_volume` (Market General Put Volume),
  `put_volume_ask_side` (Market General Put Volume Ask Side), `put_volume_bid_side` (Market General
  Put Volume Bid Side), `strike` (Option Contract Strike), `ticker` (Stock Ticker), `timestamp`
  (General UTC Timestamp)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/flow-per-strike-intraday` — Flow per strike intraday

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.flow_per_strike_intraday)

**What it does:**

Returns the options flow for a given date in one minute intervals (the one minute intervals are not
aggregated with each other).

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `filter`** (optional; type `Filter`; default=NetPremium; enum=[NetPremium, Volume,
  Trades]; example=Volume) — Retrieve the strikes with the highest filter parameter.

**Response payload**

- `200`: `application/json` → `Flow per strike` — Flow per strike; fields: `call_premium` (Market
  General Call Premium), `call_premium_ask_side` (Market General Call Premium Ask Side),
  `call_premium_bid_side` (Market General Call Premium Bid Side), `call_trades` (Market General Call
  Trades), `call_volume` (Market General Call Volume), `call_volume_ask_side` (Market General Call
  Volume Ask Side), `call_volume_bid_side` (Market General Call Volume Bid Side), `date` (Market
  General Trading day), `put_premium` (Market General Put Premium), `put_premium_ask_side` (Market
  General Put Premium Ask Side), `put_premium_bid_side` (Market General Put Premium Bid Side),
  `put_trades` (Market General Put Trades), `put_volume` (Market General Put Volume),
  `put_volume_ask_side` (Market General Put Volume Ask Side), `put_volume_bid_side` (Market General
  Put Volume Bid Side), `strike` (Option Contract Strike), `ticker` (Stock Ticker), `timestamp`
  (General UTC Timestamp)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/financials` — Full Financials

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AvFundamentalController.financials)

**What it does:**

Returns full financial data for the given ticker, including income statements, balance sheets, cash
flows, and earnings. Supports both annual and quarterly data.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Full Financials Response` — Full Financials Response; fields: `data`
  (Full Financials) — Full Financials
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/fundamental-breakdown` — Fundamental Breakdown

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.FundamentalController.show)

**What it does:**

Returns the fundamental financial data for the given ticker, including earnings per share, revenue,
dividends, share counts, RSU data, and revenue breakdowns by product and geography.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Fundamental Breakdown Response` — Fundamental Breakdown Response;
  fields: `data` (Fundamental Breakdown) — Fundamental Breakdown
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/greeks` — Greeks

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.greeks)

**What it does:**

Returns the greeks for each strike for a single expiry date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `expiry`** (required; type `Single expiry date`; example=2024-02-02) — A single expiry
  date in ISO date format.

**Response payload**

- `200`: `application/json` → `Greeks` — Greeks; fields: `call_charm` (Charm), `call_delta` (Delta),
  `call_gamma` (Gamma), `call_rho` (Rho), `call_theta` (Theta), `call_vanna` (Vanna), `call_vega`
  (Vega), `call_volatility` (Implied Volatility), `date` (General ISO Date), `expiry` (General ISO
  Date), `put_charm` (Charm), `put_delta` (Delta), `put_gamma` (Gamma), `put_rho` (Rho), `put_theta`
  (Theta), `put_vanna` (Vanna), `put_vega` (Vega), `put_volatility` (Implied Volatility), `strike`
  (Strike)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/historical-risk-reversal-skew` — Historical Risk Reversal Skew

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.historical_risk_reversal_skew)

**What it does:**

Returns the historical risk reversal skew (the difference between put and call volatility) at a
delta of 25 or 10 (colloquial for 0.25 or 0.1) for a given expiry date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `expiry`** (required; type `Single expiry date`; example=2024-02-02) — A single expiry
  date in ISO date format.
- **query `timeframe`** (optional; type `Time frame`; default=1Y; example=2M) — The timeframe of the
  data to return. Can be one of the following formats:
  - YTD
  - 1D, 2D, etc.
  - 1W, 2W, etc.
  - 1M, 2M, etc.
  - 1Y, 2Y, etc.
- **query `delta`** (required; type `Delta`; example=0.610546281537814) — The delta of the option
  trade.

**Response payload**

- `200`: `application/json` → `Historical Risk Reversal Skew` — Historical Risk Reversal Skew;
  fields: `date` (General ISO Date), `delta` (Risk Reversal Delta), `risk_reversal` (Risk Reversal),
  `ticker` (Stock Ticker)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/iv-rank` — IV Rank

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.iv_rank)

**What it does:**

Returns the IV rank data for a ticker over a period of time. IV rank is a measure of where current
implied volatility stands relative to its historical range.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `timespan`** (optional; type `OptionalTimespan`; example=1y) — Optional timespan parameter
  that can be used to specify a date range. Accepts values like '1y', '6m', '3m', '1m', '1w' or an
  ISO date (example: 2024-01-25).

**Response payload**

- `200`: `application/json` → `IV Rank` — IV Rank; fields: `close` (Stock Close Price), `date`
  (Market General Trading day), `iv_rank_1y` (Market General IV Rank), `updated_at` (General UTC
  Timestamp), `volatility` (Market General Volatility)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/volatility/term-structure` — Implied Volatility Term Structure

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.implied_volatility_term_structure)

**What it does:**

The average of the latest volatilities for the at the money call and put contracts for every expiry
date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Implied Volatility Term Structure` — Implied Volatility Term
  Structure; fields: `date` (Market General Trading day), `dte` (DTE), `expiry` (Stock Expiry),
  `implied_move` (Stock Implied Move), `implied_move_perc` (Stock Implied Move Perc), `volatility`
  (Stock Volatility)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/income-statements` — Income Statements

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AvFundamentalController.income_statements)

**What it does:**

Returns income statement data including revenue, gross profit, operating income, EBIT, EBITDA, net
income, R&D, SG&A, depreciation, interest, and tax expenses.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Income Statements Response` — Income Statements Response; fields:
  `data` (array<Income Statement>); data item `Income Statement` fields:
  `comprehensive_income_net_of_tax` (string), `cost_of_goods_and_services_sold` (string),
  `cost_of_revenue` (string), `depreciation` (string), `depreciation_and_amortization` (string),
  `ebit` (string), `ebitda` (string), `fiscal_date_ending` (string (date)), `gross_profit` (string),
  `income_before_tax` (string), `income_tax_expense` (string), `inserted_at` (string (date-time)),
  `interest_and_debt_expense` (string), `interest_expense` (string), `interest_income` (string),
  `investment_income_net` (string), `net_income` (string), `net_income_from_continuing_operations`
  (string), `net_interest_income` (string), `non_interest_income` (string), `operating_expenses`
  (string), `operating_income` (string), `other_non_operating_income` (string), `report_type`
  (string), `reported_currency` (string), `research_and_development` (string),
  `selling_general_and_administrative` (string), `ticker` (string), `total_revenue` (string),
  `updated_at` (string (date-time)) — Income Statements
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/info` — Information

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.info)

**What it does:**

Returns a information about the given ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Ticker Info` — Ticker Info; fields: `announce_time` (Stock Earnings
  time), `avg30_volume` (Stock Average 30 Day Volume), `beta` (Stock Beta), `full_name` (Stock Full
  Name), `has_dividend` (Stock Has Dividend), `has_earnings_history` (Stock Has Earnings History),
  `has_investment_arm` (Stock Has Investment Arm), `has_options` (Stock Has Options), `issue_type`
  (Stock Issue Type), `logo` (Logo), `marketcap` (Stock Marketcap AUM), `next_earnings_date` (Stock
  Next Earnings Date), `sector` (Market General Sector), `short_description` (Stock Short
  Description)
- `404`: `application/json` → `Error Message stating that the requested element was not found
  causing an empty result to be generated.` — Error Message stating that the requested element was
  not found causing an empty result to be generated.; fields: Error Message stating that the
  requested element was not found causing an empty result to be generated. — Not Found
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/insider-buy-sells` — Insider buy & sells

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.insider_buy_sell)

**What it does:**

Returns the total amount of purchases & sells as well as notional values for insider transactions
for the given ticker

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Insider statistics` — Insider statistics; fields: `filing_date`
  (Insider Trades Filing Date), `purchases` (Insider Trades Purchases), `purchases_notional`
  (Insider Trades Notional Purchases), `sells` (Insider Trades Sells), `sells_notional` (Insider
  Trades Notional Sells)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/interpolated-iv` — Interpolated IV

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.interpolated_iv)

**What it does:**

Returns the Interpolated IV for a given trading day. If there is no expiration then the data is
calcualted via linear interpolation with the next 2 closest expirations Date must be the current or
a past date. If no date is given, returns data for the current/last market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Interpolated IV` — Interpolated IV; fields: `date` (Market General
  Trading day), `days` (integer), `implied_move_perc` (string), `percentile` (string), `volatility`
  (string)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/max-pain` — Max Pain

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.max_pain)

**What it does:**

Returns the max pain for all expirations for the given ticker for the last 120 days

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Max Pain` — Max Pain; fields: `expiry` (Option Contract Expiry),
  `max_pain` (string)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/nope` — Nope

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.nope)

**What it does:**

Returns the tickers NOPE for the given market day broken down per minute. NOPE is the Net Options
Pricing Effect, which tracks the intraday net delta of any ticker, but most research has been done
on indexes. It functions under 2 assumptions:

1. MM's take short side of any call or put traded during the day
2. MM's try to minimize risk by dynamically hedging their delta-gamma exposure, and do so by
   buying/shorting the underlying stock in proportion to the total net delta being tradedBased on
   these assumptions, options trading in large amounts (re: very liquid tickers) can potentially
   drive the price of the underlying, to a certain extent. Large movements might exacerbate this
   real time hedging, and drive price movements further in respective directions. In short, NOPE
   represents a best-estimate of expected number of shares to be hedged at any given time, and will
   show a general expected direction on the underlying The original NOPE calculation was based on
   the following formula: `NOPE = (Call Delta

- Put Delta) / Stock Volume` where call/put delta is obtained by multiplying each chains volume with
  its latest delta and then summing those values up. `NOPE fill` on the other hand is based on the
  delta at the time of the transaction Date must be the current or a past date. If no date is given,
  returns data for the current/last market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Nope` — Nope; fields: `call_delta` (Call Delta), `call_fill_delta`
  (Call Fill Delta), `call_vol` (Call Volume), `nope` (Nope score), `nope_fill` (Nope Fill),
  `put_delta` (Put Delta), `put_fill_delta` (Put Fill Delta), `put_vol` (Put Volume), `stock_vol`
  (Stock Volume), `timestamp` (Timestamp)

#### `GET /api/stock/{ticker}/ohlc/{candle_size}` — OHLC

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.ohlc)

**What it does:**

Returns the Open High Low Close (OHLC) candle data for a given ticker. Results are limited to 2,500
elements even if there are more available. Note: If you select 1d or 1w as a candle_size then the
candles won't have a start & end time. For 1w, the `date` field is the Monday of the ISO week. Note:
Suppose you enter end_date value 2024-11-25 which was a Monday. Your response will include 1-2 hours
of data from Tuesday 2024-11-26 due to UTC date rollover. Rest-assured, the response data covers the
full trading day (based on Eastern time) according to your entered end_date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **path `candle_size`** (required; type `string`; enum=[1m, 5m, 10m, 15m, 30m, 1h, 4h, 1d, …])
- **query `timeframe`** (optional; type `Time frame`; default=1Y; example=2M) — The timeframe of the
  data to return. Can be one of the following formats:
  - YTD
  - 1D, 2D, etc.
  - 1W, 2W, etc.
  - 1M, 2M, etc.
  - 1Y, 2Y, etc.
- **query `end_date`** (optional; type `OHLC End Date`; example=2024-01-18) — A trading date in the
  format of YYYY-MM-DD. This is the end date for the given timeframe. If you set the timeframe to 1
  minute ticks and have an end_date of 2024-01-18 then it would start searching backwards from
  2024-01-18
- **query `date`** (optional; type `Market Date`; example=2024-01-18) — A trading date in the format
  of YYYY-MM-DD.
- **query `limit`** (optional; type `Max 2500 Min 1`; minimum=1; maximum=2500; example=10) — How
  many items to return. Max: 2500. Min: 1.

**Response payload**

- `200`: `application/json` → `Candle data` — Candle data; fields: `close` (Candle Close),
  `end_time` (Candle End time), `high` (Candle High), `low` (Candle Low), `market_time` (Market
  General Market Time), `open` (Candle Open), `start_time` (Candle Start time), `total_volume`
  (Candle Total volume), `volume` (Candle Volume) — OHLC Response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/oi-change` — OI Change

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.oi_change)

**What it does:**

Returns the tickers contracts' OI change data ordered by absolute OI change (default: descending).
Date must be the current or a past date. If no date is given, returns data for the current/last
market day.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `limit`** (optional; type `Min Limit 1`; minimum=1; example=10) — How many items to
  return. If no limit is given, returns all matching data. Min: 1.
- **query `page`** (optional; type `Page`; example=1) — Page number (use with limit). Starts on page
  0.
- **query `order`** (optional; type `OrderDirection`; default=desc; enum=[desc, asc]; example=asc) —
  Whether to sort descending or ascending. Descending by default.

**Response payload**

- `200`: `application/json` → `OI Change` — OI Change; fields: `avg_price` (ToBeDone), `curr_date`
  (General ISO Date), `curr_oi` (ToBeDone), `last_ask` (ToBeDone), `last_bid` (ToBeDone),
  `last_date` (General ISO Date), `last_fill` (ToBeDone), `last_oi` (ToBeDone), `oi_change`
  (ToBeDone), `oi_diff_plain` (ToBeDone), `option_symbol` (Option Contract Symbol),
  `percentage_of_total` (ToBeDone), `prev_ask_volume` (ToBeDone), `prev_bid_volume` (ToBeDone),
  `prev_mid_volume` (ToBeDone), `prev_multi_leg_volume` (ToBeDone), `prev_neutral_volume`
  (ToBeDone), `prev_stock_multi_leg_volume` (ToBeDone), `prev_total_premium` (ToBeDone), `rnk`
  (ToBeDone), `trades` (ToBeDone), `underlying_symbol` (ToBeDone), `volume` (ToBeDone)

#### `GET /api/stock/{ticker}/oi-per-expiry` — OI per Expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.oi_per_expiry)

**What it does:**

Returns the total open interest for calls and puts for a specific expiry date.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Open Interest per Expiry` — Open Interest per Expiry; fields:
  `call_oi` (Total Call Open interest), `date` (Market General Trading day), `expiry` (Option
  Contract Expiry), `put_oi` (Total Put Open interest) — OI per expiry response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/oi-per-strike` — OI per Strike

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.oi_per_strike)

**What it does:**

Returns the total open interest for calls and puts for a specific strike.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Open Interest per Strike` — Open Interest per Strike; fields:
  `call_oi` (Total Call Open interest), `date` (Market General Trading day), `put_oi` (Total Put
  Open interest), `strike` (Option Contract Strike) — OI per strike response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/stock-volume-price-levels` — Off/Lit Price Levels

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.stock_volume_price_level)

**What it does:**

Returns the lit & off lit stock volume per price level for the given ticker. ---- Important: The
volume does **NOT** represent the full market dialy volume. It only represents the volume of
executed trades on exchanges operated by Nasdaq and FINRA off lit exchanges.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Off Lit Price Levels` — Off Lit Price Levels; fields: `data`
  (array<Off Lit Price Level>); data item `Off Lit Price Level` fields: `lit_vol` (Stock Lit
  Volume), `off_vol` (Stock Off Lit Volume), `price` (number (double))
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/option-chains` — Option Chains

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.option_chains)

**What it does:**

Returns all option symbols for the given ticker that were present at the given day. If no date is
given, returns data for the current/last market day. You can use the following regex to extract
underlying ticker, option type, expiry & strike:
`^(?<symbol>[\w]*)(?<expiry>(\d{2})(\d{2})(\d{2}))(?<type>[PC])(?<strike>\d{8})$` Keep in mind that
the strike needs to be divided by 1,000.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `greeks`** (optional; type `boolean`) — When true, return an enriched row per contract
  (strike, expiry, type, NBBO, IV, OI, volume, delta, gamma, theta, vega, rho, last_tape_time)
  instead of the default array of option-symbol strings.

**Response payload**

- `200`: `application/json` → `Option Chains response` — Option Chains response; fields: sample
  keys: `data` — Option chains response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/option/stock-price-levels` — Option Price Levels

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.option_price_level)

**What it does:**

Returns the call and put volume per price level for the given ticker. ---- Can be used to build a
chart such as following:

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Option Price Level` — Option Price Level; fields: `call_volume`
  (Market General Call Volume), `price` (Stock Price Level), `put_volume` (Market General Put
  Volume)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/options-volume` — Options Volume

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.options_volume)

**What it does:**

Returns the options volume & premium for all trade executions that happened on a given trading date
for the given ticker. ---- This can be used to build a ticker options overview such as: ----

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `limit`** (optional; type `Default 1 Max 500 Min 1`; default=1; minimum=1; maximum=500;
  example=10) — How many items to return. Default: 1. Max: 500. Min: 1.

**Response payload**

- `200`: `application/json` → `Ticker Options Volume` — Ticker Options Volume; fields:
  `avg_30_day_call_volume` (Market General Avg 30 Day Call Volume), `avg_30_day_put_volume` (Market
  General Avg 30 Day Put Volume), `avg_3_day_call_volume` (Market General Avg 3 Day Call Volume),
  `avg_3_day_put_volume` (Market General Avg 3 Day Put Volume), `avg_7_day_call_volume` (Market
  General Avg 7 Day Call Volume), `avg_7_day_put_volume` (Market General Avg 7 Day Put Volume),
  `bearish_premium` (Market General Bearish Premium), `bullish_premium` (Market General Bullish
  Premium), `call_open_interest` (Market General Call Open Interest), `call_premium` (Market General
  Call Premium), `call_volume` (Market General Call Volume), `call_volume_ask_side` (Market General
  Call Volume Ask Side), `call_volume_bid_side` (Market General Call Volume Bid Side), `date`
  (Market General Trading day), `net_call_premium` (Market General Net Call Premium),
  `net_put_premium` (Market General Net Put Premium), `put_open_interest` (Market General Put Open
  Interest), `put_premium` (Market General Put Premium), `put_volume` (Market General Put Volume),
  `put_volume_ask_side` (Market General Put Volume Ask Side), `put_volume_bid_side` (Market General
  Put Volume Bid Side)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/ownership` — Ownership

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.ownership)

**What it does:**

Returns the institutions, insider trades and politicians with the most shares. This is an enterprise
only endpoint. Contact dev@unusualwhales.com for details about accessing this data.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `limit`** (optional; type `Default 20 Max 100 Min 1`; default=20; minimum=1; maximum=100;
  example=20) — How many items to return. Default: 20. Max: 100. Min: 1.

**Response payload**

- `200`: `application/json` → `Ownership` — Ownership; fields: `cik` (Ownership CIK), `entity_type`
  (Ownership Entity Type), `name` (Ownership Name), `ownership` (Ownership Relative),
  `ownership_perc` (Ownership Percentage), `units` (Ownership Units)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/volatility/realized` — Realized Volatility

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.realized_volatility)

**What it does:**

The implied and realized volatility of a given ticker. The implied volatility is the expected 30 day
forward looking volatility. The realized/historical volatility is the volatility of the stock price
in the last 30 days. Since IV is forward looking, the realized volatility is shifted 30 days
backwards to see if the past IV pricings were frequently underpricing or overpricing the realized
volatility risk.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `timeframe`** (optional; type `Time frame`; default=1Y; example=2M) — The timeframe of the
  data to return. Can be one of the following formats:
  - YTD
  - 1D, 2D, etc.
  - 1W, 2W, etc.
  - 1M, 2M, etc.
  - 1Y, 2Y, etc.

**Response payload**

- `200`: `application/json` → `Realized Volatility` — Realized Volatility; fields: `date` (Market
  General Trading day), `implied_volatility` (Stock IV 30d), `price` (Stock Close Price),
  `realized_volatility` (Realized Stock Volatility), `unshifted_rv_date` (Unshifted RV Date)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/flow-recent` — Recent flows

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.flow_recent)

**What it does:**

Returns the latest flows for the given ticker. Optionally a min premium and a side can be supplied
in the query for further filtering.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `side`** (optional; type `Side`; default=ALL; enum=[ALL, ASK, BID, MID]; example=ASK) —
  The side of a stock trade. Must be one of ASK, BID, MID. If not set, will return all side's
  trades.
- **query `min_premium`** (optional; type `StockTradesMinPremium`; default=0; minimum=0;
  example=50000) — The minimum premium requested trades should have.

**Response payload**

- `200`: `application/json` → `Flow per expiry` — Flow per expiry; fields: `call_otm_premium`
  (Market General Call OTM Premium), `call_otm_trades` (Market General Call OTM Trades),
  `call_otm_volume` (Market General Call OTM Volume), `call_premium` (Market General Call Premium),
  `call_premium_ask_side` (Market General Call Premium Ask Side), `call_premium_bid_side` (Market
  General Call Premium Bid Side), `call_trades` (Market General Call Trades), `call_volume` (Market
  General Call Volume), `call_volume_ask_side` (Market General Call Volume Ask Side),
  `call_volume_bid_side` (Market General Call Volume Bid Side), `date` (Market General Trading day),
  `expiry` (Option Contract Expiry), `put_otm_premium` (Market General Put OTM Premium),
  `put_otm_trades` (Market General Put OTM Trades), `put_otm_volume` (Market General Put OTM
  Volume), `put_premium` (Market General Put Premium), `put_premium_ask_side` (Market General Put
  Premium Ask Side), `put_premium_bid_side` (Market General Put Premium Bid Side), `put_trades`
  (Market General Put Trades), `put_volume` (Market General Put Volume), `put_volume_ask_side`
  (Market General Put Volume Ask Side), `put_volume_bid_side` (Market General Put Volume Bid Side),
  `ticker` (Stock Ticker)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/quote` — Stock Quote

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.StockQuoteController.show)

**What it does:**

Returns the latest trade, national best bid and ask, derived quote values, and quote change
statistics for each market session.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Stock Quote Response` — Stock Quote Response; fields: `data` (Stock
  Quote) — Latest Stock Quote
- `404`: `application/json` → `Stock Quote Error Response` — Stock Quote Error Response; fields:
  `error` (object) — Ticker is unknown
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error
- `502`: `application/json` → `Stock Quote Error Response` — Stock Quote Error Response; fields:
  `error` (object) — Stock quote service unavailable

#### `GET /api/stock/{ticker}/stock-state` — Stock State

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.last_stock_state)

**What it does:**

Returns the last stock state for the given ticker. This is the easiest way to retreive the open,
close, high, low and volume of the last trading day. For the latest trade and national best bid and
ask use `/api/stock/{ticker}/quote`. It also returns the quote age, midpoint, spread, spread in
basis points, top of book size imbalance and size-weighted midpoint.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker

**Response payload**

- `200`: `application/json` → `Stock state` — Stock state; fields: `close` (Stock Close Price),
  `high` (Stock High Price), `low` (Stock Low Price), `market_time` (Market Time), `open` (Stock
  Open Price), `prev_close` (Stock Prev Close Price), `tape_time` (Stock Last Tape Time),
  `total_volume` (Cumulative Volume), `volume` (Stock Volume) — Stock State Response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/technical-indicator/{function}` — Technical Indicator

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.AvFundamentalController.technical_indicator)

**What it does:**

Returns any technical indicator time series for a ticker. Supports international stocks and OTC.
**Supported functions:** SMA, EMA, WMA, DEMA, TEMA, TRIMA, KAMA, MAMA, T3, MACD, MACDEXT, STOCH,
STOCHF, RSI, STOCHRSI, WILLR, ADX, ADXR, APO, PPO, MOM, BOP, CCI, CMO, ROC, ROCR, AROON, AROONOSC,
MFI, TRIX, ULTOSC, DX, MINUS_DI, PLUS_DI, MINUS_DM, PLUS_DM, BBANDS, MIDPOINT, MIDPRICE, SAR,
TRANGE, ATR, NATR, AD, ADOSC, OBV, HT_TRENDLINE, HT_SINE, HT_TRENDMODE, HT_DCPERIOD, HT_DCPHASE,
HT_PHASOR, VWAP. Not all parameters apply to every function. For example, STOCH and BOP only use
`interval`, while RSI uses all three (`interval`, `time_period`, `series_type`). The `month`
parameter is only relevant for intraday intervals (1min, 5min, 15min, 30min, 60min). For
daily/weekly/monthly intervals, full history is returned and `month` is ignored.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **path `function`** (required; type `string`) — The technical indicator function name (e.g. RSI,
  SMA, MACD, BBANDS)
- **query `interval`** (optional; type `string`; default=daily; enum=[1min, 5min, 15min, 30min,
  60min, daily, weekly, monthly]) — Time interval between data points
- **query `time_period`** (optional; type `integer`; default=14) — Number of data points used to
  calculate the indicator (e.g. 14 for RSI, 20 for SMA)
- **query `series_type`** (optional; type `string`; default=close; enum=[close, open, high, low]) —
  The price type used in the calculation
- **query `month`** (optional; type `string`) — Target month for intraday data in YYYY-MM format
  (e.g. 2026-02). Only applies to intraday intervals.

**Response payload**

- `200`: `application/json` → `Technical Indicator Response` — Technical Indicator Response; fields:
  `data` (array<Technical Indicator Point>); data item `Technical Indicator Point` fields: `date`
  (string (date)), `indicator` (string), `interval` (string), `series_type` (string), `ticker`
  (string), `time_period` (integer), `values` (object/map) — Technical Indicator
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/volatility/stats` — Volatility Statistics

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.volatility_stats)

**What it does:**

Returns comprehensive volatility statistics for a ticker on a specific date, including implied
volatility data, realized volatility data, and their respective high/low values for the past year.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Volatility Statistics` — Volatility Statistics; fields: `date`
  (Market General Trading day), `iv` (Market General Volatility), `iv_high` (Market General
  Volatility), `iv_low` (Market General Volatility), `iv_rank` (Market General IV Rank), `rv`
  (Market General Volatility), `rv_high` (Market General Volatility), `rv_low` (Market General
  Volatility), `ticker` (Stock Ticker)
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/option/volume-oi-expiry` — Volume & OI per Expiry

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.TickerController.vol_oi_per_expiry)

**What it does:**

Returns the total volume and open interest per expiry for the given ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Volume & OI per Expiry` — Volume & OI per Expiry; fields: `expires`
  (Stock Expiry), `oi` (Stock Open Interest), `volume` (Stock Options Volume) — Vol & OI per expiry
  response
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Stock directory — 1

#### `GET /api/stock-directory/ticker-exchanges` — Ticker Exchange Mapping

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.StockDirectoryController.ticker_exchanges)

**What it does:**

Returns a mapping of all tickers to their exchanges.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Ticker Exchanges` — Ticker Exchanges; fields: `data` (array<Ticker
  Exchange>); data item `Ticker Exchange` fields: `exchange` (string), `ticker` (string)
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Unusual congressional trades — 4

#### `GET /api/congress/unusual-trades` — Unusual Congressional Trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.UnusualTradesController.recent)

**What it does:**

Returns congressional trades that have been flagged as unusual, optionally filtered by reason tags.
Supported `types` values include `committee_conflict`, `first_person_to_trade`, `low_marketcap`,
`unusual_industry`, `unusually_large_trade`, and `fec_donation_conflict`. This is a premium
endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **query `types`** (optional; type `string`) — Comma-separated list of unusual-activity tags to
  filter by. Returns all unusual trades when omitted.
- **query `limit`** (optional; type `integer`; default=100; minimum=1; maximum=500) — Number of
  results per page (max 500).
- **query `page`** (optional; type `integer`; default=1; minimum=1) — Page number (1-indexed).

**Response payload**

- `200`: `application/json` → `Unusual Trades` — Unusual Trades; fields: `data` (array<Unusual
  Congressional Trade>); data item `Unusual Congressional Trade` fields: `amount` (string), `asset`
  (string), `asset_type` (string), `company_name` (string), `description` (string), `filing_date`
  (string (date)), `id` (string (uuid)), `politician_id` (string (uuid)), `politician_name`
  (string), `ticker` (string), `transaction_date` (string (date)), `transaction_type` (string),
  `unusual_activity_meta` (array<object>), `unusual_activity_tags` (array<string>) — Unusual Trades
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/congress/unusual-trades/stats` — Unusual Trades Aggregate Stats

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.UnusualTradesController.stats)

**What it does:**

Returns the most recent cached overview statistics: top 30 politicians by unusual trade count,
party/chamber breakdowns, committee and industry groupings, top tickers, biggest trades, and daily
trade volume series with SPY benchmark. This is a premium endpoint. Contact dev@unusualwhales.com to
request access.

**Parameters**

- None documented.

**Response payload**

- `200`: `application/json` → `Unusual Trades Stats` — Unusual Trades Stats; fields: `data`
  (object), `date` (string (date)), `type` (string) — Unusual Trades Stats
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/congress/unusual-trades/chart-data` — Unusual Trades Chart Data

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.UnusualTradesController.chart_data)

**What it does:**

Returns trade points (with price context) and SPY daily closes over the requested date range,
suitable for plotting congressional trade activity against the broader market. Defaults to the last
~4 months when no range is supplied. This is a premium endpoint. Contact dev@unusualwhales.com to
request access.

**Parameters**

- **query `date_from`** (optional; type `string (date)`; format=date) — Inclusive lower bound on
  transaction_date (ISO 8601).
- **query `date_to`** (optional; type `string (date)`; format=date) — Inclusive upper bound on
  transaction_date (ISO 8601).

**Response payload**

- `200`: `application/json` → `Unusual Trades Chart Data` — Unusual Trades Chart Data; fields:
  `spy_prices` (array<SPY Daily Close>), `trades` (array<Unusual Trade Chart Point>) — Unusual
  Trades Chart Data
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/congress/unusual-trades/by-tickers` — Unusual Trades by Ticker

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.UnusualTradesController.by_tickers)

**What it does:**

Returns congressional trades filtered by one or more tickers, with optional date range, transaction
type, and politician name filters. Includes price-at-trade and current-price context. This is a
premium endpoint. Contact dev@unusualwhales.com to request access.

**Parameters**

- **query `tickers`** (optional; type `string`) — Comma-separated list of tickers to filter by.
- **query `transaction_type`** (optional; type `string`; enum=[buy, sell]) — Filter to either buy-
  or sell-side transactions.
- **query `date_from`** (optional; type `string (date)`; format=date) — Inclusive lower bound on
  transaction_date (ISO 8601).
- **query `date_to`** (optional; type `string (date)`; format=date) — Inclusive upper bound on
  transaction_date (ISO 8601).
- **query `politician`** (optional; type `string`) — Case-insensitive substring match against the
  politician's full name.
- **query `limit`** (optional; type `integer`; default=100; minimum=1; maximum=500)
- **query `page`** (optional; type `integer`; default=1; minimum=1)

**Response payload**

- `200`: `application/json` → `Unusual Trades` — Unusual Trades; fields: `data` (array<Unusual
  Congressional Trade>); data item `Unusual Congressional Trade` fields: `amount` (string), `asset`
  (string), `asset_type` (string), `company_name` (string), `description` (string), `filing_date`
  (string (date)), `id` (string (uuid)), `politician_id` (string (uuid)), `politician_name`
  (string), `ticker` (string), `transaction_date` (string (date)), `transaction_type` (string),
  `unusual_activity_meta` (array<object>), `unusual_activity_tags` (array<string>) — Unusual Trades
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### Volatility — 6

#### `GET /api/volatility/anomaly/top` — Top Volatility Anomalies

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.VolatilityController.anomaly_top)

**What it does:**

Screener of the top volatility anomalies on a date, by `direction` (`short_vol` or `long_vol`).
Supports the same filters as the website screener.

**Parameters**

- **query `direction`** (required; type `string`) — `short_vol` or `long_vol`.
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `limit`** (optional; type `integer`) — Max rows (default 50, max 200).

**Response payload**

- `200`: `application/json` → `Volatility Result` — Volatility Result; fields: Volatility Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/volatility/character/top` — Top Volatility Character

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.VolatilityController.character_top)

**What it does:**

Screener of tickers by volatility character (mean-reverting / persistent / moderate) on a date.

**Parameters**

- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.
- **query `limit`** (optional; type `integer`) — Max rows (default 50, max 200).
- **query `sort`** (optional; type `string`) — `half_life` (default), `hurst`, or `neg_entropy`.
- **query `dir`** (optional; type `string`) — `asc` (default) or `desc`.

**Response payload**

- `200`: `application/json` → `Volatility Result` — Volatility Result; fields: Volatility Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/volatility/vix-term-structure` — VIX Term Structure

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.VolatilityController.vix_term_structure)

**What it does:**

The latest VIX futures term structure plus history. Requires the `volatility` API add-on.

**Parameters**

- **query `history_days`** (optional; type `integer`) — Days of history (default 90, max 365).

**Response payload**

- `200`: `application/json` → `Volatility Result` — Volatility Result; fields: Volatility Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/volatility/variance-risk-premium` — Variance Risk Premium

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.VolatilityController.variance_risk_premium)

**What it does:**

The variance risk premium (implied vs realized variance) history for a ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Volatility Result` — Volatility Result; fields: Volatility Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/volatility/anomaly` — Volatility Anomaly Score

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.VolatilityController.anomaly)

**What it does:**

The volatility anomaly score (and recent history) for a ticker — a composite signal flagging
unusually rich/cheap volatility.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Volatility Result` — Volatility Result; fields: Volatility Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

#### `GET /api/stock/{ticker}/volatility/character` — Volatility Character

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.VolatilityController.character)

**What it does:**

The volatility character (Hurst exponent / mean-reversion half-life / entropy) and history for a
ticker.

**Parameters**

- **path `ticker`** (required; type `SingleTicker`; example=AAPL) — A single ticker
- **query `date`** (optional; type `Optional Market Date`; example=2024-01-18) — A trading date in
  the format of YYYY-MM-DD. This is optional and by default the last trading date.

**Response payload**

- `200`: `application/json` → `Volatility Result` — Volatility Result; fields: Volatility Result
- `422`: `application/json` → `Error Message` — Error Message; fields: `msg` (string), `path`
  (string), `query` (string), `url` (string) — Unprocessable Entity
- `500`: `text/plain` → `Error Message on an internal server error.` — Error Message on an internal
  server error.; fields: Error Message on an internal server error. — Internal Server Error

### WebSocket data — 15

#### `GET /api/socket/contract_screener` — Contract screener

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.contract_screener)

**What it does:**

**NOTE:** This is the documentation for websocket channel `contract_screener`. Websocket access for
personal use is only available through the Advanced plan. The data on this channel is the exact same
data shown on https://unusualwhales.com/options-screener. You can find fully-functional examples
that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `contract_screener` for live hot-option-contract snapshots
  (Greeks, side volumes, OI growth indicators). Payload format:

##### Field reference

Each message is a snapshot of a single option contract that's currently active. Streamed
continuously throughout the session.

| Field | Type | Description |
| --- | --- | --- |
| `option_symbol` | string | OSI option chain id. |
| `tape_time` | ISO 8601 string | Snapshot time, UTC. |
| `volume` | int | Session volume on the contract. |
| `open_interest` | int | Open interest as of `tape_time`. |
| `prev_oi` | int | Open interest at the previous session close. Compare with `open_interest` to detect OI growth. |
| `trades` | int | Number of distinct trades in the session. |
| `premium` | decimal string | Total dollar premium traded today. |
| `ask_side_volume`, `bid_side_volume`, `neutral_volume` | int | Volume broken out by side classification |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/custom_alerts` — Custom alerts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.custom_alerts)

**What it does:**

**NOTE:** This is the documentation for websocket channel `custom_alerts`. Websocket access for
personal use is only available through the Advanced plan. Unlike the other websocket channels,
`custom_alerts` is a **per-user** stream. You only receive notifications that match the alert
configurations on the Unusual Whales account that owns the API token used to connect. Configure
alerts on https://unusualwhales.com/notifications; every notification fired against your account is
also delivered on this channel. You can find fully-functional examples that stream data from many
channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel: `custom_alerts`. Payload format: `meta` is a free-form object whose shape
  depends on `noti_type`.

##### Field reference

| Field | Type | Description |
| --- | --- | --- |
| `name` | string | The alert configuration's display name. |
| `noti_type` | string | The kind of alert that fired (e.g. `"flow_alerts"`, `"price_alert"`, `"news"`). Determines the shape of `meta`. |
| `user_noti_config_id` | uuid string | Identifier of the alert configuration that produced this notification. Use this to group alerts by the rule that fired them. |
| `user_id` | uuid string | Account that received the alert (always the account associated with the API token) |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/flow_alerts` — Flow alerts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.flow_alerts)

**What it does:**

**NOTE:** This is the documentation for websocket channel `flow-alerts`. Websocket access for
personal use is only available through the Advanced plan. You can find fully-functional examples
that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `flow-alerts` for all flow alerts. Payload format:

##### Field reference

A flow alert is an aggregate of one or more individual `option_trades` that matched a rule (sweep,
repeated hits, etc.). Use `trade_ids` to look up the underlying option trades.

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid string | Unique alert identifier. |
| `rule_id` | uuid string | Identifier of the rule that produced the alert. |
| `rule_name` | string \| null | Human-readable rule name (e.g. `"RepeatedHitsDescendingFill"`). May be null if the rule definition is not present. |
| `ticker` | string | Underlying ticker. |
| `option_chain` | string | OSI option chain id matched by the alert. |
| `underlying_price` | float | Spot price of the underlying at `end_time`. |
| `volume` | int | Aggregate session volume on the contract at `end_time`. |
| `total_size` | int | Sum of `size` across the trades in this alert. |
| `total_premium` | float | Sum of `premium` across the trades in this alert |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/gex` — GEX

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.gex)

**What it does:**

**NOTE:** This is the documentation for websocket channels `gex`, `gex:<TICKER>`, `gex_strike`,
`gex_strike:<TICKER>`, `gex_strike_expiry`, and `gex_strike_expiry:<TICKER>`. Websocket access for
personal use is only available through theAdvanced plan. You can find fully-functional examples that
stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream, for example `gex:SPY` for live GEX updates for SPY,
  `gex_strike:SPY` for strike-level GEX data, or `gex_strike_expiry:SPY` for strike and expiry level
  GEX data. Omit the ticker suffix to receive the corresponding updates for every ticker. If you are
  connecting to multiple ticker specific GEX channels and just want updates for all tickers, connect
  only to the corresponding global channel: `gex`, `gex_strike`, or `gex_strike_expiry`. Payload
  format: Format for `gex:<TICKER>`: Format for `gex_strike:<TICKER>`: Format for
  `gex_strike_expiry:<TICKER>`:

##### Field reference

`gex:<TICKER>` (ticker-aggregate) Per-ticker totals of dealer Greek exposures, expressed as the
dollar move per 1% move in the underlying.

| Field | Type | Description |
| --- | --- | --- |
| `ticker` | string |  |
| `timestamp` | int (ms) | Calculation time, unix epoch milliseconds |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/lit_trades` — Lit trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.lit_trades)

**What it does:**

**NOTE:** This is the documentation for websocket channel `lit_trades`. Websocket access for
personal use is only available through the Advanced plan. You can find fully-functional examples
that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `lit_trades` for all lit (exchange-based) trades. Payload
  format:

##### Field reference

| Field | Type | Description |
| --- | --- | --- |
| `symbol` | string | Underlying ticker. |
| `price` | decimal string | Trade price. |
| `size` | int | Trade size in shares. |
| `volume` | int | Cumulative session volume on the symbol as of this trade. |
| `type` | `"lit"` | Always `"lit"` on this channel. |
| `trade_settlement` | string \| null | Settlement type, e.g. `"regular"`, `"cash"`, `"next_day"`, `"seller"`. |
| `trade_code` | string \| null | Special trade code such as `"opening_print"`, `"closing_print"`, `"odd_lot"`, etc. `null` for unflagged trades. |
| `ext_hour_sold_codes` | string \| null | Extended-hours indicator, e.g. `"extended_hours_trade"`. `null` during regular session. |
| `sale_cond_codes` | string \| null | Special sale condition, e.g. `"cross_trade"`, `"derivatively_priced"`. |
| `executed_at` | ISO 8601 string | Trade execution timestamp on the exchange (UTC) |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/market_tide` — Market tide

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.market_tide)

**What it does:**

**NOTE:** This is the documentation for websocket channel `market_tide`. Websocket access for
personal use is only available through the Advanced plan. You can find fully-functional examples
that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `market_tide` for the live market-wide net call/put premium
  and volume aggregate. The OTM-only breakout is included in the same payload via the `otm_net_*`
  fields. Payload format:

##### Field reference

Aggregates roll up every option transaction across the market. `net_*` fields cover all options;
`otm_net_*` fields are restricted to out-of-the-money contracts (where dealer hedging tends to
dominate price action).

| Field | Type | Description |
| --- | --- | --- |
| `timestamp` | ISO 8601 string | Aggregation time, UTC. |
| `net_volume` | int | Net call volume minus net put volume across the entire market. |
| `net_call_premium` | decimal string | Net call premium (`ask_side - bid_side`) across all options, in dollars. |
| `net_put_premium` | decimal string | Net put premium (`ask_side - bid_side`) across all options, in dollars. |
| `otm_net_volume` | int | Same as `net_volume` but restricted to OTM contracts |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/net_flow` — Net flow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.net_flow)

**What it does:**

**NOTE:** This is the documentation for websocket channel `net_flow:<TICKER>`. Websocket access for
personal use is only available through the Advanced plan. You can find fully-functional examples
that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream, for example `net_flow:SPY` for live net call/put premium
  and volume aggregates for SPY. Payload format:

##### Field reference

| Field | Type | Description |
| --- | --- | --- |
| `ticker` | string | Underlying ticker. Same as the `:` suffix on the channel name. |
| `net_call_prem` | decimal string | Cumulative net call premium for the session (`ask_side - bid_side`) in dollars. |
| `net_call_vol` | int | Cumulative net call volume for the session (`ask_side - bid_side`). |
| `net_put_prem` | decimal string | Cumulative net put premium for the session (`ask_side - bid_side`) in dollars. |
| `net_put_vol` | int | Cumulative net put volume for the session (`ask_side - bid_side`). |
| `time` | int (ms) | Last update time, unix epoch milliseconds (UTC). |

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/news` — News

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.news)

**What it does:**

**NOTE:** This is the documentation for websocket channel `news`. Websocket access for personal use
is only available through the Advanced plan. You can find fully-functional examples that stream data
from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `news` for all live headline news. Payload format:

##### Field reference

| Field | Type | Description |
| --- | --- | --- |
| `headline` | string | The full headline text. For Truth Social posts this is the post body verbatim. |
| `timestamp` | ISO 8601 string | Time the headline was published, UTC. |
| `source` | string | One of `"social-media"` (Truth Social posts), `"aggregator"`, or any other upstream source string passed through verbatim. |
| `tickers` | string[] | Tickers tagged on the headline. Empty for Truth Social posts. |
| `is_trump_ts` | bool | `true` only when the source is Truth Social and the author is `@realDonaldTrump`. Use this to filter the highest-signal Truth Social subset. |

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/off_lit_trades` — Off-lit trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.off_lit_trades)

**What it does:**

**NOTE:** This is the documentation for websocket channel `off_lit_trades`. Websocket access for
personal use is only available through the Advanced plan. You can find fully-functional examples
that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `off_lit_trades` for all off-lit (dark pool) trades.
  Payload format: The schema is identical to `lit_trades` but `type` is always `"off_lit"`. Off-lit
  prints are reported to the consolidated tape via a TRF (typically with a small delay), so
  `trf_executed_at` is more frequently populated on this channel and may differ from `executed_at`.

##### Field reference

See the `lit_trades` field reference above for the full schema. The two channels share the same
payload shape.

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/option_trades` — Option trades

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.option_trades)

**What it does:**

**NOTE:** This is the documentation for websocket channels `option_trades` and
`option_trades:<TICKER>`. Websocket access for personal use is only available through the Advanced
plan. You can find fully-functional examples that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel(s) you wish to stream, for example `option_trades` for all tickers or
  `option_trades:TSLA` for TSLA transactions only. Payload format:

##### Field reference

| Field | Type | Description |
| --- | --- | --- |
| `id` | uuid string | Unique trade identifier. Use this to dedupe and to cross-reference `flow-alerts.trade_ids`. |
| `underlying_symbol` | string | Underlying ticker (e.g. `AAPL`, `SPX`). |
| `executed_at` | int (ms) | Trade execution time, unix epoch milliseconds (UTC). |
| `nbbo_bid` / `nbbo_ask` | decimal string | NBBO at the time of the trade. |
| `ewma_nbbo_bid` / `ewma_nbbo_ask` | decimal string | Exponentially-weighted moving average of the NBBO; smoother reference price for noisy quotes. |
| `size` | int | Number of contracts in this trade. |
| `price` | decimal string | Per-contract price. |
| `premium` | decimal string | Total dollar premium for the trade (`size * price * multiplier`). The multiplier is `100` for most contracts; `NANOS` uses `1` and `XSP` uses `10` |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/periscope` — Periscope

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.periscope)

**What it does:**

**NOTE:** This is the documentation for websocket channels `periscope` and `periscope:<TICKER>`.
Websocket access for personal use is only available through the Advanced plan. These channels are
not available on the enterprise and enterprise startup plans; joining them returns an error. You can
find fully-functional examples that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `periscope` for every index ticker at once, or
  `periscope:SPX` for a single ticker. Payload format: Format for `periscope:<TICKER>`: The global
  `periscope` channel streams the same per-row object for every index ticker, tagged `periscope`
  (the `ticker` field distinguishes them): One frame is delivered per (strike, expiry). Fields:

| Field | Type | Description |
| --- | --- | --- |
| `ticker` | string | Index ticker (SPX, VIX, XSP, NANOS). |
| `timestamp` | int (ms) | Snapshot time, unix epoch milliseconds. |
| `strike` | decimal string | Strike price. |
| `expiry` | date string `YYYY-MM-DD` | Contract expiration. |
| `gamma`, `charm`, `vanna` | decimal string | Market-maker greek exposure at this strike & expiry. |

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/price` — Price

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.price)

**What it does:**

**NOTE:** This is the documentation for websocket channels `price` and `price:<TICKER>`. Websocket
access for personal use is only available through theAdvanced plan. You can find fully-functional
examples that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream, for example `price:SPY` for live price updates for SPY or
  `price` for live price updates for every ticker. If you are connecting to multiple
  `price:<TICKER>` channels and just want prices for all tickers, connect only to the global `price`
  channel. Payload format: Format for `price:<TICKER>`: Format for `price`:

##### Field reference

| Field | Type | Description |
| --- | --- | --- |
| `ticker` | string | Ticker symbol. |
| `close` | decimal string | Last trade price. |
| `time` | int (ms) | Trade execution time, unix epoch milliseconds (UTC). |
| `vol` | int | Cumulative session volume on the underlying as of `time`. |

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/interval_flow` — Ticker Interval flow

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.interval_flow)

**What it does:**

**NOTE:** This is the documentation for websocket channel `interval_flow`. Websocket access for
personal use is only available through the Advanced plan. The data on this channel is the exact same
data shown on https://unusualwhales.com/ticker-interval-flow. Ticker Interval flow shows you summary
stats about option transactions for a given ticker in a 5min window. It shows the options volume,
transactions count, greek exposure, net premium, implied move change and many more stats. This is
very useful if you want to build an alerting system that alerts if there is a volume spike or some
other sort of spike in a ticker in a short time frame. You can find fully-functional examples that
stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `interval_flow` for live per-interval option flow
  statistics. Each message carries a single ticker in its `ticker` field; subscribers receive
  updates for all tickers across the channel. Payload format:

##### Field reference

Each message describes a 5-minute window of option flow for a single ticker. Two messages are
emitted per ticker per interval, distinguished by `interval_type`: one rolls up all contracts, the
other restricts to OTM contracts … [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket/trading_halts` — Trading halts

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.trading_halts)

**What it does:**

**NOTE:** This is the documentation for websocket channel `trading_halts`. Websocket access for
personal use is only available through the Advanced plan. You can find fully-functional examples
that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  Connect to the websocket URI: `wss://api.unusualwhales.com/socket?token=<YOUR_API_TOKEN>` then
  `join` the channel you wish to stream: `trading_halts` for live trading-state changes (halts,
  resumes, LULD pauses, etc.) sourced from UTP/CTA. Payload format:

##### Field reference

| Field | Type | Description |
| --- | --- | --- |
| `ticker` | string | Ticker whose trading state changed. |
| `state` | string | New trading state, e.g. `"halted"`, `"resumed"`, `"paused"`. Values are passed through from the upstream feed. |
| `reason` | string | Reason code from the upstream feed (e.g. `"LUDP"` for LULD volatility pause, `"T1"` for news pending, `""` if no reason was provided). |
| `time` | ISO 8601 string | When the state change occurred, UTC. |

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

#### `GET /api/socket` — WebSocket channels

[Official operation documentation](https://api.unusualwhales.com/docs/operations/PublicApi.SocketController.channels)

**What it does:**

Returns the available WebSocket channels for connections.

##### Websocket guide

You can find fully-functional examples that stream data from many channels here:

- Python:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output

- Javascript:
  https://github.com/unusual-whales/api-examples/tree/main/examples/ws-multi-channel-multi-output-nodejs
  If you are an AI or work with AI there is an available websocket skill at:
  https://unusualwhales.com/skills/websocket.md The following channels are available:

| Channel | Description |
| --- | --- |
| option_trades | Receive live option trades throughout the trading session. Expect 6-10M records per day. |
| option_trades:TICKER | Similar to `option_trades` but receive all trades only for the specified ticker. |
| flow-alerts | Receive live flow alerts (all of them unfiltered). This data can be used to build views like https://unusualwhales.com/option-flow-alerts. |
| price | Receive live price updates for every ticker at once. |
| price:TICKER | Receive live price updates for the given ticker. |
| news | Receive live headline news (Truth Social posts + aggregator news). |
| lit_trades | Receive live lit (exchange-based) trades throughout the trading session. |
| off_lit_trades | Receive live off-lit (dark pool) trades throughout the trading session. |
| gex | Receive live gex updates for every ticker at once |

… [see the official operation documentation for full notes]

**Parameters**

- None documented.

**Response payload**

- No HTTP response schema is documented; this route provides WebSocket channel guidance rather than
  a REST payload. Live frames are `[CHANNEL_NAME, PAYLOAD]`.

## Depth-2 crawl findings

- Historical downloadable datasets: options flow, OHLC/price data, insider trades, GEX, dark-pool
  data, Market Tide, IV Rank, and 13F filings.
- MCP access to 50+ financial-data tools.
- Kafka streams for option trades, insider trades, equities, and other enterprise feeds.
- Specialized linked skills for API-usage monitoring, institutional/13F data, WebSockets,
  earnings-volatility scanning, and a full-tape options data lake.

## Public free sources and practical substitutes

This section maps the inventory to public sources that can be queried, downloaded, or collected from public pages as of 2026-08-14. “Free” does not mean equivalent to Unusual Whales: several Unusual Whales products are derived from licensed exchange feeds, proprietary classifications, historical databases, or alert logic. The sources below are the closest public inputs for rebuilding or validating portions of the inventory.

Use APIs and bulk files before scraping rendered pages. Respect each provider’s rate limits, robots.txt, terms of use, attribution requirements, and data licenses. In particular, the SEC requires a descriptive User-Agent and fair-access behavior, and Cboe’s delayed-quote pages explicitly prohibit automated extraction.

### Coverage matrix

| Inventory area | Public/free source | Access and useful data | Coverage and important gap |
|---|---|---|---|
| Equity quotes, OHLC, adjusted prices, volume, technical indicators, movers, seasonality inputs, and correlations | [Alpha Vantage documentation](https://www.alphavantage.co/documentation/) | Free API key; time series, quotes, technical indicators, fundamentals, earnings, dividends, splits, news/sentiment, ETFs, commodities, FX, and some options endpoints. | Good starting point for delayed/historical equity analytics. Free-tier limits apply; realtime US market data and several deeper historical/options features may require a premium entitlement. |
| Exchange-facing quote and market-activity pages, option-chain display, IPO pages, and short-interest display | [Nasdaq market activity](https://www.nasdaq.com/market-activity), [Nasdaq option chain](https://www.nasdaq.com/market-activity/stocks/comm/option-chain), [Nasdaq Short Interest](https://www.nasdaqtrader.com/Trader.aspx?id=ShortInterest) | Public, human-facing pages; Nasdaq short interest is published by issue twice monthly and covers a rolling historical window. | Useful for manual verification and limited collection. These pages are not a stable, documented replacement for a consolidated quote, options-tape, or websocket API; check page terms before automation. |
| Company submissions, filing history, issuer metadata, XBRL fundamentals, reported shares, revenue, income, and balance-sheet facts | [SEC EDGAR APIs](https://www.sec.gov/search-filings/edgar-application-programming-interfaces), [data.sec.gov](https://data.sec.gov/) | Public JSON APIs without an API key. Useful endpoints include `https://data.sec.gov/submissions/CIK##########.json` and `https://data.sec.gov/api/xbrl/companyfacts/CIK##########.json`. SEC also publishes bulk archives. | Strong source for the inventory’s fundamentals, filings, corporate actions, and filing-linked research. It is as-filed/regulatory data, not a low-latency market-data feed. Use a descriptive User-Agent and throttle requests. |
| Insider transactions, Form 3/4/5, insider ownership, and filing-derived insider activity | [SEC insider transaction data sets](https://www.sec.gov/data-research/sec-markets-data/insider-transactions-data-sets), [SEC EDGAR filings](https://www.sec.gov/edgar/search/) | Quarterly downloadable data sets plus the underlying filings. | Closest public source for `/insider-trades`-style data. Normalize amendments, accession numbers, transaction codes, derivative transactions, and issuer CIKs yourself; there is no free source with Unusual Whales’ enrichment and alert labels. |
| Institutional holdings, 13F filings, manager portfolios, and ownership snapshots | [SEC EDGAR 13F filings](https://www.sec.gov/edgar/searchedgar/13F) and [SEC filing search](https://www.sec.gov/edgar/search/) | Public filings and SEC datasets/downloads. | Covers reported institutional holdings on the filing schedule, not live positions, options-flow intent, or Unusual Whales’ manager/position normalization. |
| Short interest and failed-to-deliver data | [Nasdaq short interest](https://www.nasdaqtrader.com/Trader.aspx?id=ShortInterest), [SEC fails-to-deliver data](https://www.sec.gov/data-research/sec-markets-data/fails-deliver-data) | Nasdaq short interest is twice monthly. SEC FTD files are downloadable twice monthly and include date, CUSIP, ticker, issuer, price, and aggregate quantity. | Public substitute for periodic short-interest/FTD routes. SEC FTD quantity is an aggregate outstanding balance for the reporting date, not a daily trade tape and not proof of naked short selling. |
| Daily short-sale volume, ATS volume, OTC transparency, and member-firm/ATS attribution | [FINRA OTC transparency](https://www.finra.org/filing-reporting/otc-transparency), [FINRA short-sale volume](https://www.finra.org/finra-data/browse-catalog/short-sale-volume), [FINRA daily short-sale transaction data](https://www.finra.org/finra-data/daily-short-sale-volume-transaction-data) | Public delayed/aggregated OTC and ATS reports; downloadable daily/monthly short-sale-volume files. FINRA describes the short-sale data as free for non-commercial use. | Useful for `/darkpool`, `/off-lit`, `/lit-flow`, and short-sale research at aggregate level. It does not provide the same consolidated, low-latency per-print dark-pool tape, participant labels, premium classification, or “whale” scoring as Unusual Whales. FINRA short-sale volume is not short interest. |
| Options chains, implied volatility, Greeks, put/call ratios, open interest, and volume/open-interest ratios | [Alpha Vantage options documentation](https://www.alphavantage.co/documentation/), [Nasdaq option chain](https://www.nasdaq.com/market-activity/stocks/comm/option-chain), [Cboe delayed quotes](https://www.cboe.com/delayed_quotes/cboe/advanced_charts/) | Alpha Vantage documents realtime option-chain, put/call-ratio, and volume/open-interest-ratio functions; Nasdaq and Cboe provide public delayed/manual displays. | Partial substitute only. A free, unrestricted equivalent for the full OPRA options tape, historical full-chain coverage, Unusual Whales flow classification, or its precomputed option analytics was not found. Cboe states that automated extraction/download from its delayed-quote pages is prohibited; do not scrape those pages. |
| GEX, net gamma, delta/charm/vanna exposure, IV rank, volatility term structure, and options anomaly metrics | Public option-chain/underlying data from [Alpha Vantage](https://www.alphavantage.co/documentation/) plus a local Black–Scholes/Greek calculator | Query chain snapshots where available, then calculate Greeks and aggregate by strike/expiry using open interest, volume, price, and an explicit assumption for dealer-side positioning. | These metrics can be approximated, but they are not directly observed. Results depend on contract selection, stale quotes, borrow/position assumptions, dealer-sign inference, and timestamp alignment. There is no public free source for Unusual Whales’ exact GEX, Periscope, Market Tide, NOPE, Volatility Character, or anomaly algorithms. |
| ETF profile, holdings, weights, issuer metadata, and exposure | [SEC EDGAR filing search](https://www.sec.gov/edgar/search/), ETF issuer holdings pages, and the ETF/profile functions in [Alpha Vantage](https://www.alphavantage.co/documentation/) | SEC filings and issuer files are public; Alpha Vantage offers normalized ETF profile/holdings functionality subject to key and plan limits. | Good for periodic holdings and reported exposures. Holdings, creation/redemption activity, intraday flows, and Unusual Whales’ exposure normalization are not available as one complete free feed. |
| Earnings dates, company events, dividends, splits, IPO calendar, and corporate-action research | [SEC EDGAR](https://www.sec.gov/edgar/search/), [Alpha Vantage](https://www.alphavantage.co/documentation/), [Nasdaq IPO activity](https://www.nasdaq.com/market-activity/ipos), and company investor-relations pages | Public filings and issuer calendars; Alpha Vantage provides several normalized event/fundamental endpoints with a free key. | Event timing is often available, but a consistent historical earnings-calendar API and Unusual Whales’ event normalization/alerts are not guaranteed for free. Prefer filings and issuer releases as the source of record. |
| News, filing headlines, social/news activity, and event-driven text inputs | [GDELT data](https://www.gdeltproject.org/data.html), [GDELT downloads](https://data.gdeltproject.org/), official publisher RSS feeds, and SEC EDGAR | GDELT publishes downloadable event/news datasets; publishers and agencies commonly expose RSS; SEC provides filing submissions and filing feeds. | Provides public text/event inputs for sentiment and headline monitoring. It is not the same feed as Unusual Whales’ aggregator, Truth Social coverage, deduplication, ticker tagging, or alert ranking. Verify publisher terms before collecting content. |
| Economic calendar and macro indicators: rates, inflation, employment, GDP, yields, and related series | [FRED API](https://fred.stlouisfed.org/docs/api/fred/v2/), [BLS Public Data API](https://www.bls.gov/developers/api_signature_v2.htm), [EIA Open Data API](https://www.eia.gov/developer/) | FRED requires a free API key. BLS offers a public API, with registration recommended for larger requests. EIA offers a free API key and bulk downloads that do not require a key. | Strong public replacements for macro/economic series used by calendar and market-context routes. They do not reproduce Unusual Whales’ calendar UI, event scoring, or notification system. |
| FDA approvals, drug events, and healthcare catalysts | [openFDA drug API](https://open.fda.gov/apis/drug/), [FDA approvals](https://www.fda.gov/news-events/approvals-fda-regulated-products), [Drugs@FDA](https://www.accessdata.fda.gov/scripts/cder/daf/) | Public JSON API, downloadable data, approval pages, and Drugs@FDA records. | Useful for approval/event research and building a catalyst calendar. There is no single free official API that exactly reproduces a PDUFA calendar with Unusual Whales’ dates, confidence, and ticker mapping; derive it from FDA records, company releases, and SEC filings. |
| Congressional/politician holdings and transactions | [House Clerk financial disclosure search](https://disclosures-clerk.house.gov/FinancialDisclosure/ViewSearch), [Senate financial disclosure search](https://efdsearch.senate.gov/search/home/), [Senate financial disclosure information](https://www.ethics.senate.gov/public/index.cfm/financialdisclosure) | Public reports, generally PDF/HTML/downloadable records, that can be collected and normalized by filer, asset, transaction date, and filing date. | Closest public source for congressional-trade routes. Reports are delayed, amended, inconsistently formatted, and may disclose ranges rather than exact amounts; they are not a live portfolio API. Respect each site’s use restrictions. |
| Futures, commodity prices, and commodity-market context | [Alpha Vantage commodities/futures documentation](https://www.alphavantage.co/documentation/), [EIA Open Data](https://www.eia.gov/developer/), [FRED](https://fred.stlouisfed.org/docs/api/fred/v2/) | Free-key or public bulk/API series for selected commodities, energy products, macro prices, and related indicators. | Useful for historical/context data, not a universal free CME/OPRA-style tick feed. Exchange-grade futures quotes, full depth, and low-latency history generally require licensed or paid market-data access. |
| Cryptocurrency OHLC, trades, order books, tickers, and market depth | [Binance public market-data API](https://github.com/binance/binance-spot-api-docs/blob/master/faqs/market_data_only.md), [Binance API Swagger](https://binance.github.io/binance-api-swagger/), [Coinbase Exchange public API](https://docs.cdp.coinbase.com/api-reference/exchange-api/rest-api/products/get-product-candles) | Public market-data endpoints do not require authentication for the documented public routes. Examples: `https://data-api.binance.vision/api/v3/klines`, `/trades`, `/ticker/24hr`, `/depth`; Coinbase exposes public products, candles, trades, and order-book routes. | Good replacements for crypto price/flow inputs. They cover individual venues, not all exchanges, and do not provide Unusual Whales’ cross-venue whale labels, wallet attribution, or derived alerts. |
| Foreign exchange rates and currency context | [Frankfurter API](https://www.frankfurter.dev/docs/), for example `https://api.frankfurter.dev/v1/latest?from=USD`, plus [Alpha Vantage FX](https://www.alphavantage.co/documentation/) | Frankfurter is a no-key public API for ECB reference rates; Alpha Vantage supplies free-key FX functions with plan/rate limits. | Frankfurter is daily reference data, not intraday FX ticks or dealer flow. Use Alpha Vantage or a licensed feed for higher-frequency series; neither recreates Unusual Whales’ derived FX analytics. |
| Prediction-market markets, outcomes, trades, positions, holders, open interest, and order books | [Polymarket API documentation](https://docs.polymarket.com/api-reference/introduction), [Gamma API](https://gamma-api.polymarket.com/), [Data API](https://data-api.polymarket.com/), [CLOB API](https://clob.polymarket.com/) | Polymarket’s public Gamma/Data APIs expose markets, events, tags, activity, trades, positions, holders, and related data; public CLOB routes expose prices, books, spreads, and history. Public endpoints do not require trading authentication. | This is a strong substitute for the Polymarket subset of the inventory only. It is not a universal prediction-market feed and does not reproduce Unusual Whales’ cross-venue normalization or scoring. |
| Private-company profiles, funding, valuations, and private-market activity | SEC EDGAR Form D/issuer filings, company investor-relations pages, and public announcements | Public filings and company pages can be queried or collected individually; normalize company identifiers and reported funding rounds locally. | No complete free source was identified for private-market valuations, management, funding history, or secondary pricing at Unusual Whales’ breadth. Treat search-engine snippets and aggregators as discovery aids, not authoritative records. |
| Alerts, screeners, query grammar, flow summaries, Market Tide, Options Pulse, and other derived analytics | Build locally from the public inputs above | Store raw filings, prices, options snapshots, FINRA aggregates, news, and event records; then calculate alerts and screeners with documented rules. | These are derived/proprietary products, not raw public datasets. A free source can supply inputs but cannot reproduce Unusual Whales’ exact alert thresholds, trade-side classification, participant labels, historical backfill, or proprietary feature engineering. |

### Practical source details and example requests

The following no-auth or free-key examples are useful smoke tests when building a public-data replacement layer:

```text
# SEC company submissions and XBRL facts (no API key; send a descriptive User-Agent)
https://data.sec.gov/submissions/CIK0000320193.json
https://data.sec.gov/api/xbrl/companyfacts/CIK0000320193.json

# Alpha Vantage demo request (free key; demo coverage/rate limits apply)
https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&apikey=demo

# Binance public candles (no authentication)
https://data-api.binance.vision/api/v3/klines?symbol=BTCUSDT&interval=1d&limit=1000

# Coinbase public candles (no authentication)
https://api.exchange.coinbase.com/products/BTC-USD/candles?granularity=86400

# Frankfurter daily reference FX rate (no authentication)
https://api.frankfurter.dev/v1/latest?from=USD

# Polymarket public market discovery (no authentication)
https://gamma-api.polymarket.com/markets?active=true&closed=false&limit=1
```

### What cannot be obtained as a genuinely free equivalent

The following inventory capabilities should be treated as partial-rebuild projects rather than simple source substitutions:

- Full OPRA option-trade history, low-latency option-flow websockets, and complete historical chains
  with reliable trade-side classification.
- Consolidated dark-pool/off-exchange per-print data with participant identity, premium labels, and
  the same-time lit/off-lit join logic.
- Unusual Whales’ GEX/Periscope, Market Tide, NOPE, Volatility Character, Options Pulse, alert
  rules, whale labels, and other derived analytics.
- Complete cross-venue crypto whale attribution, wallet identity, and private-market
  valuation/funding history.
- A single free API combining all equities, options, futures, crypto, FX, prediction markets,
  filings, congress, FDA, macro, and news data with stable schemas and historical backfill.

For an implementation, keep provenance on every record: source URL, retrieval timestamp, source identifier (CIK, accession number, ticker, contract, or market ID), licensing/terms classification, and whether the value is reported, aggregated, or locally derived. This makes it possible to compare a locally reconstructed metric with the corresponding Unusual Whales route without implying that the two datasets are identical.
