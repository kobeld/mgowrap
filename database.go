package mgowrap

import "gopkg.in/mgo.v2"

type Database struct {
	DialString string
	Name       string
}

func NewDatabase(dialString string, name string) (db *Database) {
	db = &Database{
		DialString: dialString,
		Name:       name,
	}
	return
}

// Key: DialString
// Value: mgo Session
var ConnectedSessions map[string]*mgo.Session

func (db *Database) GetOrDialSession() (session *mgo.Session) {
	if db == nil {
		panic("mgo: Database is nil")
	}

	if db.Name == "" || db.DialString == "" {
		panic("mgo: must provide valid dialstring and database name")
	}
	if ConnectedSessions == nil {
		ConnectedSessions = make(map[string]*mgo.Session)
	}

	session = ConnectedSessions[db.DialString]
	if session == nil {
		var err error
		session, err = mgo.Dial(db.DialString)
		if err != nil {
			panic(err)
		}
		ConnectedSessions[db.DialString] = session
	}
	return
}

func (this *Database) DropDatabase() (err error) {
	this.DatabaseDo(func(db *mgo.Database) {
		err = db.DropDatabase()
	})
	return
}

func (this *Database) DatabaseDo(f func(db *mgo.Database)) {
	s := this.GetOrDialSession().Copy()
	defer s.Close()
	f(s.DB(this.Name))
}

func (db *Database) CollectionDo(name string, f func(c *mgo.Collection)) {
	s := db.GetOrDialSession().Copy()
	defer s.Close()
	f(s.DB(db.Name).C(name))
}

func (db *Database) CollectionsDo(f func(c ...*mgo.Collection), names ...string) {
	s := db.GetOrDialSession().Copy()
	defer s.Close()
	cs := []*mgo.Collection{}
	for _, name := range names {
		cs = append(cs, s.DB(db.Name).C(name))
	}
	f(cs...)
}
