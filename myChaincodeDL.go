package main


import (
//	"errors"
//	"fmt"
//	"strconv"
//	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


//********************************
//         Trade Repository
//********************************

type TradeRepository struct {
	chainArray ChainArray
}

func (self *TradeRepository) init(stub *shim.ChaincodeStub) bool {
	self.chainArray.init(stub, tradeIndex)
	
	return true
}

func (self *TradeRepository) newTrade(trade Trade) (int, error) {    //new trade returns "trade" or takes trade???
	index, err := self.chainArray.appendValue(trade)
	if err != nil {
		return -1, err
	}
	return index, nil
}

func (self *TradeRepository) getTradeByPostion(index int) (Trade, error) {
	var trade Trade
	
	err := self.chainArray.get(index, &trade)
	if err != nil {
		return trade, err
	}
	
	return trade, nil
}


func (self *TradeRepository) updateTrade(index int, trade Trade) (string, error) {
	key, err := self.chainArray.put(index, trade)
	if err != nil {
		return "", err
	}
	
	return key, nil
}


//********************************
//         Trade
//********************************

type Trade struct {
	UserID		string	`json:"userid"`
	SecurityID	string	`json:"securityid"`
	SecurityPointer	string	`json:"securitypointer"`
	TransType	string	`json:"transtype"`
	Price		float64	`json:"price"`
	Units		int		`json:"units"`
	Status		string	`json:"status"`
	Expiry		string	`json:"expiry"`
	Fulfilled	int		`json:"fulfilled"`
}

func (self *Trade) init(userID string, securityID string, securityPointer string, transType string, price float64, units int, expiry string) bool {
	self.UserID = userID
	self.SecurityID = securityID
	self.SecurityPointer = securityPointer
	self.TransType = transType
	self.Price = price
	self.Units = units
	self.Status = "Active"
	self.Expiry = expiry
	self.Fulfilled = 0
	
	return true
}



//********************************
//         User Repository
//********************************

type UserRepository struct {
	LinkedList ChainLinkedList
}

func (self *UserRepository) init(stub *shim.ChaincodeStub) bool {
	self.LinkedList.init(stub, userIndex)
	return true
}


func (self *UserRepository) newUser(userID string, ballance int) (string,error) {
	
	var user User
	user.init(userID, ballance)
	
	key, err := self.LinkedList.put(userID, user)
	if err != nil {
		return "", err
	}
	
	return key, nil
}

func (self *UserRepository) getUser(userId string) (User, error) {
	var user User
	
	err := self.LinkedList.get(userId, &user)
	if err != nil {
		return user, err
	}
	
	return user, nil
}

func (self *UserRepository) updateUser(userID string, user User) (string, error) {
	key, err := self.LinkedList.put(userID, user)
	if err != nil {
		return "", err
	}
	
	return key, nil
}

func (self *UserRepository) deleteUser(userID string) error {
	err :=  self.LinkedList.del(userID)
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