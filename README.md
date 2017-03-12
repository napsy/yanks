# yanks
A collector and agent for Go runtime stats

## Usage

1. obtain an API key
2. install the agent into your Go app and use the obtained API key
3. run your app

The app will send memory stats to the data collector every 10 seconds.

## What's Done

- agent: collecting memory stats and sending through TCP
- collector: receiving stats and parsing data

## Todo

- authentication using API keys
- write collected data into database
- data browsing web frontend
- use HTTPS instead of TCP
- integrate pprof code instead of ``os.Exec``
