package main

import (
"encoding/json"
"flag"
"io/ioutil"
"log"
"os"
"strings"
"time"

"github.com/superliminal/nhighget/nget"
)

func main() {
	log.Println("Starting...")
	startTime := time.Now()
	run()
	runTime := time.Since(startTime)
	log.Printf("Complete: Score retrieval took %v seconds\n", runTime.Seconds())
}

func run() {
	flag.Parse()
	fileName := strings.TrimSpace(flag.Arg(0))
	if fileName == "" {
		fileName = "n-scores.json"
	}

	log.Println("Downloading Scores")
	scores := nget.GetAllScores()

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
	log.Printf("Scores saved to %v", fileName)
}
