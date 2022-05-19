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

	return p.room.desc(state, p), nil
}

var takeCommand Command = func(p *Player, params ...string) (string, error) {
	if len(params) != 1 {
		return "", errors.New("неверное кол-во параметров, требуется - 1")
	}
	requestedItem := params[0]

	err := p.get(item{requestedItem, p.room})
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
		p.takeOn(&trueBackpack{})
		return "вы одели: рюкзак", nil
	}

	return "", errors.New("нечего надевать")
}

var useCommand Command = func(p *Player, params ...string) (string, error) {
	if len(params) != 2 {
		return "", errors.New("неверное кол-во параметров, требуется - 2")
	}
	tool := item{name: params[0]}
	target := params[1]
	if !p.has(tool) {
		return "", errors.New("нет предмета в инвентаре - " + tool.name)
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

var sayCommand Command = func(p *Player, params ...string) (string, error) {
	msg := p.name + " говорит: "
	for _, word := range params {
		msg += word + " "
	}
	msg = strings.TrimRight(msg, " ")
	for _, player := range p.room.players {
		if player != p {
			player.HandleOutput(msg)
		}
	}
	return msg, nil
}

var sayToPlayerCommand Command = func(p *Player, params ...string) (string, error) {
	targetPlayerName := params[0]
	params = params[1:]
	msg := ""

	if len(params) >= 1 {
		msg = p.name + " говорит вам: "
		for _, word := range params {
			msg += word + " "
		}
		msg = strings.TrimRight(msg, " ")
	} else {
		msg = p.name + " выразительно молчит, смотря на вас"
	}
	targetPlayerExists := false

	for _, player := range p.room.players {
		if player.name == targetPlayerName {
			targetPlayerExists = true
			player.HandleOutput(msg)
		}
	}
	if !targetPlayerExists {
		return "тут нет такого игрока", nil
	}
	return "", nil
}

var commands = make(map[string]Command)

func initCommands() {
	commands["идти"] = goCommand
	commands["осмотреться"] = lookAroundCommand
	commands["взять"] = takeCommand
	commands["применить"] = useCommand
	commands["одеть"] = takeOnCommand
	commands["сказать"] = sayCommand
	commands["сказать_игроку"] = sayToPlayerCommand
}

func handleCommand(c string, player *Player) string {
	params := strings.Split(c, " ")
	command := params[0]
	params = params[1:]
	if f, ok := commands[command]; ok {
		res, err := f(player, params...)
		if err != nil {
			return err.Error()
		}
		return res
	}

	return "неизвестная команда"
}
