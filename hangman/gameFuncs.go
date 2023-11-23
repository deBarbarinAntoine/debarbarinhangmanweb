package hangman

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

type Game struct {
	Name                     string
	Score                    int
	Word                     string
	Difficulty               int
	Dictionary               string
	WordDisplay, RunesPlayed []rune
	words                    []string
	nbLettersFound, NbError  int
}

var hangman []string

const (
	ALREADYPLAYED  = 1
	CORRECTRUNE    = 2
	INCORRECTRUNE  = 3
	CORRECTWORD    = 5
	INCORRECTWORD  = 6
	INCORRECTINPUT = 8

	EASY      = 4
	MEDIUM    = 7
	DIFFICULT = 10
	LEGENDARY = 13

	ONGOING = 14
	WIN     = 15
	LOOSE   = 16
)

// InitGame initialises the game with user data and all variables necessary.
func (game *Game) InitGame() {
	game.words = retrieveWords(game.Dictionary)
	game.Word = game.chooseWord()
	game.WordDisplay = []rune(strings.Repeat("_ ", len(game.Word)))
	game.hint()
}

// CheckInputFormat checks if the input's format is right and returns the result with a boolean.
func CheckInputFormat(input string) bool {
	for _, char := range input {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return false
		}
	}
	return true
}

// CheckWord checks if the player found the word. Returns true if he found it and false if not.
func (game *Game) CheckWord(try string) int {
	if try == game.Word {
		return CORRECTWORD
	}
	return INCORRECTWORD
}

// DisplayWord changes the wordDisplay to replace the '_' character with the rune played if it is in the word.
func (game *Game) DisplayWord(char rune) {
	for i, r := range game.Word {
		if r == char {
			game.WordDisplay[i*2] = char - 32
			game.nbLettersFound++
		}
	}
}

// RevealWord reveals all runes in wordDisplay.
func (game *Game) RevealWord() {
	for i, r := range game.Word {
		game.WordDisplay[i*2] = r - 32
	}
}

// NbRemainingLetters returns the number of letters still not found in the word.
func (game *Game) NbRemainingLetters() int {
	var result int
	for _, char := range game.WordDisplay {
		if char == '_' {
			result++
		}
	}
	return result
}

// retrieveWords retrieves the words from the selected dictionary.
func retrieveWords(dictionary string) []string {
	var words []string
	if dictionary == "" {
		dictionary = "../Files/Dictionaries/ods5.txt"
	}
	content, err := os.ReadFile(dictionary)
	if err != nil {
		log.Fatal(err)
	}
	words = strings.Split(string(content), "\n")
	if !checkDictionary(words) {
		fmt.Println(colorCode(Red), "Erreur d'acquisition des mots du dictionnaire", CLEARCOLOR)
		time.Sleep(time.Second * 2)
		words = retrieveWords("../Files/Dictionaries/ods5.txt")
	}
	return words
}

// retrieveHangman retrieves the hangman in /Files/hangman.txt and stores it in hangman.
func retrieveHangman() {
	hangman = append(hangman[0:0])
	content, err := os.ReadFile("../Files/hangman.txt")
	if err != nil {
		log.Fatal(err)
	}
	var line int
	var str string
	for _, char := range content {
		str += string(char)
		if char == '\n' {
			line++
		}
		if line == 8 {
			hangman = append(hangman, str)
			str = ""
			line = 0
		}
	}
}

// CheckRune checks if the rune played is already played, correct or incorrect.
func (game *Game) CheckRune(char rune) int {
	for _, r := range game.RunesPlayed {
		if r == char {
			return ALREADYPLAYED
		}
	}
	for _, r := range strings.ToUpper(game.Word) {
		if r == char {
			game.RunesPlayed = append(game.RunesPlayed, char)
			return CORRECTRUNE
		}
	}
	game.RunesPlayed = append(game.RunesPlayed, char)
	return INCORRECTRUNE
}

