package blockchain

type Security struct {
	SecurityId   int
	SecurityName string
	HoldingId    string
}

func (self *Security) RegisterSecurity(security Security) bool {

}

func (self *Security) GetSecurity(securityId int) Security {

}

func (self *Security) UpdateSecurity(securityId int, security Security) bool {

}

func (self *Security) DeleteSecurity(securityId int) bool {

}
