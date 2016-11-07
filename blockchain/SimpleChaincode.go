package blockchain

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//IBM HyperLedger is an implementation of BlockChain
type SimpleChaincode struct {
	stub *shim.ChaincodeStub
}

//implement interface func Get
func (self *SimpleChaincode) Get(key string) ([]byte, error) {
	return self.stub.GetState(key)
}

//implement interface func Push
func (self *SimpleChaincode) Push(key string, value interface{}) ([]byte, error) {
	switch reflect.TypeOf(value).String() {
	case "string":
		return self.stub.PutState(key, []byte(fmt.Sprint(value)))
	case "[]uint8":
		return self.stub.PutState(key, value)
	default: //use json to parse the object
		val, error = json.Marshal(value)
		if nil != error {
			return nil, error
		}
		return self.stub.PutState(key, val)
	}
}
