package hangman

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Save struct {
	Name       string
	Score      int
	Word       string
	Difficulty string
	Dictionary string
}

type Parameters struct {
	Name           string
	DictionaryPath string
	Difficulty     int
}

var savedGames []Save

// SaveGame saves the current game in fileName.
func (game *Game) SaveGame(fileName string, customDictionary bool) {
	var dictionary string
	if !customDictionary {
		dictionary = DictionaryName(game.Dictionary)
	} else {
		dictionary = game.Dictionary
	}
	currentGame := Save{
		Name:       game.Name,
		Score:      game.Score,
		Word:       game.Word,
		Difficulty: ToStringDifficulty(game.Difficulty),
		Dictionary: dictionary,
	}
	newEntry, err := json.Marshal(currentGame)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	previousContent, _ := os.ReadFile(fileName)
	if string(previousContent) != "" {
		newEntry = append([]byte{',', '\n'}, newEntry...)
	}
	_, err = file.Write(newEntry)
	if err != nil {
		log.Fatal(err)
	}
}

// RetrieveSavedGames retrieves all saved games present in fileName and put it in savedEntries.
func RetrieveSavedGames(fileName string) []Save {
	savedEntries, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(colorCode(Salmon), "Aucune sauvegarde détectée...", CLEARCOLOR)
		return nil
	}
	savedEntries = append([]byte{'[', '\n'}, savedEntries...)
	savedEntries = append(savedEntries, '\n', ']')
	err = json.Unmarshal(savedEntries, &savedGames)
	if err != nil {
		fmt.Println(colorCode(Red), "Erreur de récupération des données...", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Orangered), "Données récupérées :", CLEARCOLOR)
		fmt.Println(colorCode(Orange), string(savedEntries), CLEARCOLOR)
		log.Fatal(err)
	}
	return savedGames
}

// ChargeParameters retrieves the parameters present in fileName and changes all corresponding variables.
func (game *Game) ChargeParameters(fileName string) {
	var savedParameters Parameters
	savedEntries, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(colorCode(Salmon), "Aucun fichier de configuration détecté...", CLEARCOLOR)
		time.Sleep(time.Second * 1)
		return
	}
	err = json.Unmarshal(savedEntries, &savedParameters)
	if err != nil {
		fmt.Println(colorCode(Red), "Erreur de récupération des données...", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Aquamarine), "Il est conseillé de supprimer le fichier config.txt\n    dans le dossier Files afin de résoudre le problème.", CLEARCOLOR)
		fmt.Println()
		fmt.Println(colorCode(Orange), "Données récupérées :", CLEARCOLOR)
		fmt.Println(colorCode(Orange), string(savedEntries), CLEARCOLOR)
		log.Fatal(err)
	} else {
		game.Name = savedParameters.Name
		game.Dictionary = savedParameters.DictionaryPath
		game.Difficulty = savedParameters.Difficulty
	}
}

// SaveParameters saves all current parameters in fileName for later use.
func (game *Game) SaveParameters(fileName string) {
	currentParameters := Parameters{
		Name:           game.Name,
		DictionaryPath: game.Dictionary,
		Difficulty:     game.Difficulty,
	}
	newEntry, err := json.Marshal(currentParameters)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fileName, newEntry, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
