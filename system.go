package goecs

import (
	"fmt"
)

type SIDE string

const (
	SERVER SIDE = "SERVER"
	CLIENT SIDE = "CLIENT"
	HYBRID SIDE = "HYBRID"
)

type ISystem interface {
	GetName() string
	GetId() Identifier
	// Retourne si le système est un système HYBRID, CLIENT, SERVER
	GetSide() SIDE
	UpdateClient()
	UpdateServer()
	Init(*IWorld)
	Listen(string, func(...interface{}) error) error
	Call(string, ...interface{}) error
}

type System struct {
	Id        Identifier
	Name      string
	Type      SIDE
	World     *IWorld
	listening map[string]func(...interface{}) error
}

func (ss *System) GetType() SIDE {
	return ss.Type
}

func (ss *System) Listen(id string, handler func(...interface{}) error) error {
	_, ok := ss.listening[id]
	if !ok {
		ss.listening[id] = handler
		return nil
	}

	return fmt.Errorf("the listener [%s] already exist", id)
}

func (ss *System) Call(id string, args ...interface{}) error {
	listener, ok := ss.listening[id]
	if !ok {
		return fmt.Errorf("the listener [%s] don't exist", id)
	}

	if err := listener(args...); err != nil {
		return err
	}

	return nil
}

func (ss *System) Init(world *IWorld) {
	ss.World = world
}

func (ss *System) GetName() string {
	return ss.Name
}

func (ss *System) GetId() Identifier {
	return ss.Id
}
