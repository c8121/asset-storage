package metadata_db_entity

import (
	"time"
)

type Collection struct {
	Id      int64
	UUID    string    `dbe:"Size: 64, Index"`
	Name    string    `dbe:"Size: 255"`
	Created time.Time `dbe:"Index"`
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
