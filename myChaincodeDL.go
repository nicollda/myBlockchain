package main

import (
	"errors"
//	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)



//********************************
//         Holdings Repository
//********************************
type HoldingsRepository struct {
	LinkedList ChainLinkedList
}


func (self *HoldingsRepository) init(stub *shim.ChaincodeStub) bool {
	self.LinkedList.init(stub, holdingIndex)
	
	return true
}



func (self *HoldingsRepository) getHoldingID(userID string, securityID string) string {
	return userID + securityID
}



func (self *HoldingsRepository) newHolding(userID string, securityID string, units int) (string,error) {
	
	var holding Holding
	holding.init(userID, securityID, units)
	
	key, err := self.LinkedList.put(self.getHoldingID(userID, securityID), holding)
	if err != nil {
		return "", err
	}
	
	return key, nil
}



func (self *HoldingsRepository) getHolding(userID string, securityID string) (Holding, error) {
	var holding Holding
	
	holding, err := self.getHoldingByID(self.getHoldingID(userID, securityID))
	if err != nil {
		return holding, errors.New("this is the error")
	}
	
	return holding, nil
}



func (self *HoldingsRepository) getHoldingByID(holdingID string) (Holding, error) {
	var holding Holding
	
	err := self.LinkedList.get(holdingID, &holding)
	if err != nil {
		return holding, errors.New("this is the error2")
	}
	
	return holding, nil
}



func (self *HoldingsRepository) getFirstHolding() (Holding, error) {
	var holding Holding
	
	err := self.LinkedList.getFirst(&holding)
	if err != nil {
		return holding, err
	}
	
	return holding, nil
}



func (self *HoldingsRepository) getNextHolding() (Holding, error) {
	var holding Holding
	
	err := self.LinkedList.getNext(&holding)
	if err != nil {
		return holding, err
	}
	
	return holding, nil
}



func (self *HoldingsRepository) updateHolding(holding Holding) (string, error) {
	key, err := self.LinkedList.put(self.getHoldingID(holding.UserID, holding.SecurityID), holding)
	if err != nil {
		return "", err
	}
	
	return key, nil
}



func (self *HoldingsRepository) deleteHolding(userID string, securityID string,) error {
	err := self.LinkedList.del(self.getHoldingID(userID, securityID))
	if err != nil {
		return err
	}
	
	return nil
}



//********************************
//         Holding
//********************************
type Holding struct {
	SecurityID	string	`json:"securityid"`
	UserID		string	`json:"userid"`
	Units		int		`json:"units"`
}



func (self *Holding) init(userID string, securityID string, units int) bool {
	self.UserID = userID
	self.SecurityID = securityID
	self.Units = units
	
	return true
}



func (self *Holding) updateUnits(units int) bool {
	self.Units = units
	
	return true
}



//********************************
//         Security Repository
//********************************
 
 type SecurityRepository struct {
	chainArray ChainArray
}



func (self *SecurityRepository) init(stub *shim.ChaincodeStub) bool {
	self.chainArray.init(stub, securityIndex)
	
	return true
}



func (self *SecurityRepository) newSecurity(security Security) (int, error) {    //new trade returns "trade" or takes trade???
	index, err := self.chainArray.appendValue(security)
	if err != nil {
		return -1, err
	}
	return index, nil
}



func (self *SecurityRepository) getSecurityPosition(securityID string) (int, error) {
	var security Security 
	returnVal := 0
	lastIndex, err := self.chainArray.getLastIndex()
	
	for i:= 1; i<=lastIndex; i++ {
		security, err = self.getSecurityByPostion(i)
		if err != nil {
			return 0, err
		}
		
		if security.SecurityID == securityID {
			returnVal = i
			i = lastIndex + 1
		}
	}

	return returnVal, nil

}



func (self *SecurityRepository) getSecurityByPostion(index int) (Security, error) {
	var security Security
	
	err := self.chainArray.get(index, &security)
	if err != nil {
		return security, err
	}
	
	return security, nil
}



func (self *SecurityRepository) updateSecurity(index int, security Security) (string, error) {
	key, err := self.chainArray.put(index, security)
	if err != nil {
		return "", err
	}
	
	return key, nil
}



//********************************
//         Security
//********************************


type Security struct {
	SecurityID		string	`json:"securityid"`
	Description		string	`json:"description"`
	Status			string	`json:"status"`
}



func (self *Security) init(securityID string, description string) bool {
	self.SecurityID = securityID
	self.Description = description
	self.Status = "Active"	//should be in BL
	
	return true
}



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



func (self *TradeRepository) getTradeByPosition(index int) (Trade, error) {
	var trade Trade
	
	err := self.chainArray.get(index, &trade)
	if err != nil {
		return trade, err
	}
	
	return trade, nil
}



func (self *TradeRepository) getLastIndex() (int, error) {
	lastIndex, err := self.chainArray.getLastIndex()
	if err != nil {
		return -1, err
	}
	
	return lastIndex, nil
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
	UserID			string	`json:"userid"`
	SecurityID		string	`json:"securityid"`
	SecurityPointer	int		`json:"securitypointer"`
	TransType		string	`json:"transtype"`
	Price			float64	`json:"price"`
	Units			int		`json:"units"`
	Status			string	`json:"status"`
	Expiry			string	`json:"expiry"`
	Fulfilled		int		`json:"fulfilled"`
}



func (self *Trade) init(userID string, securityID string, securityPointer int, transType string, price float64, units int, expiry string) bool {
	self.UserID = userID
	self.SecurityID = securityID
	self.SecurityPointer = securityPointer
	self.TransType = transType
	self.Price = price
	self.Units = units
	self.Status = "Active"	//should be in BL
	self.Expiry = expiry
	self.Fulfilled = 0		//should be in BL
	
	return true
}



func (self *Trade) getUserID() string {
	return self.UserID
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



func (self *UserRepository) newUser(userID string, ballance int, status string) (string,error) {
	
	var user User
	user.init(userID, ballance, status)
	
	key, err := self.LinkedList.put(userID, user)
	if err != nil {
		return "", err
	}
	
	//debug code
	user, err = self.getUser(userID)
	if err != nil {
		return "", err
	}
	
	curOutByteA,err := self.LinkedList.stub.GetState("currentOutput")
	outByteA := []byte(string(curOutByteA) + ":::debug for userID " + user.UserID)
	err = self.LinkedList.stub.PutState("currentOutput", outByteA)
	
	
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



func (self *UserRepository) updateUser(user User) (string, error) {
	key, err := self.LinkedList.put(user.UserID, user)
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



func (self *User) init(userID string, ballance int, status string) bool {
	self.UserID = userID
	self.Status = status
	self.Ballance = ballance
	
	return true
}



func (self *User) getBallance() int {
	return self.Ballance
}