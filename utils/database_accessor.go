package utils

import (
	"gopkg.in/mgo.v2"
	"log"
)

type DatabaseAccessor struct {
	session			*mgo.Session
	url				string
	name			string
}

func NewDatabaseAccessor(url string, name string) *DatabaseAccessor {
	session, _ := mgo.Dial(url)
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)

	addIndexes(session, name)

	return &DatabaseAccessor{session, url, name}
}

func (dba *DatabaseAccessor) GetDB() *mgo.Database {
	db := dba.session.DB(dba.name)
	return db
}

func addIndexes(s *mgo.Session, dbname string) {
	var err error
	session := s.Copy()
	defer session.Close()

	// USERS INDEX
	userIndex := mgo.Index{
		Key: []string{"email"},
		Unique: true,
		Background: true,
		Sparse: true,
	}
	userCol := session.DB(dbname).C("users")
	err = userCol.EnsureIndex(userIndex)
	if err != nil {
		log.Fatalf("[addIndex]: %s\n", err)
	}

}