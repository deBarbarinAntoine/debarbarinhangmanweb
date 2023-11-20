package HangmanWeb

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"time"
)

var tmpl *template.Template

func fileServer() {
	fs := http.FileServer(http.Dir("../assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func runServer() {
	port := "localhost:8080"
	url := "http://" + port + "/index"
	go http.ListenAndServe(port, nil)
	fmt.Println("Server is running...")
	time.Sleep(time.Second * 5)
	cmd := exec.Command("explorer", url)
	cmd.Run()
	fmt.Println("If the navigator didn't open on its own, just go to ", url, " on your navigator.")
	isRunning := true
	for isRunning {
		fmt.Print("If you want to end the server, type 'stop' here :")
		var command string
		fmt.Scanln(&command)
		if command == "stop" {
			isRunning = false
		}
	}
}

func Run() {
	var err error
	tmpl, err = template.ParseGlob("../templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	routes()
	fileServer()
	runServer()
}
