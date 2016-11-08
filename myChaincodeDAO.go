package main


import (
//	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


//********************************
//         Data Access Objects
//********************************




//********************************************************************************
//***** HashMap
//********************************************************************************

//Hashmap does not have the ability to iterate through the list
//returns []byte array.  maybe should have a mapType = "[]String" etc... ?
type HashMap struct {
	stub *shim.ChaincodeStub
	mapName string
}

func (self *HashMap) init(stub *shim.ChaincodeStub, mapName string) bool {
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


//********************************************************************************
//***** ChainArray (cannot call it array)
//********************************************************************************

//Hashmap does not have the ability to iterate through the list
//returns []byte array.  maybe should have a mapType = "[]String" etc... ?
type ChainArray struct {
	stub *shim.ChaincodeStub
	ArrayName string
}

func (self *ChainArray) init(stub *shim.ChaincodeStub, arrayName string) bool {
	self.stub = stub
	self.ArrayName = arrayName
	
	return true
}


func (self *ChainArray) getKey(key int) string {
	return self.ArrayName + strconv.Itoa(key)
}


func (self *ChainArray) get(index int, returnVal interface{}) error {
	vByteA, err := self.stub.GetState(self.getKey(index))
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(vByteA, &returnVal)
	if err != nil {
		return err
	}
	
	return nil
}



func (self *ChainArray) put(index int, val interface{}) (string, error) {
	mKey := self.getKey(index)
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

func (self *ChainArray) getLastIndex() (int, error) {
	var lastIndex int
	
	lastIndexByteA, err := self.stub.GetState("Last" + self.ArrayName)
	if err != nil {
		return -1, err
	}
	
	lastIndex, err = strconv.Atoi(string(lastIndexByteA))
	if err != nil {
		lastIndex = 0
	}
	
	return lastIndex, nil
}


func (self *ChainArray) getNextIndex() (int, error) {
	fmt.Printf("Running getNextIndex")

	lastIndex, err := self.getLastIndex()
	if err != nil {
		lastIndex = 1
	} else { 
		lastIndex = lastIndex + 1
	}
	
	lastIndexByteA := []byte(strconv.Itoa(lastIndex))   
	err = self.stub.PutState("Last" + self.ArrayName, lastIndexByteA)
	if err != nil {
		return -1, err
	}
	
	
	fmt.Printf(strconv.Itoa(lastIndex))
	
	return lastIndex, nil
}


func (self *ChainArray) appendValue(val interface{}) (int, error) {
	index, err := self.getNextIndex()
	if err != nil {
		return -1, err
	}
	
	mKey := self.getKey(index)
	
	mByteA, err := json.Marshal(&val)
	if err != nil {
		return -1, err
	}
	
	err = self.stub.PutState(mKey, mByteA)
	if err != nil {
		return -1, err
	}
	
	return index, nil
}