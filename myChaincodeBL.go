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
//	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const separator = 		"::::"
const userIndex =		"UserIndex" + separator
const tradeIndex =		"TradeIndex" + separator
const securityIndex = 	"SecurityIndex" + separator
const holdingIndex =	"HoldingIndex" + separator
const initialCash =		1000
const payout =			5
const defaultPrice =	5
const debug =			true





//********************************************************************************************************
//****                        Query function inplimentations                                          ****
//********************************************************************************************************


func (t *SimpleChaincode) securities() ([]byte, error) {
	s := []string {"Jaime,Killed", "Jaime,Killer", "Jon,Killed", "Jon,Killer"}
	
	sByteA, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	
	return sByteA, nil
}


func (t *SimpleChaincode) users() ([]byte, error) {
	u := []string {"David", "Aaron", "Wesley"}
	
	uByteA, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	
	return uByteA, nil
}

func (t *SimpleChaincode) holdings(userID string) ([]byte, error) {
	fmt.Printf("Running holdings")
	
	user, err := t.userRep.getUser(userID)
	if err != nil {
		return nil, err
	}
	
	return []byte(strconv.Itoa(user.getBallance())), nil
	
}

func (t *SimpleChaincode) ballance(userID string) ([]byte, error) {
	fmt.Printf("Running cash")
	
	
	curOutByteA, err := t.stub.GetState("currentOutput") //tradeIndex + "5") //bankUser)//"LastTradeIndex")//bankUser)  //userIndex + "BANK")
	if err != nil {
		return nil, err
	}
	
	var shareKey Holding
	shareKey.UserID = "Aaron"
	shareKey.SecurityID = "Jaime,Killed"

	
	shareKeyByteA, err := json.Marshal(shareKey)
	if err != nil {
		return nil, err
	}
	
	tradeOutByteA, err := t.stub.GetState(string(shareKeyByteA)) //tradeIndex+ "5") //bankUser)//"LastTradeIndex")//bankUser)  //userIndex + "BANK")
	if err != nil {
		return nil, err
	}
	
	aaronByteA, err := t.stub.GetState(userIndex + "3") 
	
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




//********************************************************************************************************
//****                        Invoke function inplimentations                                          ****
//********************************************************************************************************





// initial public offering for a square
func (t *SimpleChaincode) registerTrade(tradeType string, userID string, securityID string, price float64, units int, expiry string) ([]byte, error) {
	fmt.Printf("Running registerTrade")
	
	var trade Trade
	var err error
	
	securityPointer, err := t.securitiesRep.getSecurityPosition(securityID)
	if err != nil {
		return nil, err
	}
	
	if securityPointer == 0 {
		return nil, errors.New("Security Not Found.")
	}
	
	trade.init(userID, securityID, securityPointer, tradeType, price, units, expiry)
	
	_, err = t.tradeRep.newTrade(trade)
	if err != nil {
		return nil, err
	}
	
	_, err = t.exchange()
	if err != nil {
		return nil, err
	}
	
	return nil, nil  //what should this function return?
}



// initial public offering for a square
func (t *SimpleChaincode) registerSecurity(securityID string, desc string) ([]byte, error) {
	fmt.Printf("Running registerSecurity")
	
	var security Security
	
	security.init(securityID, desc)
	_, err := t.securitiesRep.newSecurity(security)
	if err != nil {
		return nil, err
	}
	
	return nil, nil  //what should this function return?
}



// called by the moderator watson?  to specify that an event happened pay it out
func (t *SimpleChaincode) dividend(securityID string) ([]byte, error) {
	fmt.Printf("Running dividend")
	
	t.writeOut("in dividend")
	
	var shareKey Holding
	var numberUsers int
	var currentUser User
	
	shareKey.SecurityID = securityID
	
	
	//todo: need to make data abstraction
	numberUsersByteA, err := t.stub.GetState("Last" + userIndex)
	if err != nil {
		return nil, err
	}
	
	numberUsers, err = strconv.Atoi(string(numberUsersByteA))
	if err != nil {
		return nil, err
	}
	
	t.writeOut("in dividend: before for loop")
	//For each user
	for i := 1; i <= numberUsers; i++ {
		currentUserByteA, err := t.stub.GetState(userIndex + strconv.Itoa(i))
		err = json.Unmarshal(currentUserByteA, &currentUser)
		if err != nil {
			return nil, err
		}
		
		//create a string to look up the number of shares using , userid, char and event.  the result is number of shares
		shareKey.UserID = currentUser.UserID
		
		shareKeyByteA, err := json.Marshal(shareKey)
		if err != nil {
			return nil, err
		}
		
		numberSharesByteA, err := t.stub.GetState(string(shareKeyByteA))
		numberShares, err := strconv.Atoi(string(numberSharesByteA))
		
		
		if err == nil {  //means the user has stock in this security
			if currentUser.Status == "Active" && numberShares > 0 {
				
				currentUser.Ballance = currentUser.Ballance + payout		//todo:  should be transfer of funds not "creating money".  
				
				currentUserByteA,err := json.Marshal(currentUser)
				if err != nil {
					return nil, err
				}
				
				t.stub.PutState(userIndex + strconv.Itoa(i), currentUserByteA)  //should be via data layer
			}
		}
	}
	
	t.writeOut("in dividend: before return")
	return nil,nil
}


func (t *SimpleChaincode) writeOut(out string) ([]byte, error) {
	if debug {
		curOutByteA,err := t.stub.GetState("currentOutput")
		outByteA := []byte(string(curOutByteA) + ":::" + out)
		err = t.stub.PutState("currentOutput", outByteA)
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
func (t *SimpleChaincode) exchange() ([]byte, error) {
	fmt.Printf("Running exchange")
	
	t.writeOut("in exchange")
	
	var buyTrade	Trade
	var sellTrade	Trade
	
	numberTradesByteA, err := t.stub.GetState("Last" + tradeIndex)  //should be through data layer
	if err != nil {
		return nil, err
	}
	
	numberTrades, err := strconv.Atoi(string(numberTradesByteA))
	if err != nil {
		return nil, err
	}
	
	
	t.writeOut("in exchange: before matching loop")
	//trade matching loop
	for b := 1; b <= numberTrades; b++{
		bTradeByteA, err := t.stub.GetState(tradeIndex + strconv.Itoa(b))
		if err != nil {
			return nil, err
		}
		
		err = json.Unmarshal(bTradeByteA, &buyTrade)  //should be via data layer
		if err != nil {
			return nil, err
		}
		
		for s := 1; s <= numberTrades; s++ {
			sTradeByteA, err := t.stub.GetState(tradeIndex + strconv.Itoa(s))
			if err != nil {
				return nil, err
			}
			
			err = json.Unmarshal(sTradeByteA, &sellTrade)  //should be via data layer
			if err != nil {
				return nil, err
			}
			
			
			//t.writeOut(sellTrade.Status + " " + ")
			if sellTrade.Status == "Open" && buyTrade.Status == "Open" && sellTrade.TransType == "Ask" && buyTrade.TransType == "Bid" && sellTrade.SecurityID == buyTrade.SecurityID {
				t.writeOut("in exchange: before executeTrade")
				_, err := t.executeTrade(b, buyTrade, s, sellTrade)
				
				if err != nil {
					return nil, err
				}
			}
		}	
	}
	
	
	t.writeOut("in exchange: before return")
	return nil, nil
}







//actually make the trade.  does not vlidate anything.  this should be added at some point:
	//no roll back / transation management
	//doesnt check holdings
	//or users are valid
	//or expiry date etc...
func (t *SimpleChaincode) executeTrade(buyTradeIndex int, buyTrade Trade, sellTradeIndex int, sellTrade Trade) ([]byte, error) {
	fmt.Printf("Running exchange")
	
	var buyUser 	User
	var buyerIndex	int
	var sellUser 	User
	var sellerIndex	int
	var numberUsers	int
	var tempUser	User
	
	
	
	numberUsersByteA, err := t.stub.GetState("Last" + userIndex)  //should be through data layer
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
		userByteA, err := t.stub.GetState(userIndex + strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		
		err = json.Unmarshal(userByteA, &tempUser)
		if err != nil {
			return nil, err
		}
		
		if buyTrade.UserID == tempUser.UserID {
			buyUser = tempUser
			buyerIndex = i
		}
		
		if sellTrade.UserID == tempUser.UserID {
			sellUser = tempUser
			sellerIndex = i
		}
	}
	
	if sellTrade.UserID == "BANK" {
		sellerUserByteA, err := t.stub.GetState(userIndex + "BANK")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(sellerUserByteA, &sellUser)
		if err != nil {
			return nil, err
		}
	}
	
	
	
	//transfers funds and closes the trades
	buyUser.Ballance = buyUser.Ballance - defaultPrice
	sellUser.Ballance = sellUser.Ballance + defaultPrice
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
	err = t.stub.PutState(userIndex + strconv.Itoa(buyerIndex), buyUserByteA)
	if err != nil {
		return nil, err
	}
	
	
	//Saves the changes to the seller
	err = t.stub.PutState(userIndex + strconv.Itoa(sellerIndex), sellUserByteA)
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
	err = t.stub.PutState(tradeIndex + strconv.Itoa(buyTradeIndex), buyTradeByteA)
	if err != nil {
		return nil, err
	}
	
	
	//saves the changes to the sell trade
	err = t.stub.PutState(tradeIndex + strconv.Itoa(sellTradeIndex), sellTradeByteA)
	if err != nil {
		return nil, err
	}
	
	var shareKey Holding
	shareKey.UserID = buyTrade.UserID
	shareKey.SecurityID = buyTrade.SecurityID
	
	//adjust Holdings for buyer and seller
	shareKeyByteA, err := json.Marshal(shareKey)
	if err != nil {
		return nil, err
	}
	
	unitsString:= strconv.Itoa(buyTrade.Units)
	err = t.stub.PutState(string(shareKeyByteA), []byte(unitsString))
	if err != nil {
		return nil, err
	}
	
	if sellTrade.UserID != "BANK" {
		shareKey.UserID = sellTrade.UserID
		
		//remove holdings from seller.  should really only delstate if units = 0 but this is not inplimented
		err = t.stub.DelState(string(shareKeyByteA))
		if err != nil {
			return nil, err
		}
	}
	
	return nil, nil
}



// register user
func (t *SimpleChaincode) registerUser(userID string) ([]byte, error) {
	fmt.Printf("Running registerUser")
	//need to make sure the user is not already registered
	newCash := initialCash
	
	if userID == "BANK" {
		newCash = 100000
	} 
	
	index, err := t.userRep.newUser(userID, newCash, "Active")
	if err != nil {
		return nil, err
	}
	
	return []byte(index), nil
}


//   curently not used but should be used in place of taking the user id via the interface.  user id should come from the security model
func (t *SimpleChaincode) getUserID(args []string) ([]byte, error) {
	//returns the user's ID 
	
	return nil, nil  //dont know how to get the current user
}