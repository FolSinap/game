package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type item string

type backpack interface {
    get(item) error
	has(item) bool
}

type nullBackpack struct {}

func (n *nullBackpack) get(item) error {
	return errors.New("некуда класть")
}

func (n *nullBackpack) has(item) bool {
	return false
}

type trueBackpack []item

func (t *trueBackpack) get(i item) error {
	*t = append(*t, i)
	return nil
}

func (t *trueBackpack) has(i item) bool {
	for _, item := range *t {
		if item == i {
			return true
		}
	}
	return false
}

const (
	noBackpack = iota
	hasBackpack
	hasStuff
	doorIsOpened
)

var (
	kitchen Room
	hall Room
	room Room
	outside Room
	players []*Player
	state int
)

func initGame() {
	state = noBackpack
	initCommands()
	initRooms()
}

func main() {
	initGame()
	addPlayer(NewPlayer("player1"))
	reader := bufio.NewReader(os.Stdin)
	command := ""

	for {
		command, _ = reader.ReadString('\n')
		fmt.Println(handleCommand(strings.TrimSpace(command), players[0]))
	}

}
