package blockchain

import (
	"strconv"
	"encoding/json"
)

type UserRepo struct {
	sc *SimpleChaincode
}

func (self *UserRepo) RegisterUser(user User) bool {
	defer return nil == err
	lastUserIndex, err = self.sc.Get("Last" + USERINDEX)
	lastUserIndex += 1;
	_, err = self.sc.Push(strconv.Itoa(lastUserIndex), User)
}

func (self *UserRepo) GetUser(userId int) User {
	defer {
		if nil != err {
			return nil
		}
	}
	jsonstr, err = self.sc.Get(USERINDEX + strconv.Itoa(userId))
	var user User
	json.Unmarshal(jsonstr, user)
	return user;
}

func (self *UserRepo) UpdateUser(userId int, user User) bool {

}

func (self *UserRepo) DeleteUser(userId int) bool {

}

func (self *UserRepo) GetUserList() []User {

}
