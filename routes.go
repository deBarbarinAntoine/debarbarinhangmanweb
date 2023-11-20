package HangmanWeb

import "net/http"

func routes() {
	//http.HandleFunc("/", rootHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/user/create", createUserHandler)
	http.HandleFunc("/user/creatingUser", creatingUserHandler)
	http.HandleFunc("/user/login", loginHandler)
	http.HandleFunc("/user/login/check", openSessionHandler)
	http.HandleFunc("/user/home", homeHandler)
	http.HandleFunc("/user/home/scores", scoresHandler)
	http.HandleFunc("/user/logout", logoutHandler)
	http.HandleFunc("/hangman/game", hangmanHandler)
	http.HandleFunc("/hangman/game/reset", resetHandler)
	http.HandleFunc("/hangman/init", initHandler)
	http.HandleFunc("/hangman/init/treatment", treatmentHandler)
	http.HandleFunc("/hangman/checkInput", checkInputHandler)
	http.HandleFunc("/hangman/endgame", endgameHandler)
}
