package user

import (
	"fmt"
	"sort"
)

type User struct {
	Id       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

type DB struct {
	users      map[int]*User
	idIterator int
}

var Admin = &User{
	Id:       0,
	Name:     "admin",
	Password: "admin",
}

func NewDB() *DB {
	db := &DB{
		users:      map[int]*User{},
		idIterator: 1,
	}
	db.CreateUser(Admin)
	db.CreateUser(&User{
		Name:     "Alice",
		Password: "123456",
	})
	db.CreateUser(&User{
		Name:     "Bob",
		Password: "123456",
	})
	return db
}

func (db *DB) CreateUser(u *User) (*User, error) {
	for _, user := range db.users {
		if user.Name == u.Name {
			return nil, fmt.Errorf("user name %s already exists", u.Name)
		}
	}
	u.Id = db.idIterator
	db.idIterator++
	db.users[u.Id] = u
	return u, nil
}

func (db *DB) DeleteUser(id int) {
	delete(db.users, id)
}

func (db *DB) UpdateUser(u *User) (*User, error) {
	if _, ok := db.users[u.Id]; !ok {
		return nil, fmt.Errorf("user id %d not found", u.Id)
	}
	db.users[u.Id] = u
	return u, nil
}

func (db *DB) Login(name, password string) *User {
	for _, u := range db.users {
		if u.Name == name && u.Password == password {
			return u
		}
	}
	return nil
}

func (db *DB) GetUser(id int) *User {
	return db.users[id]
}

func (db *DB) GetUsers() []*User {
	users := make([]*User, 0, len(db.users))
	for _, u := range db.users {
		users = append(users, u)
	}
	sort.Sort(sortUsers(users))
	return users
}

type sortUsers []*User

func (s sortUsers) Len() int {
	return len(s)
}

func (s sortUsers) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

func (s sortUsers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
