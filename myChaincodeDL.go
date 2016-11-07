package main


import (
//	"errors"
//	"fmt"
//	"strconv"
//	"encoding/json"
//	"github.com/hyperledger/fabric/core/chaincode/shim"
)


//********************************
//         User Repository
//********************************

type UserRepo struct {
	sc *SimpleChaincode
	tempUser User
}

func (self *UserRepo) Init(sc *SimpleChaincode) bool {
	self.sc = sc
	return true
}


func (self *UserRepo) RegisterUser(user User) bool {
	self.tempUser = user
	return true
}

func (self *UserRepo) GetUser(userId int) User {
	return self.tempUser
}

func (self *UserRepo) UpdateUser(userId int, user User) bool {
	self.tempUser = user
	return true
}

func (self *UserRepo) DeleteUser(userId int) bool {
	//self.tempUser = nil
	return true
}

func (self *UserRepo) GetUserList() []User {
	ul := []User {self.tempUser}
	
	return ul
}