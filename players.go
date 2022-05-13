package main

import "errors"

func NewPlayer(name string) *Player {
	return &Player{&nullBackpack{}, name, &kitchen, make(chan string)}
}

func addPlayer(player *Player) {
	kitchen.addPlayer(player)
	players = append(players, player)
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

func (p *Player) HandleInput(msg string) {
	p.output <- handleCommand(msg, p)
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
