package db

import "time"

type LevelHighscore struct {
	Id         int       `gorm:"primary_key;AUTO_INCREMENT"`
	LevelId    int       `gorm:"unique_index:uniq0level0player"`
	PlayerId   int       `gorm:"unique_index:uniq0level0player"`
	Score      int
	Rank       int
	CreatedAt  time.Time
	PlayerName string
}
func (LevelHighscore) TableName() string {
	return "level_highscores"
}

type EpisodeHighscore struct {
	Id         int       `gorm:"primary_key;AUTO_INCREMENT"`
	EpisodeId  int       `gorm:"unique_index:uniq0episode0player"`
	PlayerId   int       `gorm:"unique_index:uniq0episode0player"`
	Score      int
	Rank       int
	CreatedAt  time.Time
	PlayerName string
}
func (EpisodeHighscore) TableName() string {
	return "episode_highscores"
}

// For future use if M&R implement changeable player names
type Player struct {
	Id        int
	Name      string
	IsIgnored bool
}
func (Player) TableName() string {
	return "players"
}

// For future use if we want to display level/episode names
type Episode struct {
	Id       int
	Name     string
	LongName string
}
func (Episode) TableName() string {
	return "episodes"
}

type Level struct {
	Id        int
	EpisodeId int
	Name      string
	LongName  string
}
func (Level) TableName() string {
	return "levels"
}
