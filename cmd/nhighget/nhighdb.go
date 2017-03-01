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
	log.Printf("%v", config.Conf.IgnoreList)
	log.Printf("%v", config.Conf.IgnoreIds)

	log.Println("Downloading Scores")
	scoresList := nget.GetAllScores()
	log.Println("Download Complete: Updating Database")

	DB, err := db.GetDB(config.Conf.Db.GetConnectionString())
	if err != nil {
		log.Printf("Cound not connect to database: %v", err.Error())
		return
	}

	// TODO: Concurrent db updates
	for _, scores := range scoresList {
		if scores.IsEpisodeScore() {
			updateEpisodeScores(DB, scores)
		} else {
			updateLevelScores(DB, scores)
		}
	}
}

// Updates the highscores for the given level
func updateLevelScores(s *db.DB, levelScores nget.NScoresResponse) {

	// Used to offset ranks due to cheaters
	rankOffset := 0

	for _, highscore := range levelScores.Scores {
		// Do we ignore this player's scores?
		if isPlayerIgnored(highscore) {
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
				dbScore.Score = newScore.Score
				dbScore.CreatedAt = time.Now()
				s.Save(dbScore)
			}
		}
	}
}

// Updates the highscores for the given episode
func updateEpisodeScores(s *db.DB, episodeScores nget.NScoresResponse) {

	// Used to offset ranks due to cheaters
	rankOffset := 0

	for _, highscore := range episodeScores.Scores {
		// Do we ignore this player's scores?
		if isPlayerIgnored(highscore) {
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
				dbScore.Score = newScore.Score
				s.Save(dbScore)
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