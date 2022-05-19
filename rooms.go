package main

import (
	"errors"
	"strings"
)

type description func (state int, p *Player) string

type direction struct {
	name string
	isLocked bool
	room *Room
}

type Room struct {
	greetings string
	desc description
	nextRooms []*direction
	items []item
	players []*Player
}

func (r *Room) addPlayer(p *Player) {
	r.players = append(r.players, p)
}

func (r *Room) removePlayer(p *Player) {
	for i, player := range r.players {
		if player == p {
			r.players = append(r.players[:i], r.players[i + 1:]...)
		}
	}
}

func (r *Room) msgPostfix(p *Player) string {
	return r.showDirections() + r.showPlayers(p)
}

func (r *Room) showPlayers(except *Player) string {
	msg := ""
	if len(r.players) > 1 {
		msg = ". Кроме вас тут ещё "
		for _, player := range r.players {
			if player != except {
				msg += player.name + ", "
			}
		}
		msg = strings.TrimRight(msg, ", ")
	}
	return msg
}

func (r *Room) greet() string {
	return r.greetings + " " + r.showDirections()
}

func (r *Room) showDirections() string {
	var directions = "можно пройти - "
	for _, direction := range r.nextRooms {
		directions += direction.name + ", "
	}
	return strings.TrimRight(directions, ", ")
}

func (r *Room) isAvailable(room string) bool {
	for _, direction := range r.nextRooms {
		if direction.name == room {
			return true
		}
	}
	return false
}

func (r *Room) canGo(room string) bool {
	for _, direction := range r.nextRooms {
		if direction.name == room {
			if direction.isLocked {
				return false
			}
			return true
		}
	}
	return false
}

func (r *Room) isEmpty() bool {
	if len(r.items) > 0 {
		return false
	}
	return true
}

func (r *Room) has(i item) bool {
	for _, item := range r.items {
		if item.name == i.name {
			return true
		}
	}
	return false
}

func (r *Room) add(i item) {
	r.items = append(r.items, i)
}

func (r *Room) remove(i item) {
	for index, item := range r.items {
		if item.name == i.name {
			r.items = append(r.items[:index], r.items[index + 1:]...)
		}
	}
}

func (r *Room) getNextRoomByName(name string) (error, *Room) {
	for _, direction := range r.nextRooms {
		if direction.name == name {
			return nil, direction.room
		}
	}
	return errors.New("нет пути в " + name), nil
}

func initRooms() {
	var kitchenDesc description = func(state int, p *Player) string {
		if state < hasStuff {
			return "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. " + kitchen.msgPostfix(p)
		}
		return "ты находишься на кухне, на столе чай, надо идти в универ. " + kitchen.msgPostfix(p)
	}
	var hallDesc description = func(i int, p *Player) string {
		return "ничего интересного." + hall.msgPostfix(p)
	}
	var roomDesc description = func(i int, p *Player) string {
		directions := room.msgPostfix(p)
		switch state {
		case noBackpack:

			return "на столе: ключи, конспекты, на стуле - рюкзак. " + directions
		case hasBackpack:
			desc := "на столе: "
			for _, item := range room.items {
				desc += item.name + ", "
			}
			return strings.TrimRight(desc, ", ") + ". " + directions
		default:
			return "пустая комната. "  + directions
		}
	}
	var outsideDesc description = func(i int, p *Player) string {
		return "на улице весна. " + outside.msgPostfix(p)
	}

	kitchen = Room{
		"кухня, ничего интересного.",
		kitchenDesc,
		[]*direction{{"коридор", false, &hall}},
		[]item{},
		[]*Player{},
	}
	hall = Room{
		"ничего интересного.",
		hallDesc,
		[]*direction{
			{"кухня", false, &kitchen},
			{"комната", false, &room},
			{"улица", true, &outside},
		},
		[]item{},
		[]*Player{},
	}
	room = Room{
		"ты в своей комнате.",
		roomDesc,
		[]*direction{{"коридор", false, &hall}},
		[]item{{"ключи", &room}, {"конспекты", &room}},
		[]*Player{},
	}
	outside = Room{
		"на улице весна.",
		outsideDesc,
		[]*direction{{"домой", false, &hall}},
		[]item{},
		[]*Player{},
	}
}
