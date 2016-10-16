package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"
	"os"

	"nhighget/nget"
)

const TOTAL_EPISODE_COUNT = 265
const TOTAL_LEVEL_COUNT = 1445
const TOTAL_SCORE_COUNT = TOTAL_EPISODE_COUNT + TOTAL_LEVEL_COUNT

func main() {
	log.Println("Downloading Scores...")
	startTime := time.Now()
	run()
	runTime := time.Since(startTime)
	log.Printf("Complete: Score retrieval took %v seconds\n", runTime.Seconds())
}

func run() {
	fileName := flag.Arg(0)
	if fileName == "" {
		fileName = "n_scores.json"
	}
	var scores []nget.NScoresResponse
	resultChan := make(chan *nget.NScoresResponse)
	go nget.GetAllScores(resultChan)
	for i:=0; i<TOTAL_SCORE_COUNT; i++ {
		resp := <-resultChan
		if resp.Err != nil {
			log.Printf("ERROR: %v", resp.Err)
			continue;
		}
		scores = append(scores, *resp)
	}

	scoresJson, err := json.Marshal(scores)
	if err != nil {
		log.Printf("Error: Failed to encode the scores as json: %v", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(fileName, scoresJson, 0644)
	if err != nil {
		log.Printf("Error: Failed to write scores to file %v: %v", fileName, err)
		os.Exit(1)
	}
}
