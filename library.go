package ecs

import "errors"

type ILibrary interface {
	GetId() Identifier
	GetStruct() Library

	GetCompositions() []Composition
	GetComponents() []Component
	GetSystems() []ISystem

	RegisterSystem(ISystem) error
	RegisterSystems([]ISystem) error

	RegisterComponent(Component) error
	RegisterComponents([]Component) error
	RegisterComposition(Composition) error
	RegisterCompositions([]Composition) error
}
type Library struct {
	Id           Identifier
	components   []Component
	systems      []ISystem
	compositions []Composition
}

func (lib Library) GetId() Identifier {
	return lib.Id
}

func (lib Library) GetStruct() Library {
	return lib
}

func (lib Library) GetCompositions() []Composition {
	return lib.compositions
}

func (lib Library) GetComponents() []Component {
	return lib.components
}

func (lib Library) GetSystems() []ISystem {
	return lib.systems
}

func (lib *Library) RegisterSystem(system ISystem) error {
	if lib.systemExists(system) {
		return errors.New("system already exists")
	}
	lib.systems = append(lib.systems, system)
	return nil
}

func (lib *Library) RegisterSystems(systems []ISystem) error {
	for _, system := range systems {
		if lib.systemExists(system) {
			return errors.New("system already exists")
		}
	}
	lib.systems = append(lib.systems, systems...)
	return nil
}

func (lib *Library) RegisterComponent(component Component) error {
	if lib.componentExists(component.Id) {
		return errors.New("component already exists")
	}
	lib.components = append(lib.components, component)
	return nil
}

func (lib *Library) RegisterComponents(components []Component) error {
	for _, component := range components {
		if lib.componentExists(component.Id) {
			return errors.New("component already exists")
		}
	}
	lib.components = append(lib.components, components...)
	return nil
}

func (lib *Library) RegisterComposition(composition Composition) error {
	if lib.compositionExists(composition.Id) {
		return errors.New("composition already exists")
	}
	lib.compositions = append(lib.compositions, composition)
	return nil
}

func (lib *Library) RegisterCompositions(compositions []Composition) error {
	for _, composition := range compositions {
		if lib.compositionExists(composition.Id) {
			return errors.New("composition already exists")
		}
		lib.compositions = append(lib.compositions, composition)
	}
	return nil
}

func (lib Library) componentExists(id Identifier) bool {
	for _, component := range lib.components {
		if component.Id == id {
			return true
		}
	}
	return false
}

func (lib Library) systemExists(system ISystem) bool {
	for _, s := range lib.systems {
		if s.GetId() == system.GetId() {
			return true
		}
	}
	return false
}

func (lib Library) compositionExists(id Identifier) bool {
	for _, composition := range lib.compositions {
		if composition.Id == id {
			return true
		}
	}

	return false
}
