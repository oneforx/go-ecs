// Package that describe an Entity-Component-System World
package goecs

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type IWorld interface {
	GetId() Identifier

	// Ajoute une entité
	AddEntity(*IEntity) *FeedBack
	// Ajoute une liste d'entité
	AddEntities([]Entity) *FeedBack

	// Retourne une entité par son instance ID
	GetEntity(uuid.UUID) *IEntity
	// Get all entities instantiated in the world
	GetEntities() []*IEntity

	// Return des une copie de toute les entité sans cycle (sans référence au monde qui référence encore toute les entité, c'est infini pour type *IEntity)
	GetEntitiesNoCycle() []EntityNoCycle
	// Add entities in the world from an array of Entity without cycle reference
	AddEntitiesFromEntitiesNoCycle(entitiesNoCycle []EntityNoCycle)
	// UUID should be an external ID maybe a ClientID
	GetEntitiesPossessedBy(uuid.UUID) []*IEntity
	GetEntitiesByComponentId(Identifier) []*IEntity
	// EN: Get entities which have at least the specified components name.
	//     Returns a slice of IEntity pointers.
	// FR: Obtiens les entités qui ont au moins les composants spécifiés.
	//     Retourne une tranche de pointeurs IEntity.
	GetEntitiesWithComponents(...Identifier) []*IEntity
	GetEntitiesWithComponentsIdString(...string) []*IEntity
	// Get entities wich have the components specified by the array of component name.
	GetEntitiesWithStrictComposition(Composition) []*IEntity
	// EN: Get entities which at least the specified composition
	// 		Returns a slice of IEntity pointers
	// FR: Obtiens les entités qui ont au moins la composition spécifié
	//     Retourne une tranche de pointeurs IEntity.
	GetEntitiesWithComposition(Composition) []*IEntity
	// Update a list of component of the entity
	UpdateEntityComponents(uuid.UUID, []*Component) *FeedBack
	// Remove an entity by his instanceId
	RemoveEntity(uuid.UUID) *FeedBack

	// Ajoute un système, veuillez ajouter une structure Système qui implémente l'interface ISystem
	AddSystem(*ISystem)
	//
	GetSystemById(Identifier) *ISystem
	RemoveSystem(Identifier) *FeedBack

	// Execute toutes les methods UpdateClient des systèmes lié à ce monde
	UpdateClient()
	// Execute toutes les methods UpdateServer des systèmes lié à ce monde
	UpdateServer()
}

type World struct {
	Id            Identifier
	Entities      []*IEntity
	Systems       []*ISystem
	entitiesMutex sync.RWMutex
}

func (world *World) GetId() Identifier {
	return world.Id
}

func (world *World) AddEntity(entity *IEntity) (err *FeedBack) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
	var found bool = false
	for _, ent := range world.Entities {
		entityLocalised := *entity
		entLocalised := *ent
		if entLocalised.GetId() == entityLocalised.GetId() {
			found = true
		}
	}
	if found {
		return &FeedBack{
			Host:    "AddEntity",
			Job:     "for",
			Label:   "ENTITY_SAME_ID_EXIST",
			Comment: "Une entité avec le même identifiant existe déjà",
		}
	} else {
		world.Entities = append(world.Entities, entity)
	}
	return err
}

func (world *World) AddEntities(entities []Entity) (fb *FeedBack) {
	for _, entity := range entities {
		if world.GetEntity(entity.Id) != nil {
			if fb == nil {
				fb = &FeedBack{
					Host:    "AddEntities",
					Job:     "GetEntity",
					Label:   "ENTITY_SAME_ID_EXIST",
					Comment: "Les entités dans la liste sont ceux qui n'ont pas été ajouté car des entités avec le même id existe déjà",
					Data:    []uuid.UUID{},
				}
			} else {
				fbData, ok := fb.Data.([]uuid.UUID)
				if !ok {
					log.Println("Could not parse feedback data to []uuid.UUID")
				}
				fb.Data = append(fbData, entity.Id)
			}
			continue
		}
		var entityLocation IEntity = &entity
		world.AddEntity(&entityLocation)
	}

	return fb
}

