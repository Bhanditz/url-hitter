package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)
import "github.com/BurntSushi/toml"

type Config struct {
	BaseURL     string
	SitePaths   []string
	Duration    int
	DefaultHits int
}

func ReadConfig() Config {
	var configfile = "hiturl.toml"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func hitURL(url string, period int, hits int, wg *sync.WaitGroup) {
	sleepTime := time.Duration(period * 1000000000 / hits) //milliseconds
	fmt.Println("Sleep time(ms): ", sleepTime)
	for hits > 0 {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error: ", err)
			wg.Done()
			return
		}

		fmt.Printf("Hit: %d, URL: %s, Status: %s\n", hits, url, resp.Status)
		resp.Body.Close()
		hits--
		time.Sleep(sleepTime)
	}
	wg.Done()
}

func main() {
	var config = ReadConfig()

	if len(config.SitePaths) == 0 {
		fmt.Println("Error: no site parts loaded.")
		return
	}
	var wg sync.WaitGroup
	for _, sitePath := range config.SitePaths {
		url := config.BaseURL + sitePath
		wg.Add(1)
		go hitURL(url, config.Duration, config.DefaultHits, &wg)
	}
	wg.Wait()
}
