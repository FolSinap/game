package main

import (
	"errors"
	"strings"
)

type description func (state int) string

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
		if item == i {
			return true
		}
	}
	return false
}

func (r *Room) remove(i item) {
	for index, val := range r.items {
		if val == i {
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
	var kitchenDesc description = func(state int) string {
		if state < hasStuff {
			return "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. " + kitchen.showDirections()
		}
		return "ты находишься на кухне, на столе чай, надо идти в универ. " + kitchen.showDirections()
	}
	var hallDesc description = func(int) string {
		return "ничего интересного." + hall.showDirections()
	}
	var roomDesc description = func(int) string {
		directions := room.showDirections()
		switch state {
		case noBackpack:

			return "на столе: ключи, конспекты, на стуле - рюкзак. " + directions
		case hasBackpack:
			desc := "на столе: "
			for _, item := range room.items {
				desc += string(item) + ", "
			}
			return strings.TrimRight(desc, ", ") + ". " + directions
		default:
			return "пустая комната. "  + directions
		}
	}
	var outsideDesc description = func(int) string {
		return "на улице весна. " + outside.showDirections()
	}

	kitchen = Room{
		"кухня, ничего интересного.",
		kitchenDesc,
		[]*direction{{"коридор", false, &hall}},
		[]item{},
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
	}
	room = Room{
		"ты в своей комнате.",
		roomDesc,
		[]*direction{{"коридор", false, &hall}},
		[]item{"ключи", "конспекты"},
	}
	outside = Room{
		"на улице весна.",
		outsideDesc,
		[]*direction{{"домой", false, &hall}},
		[]item{},
	}
}
