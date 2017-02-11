package nget

import (
	"log"
	"sort"
)

const (
	GETTER_COUNT = 100 // The number of worker go-routines to use
)

func GetAllScores() []NScoresResponse {
	var scores NScoresResponseList

	requests := getAllScoreRequests()
	requestCount := len(requests)

	requestChannel := getRequestChannel(requests)
	resultChan := GetScores(requestChannel, GETTER_COUNT)

	for i := 0; i < requestCount; i++ {
		resp := <-resultChan
		if resp.Err != nil {
			log.Printf("ERROR: %v", resp.Err)
			if resp.Request.attempts < 5 {
				resp.Request.attempts += 1
				requestChannel <- resp.Request
				i--
			} else {
				log.Print("ERROR: Max retries exceeded")
			}
			continue;
		}
		scores = append(scores, *resp)
	}
	close(requestChannel)
	sort.Sort(scores)
	return scores
}

func getAllScoreRequests() []*NScoresRequest {
	var reqSlice []*NScoresRequest
	reqSlice = append(reqSlice, getIntroEpisodeScoreRequests()...)
	reqSlice = append(reqSlice, getStandardEpisodeScoreRequests()...)
	reqSlice = append(reqSlice, getLegacyEpisodeScoreRequests()...)
	reqSlice = append(reqSlice, getIntroLevelScoreRequests()...)
	reqSlice = append(reqSlice, getStandardLevelScoreRequests()...)
	reqSlice = append(reqSlice, getLegacyLevelScoreRequests()...)
	reqSlice = append(reqSlice, getSecretLevelScoreRequests()...)
	return reqSlice
}

func getIntroEpisodeScoreRequests() []*NScoresRequest {
	return getScoreRequests("episode", 0, 25)
}

func getStandardEpisodeScoreRequests() []*NScoresRequest {
	return getScoreRequests("episode", 120, 120)
}

func getLegacyEpisodeScoreRequests() []*NScoresRequest {
	return getScoreRequests("episode", 240, 120)
}

func getIntroLevelScoreRequests() []*NScoresRequest {
	return getScoreRequests("level", 0, 125)
}

func getStandardLevelScoreRequests() []*NScoresRequest {
	return getScoreRequests("level", 600, 600)
}

func getLegacyLevelScoreRequests() []*NScoresRequest {
	return getScoreRequests("level", 1200, 600)
}

func getSecretLevelScoreRequests() []*NScoresRequest {
	return getScoreRequests("level", 1800, 120)
}

func getRequestChannel(requests []*NScoresRequest) chan *NScoresRequest {
	out := make(chan *NScoresRequest)
	go func() {
		for _, request := range requests {
			out <- request
		}
	}()
	return out
}

func getScoreRequests(scoreType string, startId, count int) []*NScoresRequest {
	var reqs []*NScoresRequest
	for i := startId; i < startId + count; i++ {
		reqs = append(reqs, NewNScoresRequest(scoreType, i))
	}
	return reqs
}