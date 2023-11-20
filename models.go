package HangmanWeb

import hangman "github.com/debarbarinantoine/hangmancli"

type Session struct {
	isOpen    bool
	isPlaying bool
	user      User
	gameData  GameData
}

type User struct {
	Name       string
	Password   string
	Dictionary string
	Difficulty string
}

type GameData struct {
	status       int
	Message      string
	MessageClass string
	RunesPlayed  string
	WordDisplay  string
	Game         hangman.Game
}
