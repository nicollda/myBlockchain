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
	hashmap HashMap
}

func (self *UserRepository) Init(stub *shim.ChaincodeStub) bool {
	self.hashmap.Init(stub, userIndex)
	return true
}


func (self *UserRepository) NewUser(userID string, ballance int) bool {
	var user User
	
	user.UserID = userID
	user.Status = "Active"
	user.Ballance = ballance
	
	self.hashmap.put(userID, user)

	return true
}

func (self *UserRepository) GetUser(userId string) (User, error) {
	var user User
	err := self.hashmap.get(userId, user)
	if err != nil {
		return user, err
	}
	
	return user, nil
}

func (self *UserRepository) UpdateUser(userID string, user User) bool {
	self.hashmap.put(userID, user)
	return true
}

func (self *UserRepository) DeleteUser(userID string) bool {
	//self.tempUser = nil
	self.hashmap.del(userID)
	return true
}