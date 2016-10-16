package nget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

const BASE_URL = "https://dojo.nplusplus.ninja/prod/steam"
const STEAM_ID = 76561198026910809

const TOTAL_EPISODE_COUNT = 265
const TOTAL_LEVEL_COUNT = 1445
const TOTAL_SCORE_COUNT = TOTAL_EPISODE_COUNT + TOTAL_LEVEL_COUNT

type NScoresResponseList []NScoresResponse

func (l NScoresResponseList) Len() int { return len(l) }
func (l NScoresResponseList) Less(i, j int) bool {
	if l[i].EpisodeId != nil {
		if l[j].EpisodeId == nil {
			return true
		}
		return *l[i].EpisodeId < *l[j].EpisodeId
	}
	if l[j].EpisodeId != nil {
		return false
	}
	return *l[i].LevelId < *l[j].LevelId
}
func (l NScoresResponseList) Swap(i, j int){ l[i], l[j] = l[j], l[i] }

type NScoresResponse struct {
	EpisodeId *int `json:"episode_id,omitempty"`
	LevelId   *int `json:"level_id,omitempty"`
	Scores    []NScore `json:"scores"`
	queryType int `json:"query_type,omitempty"`
	Err       error `json:"-"`
}

type NScore struct {
	Rank     int `json:"rank"`
	Score    int `json:"score"`
	UserName string `json:"user_name"`
	userId   int `json:"user_id,omitempty"`
	replayId int `json:"replay_id,omitempty"`
}

type NUserInfo struct {
	MyDisplayName string `json:"my_display_name"`
	MyRank        int `json:"my_rank"`
	MyScore       int `json:"my_score"`
	MyReplayId    int `json:"my_replay_id"`
}

func GetAllScores() []NScoresResponse {
	var scores NScoresResponseList
	resultChan := make(chan *NScoresResponse)

	getAllScores(resultChan)

	for i:=0; i<TOTAL_SCORE_COUNT; i++ {
		resp := <-resultChan
		if resp.Err != nil {
			log.Printf("ERROR: %v", resp.Err)
			continue;
		}
		scores = append(scores, *resp)
	}
	sort.Sort(scores)
	return scores
}

func getAllScores(resultChan chan *NScoresResponse) {
	go GetIntroEpisodeScores(resultChan)
	go GetStandardEpisodeScores(resultChan)
	go GetLegacyEpisodeScores(resultChan)
	go GetIntroLevelScores(resultChan)
	go GetStandardLevelScores(resultChan)
	go GetLegacyLevelScores(resultChan)
	go GetSecretLevelScores(resultChan)
}

func GetIntroEpisodeScores(resultChan chan *NScoresResponse) {
	getScores(resultChan, "episode", 0, 25)
}

func GetStandardEpisodeScores(resultChan chan *NScoresResponse) {
	getScores(resultChan, "episode", 120, 120)
}

func GetLegacyEpisodeScores(resultChan chan *NScoresResponse) {
	getScores(resultChan, "episode", 240, 120)
}

func GetIntroLevelScores(resultChan chan *NScoresResponse) {
	getScores(resultChan, "level", 0, 125)
}

func GetStandardLevelScores(resultChan chan *NScoresResponse) {
	getScores(resultChan, "level", 600, 600)
}

func GetLegacyLevelScores(resultChan chan *NScoresResponse) {
	getScores(resultChan, "level", 1200, 600)
}

func GetSecretLevelScores(resultChan chan *NScoresResponse) {
	getScores(resultChan, "level", 1800, 120)
}


func getScores(resultChan chan *NScoresResponse, scoreType string, startId, count int) {
	for i:=startId; i<startId+count; i++ {
		go getScore(resultChan, scoreType, i)
		time.Sleep(time.Millisecond * 50) // Metanet server was returning internal server errors when I sent the requests any quicker
	}
}

func getScore(resultChan chan *NScoresResponse, scoreType string, id int) {
	url := fmt.Sprintf("%v/get_scores?app_id=&steam_id=%v&user_id=&steam_auth=&qt=0&%v_id=%v", BASE_URL, STEAM_ID, scoreType, id)

	resp, err := http.Get(url)
	if err != nil {
		resultChan <- &NScoresResponse{Err: err}
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultChan <- &NScoresResponse{Err: err}
		return
	}
	if resp.StatusCode != 200 {
		resultChan <- &NScoresResponse{Err: err}
		return
	}

	nScoresResponse := &NScoresResponse{}
	err = json.Unmarshal(body, nScoresResponse)
	if err != nil {
		resultChan <- &NScoresResponse{Err: err}
		return
	}

	resultChan <- nScoresResponse
}
