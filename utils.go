package goecs

func EntityToNoCycle(entity *IEntity) EntityNoCycle {
	entityLocalised := *entity
	return EntityNoCycle{
		Id:          entityLocalised.GetId(),
		OwnerID:     entityLocalised.GetOwnerID(),
		PossessedID: entityLocalised.GetPossessedID(),
		Components:  entityLocalised.GetComponents(),
	}
}

func EntityNoCycleToEntity(world *IWorld, entityNoCycle EntityNoCycle) *IEntity {
	var entity IEntity = &Entity{
		Id:          entityNoCycle.Id,
		OwnerID:     entityNoCycle.OwnerID,
		PossessedID: entityNoCycle.OwnerID,
		World:       world,
	}
	return &entity
}

func EntityNoCycleArrayToEntityArray(world IWorld, entitiesNoCycle []EntityNoCycle) (entities []Entity) {
	for _, entityNoCycle := range entitiesNoCycle {
		entities = append(entities, Entity{
			Id:          entityNoCycle.Id,
			OwnerID:     entityNoCycle.OwnerID,
			PossessedID: entityNoCycle.OwnerID,
			World:       &world,
		})
	}
	return entities
}
