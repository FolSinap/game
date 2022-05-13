package main

import (
	"errors"
	"strings"
)

type Command func(p *Player, params ...string) (string, error)

var goCommand Command = func(p *Player, params ...string) (string, error) {
	if len(params) != 1 {
		return "", errors.New("неверное кол-во параметров, требуется - 1")
	}
	direction := params[0]
	err := p.goTo(direction)
	if err != nil {
		return "", err
	}
	return p.room.greet(), nil
}

var lookAroundCommand Command = func(p *Player, params ...string) (string, error) {
	if len(params) != 0 {
		return "", errors.New("неверное кол-во параметров, команда не требует параметров")
	}

	return player.room.desc(state), nil
}

var takeCommand Command = func(p *Player, params ...string) (string, error) {
	if len(params) != 1 {
		return "", errors.New("неверное кол-во параметров, требуется - 1")
	}
	requestedItem := params[0]

	err := p.get(item(requestedItem))
	if err != nil {
		return "", err
	}

	if room.isEmpty() {
		state = hasStuff
	}
	return "предмет добавлен в инвентарь: " + requestedItem, nil
}

var takeOnCommand Command = func(p *Player, params ...string) (string, error) {
	if len(params) != 1 {
		return "", errors.New("неверное кол-во параметров, требуется - 1")
	}
	requestedItem := params[0]

	if state == noBackpack && p.room == &room && requestedItem == "рюкзак" {
		state = hasBackpack
		player.takeOn(&trueBackpack{})
		return "вы одели: рюкзак", nil
	}

	return "", errors.New("нечего надевать")
}

var useCommand Command = func(p *Player, params ...string) (string, error) {
	if len(params) != 2 {
		return "", errors.New("неверное кол-во параметров, требуется - 2")
	}
	tool := item(params[0])
	target := params[1]
	if !p.has(tool) {
		return "", errors.New("нет предмета в инвентаре - " + string(tool))
	}

	if target == "дверь" && p.room == &hall {
		state = doorIsOpened
		for _, direction := range p.room.nextRooms {
			direction.isLocked = false
		}
		return "дверь открыта", nil
	}

	return "", errors.New("не к чему применить")
}

var commands = make(map[string]Command)

func initCommands() {
	commands["идти"] = goCommand
	commands["осмотреться"] = lookAroundCommand
	commands["взять"] = takeCommand
	commands["применить"] = useCommand
	commands["одеть"] = takeOnCommand
}

func handleCommand(c string) string {
	params := strings.Split(c, " ")
	command := params[0]
	params = params[1:]
	if f, ok := commands[command]; ok {
		res, err := f(&player, params...)
		if err != nil {
			return err.Error()
		}
		return res
	}

	return "неизвестная команда"
}
