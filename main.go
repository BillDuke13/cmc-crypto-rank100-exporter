package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	promAddr = ":8080"
	apiURL   = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
)

type cryptoData struct {
	Data []struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  struct {
			USD struct {
				Price float64 `json:"price"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

var (
	cryptoMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cryptocurrency_price_usd",
			Help: "Current price of cryptocurrencies in USD",
		},
		[]string{"name", "symbol"},
	)
)

func updateMetrics() {
	data, err := getCryptoData()
	if err != nil {
		log.Printf("Failed to retrieve crypto data: %v\n", err)
		return
	}

	var cryptos cryptoData
	if err := json.Unmarshal(data, &cryptos); err != nil {
		log.Printf("Failed to parse crypto data: %v\n", err)
		return
	}

	for _, crypto := range cryptos.Data {
		cryptoMetrics.With(prometheus.Labels{
			"name":   crypto.Name,
			"symbol": crypto.Symbol,
		}).Set(crypto.Quote.USD.Price)
	}
}

func getCryptoData() ([]byte, error) {
	apiKey := os.Args[1]
	if apiKey == "" {
		log.Fatal("CMC API Key not set: ")
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func main() {

	prometheus.MustRegister(cryptoMetrics)
	http.Handle("/metrics", promhttp.Handler())

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			updateMetrics()
			<-ticker.C
		}
	}()

	log.Fatal(http.ListenAndServe(promAddr, nil))
}