func (world *World) AddEntitiesFromEntitiesNoCycle(entitiesNoCycle []EntityNoCycle) {
	var w IWorld = world
	for _, entityNoCycle := range entitiesNoCycle {
		var entity IEntity = &Entity{
			Id:          entityNoCycle.Id,
			OwnerID:     entityNoCycle.OwnerID,
			PossessedID: entityNoCycle.OwnerID,
			World:       &w,
			Components:  entityNoCycle.Components,
		}
		world.AddEntity(&entity)
	}
}

func (world *World) GetEntity(id uuid.UUID) (ent *IEntity) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
	for _, entity := range world.Entities {
		entityLocalised := *entity
		if entityLocalised.GetId() == id {
			ent = entity
		}
	}
	return ent
}

func (world *World) GetEntities() (entities []*IEntity) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
	return world.Entities
}

func (world *World) GetEntitiesNoCycle() (entities []EntityNoCycle) {
	for _, entity := range world.GetEntities() {
		entityLocalised := *entity

		entities = append(entities, EntityNoCycle{
			Id:          entityLocalised.GetId(),
			OwnerID:     entityLocalised.GetOwnerID(),
			PossessedID: entityLocalised.GetPossessedID(),
			Components:  entityLocalised.GetComponents(),
		})
	}

	return entities
}

func (world *World) GetEntitiesPossessedBy(possessedId uuid.UUID) (entities []*IEntity) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
	for _, entity := range world.Entities {
		entityLocalised := *entity
		if entityLocalised.GetPossessedID() == possessedId {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (world *World) GetEntitiesByComponentId(id Identifier) (entities []*IEntity) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
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

func (world *World) GetEntitiesWithComponents(v ...Identifier) (entities []*IEntity) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
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

func (world *World) GetEntitiesWithComponentsIdString(v ...string) (entities []*IEntity) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
	for _, entity := range world.Entities {
		entityLocalised := *entity
		var checkeds int = 0
		for _, cmpName := range v {
			if entityLocalised.HaveComponentByIdString(cmpName) {
				checkeds = checkeds + 1
			}
		}
		if checkeds == len(v) {
			entities = append(entities, entity)
		}
	}
	return entities
}
func (world *World) GetEntitiesWithComposition(composition Composition) (entities []*IEntity) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
	for _, entity := range world.GetEntities() {
		entityLocalised := *entity
		if entityLocalised.HaveComposition(composition.Value) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (world *World) GetEntitiesWithStrictComposition(composition Composition) (entities []*IEntity) {
	for _, entity := range world.GetEntities() {
		entityLocalised := *entity
		if len(composition.Value) == len(entityLocalised.GetComponents()) && entityLocalised.HaveComposition(composition.Value) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (world *World) UpdateEntityComponents(id uuid.UUID, components []*Component) *FeedBack {
	entity := world.GetEntity(id)
	if entity != nil {
		entityLocalised := *entity
		entityLocalised.UpdateComponents(components)
	}
	return nil
}

func (world *World) RemoveEntity(id uuid.UUID) (err *FeedBack) {
	world.entitiesMutex.Lock()
	defer world.entitiesMutex.Unlock()
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
		return &FeedBack{
			Host:    "RemoveEntity",
			Job:     "!entityFound",
			Label:   "ENTITY_DOES_NOT_EXIST",
			Comment: "The entity " + id.String() + " doesn't exist.",
		}
	}

	world.Entities = newEntities

	return nil
}

func (world *World) AddSystem(sys *ISystem) {
	world.Systems = append(world.Systems, sys)
}

func (world *World) GetSystemById(id Identifier) *ISystem {
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

func (world *World) RemoveSystem(id Identifier) (err *FeedBack) {
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
		return &FeedBack{
			Host:    "RemoveSystem",
			Job:     "!systemFound",
			Label:   "SYSTEM_DOES_NOT_EXIST",
			Comment: "The system " + id.String() + " doesn't exist.",
		}
	} else {
		world.Systems = newSystems
	}

	return nil
}

func (world *World) UpdateClient() {
	for _, system := range world.GetSystems() {
		systemLocalised := *system
		systemLocalised.UpdateClient()
	}
}

func (world *World) UpdateServer() {
	for _, system := range world.GetSystems() {
		systemLocalised := *system
		systemLocalised.UpdateServer()
	}
}
