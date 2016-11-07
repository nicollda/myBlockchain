package main


import (
//	"errors"
//	"fmt"
//	"strconv"
//	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


//********************************
//         User Repository
//********************************

type UserRepository struct {
	stub *shim.ChaincodeStub
	tempUser User
}

func (self *UserRepository) Init(stub *shim.ChaincodeStub) bool {
	self.stub = stub
	return true
}


func (self *UserRepository) NewUser(userID string) bool {
	self.tempUser.UserID = userID
	self.tempUser.Status = "Active"
	self.tempUser.Ballance = initialCash

	return true
}

func (self *UserRepository) GetUser(userId int) User {
	return self.tempUser
}

func (self *UserRepository) UpdateUser(userId int, user User) bool {
	self.tempUser = user
	return true
}

func (self *UserRepository) DeleteUser(userId int) bool {
	//self.tempUser = nil
	return true
}

func (self *UserRepository) GetUserList() []User {
	ul := []User {self.tempUser}
	
	return ul
}