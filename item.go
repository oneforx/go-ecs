package ecs

type IItem interface {
	GetId() string
	GetName() string
	GetDescription() string
	Use(*IWorld) string
}

const (
	SLOT_LEFTHAND  = "SLOT_LEFTHAND"
	SLOT_RIGHTHAND = "SLOT_RIGHTHAND"
	SLOT_CHEST     = "SLOT_CHEST"
	SLOT_FOOT      = "SLOT_FOOT"
)

type Item struct {
	id          string
	name        string
	description string
	slotType    string
}

type ModelItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SlotType    string `json:"slotType"`
}

func (item *Item) GetId() string {
	return item.id
}

func (item *Item) GetName() string {
	return item.name
}

func (item *Item) GetDescription() string {
	return item.description
}

func (item *Item) GetSlotType() string {
	return item.slotType
}
