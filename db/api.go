package db

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type DB struct {
	*gorm.DB
}

func GetDB(connString string) (*DB, error) {
	if connString == "" {
	return nil, errors.New("Bad connection string")
	}
	gormDB, err := gorm.Open("mysql", connString)
	return &DB{gormDB}, err
}

func (db DB) CreateTables() {
	db.AutoMigrate(EpisodeHighscore{}, LevelHighscore{})
}

func (db DB) GetAllLevelScores() ([]LevelHighscore, error) {
	var highscores []LevelHighscore
	db.Find(&highscores)
	return highscores, db.Error
}

func (db DB) GetLevelScore(levelId int, playerId int) (LevelHighscore, error) {
	var highscore LevelHighscore
	db.Where("level_id = ? AND player_id = ?", levelId, playerId).First(&highscore)
	return highscore, db.Error
}

func (db DB) CreateLevelScore() {

}

func (db DB) UpdateLevelScore(score *LevelHighscore, newScore int) {
	score.Score = newScore
	score.CreatedAt = time.Now()
	db.Save(score)
}

func (db DB) GetEpisodeScore(episodeId int, playerId int) (EpisodeHighscore, error) {
	var highscore EpisodeHighscore
	db.Where("episode_id = ? AND player_id = ?", episodeId, playerId).First(&highscore)
	return highscore, db.Error
}

func (db DB) CreateEpisodeScore() {

}

func (db DB) UpdateEpisodeScore(score *EpisodeHighscore, newScore int) {
	score.Score = newScore
	score.CreatedAt = time.Now()
	db.Save(score)
}

func (db DB) GetAllPlayers() ([]Player, error) {
	var players []Player
	db.Find(&players)
	return players, db.Error
}

func (db DB) CreatePlayer(id int, name string, isIgnored bool) (*Player) {
	player := &Player{
		Id: id,
		Name: name,
		IsIgnored: isIgnored,
	}
	db.Create(player)
	return player
}