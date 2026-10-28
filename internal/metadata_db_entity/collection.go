package metadata_db_entity

import (
	"time"
)

type Collection struct {
	Id      int64
	UUID    string
	Name    string
	Created time.Time
}

func (c *Collection) GetEntityName() string {
	return "Collection"
}

func (c *Collection) GetId() int64 {
	return c.Id
}

func (c *Collection) SetId(id int64) {
	c.Id = id
}
