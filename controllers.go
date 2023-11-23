package HangmanWeb

import (
	hangman "HangmanWeb/hangman"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

var mySession Session
var dictionaries = []string{"français", "anglais", "italien"}

const (
	INGAME    = 30
	INSESSION = 31
	OUT       = 32
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusMovedPermanently)
}

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
	fmt.Println("home handler")
	if !redirect(INSESSION, w, r) {
		return
	}
	var nbGames, maxScore, totalScore int
	dictionaryUse := make(map[string]int)
	allGames := hangman.RetrieveSavedGames("../Files/scores.txt")
	for _, game := range allGames {
		if game.Name == mySession.MyUser.Name {
			nbGames++
			totalScore += game.Score
			if game.Score > maxScore {
				maxScore = game.Score
			}
			dictionaryUse[game.Dictionary]++
		}
	}
	var err error
	if nbGames != 0 {
		averageScore := totalScore / nbGames
		var favoriteDictionary string
		var maxUsed int
		for dictionary, nbUsed := range dictionaryUse {
			if nbUsed > 0 && nbUsed > maxUsed {
				favoriteDictionary = dictionary
				maxUsed = nbUsed
			}
		}
		homeData := struct {
			Name               string
			NbGames            int
			MaxScore           int
			AverageScore       int
			FavoriteDictionary string
		}{
			Name:               mySession.MyUser.Name,
			NbGames:            nbGames,
			MaxScore:           maxScore,
			AverageScore:       averageScore,
			FavoriteDictionary: favoriteDictionary,
		}
		err = tmpl.ExecuteTemplate(w, "home", homeData)
	} else {
		homeData := struct {
			Name               string
			NbGames            int
			MaxScore           int
			AverageScore       int
			FavoriteDictionary string
		}{
			Name:               mySession.MyUser.Name,
			NbGames:            nbGames,
			MaxScore:           maxScore,
			AverageScore:       0,
			FavoriteDictionary: "aucun",
		}
		err = tmpl.ExecuteTemplate(w, "home", homeData)
	}
	if err != nil {
		log.Fatal(err)
	}
	return
}

func scoresHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	savedGames := hangman.RetrieveSavedGames("../Files/scores.txt")
	sort.SliceStable(savedGames, func(i, j int) bool { return savedGames[i].Score > savedGames[j].Score })
	if savedGames != nil {
		err := tmpl.ExecuteTemplate(w, "scores", savedGames)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}
}

func modifyUserHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	fmt.Println("modify user")
	tmpl.ExecuteTemplate(w, "modify", mySession)
	mySession.MyGameData.Message = ""
	mySession.MyGameData.MessageClass = ""
	return
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	if r.Method == http.MethodPost {
		fmt.Println("update user")
		username := r.FormValue("username")
		newPassword := r.FormValue("newPassword")
		if (checkUsername(username) || username == mySession.MyUser.Name) && mySession.MyUser.Password == r.FormValue("password") && len(newPassword) > 5 {
			fmt.Println("Previous name: ", mySession.MyUser.Name)
			fmt.Println("Previous password: ", mySession.MyUser.Password)
			fmt.Println()
			newUser := User{Name: username, Password: newPassword}
			err := mySession.MyUser.modifyUser(newUser)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("New name: ", mySession.MyUser.Name)
			fmt.Println("New password: ", mySession.MyUser.Password)
			fmt.Println()
			http.Redirect(w, r, "/user/home", http.StatusMovedPermanently)
			return
		} else {
			mySession.MyGameData.Message = "Données invalides !"
			mySession.MyGameData.MessageClass = "message red"
			fmt.Println("Error: ", mySession.MyGameData.Message)
			http.Redirect(w, r, "/user/modify", http.StatusSeeOther)
			return
		}
	} else {
		http.Redirect(w, r, "/user/modify", http.StatusMovedPermanently)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout")
	mySession.Close()
	fmt.Println("From handler:")
	fmt.Printf("%#v\n", mySession)
	http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	return
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	err := tmpl.ExecuteTemplate(w, "init", dictionaries)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func treatmentHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INSESSION, w, r) {
		return
	}
	if r.Method == http.MethodPost {
		var difficulty int
		switch r.FormValue("difficulty") {
		case "easy":
			mySession.MyUser.Difficulty = "Facile"
			difficulty = hangman.EASY
		case "medium":
			mySession.MyUser.Difficulty = "Intermédiaire"
			difficulty = hangman.MEDIUM
		case "difficult":
			mySession.MyUser.Difficulty = "Difficile"
			difficulty = hangman.DIFFICULT
		case "legendary":
			mySession.MyUser.Difficulty = "Légendaire"
			difficulty = hangman.LEGENDARY
		}
		mySession.isPlaying = true
		mySession.MyGameData.Game.Name = mySession.MyUser.Name
		mySession.MyGameData.Game.Difficulty = difficulty
		mySession.MyGameData.Game.Dictionary = "../Files/Dictionaries/" + r.FormValue("dictionary") + ".txt"
		mySession.MyUser.Dictionary = r.FormValue("dictionary")
		mySession.MyUser.modifyUser(mySession.MyUser)
		mySession.MyGameData.Game.InitGame()
		mySession.MyGameData.Game.Dictionary = mySession.MyUser.Dictionary
		mySession.MyGameData.RunesPlayed = string(mySession.MyGameData.Game.RunesPlayed)
		mySession.MyGameData.WordDisplay = string(mySession.MyGameData.Game.WordDisplay)
		http.Redirect(w, r, "/hangman/game", http.StatusMovedPermanently)
		return
	} else {
		http.Redirect(w, r, "/hangman/init", http.StatusMovedPermanently)
		return
	}
}

func hangmanHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INGAME, w, r) {
		return
	}
	if mySession.MyGameData.status == hangman.WIN || mySession.MyGameData.status == hangman.LOOSE {
		http.Redirect(w, r, "/hangman/endgame", http.StatusMovedPermanently)
		return
	}
	switch mySession.MyGameData.status {
	case 0:
		mySession.MyGameData.Message = ""
		mySession.MyGameData.MessageClass = "none"
	case hangman.ALREADYPLAYED:
		mySession.MyGameData.Message = "Vous avez déjà joué cette lettre !"
		mySession.MyGameData.MessageClass = "alert"
	case hangman.CORRECTRUNE:
		mySession.MyGameData.Message = ""
		mySession.MyGameData.MessageClass = "none"
	case hangman.INCORRECTRUNE:
		mySession.MyGameData.Message = "Lettre incorrecte !"
		mySession.MyGameData.MessageClass = "alert"
	case hangman.CORRECTWORD:
		mySession.MyGameData.Message = ""
		mySession.MyGameData.MessageClass = "none"
	case hangman.INCORRECTWORD:
		mySession.MyGameData.Message = "Mot incorrect !"
		mySession.MyGameData.MessageClass = "alert"
	case hangman.INCORRECTINPUT:
		mySession.MyGameData.Message = "Saisie incorrecte !"
		mySession.MyGameData.MessageClass = "alert"
	}
	mySession.MyGameData.RunesPlayed = string(mySession.MyGameData.Game.RunesPlayed)
	mySession.MyGameData.WordDisplay = string(mySession.MyGameData.Game.WordDisplay)
	err := tmpl.ExecuteTemplate(w, "hangman", mySession)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func checkInputHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INGAME, w, r) {
		return
	}
	if r.Method != http.MethodPost {
		fmt.Println("r.Method != http.MethodPost")
		http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
		return
	}
	input := r.FormValue("input")
	fmt.Println("input: ", input)
	if len(input) > 1 && hangman.CheckInputFormat(input) {
		mySession.MyGameData.status = mySession.MyGameData.Game.CheckWord(input)
		mySession.MyGameData.Game.CountScore(mySession.MyGameData.status)
		if mySession.MyGameData.status == hangman.CORRECTWORD {
			mySession.MyGameData.Game.RevealWord()
		}
		fmt.Println("mySession.gameData.status = ", mySession.MyGameData.status)
	} else if len(input) == 1 && hangman.CheckInputFormat(input) {
		input = strings.ToUpper(input)
		mySession.MyGameData.status = mySession.MyGameData.Game.CheckRune([]rune(input)[0])
		mySession.MyGameData.Game.DisplayWord([]rune(strings.ToLower(input))[0])
		mySession.MyGameData.Game.CountScore(mySession.MyGameData.status)
		fmt.Println("mySession.gameData.status = ", mySession.MyGameData.status)
	} else {
		mySession.MyGameData.status = hangman.INCORRECTINPUT
		fmt.Println("mySession.gameData.status = ", mySession.MyGameData.status)
	}
	if endgameStatus, endgame := mySession.MyGameData.Game.CheckEndgame(9); endgame {
		mySession.MyGameData.status = endgameStatus
		http.Redirect(w, r, "/hangman/endgame", http.StatusMovedPermanently)
		return
	} else {
		http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
		return
	}
}

func endgameHandler(w http.ResponseWriter, r *http.Request) {
	if !redirect(INGAME, w, r) {
		return
	}
	if mySession.MyGameData.status == hangman.WIN || mySession.MyGameData.status == hangman.LOOSE {
		var imgSrc, imgAlt, message, messageClass string
		if mySession.MyGameData.status == hangman.WIN {
			mySession.MyGameData.Game.SaveGame("../Files/scores.txt", true)
			imgSrc = "win.png"
			imgAlt = "vous avez gagné"
			message = "Félicitations " + mySession.MyGameData.Game.Name + ", vous avez gagné !"
			messageClass = "victory"
		} else {
			imgSrc = "hangman9.png"
			imgAlt = "image de pendu - vous avez perdu"
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
			Data:         mySession.MyGameData.Game,
		}
		mySession.MyGameData = GameData{}
		mySession.isPlaying = false
		err := tmpl.ExecuteTemplate(w, "endgame", endGameData)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else {
		http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
		return
	}
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("resetHandler")
	fmt.Println("isPlaying: ", mySession.isPlaying)
	if !redirect(INGAME, w, r) {
		return
	}
	fmt.Println("exiting game...")
	mySession.MyGameData = GameData{}
	mySession.isPlaying = false
	fmt.Println("isPlaying: ", mySession.isPlaying)
	http.Redirect(w, r, "/user/home", http.StatusSeeOther)
	return
}

func redirectSession(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/user/home", http.StatusSeeOther)
	return
}

func redirectGame(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/hangman/game", http.StatusSeeOther)
	return
}

func redirectIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusSeeOther)
	return
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
			//fmt.Println("Correct game route")
			return true
		} else {
			//fmt.Println("Incorrect game route")
			redirectGame(w, r)
			return false
		}
	} else if mySession.isOpen {
		if isSessionRoute {
			//fmt.Println("Correct session route")
			return true
		} else {
			//fmt.Println("Incorrect session route")
			redirectSession(w, r)
			return false
		}
	} else if !isGameRoute && !isSessionRoute {
		//fmt.Println("Correct visitor route")
		return true
	} else {
		//fmt.Println("Incorrect visitor route")
		redirectIndex(w, r)
		return false
	}
}
