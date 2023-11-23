# Hangman GoWeb

## Introduction

This hangman project is an assignment done in our first semester studying development. It uses an existing ```hangmancli``` project present in [this repository](https://github.com/debarbarinantoine/hangmancli.git) and runs it in a golang server using ```net/http``` package.

## How to execute it

To run it on your pc, you need to download the **repository** as a **.zip** file clicking on **Code** and then on **Download ZIP**.

Unzip it and then go to ```debarbarinantoinehangmanweb/exec``` and right click in it. Then, click on **Open in the terminal**.

In the terminal, write that line and press ```Enter```:
``` powershell
go run ./main.go
```

Then, it should display that:
```
Server is running...
If the navigator didn't open on its own, just go to http://localhost:8080/index on your browser.
If you want to end the server, type 'stop' here :
```
_Don't close the terminal, otherwise, the server will stop automatically._

When you open your browser at the right *url*, you can then access to the game.
If you want to begin with a clean user list and score list, you can just erase the ```accounts.thorg``` and the ```scores.txt``` in ```debarbarinantoinehangmanweb/Files```.

To stop the server, return to the terminal and type ```stop```, then press ```Enter```.

A lot of information should be displayed in the terminal, because the program is currently quite verbose. Just don't worry about it.

## Project management

This project was made by a team of three people: **Antoine**, **Nicolas** and **Sabrina**.
It doesn't show on the repository nor on the commit history because of our project management.

Because of the willingness of all team members, we decided to do each a different program using our own ````hangmancli```` previously made.

We didn't really work alone though, because we were always ready to help or advise each other if necessary and sometime even coding together during class or on Discord.

Doing that, we wanted to challenge ourselves to learn more and to emulate each other to do better. In the end, we wanted to choose the best project to send it to evaluation: so here it is.

The other projects are here:
- [Hangman GoWeb de Nicolas](https://github.com/Nicolas13100/Hangman_Web_Duo.git)
- [Hangman GoWeb de Sabrina]()