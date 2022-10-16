package ecs

import (
	"errors"
	"fmt"
)

type IEntity interface {
	GetWorld() *IWorld
	GetId() string
	AddComponent(cmp IComponent) error
	HaveComponent(cn string) bool
	GetComponent(id string) (*IComponent, error)
	GetComponents() []*IComponent
	GetComposition() []string
	HaveComposition([]string) bool
}

type ModelEntity struct {
	Id         string                    `json:"id"`
	Components map[string]ModelComponent `json:"components"`
}

type Entity struct {
	Id         string
	World      *IWorld
	Components []*IComponent
}

func (entity *Entity) GetId() string {
	return entity.Id
}

func CEntity(world *IWorld, id string, components []*IComponent) *Entity {
	return &Entity{
		Id:         id,
		World:      world,
		Components: components,
	}
}

func CEntityFromData(world *IWorld, data ModelEntity) *Entity {
	var components []*IComponent

	for k, v := range data.Components {
		components = append(components, CreateComponent(k, v))
	}

	return &Entity{
		Id:         data.Id,
		World:      world,
		Components: components,
	}
}

func (entity *Entity) AddComponent(cmp IComponent) error {
	// Check if we already have a component with same id
	var foundId int = -1

	for idx, component := range entity.GetComponents() {
		componentLocalised := *component
		if componentLocalised.GetId() == cmp.GetId() {
			foundId = idx
		}
	}
	if foundId != -1 {
		return errors.New("Component with same id already exist")
	} else {
		entity.Components = append(entity.Components, &cmp)
		return nil
	}
}

func (entity *Entity) HaveComponent(cn string) bool {
	for _, component := range entity.Components {
		componentLocalised := *component
		if componentLocalised.GetId() == cn {
			return true
		}
	}
	return false
}

func (entity *Entity) GetComponent(id string) (cmp *IComponent, err error) {
	var foundId int = -1

	for idx, component := range entity.GetComponents() {
		componentLocalised := *component
		if componentLocalised.GetId() == id {
			foundId = idx
		}
	}
	if foundId != -1 {
		cmp = entity.Components[foundId]
	} else {
		err = fmt.Errorf("Component %s not found", id)
	}

	return cmp, err
}

func (entity *Entity) GetComponents() (components []*IComponent) {
	return entity.Components
}

func (entity *Entity) GetComposition() (composition []string) {
	for _, component := range entity.Components {
		cmp := *component
		composition = append(composition, cmp.GetId())
	}
	return composition
}

func (entity *Entity) HaveComposition(composition []string) bool {
	entityComposition := entity.GetComposition()
	haveComponent := 0
	for _, componentName := range entityComposition {
		for _, targetComponentName := range composition {
			if componentName == targetComponentName {
				haveComponent++
			}
		}
	}
	return len(composition) == haveComponent
}

func (entity *Entity) GetWorld() *IWorld {
	return entity.World
}
