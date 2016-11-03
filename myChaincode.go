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
//todo: data abstraction layer
//todo: put funtions etc.  in the right order
//todo: add error trapping
//todo: add security to get user names




package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)



const separator = 		"::::"
const userIndex =		"UserIndex" + separator
const tradeIndex =		"TradeIndex" + separator
const happeningIndex = 	"HappeningIndex" + separator
const initialCash =		1000
const payout =			5
const defaultPrice =	5
const bankUser =		userIndex + "BANK"
const debug = 			true



type Trade struct {
	Entity		string	`json:"entity"`
	Char		string	`json:"char"`
	Event		string	`json:"event"`
	Action		string	`json:"action"`
	Price		float64	`json:"price"`
	Units		int		`json:"units"`
	Status		string	`json:"status"`
	Expiry		string	`json:"expiry"`
	Fulfilled	int		`json:"fulfilled"`
}


type HappeningRegister struct {
	Char		string	`json:"char"`
	Event		string	`json:"event"`
	User		string	`json:"user"`
}

type User struct {
	UserID		string	`json:"userID"`
	Status		string	`json:"status"`
	Cash		int		`json:"cash"`
}



// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}





//Init the blockchain.  populate a 2x2 grid of potential events for users to buy
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Printf("Init called, initializing chaincode")
	
	t.writeOut(stub, "in init")
	
	//create a bank with some money
	var user User
	user.UserID = "BANK"
	user.Status = "Active"
	user.Cash	= 1000
	
	u, err := json.Marshal(user)
	if err != nil {
		
		return nil, err
	}
	
	
	err = stub.PutState(bankUser, u)
	if err != nil {
		return nil, err
	}
	
	// initially offer some happenings
	
	a := []string{"Jaime", "Killed",strconv.Itoa(defaultPrice), "100", "", "BANK"}
	
	a[0] = "Jaime"	//character
	a[1] = "Killed" //action
	a[2] = strconv.Itoa(defaultPrice)		//cost per share
	a[3] = "100"		//number of shares
	a[4] = ""
	a[5] = "BANK"
	
	_, err = t.registerTrade(stub, "IPO", a)
	if err != nil {
		return nil, err
	}
	
	a[1] = "Killer"
	_, err = t.registerTrade(stub, "IPO", a)
	if err != nil {
		return nil, err
	}
	
	a[0] = "Jon"
	_, err = t.registerTrade(stub, "IPO", a)
	if err != nil {
		return nil, err
	}
	
	a[1] = "Killed"
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
	
	_, err = t.registerUser(stub, "Aaron")
	if err != nil {
		return nil, err
	}
	
	//register some trades
	b := []string{"Jaime", "Killed", strconv.Itoa(defaultPrice), "100", "", "Aaron"}
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
	
	
	c := []string{"Jamie", "Killed"}
	c[0] = c[0]
	
	_, err = t.registerHappening(stub, c)
	
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
	} else if function == "registerHappening" {
		// enter an an character event happening in the show.  pays out to users holding squares
		fmt.Printf("Function is ask")
		return t.registerHappening(stub, args)
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
	//	return t.holdings(stub, args)
	
	} else if function == "cash" {
		// query a users cash on hand
		return t.cash(stub, args)
	
	} else if function == "market" {
		// query a users total market value based on cash + market value of holdings
	//	return t.market(stub, args)
		
	} else {
		fmt.Printf("Function is query")
		return nil, errors.New("Invalid query function name. Expecting holdings, cash or market")
	}	
	
	return nil, nil
}





func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}






func (t *SimpleChaincode) cash(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Printf("Running cash")
	
	
	curOutByteA, err := stub.GetState("currentOutput") //tradeIndex + "5") //bankUser)//"LastTradeIndex")//bankUser)  //userIndex + "BANK")
	if err != nil {
		return nil, err
	}
	
	var shareKey HappeningRegister
	shareKey.User = "Aaron"
	shareKey.Char = "Jaime"
	shareKey.Event = "Killed"
	
	shareKeyByteA, err := json.Marshal(shareKey)
	if err != nil {
		return nil, err
	}
	
	tradeOutByteA, err := stub.GetState(string(shareKeyByteA)) //tradeIndex+ "5") //bankUser)//"LastTradeIndex")//bankUser)  //userIndex + "BANK")
	if err != nil {
		return nil, err
	}
	
	aaronByteA, err := stub.GetState(userIndex +"3") 
	
	//json.Unmarshal(bank, &bankString)
/*	
	var user User
	
	err = json.Unmarshal(bank, &user)
	if err != nil {
		return nil, err
	}
	
*/	
	
	
	return []byte(string(curOutByteA) + "        " + string(tradeOutByteA) + "         " + string(aaronByteA)), nil//[]byte(strconv.Itoa(user.Cash)), nil
	
}







