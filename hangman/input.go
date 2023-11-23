package hangman

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func (game *Game) input() int {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println(colorCode(Red), "Erreur de saisie !", CLEARCOLOR)
		time.Sleep(time.Second * 2)
	}
	if len(input) > 1 && CheckInputFormat(input) {
		result := game.CheckWord(input)
		game.CountScore(result)
		if result == CORRECTWORD {
			game.RevealWord()
		}
		return result
	} else if len(input) == 1 && CheckInputFormat(input) {
		input = strings.ToUpper(input)
		result := game.CheckRune([]rune(input)[0])
		game.DisplayWord([]rune(strings.ToLower(input))[0])
		game.CountScore(result)
		return result
	} else {
		return INCORRECTINPUT
	}
}

func nameInput() (bool, string) {
	var name string
	_, err := fmt.Scanln(&name)
	if err != nil {
		log.Fatal(err)
	}
	if len(name) < 3 || len(name) > 15 || !CheckInputFormat(name) {
		return false, ""
	}
	return true, name
}

func menu(title string, options ...string) string {
	for {
		clearTerminal()
		fmt.Println(colorCode(Deepskyblue), title, CLEARCOLOR)
		fmt.Println()
		for i, option := range options {
			fmt.Println(colorCode(Forestgreen), i+1, ". ", option, CLEARCOLOR)
		}
		var selection string
		_, err := fmt.Scanln(&selection)
		if err != nil {
			log.Fatal(err)
		}
		for i := range options {
			if selection == strconv.Itoa(i+1) {
				return options[i]
			}
		}
	}
}
