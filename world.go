package ecs

import (
	"errors"
)

type IWorld interface {
	GetId() string
	AddEntity(*IEntity) error
	GetEntity(string) *IEntity
	GetEntities() []*IEntity
	GetEntitiesByComponentId(string) []*IEntity
	GetEntitiesWithComponents(...string) []*IEntity
	GetEntitiesWithStrictComposition([]string) []*IEntity
	GetEntitiesWithComposition([]string) []*IEntity
	RemoveEntity(string) error
	AddSystem(*ISystem)
	GetSystemById(string) *ISystem
	RemoveSystem(string) error
	Update()
}

type World struct {
	Id       string
	Entities []*IEntity
	Systems  []*ISystem
}

func (world *World) GetId() string {
	return world.Id
}

func (world *World) AddEntity(entity *IEntity) (err error) {
	var found bool = false
	for _, ent := range world.Entities {
		entityLocalised := *entity
		entLocalised := *ent
		if entLocalised.GetId() == entityLocalised.GetId() {
			found = true
		}
	}
	if found {
		err = errors.New("Cannot add entity because an entity with same id already exist")
		return err
	} else {
		world.Entities = append(world.Entities, entity)
	}
	return err
}

func (world *World) GetEntity(id string) (ent *IEntity) {
	for _, entity := range world.Entities {
		entityLocalised := *entity
		if entityLocalised.GetId() == id {
			ent = entity
		}
	}
	return ent
}

func (world *World) GetEntities() (entities []*IEntity) {
	return world.Entities
}

func (world *World) GetEntitiesByComponentId(id string) (entities []*IEntity) {
	for _, entity := range world.Entities {
		entityLocalised := *entity
		for _, component := range entityLocalised.GetComponents() {
			cmp := *component
			if cmp.GetId() == id {
				entities = append(entities, entity)
			}
		}
	}
	return entities
}

func (world *World) GetEntitiesWithComponents(v ...string) (entities []*IEntity) {
	for _, entity := range world.Entities {
		entityLocalised := *entity
		var checkeds int = 0
		for _, cmpName := range v {
			if entityLocalised.HaveComponent(cmpName) {
				checkeds = checkeds + 1
			}
		}
		if checkeds == len(v) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (world *World) GetEntitiesWithComposition(composition []string) (entities []*IEntity) {
	for _, entity := range world.GetEntities() {
		entityLocalised := *entity
		if entityLocalised.HaveComposition(composition) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (world *World) GetEntitiesWithStrictComposition(composition []string) (entities []*IEntity) {
	for _, entity := range world.GetEntities() {
		entityLocalised := *entity
		if len(composition) == len(entityLocalised.GetComponents()) && entityLocalised.HaveComposition(composition) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (world *World) RemoveEntity(id string) (err error) {
	var newEntities []*IEntity = []*IEntity{}
	var entityFound bool = false

	for _, entity := range world.Entities {
		localisedEntity := *entity
		if localisedEntity.GetId() != id {
			newEntities = append(newEntities, entity)
		} else {
			entityFound = true
		}
	}

	if !entityFound {
		return errors.New("Cannot delete entity because it doesn't exist")
	}

	world.Entities = newEntities

	return nil
}

func (world *World) AddSystem(sys *ISystem) {
	world.Systems = append(world.Systems, sys)
}

func (world *World) GetSystemById(id string) *ISystem {
	var systemFound *ISystem
	for _, system := range world.GetSystems() {
		systemLocalised := *system
		if systemLocalised.GetId() == id {
			systemFound = system
		}
	}
	return systemFound
}

func (world *World) GetSystems() []*ISystem {
	return world.Systems
}

func (world *World) RemoveSystem(id string) (err error) {
	var newSystems []*ISystem = []*ISystem{}
	var systemFound bool = false

	for _, system := range world.GetSystems() {
		systemLocalised := *system
		if systemLocalised.GetId() != id {
			newSystems = append(newSystems, system)
		} else {
			systemFound = true
		}
	}

	if !systemFound {
		err = errors.New("Cannot delete entity because it doesn't exist")
		return err
	} else {
		world.Systems = newSystems
	}

	return nil
}

func (world *World) Update() {
	for _, system := range world.GetSystems() {
		systemLocalised := *system
		systemLocalised.Update()
	}
}
