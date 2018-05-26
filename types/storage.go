package types

import (
	"github.com/globalsign/mgo/bson"
)

// Storer declares a set of CRUD functions for persisting and accessing data
type Storer interface {
	// Player account interface
	PlayerCreate(player Player) (id bson.ObjectId, err error)
	PlayerGetByName(name string) (player Player, err error)
	PlayerGetByID(id bson.ObjectId) (player Player, err error)
	PlayerUpdate(id bson.ObjectId, player Player) (err error)
	PlayerRemove(id bson.ObjectId) (err error)

	// Admin interface
	AdminSetLevel(id bson.ObjectId, level int32) (err error)
	AdminGetList() (result []Admin, err error)

	// misc
	DeleteEverythingPermanently() error
}
