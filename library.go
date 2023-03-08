package ecs

import "fmt"

type ILibrary interface {
	GetId() Identifier
	GetStruct() Library
	GetComponents() []Component
	GetSystems() []ISystem
	AddComponent(Component)
	AddSystem(ISystem)
}

type Library struct {
	Id           Identifier
	components   []Component
	systems      []ISystem
	compositions map[Identifier]Composition
}

func (library *Library) GetId() Identifier {
	return library.Id
}

func (library *Library) AddSystem(system ISystem) {
	library.systems = append(library.systems, system)
}

func (library *Library) AddComposition(id Identifier, composition Composition) error {
	compositionFound, ok := library.compositions[id]
	if !ok {
		library.compositions[id] = composition
		return nil
	} else {
		return fmt.Errorf("composition %s found", compositionFound)
	}
}

func (library *Library) SetCompositions(compositions map[Identifier]Composition) {
	library.compositions = compositions
}

func (library *Library) AddComponent(component Component) {
	library.components = append(library.components, component)
}

func (library *Library) GetComponents() []Component {
	return library.components
}

func (library *Library) GetSystems() []ISystem {
	return library.systems
}

func (library *Library) GetStruct() Library {
	return *library
}
