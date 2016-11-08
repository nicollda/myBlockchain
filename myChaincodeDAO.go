package main


import (
//	"errors"
//	"fmt"
//	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


//********************************
//         Data Access Objects
//********************************


//Hashmap does not have the ability to iterate through the list
//returns []byte array.  maybe should have a mapType = "[]String" etc... ?
type HashMap struct {
	stub *shim.ChaincodeStub
	mapName string
}

func (self *HashMap) Init(stub *shim.ChaincodeStub, mapName string) bool {
	self.stub = stub
	self.mapName = mapName
	
	return true
}


func (self *HashMap) getKey(key string) string {
	return self.mapName + key
}



func (self *HashMap) get(key string, returnVal interface{}) error {
	mByteA, err := self.stub.GetState(self.getKey(key))
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(mByteA, &returnVal)
	if err != nil {
		return err
	}
	
	return nil
}



func (self *HashMap) put(key string, val interface{}) (string, error) {
	mKey := self.getKey(key)
	mByteA, err := json.Marshal(&val)
	if err != nil {
		return "", err
	}
	
	err = self.stub.PutState(mKey, mByteA)
	if err != nil {
		return "", err
	}
	
	return mKey, nil
}

func (self *HashMap) del(key string) error {
	err := self.stub.DelState(self.getKey(key))
	if err != nil {
		return err
	}
	
	return nil
}