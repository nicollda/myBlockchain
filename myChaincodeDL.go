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

func (self *UserRepository) init(stub *shim.ChaincodeStub) bool {
	self.hashmap.init(stub, userIndex)
	return true
}


func (self *UserRepository) newUser(userID string, ballance int) (string,error) {
	
	var user User
	user.init(userID, ballance)
	
	key, err := self.hashmap.put(userID, user)
	if err != nil {
		return "", err
	}
	
	return key, nil
}

func (self *UserRepository) getUser(userId string) (User, error) {
	var user User
	
	err := self.hashmap.get(userId, &user)
	if err != nil {
		return user, err
	}
	
	return user, nil
}

func (self *UserRepository) updateUser(userID string, user User) (string, error) {
	key, err := self.hashmap.put(userID, user)
	if err != nil {
		return "", err
	}
	
	return key, nil
}

func (self *UserRepository) deleteUser(userID string) error {
	err :=  self.hashmap.del(userID)
	if err != nil {
		return err
	}
	
	return nil
}


//********************************
//         User
//********************************

type User struct {
	UserID		string	`json:"userID"`
	Status		string	`json:"status"`
	Ballance	int		`json:"ballance"`
}

func (self *User) init(userID string, ballance int) bool {
	self.UserID = userID
	self.Status = "Active"
	self.Ballance = ballance
	
	return true
}