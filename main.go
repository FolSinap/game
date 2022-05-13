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

type Player struct {
	backpack
	room *Room
}

func (p *Player) takeOn(b backpack) {
	p.backpack = b
}

func (p *Player) get(i item) error {
	if p.room.has(i) {
		err := p.backpack.get(i)
		if err == nil {
			p.room.remove(i)
		}
		return err
	}
	return errors.New("нет такого")
}

func (p *Player) goTo(room string) error {
	if p.room.isAvailable(room) {
		if p.room.canGo(room) {
			_, p.room = p.room.getNextRoomByName(room)
			return nil
		}
		return errors.New("дверь закрыта")
	}
	return errors.New("нет пути в " + room)
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
	player Player
	state int
)

func initGame() {
	state = noBackpack
	initCommands()
	initRooms()
	player = Player{&nullBackpack{}, &kitchen}
}

func main() {
	initGame()
	reader := bufio.NewReader(os.Stdin)
	command := ""

	for {
		command, _ = reader.ReadString('\n')
		fmt.Println(handleCommand(strings.TrimSpace(command)))
	}

}
