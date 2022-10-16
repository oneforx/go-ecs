package ecs

type ISystem interface {
	GetId() string
	Update()
}
