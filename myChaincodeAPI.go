/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/



//todo:  add constants for all string litterals
//todo:  need to make consitent status.  need better way to take them out of the process when closed
//todo: data abstraction layer, abstract persistance
//todo: add security to get user names




package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)




// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
	userRep UserRepository 
	holdingsRep HoldingsRepository
	securitiesRep SecurityRepository
	tradeRep TradeRepository
}


//********************************************************************************************************
//****      Blockchain API functions                                                                  ****
//********************************************************************************************************




//Init the blockchain.  populate a 2x2 grid of potential events for users to buy
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")
	
	t.writeOut(stub, "in init")

	t.userRep.init(stub)
	t.holdingsRep.init(stub)
	t.securitiesRep.init(stub)
	t.tradeRep.init(stub)
	
	
	
	//create a bank with some money
	var user User
	user.UserID = "BANK"
	user.Status = "Active"
	user.Ballance = 1000
	
	t.userRep.newUser(user.UserID, user.Ballance)
	
	user.Ballance = 2000
	
	t.userRep.updateUser(user.UserID, user)
	
	user2, err := t.userRep.getUser(user.UserID)
	
	t.writeOut(stub, "New Ballance: " + strconv.Itoa(user2.Ballance))
	
	u, err := json.Marshal(user)
	if err != nil {
		
		return nil, err
	}
	
	
	err = stub.PutState(bankUser, u)
	if err != nil {
		return nil, err
	}
	
	// initially offer some happenings
	
	a := []string{"Jaime,Killed",strconv.Itoa(defaultPrice), "100", "", "BANK"}
	
	_, err = t.registerTrade(stub, "IPO", a)
	if err != nil {
		return nil, err
	}
	
	a[0] = "Jaime,Killer"
	_, err = t.registerTrade(stub, "IPO", a)
	if err != nil {
		return nil, err
	}
	
	a[0] = "Jon,Killed"
	_, err = t.registerTrade(stub, "IPO", a)
	if err != nil {
		return nil, err
	}
	
	a[0] = "Jon,Killer"
	_, err = t.registerTrade(stub, "IPO", a)
	if err != nil {
		return nil, err
	}
	
	//Register some users.  this would normally happen via the UI but we will do it here to simplify
	_, err = t.registerUser(stub, "David")
	if err != nil {
		return nil, err
	}
	
	_, err = t.registerUser(stub, "Wesley")
	if err != nil {
		return nil, err
	}
	
	t.userRep.newUser("Aaron", 1000)
	_, err = t.registerUser(stub, "Aaron")
	if err != nil {
		return nil, err
	}
	
	
	
	//register some trades
	b := []string{"Jaime,Killed", strconv.Itoa(defaultPrice), "100", "", "Aaron"}
	_, err = t.registerTrade(stub, "Bid", b)
	if err != nil {
		return nil, err
	}
	
	
	t.writeOut(stub, "in init: before exchange")
	_, err = t.exchange(stub)
	if err != nil {
		t.writeOut(stub, "in init: after exhcnage in err != nil")
		return nil, err
	}
	
	
	c := []string{"Jaime,Killed"}
	
	_, err = t.dividend(stub, c)
	
	if err != nil {
		t.writeOut(stub, "in init: after registerHeppeing in err != nil")
		//return nil, err
	}
	
	
	return nil, nil
}





// Invoke callback representing the invocation of a chaincode
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke called, determining function")
	
	// Handle different functions
	if function == "init" {
		fmt.Printf("Function is init")
		return t.Init(stub, function, args)
	} else if function == "ipo" {
		// offers squares up for sale as initial public offering
		fmt.Printf("Function is ipo")
		args[4] = ""
		return t.registerTrade(stub, function, args)
	} else if function == "bid" {
		// enter a bid to buy a square
		fmt.Printf("Function is bid")
		return t.registerTrade(stub, function, args)
	} else if function == "ask" {
		// enter an ask to sell a square
		fmt.Printf("Function is ask")
		return t.registerTrade(stub,  function, args)
	} else if function == "dividend" {
		// enter an an character event happening in the show.  pays out to users holding squares
		fmt.Printf("Function is ask")
		return t.dividend(stub, args)
	} else if function == "exchange" {
		// matches trades and excecutes any matches
		fmt.Printf("Function is exchange")
		return t.exchange(stub)
	} else if function == "registeruser" {
		// matches trades and excecutes any matches
		fmt.Printf("Function is registeruser")
		userID := args[0]
		return t.registerUser(stub, userID)
	}
	
	return nil, errors.New("Received unknown function invocation")
}




// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Printf("Query called, determining function")
	
	// Handle different functions
	if function == "holdings" {
		// query a users holdings
		return t.holdings(stub, args)
	} else if function == "cash" {
		// query a users cash on hand
		return t.ballance(stub, args)
	} else if function == "users" {
		// query for list of users
		return t.users(stub)
	} else if function == "securities" {
		// query for list of securities
		return t.securities(stub)
	} else {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting holdings, cash, users or securities")
	}	
	
	return nil, nil
}





func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}