// initial public offering for a square
func (t *SimpleChaincode) registerTrade(stub *shim.ChaincodeStub, tradeType string, args []string) ([]byte, error) {
	fmt.Printf("Running registerTrade")
	
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting registerTrade(character, event, price, units, expiry, user)")
	}
	
	var trade Trade
	var err error
	
	if tradeType == "IPO" {
		trade.Entity = "BANK"  // who is the source user
		trade.Action = "Ask"
	} else { 
		trade.Entity = args[5]  //should get this from the security mechanism...  dont know how that works
		trade.Action = tradeType
	}
	
	
	trade.Char = args[0]
	trade.Event = args[1]
	trade.Price, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, err
	}
	trade.Units, err = strconv.Atoi(args[3])
	if err != nil {
		return nil, err
	}
	trade.Status = "Open"
	trade.Expiry = args[4]
	trade.Fulfilled = 0
	
	
	// Write the state back to the ledger
	temp, err := json.Marshal(trade)
	if err != nil {
		return nil, err
	}
	
	
	index, err := t.push(stub, tradeIndex, temp)
	if err != nil {
		return nil, err
	}
	
	return index, nil
}




// need to make a persistance class / data abstraction
func (t *SimpleChaincode) push(stub *shim.ChaincodeStub, structureName string, value []byte) ([]byte, error) {
	fmt.Printf("Running Push")
	
	index, err := t.getNextIndex(stub, structureName)
	if err != nil {
		return nil, err
	}
	
	// Write the state back to the ledger
	var key string
	
	key = structureName + string(index)
		
	err = stub.PutState(key, value)
	if err != nil {
		return nil, err
	}
	
	return index, nil
}	




// user offers a square for sale asking for x for y units
func (t *SimpleChaincode) getNextIndex(stub *shim.ChaincodeStub, structureName string) ([]byte, error) {
	fmt.Printf("Running getNextIndex")
	
	var lastID int
	
	lastIDByteA, err := stub.GetState("Last" + structureName)
	lastID, err = strconv.Atoi(string(lastIDByteA))
	if err != nil {
		lastID = 1
	} else { 
		lastID = lastID + 1
	}
	
	lastIDByteA = []byte(strconv.Itoa(lastID))   
	err = stub.PutState("Last" + structureName, lastIDByteA)
	if err != nil {
		return nil, err
	}
	
	
	fmt.Printf(strconv.Itoa(lastID))
	
	return lastIDByteA, nil
}






// called by the moderator watson?  to specify that an event happened pay it out
func (t *SimpleChaincode) registerHappening(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Printf("Running registerHappening")
	
	t.writeOut(stub, "in registerHappening")
	
	var shareKey HappeningRegister
	var numberUsers int
	var currentUser User
	
	shareKey.Char = args[0]
	shareKey.Event = args[1]
	
	//todo: need to make data abstraction
	numberUsersByteA, err := stub.GetState("Last" + userIndex)
	if err != nil {
		return nil, err
	}
	
	numberUsers, err = strconv.Atoi(string(numberUsersByteA))
	if err != nil {
		return nil, err
	}
	
	t.writeOut(stub, "in registerHappening: before for loop")
	//For each user
	for i := 1; i <= numberUsers; i++ {
		currentUserByteA, err := stub.GetState(userIndex + strconv.Itoa(i))
		err = json.Unmarshal(currentUserByteA, &currentUser)
		if err != nil {
			return nil, err
		}
		
		//create a string to look up the number of shares using , userid, char and event.  the result is number of shares
		shareKey.User = currentUser.UserID
		
		shareKeyByteA, err := json.Marshal(shareKey)
		if err != nil {
			return nil, err
		}
		
		t.writeOut(stub, "in registerHappening: sharesKeyByteA=" + string(shareKeyByteA))
		numberSharesByteA, err := stub.GetState(string(shareKeyByteA))
		numberShares, err := strconv.Atoi(string(numberSharesByteA))
		
		t.writeOut(stub, "in registerHappening: before big if_" + string(numberSharesByteA))
		
		if err == nil {  //means the user has stock in this security
			t.writeOut(stub, "in registerHappening: after big if")
			
			if currentUser.Status == "Active" && numberShares > 0 {
				t.writeOut(stub, "in registerHappening: in middle if: " + currentUser.UserID)
				
				currentUser.Cash = currentUser.Cash + payout				//todo:  should be transfer of funds not "creating money".  
				
				currentUserByteA,err := json.Marshal(currentUser)
				if err != nil {
					return nil, err
				}
				
				stub.PutState(userIndex + strconv.Itoa(i), currentUserByteA)  //should be via data layer
			}
		}	
	}
	
	t.writeOut(stub, "in registerHappening: before return")
	return nil,nil
}


