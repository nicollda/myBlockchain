package blockchain

type Happening struct {
	Char  string `json:"char"`
	Event string `json:"event"`
	User  string `json:"user"`
}

func (self *Happening) RegisterHappening(args []string) ([]byte, error) {

}
