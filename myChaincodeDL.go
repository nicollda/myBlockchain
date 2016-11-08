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


func (self *UserRepository) NewUser(userID string, ballance int) (string,error) {
	var user User
	
	user.UserID = userID
	user.Status = "Active"
	user.Ballance = ballance
	
	key, err := self.hashmap.put(userID, user)
	if err != nil {
		return "", err
	}
	
	return key, nil
}

func (self *UserRepository) GetUser(userId string) (User, error) {
	var user User
	
	err := self.hashmap.get(userId, &user)
	if err != nil {
		return user, err
	}
	
	return user, nil
}

func (self *UserRepository) UpdateUser(userID string, user User) (string, error) {
	key, err := self.hashmap.put(userID, user)
	if err != nil {
		return "", err
	}
	
	return key, nil
}

func (self *UserRepository) DeleteUser(userID string) error {
	err :=  self.hashmap.del(userID)
	if err != nil {
		return err
	}
	
	return nil
}