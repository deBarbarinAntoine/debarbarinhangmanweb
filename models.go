package HangmanWeb

import "hangman/hangman"

type Session struct {
	isOpen     bool
	isPlaying  bool
	MyUser     User
	MyGameData GameData
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
