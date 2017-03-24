package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/superliminal/nhighget/nget"
	"github.com/superliminal/nhighget/config"
	"github.com/superliminal/nhighget/db"
)

func main() {
	log.Println("Starting...")
	startTime := time.Now()
	toDb()
	runTime := time.Since(startTime)
	log.Printf("Complete: Update took %v seconds\n", runTime.Seconds())
}

func toDb() {
	flag.Parse()
	err := config.Init(*config.File)
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v\n", err)
		return
	}
	log.Printf("Ignore List: %v", config.Conf.IgnoreList)
	log.Printf("Ignore IDs: %v", config.Conf.IgnoreIds)
	log.Printf("Connection String: %v", config.Conf.Db.GetConnectionString())

	log.Println("Downloading Scores")
	scoresList := nget.GetAllScores()
	log.Println("Download Complete: Updating Database")

	updateDatabase(scoresList)
}

// TODO: Concurrent db updates, pipe in from score getter
func updateDatabase(scoresList []nget.NScoresResponse) {
	DB, err := db.GetDB(config.Conf.Db.GetConnectionString())
	if err != nil {
		log.Printf("Cound not connect to database: %v", err.Error())
		return
	}

	players, err := DB.GetAllPlayers()
	if err != nil {
		log.Printf("Cound not connect to database: %v", err.Error())
		return
	}
	playerMap := make(map[int]*db.Player)
	for _, player := range players {
		playerMap[player.Id] = &player
	}

	for _, scores := range scoresList {
		if scores.IsEpisodeScore() {
			updateEpisodeScores(DB, scores, playerMap)
		} else {
			updateLevelScores(DB, scores, playerMap)
		}
	}
}

// Updates the highscores for the given level
func updateLevelScores(s *db.DB, levelScores nget.NScoresResponse, players map[int]*db.Player) {

	// Used to offset ranks due to cheaters
	rankOffset := 0

	for _, highscore := range levelScores.Scores {
		// Find the player, check they aren't ignored
		player := players[highscore.UserId]
		if player == nil {
			// Player doesn't exist in db, create them
			player = s.CreatePlayer(highscore.UserId, highscore.UserName, isPlayerIgnored(highscore))
			players[player.Id] = player
		}
		if player.IsIgnored {
			rankOffset += 1
			continue
		}

		newScore := db.LevelHighscore{
			LevelId: *(levelScores.LevelId),
			PlayerId: highscore.UserId,
			Score: highscore.Score,
			Rank: highscore.Rank - rankOffset,
			PlayerName: highscore.UserName,
		}

		// Does this score already exists?
		dbScore, _ := s.GetLevelScore(newScore.LevelId, newScore.PlayerId)
		if s.NewRecord(dbScore) {
			s.Create(&newScore)
		} else {
			if (newScore.Score != dbScore.Score) {
				s.UpdateLevelScore(&dbScore, newScore.Score)
			}
		}
	}
}

// Updates the highscores for the given episode
func updateEpisodeScores(s *db.DB, episodeScores nget.NScoresResponse, players map[int]*db.Player) {

	// Used to offset ranks due to cheaters
	rankOffset := 0

	for _, highscore := range episodeScores.Scores {
		// Find the player, check they aren't ignored
		player := players[highscore.UserId]
		if player == nil {
			player = s.CreatePlayer(highscore.UserId, highscore.UserName, isPlayerIgnored(highscore))
			players[player.Id] = player
		}
		if player.IsIgnored {
			rankOffset += 1
			continue
		}

		newScore := db.EpisodeHighscore{
			EpisodeId: *(episodeScores.EpisodeId),
			PlayerId: highscore.UserId,
			Score: highscore.Score,
			Rank: highscore.Rank - rankOffset,
			PlayerName: highscore.UserName,
		}

		// Does this score already exists?
		dbScore, _ := s.GetEpisodeScore(newScore.EpisodeId, newScore.PlayerId)
		if s.NewRecord(dbScore) {
			s.Create(&newScore)
		} else {
			if (newScore.Score != dbScore.Score) {
				s.UpdateEpisodeScore(&dbScore, newScore.Score)
			}
		}
	}
}

// Check if this user's id/name match that of a cheater
func isPlayerIgnored(score nget.NScore) bool {
	for _, username := range config.Conf.IgnoreList {
		username = strings.ToLower(username)
		if strings.ToLower(score.UserName) == username {
			return true
		}
	}
	for _, userId := range config.Conf.IgnoreIds {
		if score.UserId == userId {
			return true
		}
	}
	return false
}