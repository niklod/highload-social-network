package user

import "github.com/tarantool/go-tarantool"

type TarantoolRepository struct {
	client *tarantool.Connection
}

func (t *TarantoolRepository) Create(user *User) (*User, error) {
	panic("implement me")
}
func (t *TarantoolRepository) List() ([]User, error) {
	panic("implement me")
}
func (t *TarantoolRepository) GetByID(id int) (*User, error) {
	panic("implement me")
}
func (t *TarantoolRepository) GetByLogin(login string) (*User, error) {
	panic("implement me")
}
func (t *TarantoolRepository) GetByFirstAndLastName(firstname, lastname string) ([]User, error) {
	panic("implement me")
}
func (t *TarantoolRepository) AddFriend(userId int, friendId int) error {
	panic("implement me")
}
func (t *TarantoolRepository) DeleteFriend(userId int, friendId int) error {
	panic("implement me")
}
func (t *TarantoolRepository) Friends(userId int) ([]User, error) {
	panic("implement me")
}
