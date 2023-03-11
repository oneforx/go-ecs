package ecs

import "fmt"

type FeedBackType string

var (
	FB_ERROR   FeedBackType = "FB_ERROR"
	FB_SUCCESS FeedBackType = "FB_SUCCESS"
)

type FeedBack struct {
	Host    string // Le nom du bloc ou plus précisément le nom de la fonction qui host les jobs
	Job     string // Le nom de la fonction qui a provoqué le feedback
	Label   string
	Type    FeedBackType
	Comment string
	Data    interface{}
}

func (fb *FeedBack) String() string {
	return fmt.Sprintf("[%v][%v][%v][%v]: %v", fb.Type, fb.Host, fb.Job, fb.Label, fb.Comment)
}

type Composition struct {
	Id    Identifier
	Value []string
}

func (source Composition) Len() int {
	return len(source.Value)
}

func (source Composition) Equals(target Composition) bool {
	if source.Len() != target.Len() {
		return false
	}
	seen := make(map[string]bool)
	for _, s := range source.Value {
		seen[s] = true
	}
	for _, t := range target.Value {
		if !seen[t] {
			return false
		}
	}
	return true
}

// Identifier enable the possibility to have two object with the same path but with a different namespace
// Example, a mod could add "mymod:position" and "anothermod:position"
// You can call World.GetEntityByComponentId("mymod:position")
type Identifier struct {
	Namespace string `json:"namespace"`
	Path      string `json:"path"`
}

func (id Identifier) String() string {
	return fmt.Sprintf("%s:%s", id.Namespace, id.Path)
}

func (id Identifier) Equals(other Identifier) bool {
	return id.String() == other.String()
}
