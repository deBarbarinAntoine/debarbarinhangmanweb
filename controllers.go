package HangmanWeb

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	hangman "github.com/debarbarinantoine/hangmancli"
)

// var incorrectLogin bool
var mySession Session

const (
	INGAME    = 30
	INSESSION = 31
	OUT       = 32
)

// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	if mySession.isPlaying {
// 		http.Redirect(w, r, "/hangman", http.StatusMovedPermanently)
// 	} else if mySession.isOpen {
// 		http.Redirect(w, r, "/user/home", http.StatusMovedPermanently)
// 	} else {
// 		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
// 	}
// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(OUT, w, r) {
		return
	}
	err := tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(OUT, w, r) {
		return
	}
	err := tmpl.ExecuteTemplate(w, "createUser", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func creatingUserHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(OUT, w, r) {
		return
	}
	var user = User{
		Name:     r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	if checkUsername(user.Name) {
		user.addUser()
		http.Redirect(w, r, "/user/login", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/user/create", http.StatusMovedPermanently)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(OUT, w, r) {
		return
	}
	err := tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func openSessionHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(OUT, w, r) {
		return
	}
	if !login(r.FormValue("username"), r.FormValue("password")) {
		// incorrectLogin = true
		fmt.Println("incorrect login")
		http.Redirect(w, r, "/user/login", http.StatusMovedPermanently)
	} else {
		fmt.Println("correct login: welcome ", r.FormValue("username"), "!")
		// incorrectLogin = false
		http.Redirect(w, r, "/user/home", http.StatusMovedPermanently)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	err := tmpl.ExecuteTemplate(w, "home", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout")
	mySession.Close()
	fmt.Println("From handler:")
	fmt.Printf("%#v\n", mySession)
	http.Redirect(w, r, "/index", http.StatusMovedPermanently)
}

func scoresHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	savedGames := hangman.RetreiveSavedGames("../Files/scores.txt")
	if savedGames != nil {
		err := tmpl.ExecuteTemplate(w, "scores", savedGames)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	}
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	err := tmpl.ExecuteTemplate(w, "init", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func treatmentHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	if r.Method == http.MethodPost {
		var difficulty int
		switch r.FormValue("difficulty") {
		case "easy":
			difficulty = hangman.EASY
		case "medium":
			difficulty = hangman.MEDIUM
		case "difficult":
			difficulty = hangman.DIFFICULT
		case "legendary":
			difficulty = hangman.LEGENDARY
		}
		mySession.isPlaying = true
		mySession.gameData.Game.Name = r.FormValue("username")
		mySession.gameData.Game.Difficulty = difficulty
		mySession.gameData.Game.Dictionary = "../Files/Dictionaries/ods5.txt"
		mySession.gameData.Game.InitGame()
		mySession.gameData.RunesPlayed = string(mySession.gameData.Game.RunesPlayed)
		mySession.gameData.WordDisplay = string(mySession.gameData.Game.WordDisplay)
		http.Redirect(w, r, "/hangman/game", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/hangman/init", http.StatusMovedPermanently)
	}
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INGAME, w, r) {
		return
	}
	if mySession.gameData.status == hangman.WIN || mySession.gameData.status == hangman.LOOSE {
		http.Redirect(w, r, "/hangman/endgame", http.StatusMovedPermanently)
	}
	switch mySession.gameData.status {
	case 0:
		mySession.gameData.Message = ""
		mySession.gameData.MessageClass = "none"
	case hangman.ALREADYPLAYED:
		mySession.gameData.Message = "Vous avez déjà joué cette lettre !"
		mySession.gameData.MessageClass = "alert"
	case hangman.CORRECTRUNE:
		mySession.gameData.Message = ""
		mySession.gameData.MessageClass = "none"
	case hangman.INCORRECTRUNE:
		mySession.gameData.Message = "Lettre incorrecte !"
		mySession.gameData.MessageClass = "alert"
	case hangman.CORRECTWORD:
		mySession.gameData.Message = ""
		mySession.gameData.MessageClass = "none"
	case hangman.INCORRECTWORD:
		mySession.gameData.Message = "Mot incorrect !"
		mySession.gameData.MessageClass = "alert"
	case hangman.INCORRECTINPUT:
		mySession.gameData.Message = "Saisie incorrecte !"
		mySession.gameData.MessageClass = "alert"
	}
	mySession.gameData.RunesPlayed = string(mySession.gameData.Game.RunesPlayed)
	mySession.gameData.WordDisplay = string(mySession.gameData.Game.WordDisplay)
	err := tmpl.ExecuteTemplate(w, "hangman", mySession.gameData)
	if err != nil {
		log.Fatal(err)
	}
}

func checkInputHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INGAME, w, r) {
		return
	}
	if r.Method != http.MethodPost {
		fmt.Println("r.Method != http.MethodPost")
		http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
	}
	input := r.FormValue("input")
	fmt.Println("input: ", input)
	if len(input) > 1 && hangman.CheckInputFormat(input) {
		mySession.gameData.status = mySession.gameData.Game.CheckWord(input)
		mySession.gameData.Game.CountScore(mySession.gameData.status)
		if mySession.gameData.status == hangman.CORRECTWORD {
			mySession.gameData.Game.RevealWord()
		}
		fmt.Println("mySession.gameData.status = ", mySession.gameData.status)
	} else if len(input) == 1 && hangman.CheckInputFormat(input) {
		input = strings.ToUpper(input)
		mySession.gameData.status = mySession.gameData.Game.CheckRune([]rune(input)[0])
		mySession.gameData.Game.DisplayWord([]rune(strings.ToLower(input))[0])
		mySession.gameData.Game.CountScore(mySession.gameData.status)
		fmt.Println("mySession.gameData.status = ", mySession.gameData.status)
	} else {
		mySession.gameData.status = hangman.INCORRECTINPUT
		fmt.Println("mySession.gameData.status = ", mySession.gameData.status)
	}
	if endgameStatus, endgame := mySession.gameData.Game.CheckEndgame(9); endgame {
		mySession.gameData.status = endgameStatus
		http.Redirect(w, r, "/hangman/endgame", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
	}
}

func endgameHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INGAME, w, r) {
		return
	}
	if mySession.gameData.status == hangman.WIN || mySession.gameData.status == hangman.LOOSE {
		var imgSrc, imgAlt, message, messageClass string
		if mySession.gameData.status == hangman.WIN {
			mySession.gameData.Game.SaveGame("../Files/scores.txt")
			imgSrc = "victory"
			imgAlt = "victory"
			message = "Félicitations " + mySession.gameData.Game.Name + ", vous avez gagné !"
			messageClass = "victory"
		} else {
			imgSrc = "loss"
			imgAlt = "loss"
			message = "GAME OVER !"
			messageClass = "loss"
		}
		endGameData := struct {
			ImgSrc       string
			ImgAlt       string
			Message      string
			MessageClass string
			Data         hangman.Game
		}{
			ImgSrc:       imgSrc,
			ImgAlt:       imgAlt,
			Message:      message,
			MessageClass: messageClass,
			Data:         mySession.gameData.Game,
		}
		mySession.gameData = GameData{}
		mySession.isPlaying = false
		err := tmpl.ExecuteTemplate(w, "endgame", endGameData)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
	}
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("resetHandler")
	fmt.Println("isPlaying: ", mySession.isPlaying)
	if !redirect(INGAME, w, r) {
		return
	}
	mySession.gameData = GameData{}
	mySession.isPlaying = false
	http.Redirect(w, r, "/user/home", http.StatusSeeOther)
}

func redirectSession(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/user/home", http.StatusSeeOther)
}

func redirectGame(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
}

func redirectIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

func redirect(status int, w http.ResponseWriter, r *http.Request) bool {
	var isSessionRoute, isGameRoute bool
	switch status {
	case INGAME:
		isGameRoute = true
		isSessionRoute = false
	case INSESSION:
		isGameRoute = false
		isSessionRoute = true
	case OUT:
		isGameRoute = false
		isSessionRoute = false
	}
	if mySession.isPlaying {
		if isGameRoute {
			fmt.Println("Correct game route")
			return true
		} else {
			fmt.Println("Incorrect game route")
			redirectGame(w, r)
			return false
		}
	} else if mySession.isOpen {
		if isSessionRoute {
			fmt.Println("Correct session route")
			return true
		} else {
			fmt.Println("Incorrect session route")
			redirectSession(w, r)
			return false
		}
	} else if !isGameRoute && !isSessionRoute {
		fmt.Println("Correct visitor route")
		return true
	} else {
		fmt.Println("Incorrect visitor route")
		redirectIndex(w, r)
		return false
	}
}
