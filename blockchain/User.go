package blockchain

import (
	"errors"
	"fmt"
)

var INVOKE_FUNCTIONS = []string{"init", "ipo", "bid", "ask", "exchange"}

var QUERY_FUNCTIONS = []string{"balance", "holdings", "market", "trades"}

type User struct {
	UserID   string `json:"userID"`
	Status   string `json:"status"`
	Balance  int    `json:"cash"`
	Holdings []*Holding
	Trades   []*Trade
}

func (self *User) GetBalance() int {
	return self.Balance
}

func (self *User) GetHoldings() []Holding {
	var holdingList []Holding
	for _, v := range self.Holdings {
		append(holdingList, v)
	}
	return holdingList
}

// not implement
func (self *User) GetMarket() int {

}

func (self *User) GetTrades() []Trade {
	var tradeList []Trade
	for _, v := range self.Trades {
		append(tradeList, v)
	}
	return tradeList
}

func (self *User) Init(args []string) ([]byte, error) {

}

// not implement
func (self *User) Ipo(args []string) ([]byte, error) {

}

func (self *User) Bid(args []string) ([]byte, error) {

}

func (self *User) Ask(args []string) ([]byte, error) {

}

func (self *User) Exchange(args []string) ([]byte, error) {

}

func (self *User) Invoke(function string, args []string) ([]byte, error) {
	fmt.Printf("Invoke function is %s\n", function)
	if !arrayContains(function, INVOKE_FUNCTIONS) {
		return nil, errors.New("Received unknown function invocation - " + function)
	}
	switch function {
	case "init":
		return self.Init(args)
	case "ipo":
		return self.Ipo(args)
	case "bid":
		return self.Bid(args)
	case "ask":
		return self.Ask(args)
	case "exchange":
		return self.Exchange(args)
	}
}

func (self *User) Query(function string, args []string) ([]byte, error) {
	fmt.Printf("Query function is %s\n", function)
	if !arrayContains(function, QUERY_FUNCTIONS) {
		return nil, errors.New("Invalid query function name. Expecting blance, holdings, market or trades")
	}
	switch function {
	case "balance":
		return self.GetBalance(args)
	case "holdings":
		return self.GetHoldings(args)
	case "market":
		return self.GetMarket(args)
	case "trades":
		return self.GetTrades(args)
	}
}