func (t *SimpleChaincode) writeOut(stub *shim.ChaincodeStub, out string) ([]byte, error) {
	
	if debug {
		curOutByteA,err := stub.GetState("currentOutput")		
		outByteA := []byte(string(curOutByteA) + ":::" + out)
		err = stub.PutState("currentOutput", outByteA)
		return nil, err
	}
	return nil, nil
}



// run on a schedule to execute any pending trades. matching asks with bids and updating the ledger
// first iteration will:
//		only match buyer and seller based on ticker and not on bid and ask prices.  this will simplify and elimiate items we are not trying to prove
//		assume one one share per trade.  this will elimate having to match the number of trades from buy with the sell
//		ignore expiry
//		ignore if the counterparties have the security
//		or if user is active
func (t *SimpleChaincode) exchange(stub *shim.ChaincodeStub) ([]byte, error) {
	fmt.Printf("Running exchange")
	
	t.writeOut(stub, "in exchange")
	
	var buyTrade	Trade
	var sellTrade	Trade
	
	numberTradesByteA, err := stub.GetState("Last" + tradeIndex)  //should be through data layer
	if err != nil {
		return nil, err
	}
	
	numberTrades, err := strconv.Atoi(string(numberTradesByteA))
	if err != nil {
		return nil, err
	}
	
	
	t.writeOut(stub, "in exchange: before matching loop")
	//trade matching loop
	for b := 1; b <= numberTrades; b++{
		bTradeByteA, err := stub.GetState(tradeIndex + strconv.Itoa(b))
		if err != nil {
			return nil, err
		}
		
		err = json.Unmarshal(bTradeByteA, &buyTrade)  //should be via data layer
		if err != nil {
			return nil, err
		}
		
		for s := 1; s <= numberTrades; s++ {
			sTradeByteA, err := stub.GetState(tradeIndex + strconv.Itoa(s))
			if err != nil {
				return nil, err
			}
			
			err = json.Unmarshal(sTradeByteA, &sellTrade)  //should be via data layer
			if err != nil {
				return nil, err
			}
			
			
			//t.writeOut(stub, sellTrade.Status + " " + ")
			if sellTrade.Status == "Open" && buyTrade.Status == "Open" && sellTrade.Action == "Ask" && buyTrade.Action == "Bid" && sellTrade.Char == buyTrade.Char && sellTrade.Event == buyTrade.Event {
				t.writeOut(stub, "in exchange: before executeTrade")
				_, err := t.executeTrade(stub, b, buyTrade, s, sellTrade)
				
				if err != nil {
					return nil, err
				}
			}
		}	
	}
	
	
	t.writeOut(stub, "in exchange: before return")
	return nil, nil
}







//actually make the trade.  does not vlidate anything.  this should be added at some point:
	//no roll back / transation management
	//doesnt check holdings
	//or users are valid
	//or expiry date etc...
