package goecs

import "fmt"

type LibraryManager struct {
	Libraries map[Identifier]ILibrary
}

func (lm *LibraryManager) GetLibrary(id Identifier) *ILibrary {
	library, ok := lm.Libraries[id]
	if !ok {
		return nil
	}
	return &library
}

// Return false if we couldn't add the library
func (lm *LibraryManager) AddLibrary(id Identifier, newLibrary ILibrary) bool {
	if lm.GetLibrary(id) != nil {
		return false
	}

	lm.Libraries[id] = newLibrary
	return true
}

func (lm *LibraryManager) GetComponents() []Component {
	var components []Component = []Component{}

	for _, lib := range lm.Libraries {
		components = append(components, lib.GetComponents()...)
	}

	return components
}

func (lm *LibraryManager) GetComponent(id Identifier) (cmp Component) {
	for _, library := range lm.Libraries {
		if library.GetId().Namespace == id.Namespace {
			for _, component := range library.GetComponents() {
				if component.Id.Path == id.Path {
					cmp = component
					break
				}
			}
			break
		}
	}
	return cmp
}

func (lm *LibraryManager) GetSystem(id Identifier) (sys ISystem, err error) {
	for _, library := range lm.Libraries {
		if library.GetId().Namespace == id.Namespace {
			for _, system := range library.GetSystems() {
				if system.GetId().String() == id.String() {
					sys = system
					break
				}
			}
			break
		}
	}

	if sys == nil {
		return nil, fmt.Errorf("System not found")
	}

	return sys, nil
}

func (lm *LibraryManager) InstantiateSystem(id Identifier, world *IWorld) (*ISystem, error) {
	var systemLocation *ISystem

	system, err := lm.GetSystem(id)
	if err != nil {
		return nil, err
	}

	system.Init(world)

	systemLocation = &system

	return systemLocation, nil
}

func (lm *LibraryManager) InstantiateComponent(id Identifier, data interface{}) *Component {
	var componentLocation *Component

	component := lm.GetComponent(id)

	component.SetData(data)

	componentLocation = &component

	return componentLocation
}

func (lm *LibraryManager) LoadLibrary(library ILibrary) error {
	_, ok := lm.Libraries[library.GetId()]
	if !ok {
		lm.Libraries[library.GetId()] = library
	}

	return fmt.Errorf("CONFLICT: a library with same id already exist")
}
