# CMC Crypto Rank 100 Exporter
A Prometheus exporter for rank 100 cryptocurrency from CoinMarketCap(CMC).

## Prerequisites

Before you begin, ensure that you have obtained your CoinMarketCap API key. If you haven't, go to the [official CoinMarketCap website](https://coinmarketcap.com/api/) to get it.

## Usage

1. Update the "Dockerfile" and replace `YOUR_CMC_API_KEY` with your actual CMC API key.

2. Build the Docker image with the following command:

```bash
docker build -t IMAGE_NAME . `
```

### Metrics Format

Metric name is cryptocurrency_price_usd, labels are name and symbol.

`cryptocurrency_price_usd{name="Bitcoin", symbol="BTC"} 30000`

### Contributing
Welcome to contribute via issues and PRs!