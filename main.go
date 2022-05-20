package main

import (
	"errors"
	"fmt"
)

type item struct {
	name string
	from *Room
}

func (i item) returnToRoom() {
	i.from.add(i)
}

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
		if item.name == i.name {
			return true
		}
	}
	return false
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
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
	players  = make(map[string]*Player)
	state int
)

func initGame() {
	state = noBackpack
	initCommands()
	initRooms()
	initBot()
	panicOnError(initLogs())
}

func main() {
	initGame()
	//addPlayer(NewPlayer("player1"))
	//reader := bufio.NewReader(os.Stdin)
	//command := ""
	//
	//for {
	//	command, _ = reader.ReadString('\n')
	//	fmt.Println(handleCommand(strings.TrimSpace(command), players["player1"]))
	//}
	fmt.Scanln()
}
