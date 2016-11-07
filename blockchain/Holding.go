package blockchain

type Holding struct {
	HoldingId string
	UserId    string
	Unit      int
}

func (self *Holding) RegisterHolding(holding Holding) bool {

}

func (self *Holding) GetHolding(holdingId string) Holding {

}

func (self *Holding) UpdateHolding(holdingId string, holding Holding) bool {

}

func (self *Holding) DeleteHolding(holdingId string) bool {

}

func (self *Holding) GetHoldingList() []Holding {

}
