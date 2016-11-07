// BlockChain project IBlockChain.go
package blockchain

//This is an interface of BlockChain
type IBlockChain interface {
	Get(key string) ([]byte, error)
	Push(key string, value interface{}) ([]byte, error)
	//Init(args []string) ([]byte, error)
	//Ipo(args []string) ([]byte, error)
	//RegisterUser(args []string) ([]byte, error)
	//RegisterHappening() ([]byte, error)
	//RegisterHolding() ([]byte, error)
	//RegisterSecurity() ([]byte, error)
	//RegisterTrade() []byte, error)
}
