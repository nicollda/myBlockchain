package blockchain

var TRADE_FUNCTIONS = []string{"bid", "ask", "exchange"}

type Trade struct {
	Entity    string  `json:"entity"`
	Char      string  `json:"char"`
	Event     string  `json:"event"`
	Action    string  `json:"action"`
	Price     float64 `json:"price"`
	Units     int     `json:"units"`
	Status    string  `json:"status"`
	Expiry    string  `json:"expiry"`
	Fulfilled int     `json:"fulfilled"`
}

func (self *User) RegisterTrade(function string, args []string) ([]byte, error) {
	if !arrayContains(function, TRADE_FUNCTIONS) {
		return nil, errors.New("Received unknown trade function - " + function)
	}
	//TODO
}

func (t *Trade) GetTrade(tradeId string) Trade {

}