func (t *SimpleChaincode) executeTrade(stub *shim.ChaincodeStub, buyTradeIndex int, buyTrade Trade, sellTradeIndex int, sellTrade Trade) ([]byte, error) {
	fmt.Printf("Running exchange")
	
	var buyUser 	User
	var buyerIndex	int
	var sellUser 	User
	var sellerIndex	int
	var numberUsers	int
	var tempUser	User
	
	
	
	numberUsersByteA, err := stub.GetState("Last" + userIndex)  //should be through data layer
	if err != nil {
		return nil, err
	}
	
	numberUsers, err = strconv.Atoi(string(numberUsersByteA))
	if err != nil {
		return nil, err
	}
	
	//finds the counterparties involved
	//this is no good.  should have another hash to get the index of the user in the array or maybe store it in the trades?
	for i := 1; i <= numberUsers; i++ {
		userByteA, err := stub.GetState(userIndex + strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		
		err = json.Unmarshal(userByteA, &tempUser)
		if err != nil {
			return nil, err
		}
		
		if buyTrade.Entity == tempUser.UserID {
			buyUser = tempUser
			buyerIndex = i
		}
		
		if sellTrade.Entity == tempUser.UserID {
			sellUser = tempUser
			sellerIndex = i
		}
	}
	
	if sellTrade.Entity == "BANK" {
		sellerUserByteA, err := stub.GetState(userIndex + "BANK")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(sellerUserByteA, &sellUser)
		if err != nil {
			return nil, err
		}
	}
	
	
	
	//transfers funds and closes the trades
	buyUser.Cash = buyUser.Cash - defaultPrice
	sellUser.Cash = sellUser.Cash + defaultPrice
	buyTrade.Status = "Closed"
	sellTrade.Status = "Closed"
	
	//no transaction rolling back etc...  dont know how best to handle
	buyUserByteA, err := json.Marshal(buyUser)
	if err != nil {
		return nil, err
	}
	
	sellUserByteA, err := json.Marshal(sellUser)
	if err != nil {
		return nil, err
	}
	
	//Saves the changes to the buyer
	err = stub.PutState(userIndex + strconv.Itoa(buyerIndex), buyUserByteA)
	if err != nil {
		return nil, err
	}
	
	
	//Saves the changes to the seller
	err = stub.PutState(userIndex + strconv.Itoa(sellerIndex), sellUserByteA)
	if err != nil {
		return nil, err
	}
	
	
	buyTradeByteA, err := json.Marshal(buyTrade)
	if err != nil {
		return nil, err
	}
	
	sellTradeByteA, err := json.Marshal(sellTrade)
	if err != nil {
		return nil, err
	}
	
	
	//Saves the changes to the buy trade
	err = stub.PutState(tradeIndex + strconv.Itoa(buyTradeIndex), buyTradeByteA)
	if err != nil {
		return nil, err
	}
	
	
	//saves the changes to the sell trade
	err = stub.PutState(tradeIndex + strconv.Itoa(sellTradeIndex), sellTradeByteA)
	if err != nil {
		return nil, err
	}
	
	var shareKey HappeningRegister
	shareKey.User = buyTrade.Entity
	shareKey.Char = buyTrade.Char
	shareKey.Event = buyTrade.Event
	
	//adjust Holdings for buyer and seller
	shareKeyByteA, err := json.Marshal(shareKey)
	if err != nil {
		return nil, err
	}
	
	unitsString:= strconv.Itoa(buyTrade.Units)
	err = stub.PutState(string(shareKeyByteA), []byte(unitsString))
	if err != nil {
		return nil, err
	}
	
	if sellTrade.Entity != "BANK" {
		shareKey.User = sellTrade.Entity
		
		//remove holdings from seller.  should really only delstate if units = 0 but this is not inplimented
		err = stub.DelState(string(shareKeyByteA))
		if err != nil {
			return nil, err
		}
	}
	
	return nil, nil
}



// register user
func (t *SimpleChaincode) registerUser(stub *shim.ChaincodeStub, userID string) ([]byte, error) {
	fmt.Printf("Running registerUser")
//need to make sure the user is not already registered
//need to make another hash to hold the users' id and return thier index
	
	var user User
	
	user.UserID = userID
	user.Status = "Active"
	user.Cash	= initialCash
	
	userByteA, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	
	index, err := t.push(stub, userIndex, userByteA)
	if err != nil {
		return nil, err
	}
	
	return index, nil
}


//   curently not used but should be used in place of taking the user id via the interface.  user id should come from the security model
func (t *SimpleChaincode) getUserID(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//returns the user's ID 
	
	return nil, nil  //dont know how to get the current user
}