// CountScore changes the score according to the result given in the parameters.
func (game *Game) CountScore(result int) {
	switch result {
	case ALREADYPLAYED:
		break
	case CORRECTRUNE:
		game.Score += 10
	case INCORRECTRUNE:
		game.Score -= 5
		game.NbError++
	case CORRECTWORD:
		game.Score += 11 * game.NbRemainingLetters()
	case INCORRECTWORD:
		game.NbError += 2
		game.Score -= 5
	}
}

// ClearGameData clears the previous' game's data to start a new one.
func (game *Game) ClearGameData() {
	game.Score = 0
	game.Word = ""
	game.WordDisplay = append(game.WordDisplay[0:0])
	game.RunesPlayed = append(game.RunesPlayed[0:0])
	game.words = append(game.words[0:0])
	game.nbLettersFound = 0
	game.NbError = 0
}

// hint reveal a random rune in wordDisplay.
func (game *Game) hint() {
	if game.Difficulty != LEGENDARY {
		i := rand.Intn(len(game.Word) - 1)
		char := []rune(game.Word)[i]
		game.WordDisplay[i*2] = char - 32
	}
}

// chooseWord chooses randomly a word from words (the dictionary's words' list) according to the difficulty set previously.
func (game *Game) chooseWord() string {
	var possibleWords []string
	for _, str := range game.words {
		if len(str) >= game.Difficulty-2 && len(str) <= game.Difficulty {
			possibleWords = append(possibleWords, str)
		}
		if game.Difficulty == LEGENDARY {
			if len(str) > game.Difficulty {
				possibleWords = append(possibleWords, str)
			}
		}
	}
	if len(possibleWords) < 10 {
		var i int
		for _, str := range game.words {
			i++
			if len(str) == game.Difficulty-i-2 || len(str) == game.Difficulty+i {
				possibleWords = append(possibleWords, str)
				if len(possibleWords) > 10 {
					break
				}
			}
		}
	}
	return possibleWords[rand.Intn(len(possibleWords)-1)]
}

// DictionaryName returns the name of the dictionary.
func DictionaryName(dictionary string) string {
	switch dictionary {
	case "../Files/Dictionaries/ods5.txt":
		return "Scrabble"
	case "../Files/Dictionaries/ospd3_expurgated.txt":
		return "Anglais"
	case "../Files/Dictionaries/italiano.txt":
		return "Italien"
	default:
		return "Personnalisé"
	}
}

// SortTopTenGames sort the saved games by score in decreasing order.
func SortTopTenGames() []Save {
	sort.SliceStable(savedGames, func(i, j int) bool { return savedGames[i].Score > savedGames[j].Score })
	return savedGames
}

// ToStringDifficulty returns the name of the difficulty.
func ToStringDifficulty(difficulty int) string {
	switch difficulty {
	case EASY:
		return "Facile"
	case MEDIUM:
		return "Intermédiaire"
	case DIFFICULT:
		return "Difficile"
	case LEGENDARY:
		return "Légendaire"
	default:
		return "Inconnu"
	}
}

// checkDictionary checks if the dictionary is usable or not and changes the case.
func checkDictionary(words []string) bool {
	for i, str := range words {
		words[i] = strings.ToLower(str)
		str = words[i]
		for _, char := range str {
			if char < 'a' || char > 'z' {
				return false
			}
		}
	}
	return true
}

// CheckEndgame verify if the game is finished and returns the status and true if the game is finished or false is the game is still ongoing.
func (game *Game) CheckEndgame(nbErrorMax int) (int, bool) {
	if game.NbError >= nbErrorMax {
		return LOOSE, true
	}
	if strings.Join(strings.Split(string(game.WordDisplay), " "), "") == strings.ToUpper(game.Word) {
		return WIN, true
	}
	return ONGOING, false
}

// Run executes the game, obviously... :)
func Run() {
	principalMenu()
}
