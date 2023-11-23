package hangman

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var dictionary string

// runCmd executes the command and arguments put in the parameters.
func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// clearTerminal clears the terminal using the corresponding command.
func clearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func principalMenu() {
	for {
		switch menu("------- MENU PRINCIPAL -------", "Nouvelle partie", "Dictionnaire", "Meilleurs scores", "Quitter") {
		case "Nouvelle partie":
			var myGame Game
			myGame.setGame()
		case "Dictionnaire":
			changeDictionary()
		case "Meilleurs scores":
			topScores()
		case "Quitter":
			os.Exit(0)
		}
	}
}

func (game *Game) setGame() {
	game.ChargeParameters("../Files/config.txt")
	if dictionary != "" {
		game.Dictionary = dictionary
	}
	correctInput := true
	for {
		clearTerminal()
		fmt.Println(colorCode(Deepskyblue), "------- INITIALISATION DU JEU -------", CLEARCOLOR)
		fmt.Println()
		if !correctInput {
			fmt.Println(colorCode(Red), "Nom saisi incorrect (entre 3 et 15 caractères, sans nombres ni signes)", CLEARCOLOR)
		}
		fmt.Print(colorCode(Forestgreen), "Saisissez votre nom : ", colorCode(Aquamarine))
		if correctInput, game.Name = nameInput(); correctInput {
			break
		}
	}
	game.setDifficulty()
}

func (game *Game) setDifficulty() {
	switch menu("------- DIFFICULTÉ -------", "Facile", "Intermédiaire", "Difficile", "Légendaire") {
	case "Facile":
		game.Difficulty = EASY
	case "Intermédiaire":
		game.Difficulty = MEDIUM
	case "Difficile":
		game.Difficulty = DIFFICULT
	case "Légendaire":
		game.Difficulty = LEGENDARY
	}
	game.play()
}

func (game *Game) play() {
	retrieveHangman()
	game.InitGame()
	var status int
	var gameHasEnded bool
	var previousResult int
	for {
		clearTerminal()
		fmt.Println(colorCode(Deepskyblue), "------------------------- HANGMAN -------------------------", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Forestgreen), "Nom : ", colorCode(Aquamarine), game.Name, colorCode(Forestgreen), "\tDifficulté : ", colorCode(Aquamarine), ToStringDifficulty(game.Difficulty), colorCode(Forestgreen), "\tDictionnaire : ", colorCode(Aquamarine), DictionaryName(game.Dictionary), colorCode(Forestgreen), "\tScore : ", colorCode(Aquamarine), game.Score, CLEARCOLOR)
		fmt.Println()
		fmt.Println(hangman[game.NbError])
		fmt.Println()
		fmt.Println(colorCode(Aquamarine), string(game.WordDisplay), CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Forestgreen), "Lettres déjà jouées : ", colorCode(Orange), string(game.RunesPlayed), CLEARCOLOR)
		fmt.Println()
		if status, gameHasEnded = game.CheckEndgame(len(hangman) - 1); gameHasEnded {
			time.Sleep(time.Second * 2)
			break
		}
		switch previousResult {
		case ALREADYPLAYED:
			fmt.Println(colorCode(Orangered), "Cette lettre a déjà été jouée !", CLEARCOLOR)
			fmt.Println()
		case INCORRECTINPUT:
			fmt.Println(colorCode(Orangered), "Saisie invalide !", CLEARCOLOR)
			fmt.Println()
		}
		fmt.Print(colorCode(Deepskyblue), "Proposez une lettre ou un mot : ", colorCode(Aquamarine))
		previousResult = game.input()
	}
	game.endGame(status)
}

func (game *Game) endGame(status int) {
	for {
		clearTerminal()
		if status == WIN {
			game.SaveGame("../Files/scores.txt", false)
			fmt.Println(colorCode(Cyan), "\tFÉLICITATIONS, VOUS AVEZ GAGNÉ !", CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Aquamarine), "Le mot était ", strings.ToUpper(game.Word), CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Aquamarine), "Votre score est : ", game.Score, CLEARCOLOR)
		} else if status == LOOSE {
			fmt.Println(colorCode(Orange), "\tGAME OVER !", CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Red), hangman[game.NbError], CLEARCOLOR)
			fmt.Println()
			fmt.Println(colorCode(Aquamarine), "Le mot était ", game.Word, CLEARCOLOR)
		}

		fmt.Println()
		fmt.Println(colorCode(Aquamarine), "0.  Retour au menu principal")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		if input == "0" {
			break
		}
	}
	game.SaveParameters("../Files/config.txt")
	game.ClearGameData()
}

func changeDictionary() {
	switch menu("------- CHANGER DE DICTIONNAIRE -------", "Scrabble français", "Scabble Anglais", "Italien", "Retour") {
	case "Scrabble français":
		dictionary = "../Files/Dictionaries/ods5.txt"
	case "Scabble Anglais":
		dictionary = "../Files/Dictionaries/ospd3_expurgated.txt"
	case "Italien":
		dictionary = "../Files/Dictionaries/italiano.txt"
	case "Retour":
		break
	}
}

func topScores() {
	RetrieveSavedGames("../Files/scores.txt")
	SortTopTenGames()
	for {
		clearTerminal()
		fmt.Println(colorCode(Deepskyblue), "-------------------- MEILLEURS SCORES --------------------", CLEARCOLOR)
		for i, game := range savedGames {
			fmt.Println(colorCode(Forestgreen), i+1, "\t", game.Name, "\t", game.Score, "\t", game.Difficulty, "\t", game.Dictionary, "\t", CLEARCOLOR)
			if i > 10 {
				break
			}
		}
		fmt.Println()
		fmt.Println(colorCode(Aquamarine), "0.  Retour au menu principal")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		if input == "0" {
			break
		}
	}
}
