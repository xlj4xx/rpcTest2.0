package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	Endpoints []string `json:"endpoints"`
}

func main() {
	var config Config
	configFile, err := os.Open("settings.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		fmt.Println(err)
	}

	requestTimes := make(map[string][]float64)

	fmt.Println("Enter the amount of request")
	var numCalls int
	fmt.Println("(Public nodes: 100-150 // Private Node, 800-1000) : ")
	if _, err := fmt.Scan(&numCalls); err != nil {
		fmt.Println(err)
	}
	fmt.Println("")

	for _, endpointUrl := range config.Endpoints {
		requestTimes[endpointUrl] = []float64{}
		var client *ethclient.Client
		var err error

		if strings.HasPrefix(endpointUrl, "http") || strings.HasPrefix(endpointUrl, "ws") {
			client, err = ethclient.Dial(endpointUrl)
		} else {
			fmt.Printf("Unsupported endpoint: %s\n", endpointUrl)
			continue
		}

		if err != nil {
			fmt.Printf("Failed to connect to endpoint: %s\n", endpointUrl)
			continue
		}

		for i := 0; i < numCalls; i++ {
			startTime := time.Now()
			_, err := client.CodeAt(context.Background(), common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"), nil)
			endTime := time.Now()
			if err != nil {
				fmt.Printf("Failed to call get_code function at endpoint %s: %s\n", endpointUrl, err)
				continue
			}
			requestTimes[endpointUrl] = append(requestTimes[endpointUrl], endTime.Sub(startTime).Seconds())
			fmt.Printf("Request #: %d - endpoint %s: %.3f millisecondes.\r", i+1, endpointUrl, (endTime.Sub(startTime)).Seconds()*1000)
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}
		fmt.Println("")
	}

	fmt.Println("")
	fmt.Println("")
	for endpoint, times := range requestTimes {
		var totalTime float64
		for _, time := range times {
			totalTime += time
		}
		averageTime := totalTime / float64(len(times))

		fmt.Printf("Average latency time for %s is \033[38;2;0;255;0m%.3f\033[0m ms\n", endpoint, averageTime*1000)
	}

}
