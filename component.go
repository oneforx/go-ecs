package ecs

type IComponent interface {
	GetId() string
	GetData(name string) interface{}
	SetData(name string, v interface{})
}

type ModelComponent map[string]interface{}

type Component struct {
	Id   string
	Data map[string]interface{}
}

func CreateComponent(id string, data map[string]interface{}) *IComponent {
	var component IComponent = &Component{Id: id, Data: data}
	return &component
}

func (p *Component) GetId() string {
	return p.Id
}

func (p *Component) GetData(name string) interface{} {
	return p.Data[name]
}

func (p *Component) SetData(name string, v interface{}) {
	p.Data[name] = v
}
