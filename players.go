package main

import "errors"

func NewPlayer(name string) *Player {
	return &Player{&nullBackpack{}, name, &kitchen, make(chan string)}
}

func addPlayer(player *Player) {
	kitchen.addPlayer(player)
	players[player.name] = player
}

func findPlayer(name string) *Player {
	if player, ok := players[name]; ok {
		return player
	}
	return nil
}

func removePlayer(p *Player) {
	delete(players, p.name)
	switch p.backpack.(type) {
	case *trueBackpack:
		backpack := p.backpack.(*trueBackpack)
		for _, item := range *backpack {
			item.returnToRoom()
		}
		state = noBackpack
	}
	p.room.removePlayer(p)
}

type Player struct {
	backpack
	name string
	room *Room
	output chan string
}

func (p *Player) GetOutput() chan string {
	return p.output
}

func (p *Player) HandleInput(command string) {
	p.output <- handleCommand(command, p)
}

func (p *Player) HandleOutput(msg string) {
	p.output <- msg
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
			p.room.removePlayer(p)
			_, p.room = p.room.getNextRoomByName(room)
			p.room.addPlayer(p)
			return nil
		}
		return errors.New("дверь закрыта")
	}
	return errors.New("нет пути в " + room)
}
