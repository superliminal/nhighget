package nget

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	BASE_URL = "https://dojo.nplusplus.ninja/prod/steam"
	STEAM_ID = 76561198026910809
	READ_TIMEOUT = time.Millisecond * 2000
)

func GetScore(in <-chan *NScoresRequest) <-chan *NScoresResponse {
	out := make(chan *NScoresResponse)
	go func() {
		for r := range in {
			out <- handleRequest(r)
		}
		close(out)
	}()
	return out
}

func GetScores(in <-chan *NScoresRequest, getterCount int) <-chan *NScoresResponse {
	var outChs []<-chan *NScoresResponse
	for i := 0; i<= getterCount; i++ {
		outChs = append(outChs, GetScore(in))
	}
	return mergeChannels(outChs...)
}

func mergeChannels(cs... <-chan *NScoresResponse) <-chan *NScoresResponse {
	var wg sync.WaitGroup
	out := make(chan *NScoresResponse)

	// Starts a go-routine for the given channel to send its output to the out channel
	output := func(c <-chan *NScoresResponse) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	// Start output for each channel to be merged
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Wait for all in channels to finish sending before we close
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func handleRequest(request *NScoresRequest) *NScoresResponse {
	nScoresResponse := &NScoresResponse{Request: request}

	url := fmt.Sprintf("%v/get_scores?app_id=&steam_id=%v&user_id=&steam_auth=&qt=0&%v_id=%v", BASE_URL, STEAM_ID, request.scoreType, request.scoreId)

	resp, err := http.Get(url)
	if err != nil {
		nScoresResponse.Err = err
		return nScoresResponse
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		nScoresResponse.Err = err
		return nScoresResponse
	}
	if resp.StatusCode != 200 {
		nScoresResponse.Err = err
		return nScoresResponse
	}

	err = json.Unmarshal(body, nScoresResponse)
	if err != nil {
		nScoresResponse.Err = err
		return nScoresResponse
	}

	return nScoresResponse
}
