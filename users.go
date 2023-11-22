package HangmanWeb

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Admin info: {"Name":"Thorgan","Password":"Th0r9An","Dictionary":"../Files/Dictionaries/ods5.txt","Difficulty":"difficult"}
// Other info: {"Name":"Antoine","Password":"Ant01n3","Dictionary":"../Files/Dictionaries/ods5.txt","Difficulty":"difficult"}

var fileName = "../Files/accounts.thorg"

type UserNotFoundError struct{}

func (m *UserNotFoundError) Error() string {
	return "user not found in " + fileName
}

func iterativeDecrypt(slice []string) []string {
	var result []string
	for {
		i := len(slice) - 1
		if len(slice) == 2 {
			result = append([]string{Decrypt(slice[0] + "\n" + slice[1])}, result...)
			break
		}
		result = append([]string{Decrypt(slice[0] + "\n" + slice[i])}, result...)
		slice = append(slice[:i])
	}
	return result
}

func retreiveUsers() ([]User, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	usersJson := []byte(strings.Join(iterativeDecrypt(strings.Split(string(data), "\n")), ",\n"))

	var users []User

	usersJson = append([]byte{'[', '\n'}, usersJson...)
	usersJson = append(usersJson, '\n', ']')

	fmt.Printf("usersJson: %v\n", string(usersJson))

	err = json.Unmarshal(usersJson, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func checkUsername(username string) bool {
	users, err := retreiveUsers()
	if err != nil {
		log.Fatal(err)
	}

	if users == nil {
		return true
	}

	for _, user := range users {
		if user.Name == username {
			return false
		}
	}
	return true
}

func login(username, password string) bool {
	fmt.Println("login function called")
	users, err := retreiveUsers()
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		if user.Name == username && user.Password == password {
			mySession.MyUser = user
			mySession.isOpen = true
			return true
		}
	}
	return false
}

func (user *User) addUser() {

	newEntry, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	previousContent, _ := os.ReadFile(fileName)
	if len(previousContent) == 0 {
		_, err = file.Write([]byte(firstContent()))
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.Write([]byte{'\n'})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("writing newline...")
		_, err = file.Write([]byte{'\n'})
		if err != nil {
			log.Fatal(err)
		}
	}
	newEntry = append([]byte{'\n'}, newEntry...)
	newEntry = append([]byte(firstContent()), newEntry...)
	fmt.Println(string(newEntry))
	newEntry = []byte(Encrypt(string(newEntry)))
	_, err = file.Write(newEntry)
	if err != nil {
		log.Fatal(err)
	}

}

func (user *User) modifyUser(newUserData User) error {
	if checkUsername(user.Name) {
		return &UserNotFoundError{}
	}
	users, err := retreiveUsers()
	if err != nil {
		return err
	}
	var userIndex int
	for i, singleUser := range users {
		if singleUser.Name == user.Name {
			userIndex = i
			break
		}
	}
	users[userIndex] = newUserData
	var newData []byte
	for i, user := range users {
		data, err := json.Marshal(user)
		fmt.Println(string(data))
		data = append([]byte{'\n'}, data...)
		data = append([]byte(firstContent()), data...)
		data = []byte(Encrypt(string(data)))
		if i < len(users)-1 {
			data = append(data, ',', '\n')
		}
		fmt.Println()
		fmt.Println(string(data))
		newData = append(newData, data...)
		if err != nil {
			return err
		}
	}
	os.WriteFile(fileName, []byte(firstContent()), 0666)

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte{'\n'})
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(newData)
	user.Name = newUserData.Name
	user.Password = newUserData.Password
	user.Dictionary = newUserData.Dictionary
	user.Difficulty = newUserData.Difficulty
	return err
}

func (session *Session) Close() {
	session.isOpen = false
	session.isPlaying = false
	session.MyUser.Name = ""
	session.MyUser.Password = ""
	session.MyUser.Dictionary = ""
	session.MyUser.Difficulty = ""
	session.MyGameData.Game.Name = ""

	fmt.Println("session closed")
	fmt.Printf("%#v\n", session)
}

func firstContent() string {
	return hex.EncodeToString([]byte{0xb5, 0xd4, 0xe6, 0x8a, 0xa5, 0x3d, 0x54, 0x53, 0xc8, 0xd5, 0x77, 0x66, 0x31, 0xf5, 0x5, 0xf0, 0x99, 0xce, 0x5a, 0xc6, 0x10, 0x5e, 0xd8, 0xc6, 0xaf, 0x4a, 0xd5, 0xad, 0xc4, 0x47, 0x4e, 0xf8})
}
