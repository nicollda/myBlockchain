package main


import (
//	"errors"
	"fmt"
//	"strings"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


//********************************
//         Data Access Objects
//********************************




//********************************************************************************
//***** LinkedList
//********************************************************************************
//wraps the payload (struct we want to save) in a new struct that acts as a linked list
type LinkedListNode struct {
	NextNode string
	Payload []byte
}

type ChainLinkedList struct {
	stub *shim.ChaincodeStub
	mapName string
	originKey string
	nextNode string
}

func (self *ChainLinkedList) init(stub *shim.ChaincodeStub, mapName string) error {
	self.stub = stub
	self.mapName = mapName
	
	return nil
}


func (self *ChainLinkedList) getKey(key string) string {
	return self.mapName + separator + key
}



func (self *ChainLinkedList) getFirst(returnVal interface{}) error {
	self.nextNode = self.originKey
	
	return self.getNext(&returnVal)
}



func (self *ChainLinkedList) getNext(returnVal interface{}) error {
	lByteA, err := self.stub.GetState(self.nextNode)
	if err != nil {
		return err
	}
	
	var llNode LinkedListNode
	err = json.Unmarshal(lByteA, &llNode)
	if err != nil {
		return err
	}
	
	self.nextNode = llNode.NextNode
	
	return json.Unmarshal(llNode.Payload, &returnVal)
}


func (self *ChainLinkedList) get(key string, returnVal interface{}) error {
	lByteA, err := self.stub.GetState(self.getKey(key))
	if err != nil {
		return err
	}
	
	if lByteA == nil {
		return nil
	}
	
	var llNode LinkedListNode
	err = json.Unmarshal(lByteA, &llNode)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(llNode.Payload, &returnVal)
}



func (self *ChainLinkedList) put(key string, val interface{}) (string, error) {
	var newNode LinkedListNode
	
	
	mKey := self.getKey(key)
	
	mByteA, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	newNode.Payload = mByteA
		
	//Does the data already have a node?
	existingNodeByteA, err := self.stub.GetState(mKey)
	if err != nil {
		return "", err
	}
	
	if existingNodeByteA == nil {
		//No, the data does is not already persisted
		if self.originKey != "" {					//is our list empty? 
			newNode.NextNode = self.originKey
		}
		
		self.originKey = mKey
	} else {
		var existingNode LinkedListNode
		
		err := json.Unmarshal(existingNodeByteA, &existingNode)
		if err != nil {
			return "", err
		}
		
		newNode.NextNode = existingNode.NextNode
	}
	
	
	newNodeByteA, err := json.Marshal(newNode) 
	if err != nil {
		return "", err
	}
	
	err = self.stub.PutState(mKey, newNodeByteA)
	if err != nil {
		return "", err
	}
	
	return mKey, nil
}

func (self *ChainLinkedList) del(key string) error {
	mKey := self.getKey(key)
	
	//Does the data already have a node?
	existingNodeByteA, err := self.stub.GetState(mKey)
	if err != nil {
		return err
	}
	
	//Need to get the pointer from our deleted node so we can set its next door neighbour's pointer
	var existingNode LinkedListNode
	err = json.Unmarshal(existingNodeByteA, &existingNode)  //we can now retrieve the pointer we are about to delete
	if err != nil {
		return err
	}
	
	
	//Need to walk down the list until we find a node that points to the node that needs to be deleted
	for nextPointer := self.originKey; nextPointer != "";  {
		var llNode LinkedListNode
		llNodeByteA, err := self.stub.GetState(nextPointer)
		if err != nil {
			return err
		}
		
		err = json.Unmarshal(llNodeByteA, &llNode) 
		if err != nil {
			return err
		}
		
		if llNode.NextNode == mKey {
			//we found the node pointing to our deletion candidate.
			//edit its pointer and save it back
			llNode.NextNode = existingNode.NextNode
			
			llNodeByteA, err = json.Marshal(llNode)
			if err != nil {
				return err
			}
			
			//Save back our change
			err = self.stub.PutState(nextPointer, llNodeByteA)
			if err != nil {
				return err
			}
			nextPointer = ""  //ditch out of our for loop
		} else {
			nextPointer = llNode.NextNode
		}
	}
	
	//we have updated the links now we can just delete the node
	err = self.stub.DelState(self.getKey(key))
	if err != nil {
		//at this point if this fails it won't be in the linked list but it will still exist.  not good
		return err
	}
	
	return nil
}


//********************************************************************************
//***** ChainArray (cannot call it array)
//********************************************************************************

type ChainArray struct {
	stub *shim.ChaincodeStub
	ArrayName string
}

func (self *ChainArray) init(stub *shim.ChaincodeStub, arrayName string) error {
	self.stub = stub
	self.ArrayName = arrayName
	
	return nil
}


func (self *ChainArray) getKey(key int) string {
	return self.ArrayName + separator + strconv.Itoa(key)
}


func (self *ChainArray) get(index int, returnVal interface{}) error {
	vByteA, err := self.stub.GetState(self.getKey(index))
	if err != nil {
		return err
	}
	
	return json.Unmarshal(vByteA, &returnVal)
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
		return -1, err
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