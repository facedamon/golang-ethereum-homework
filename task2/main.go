package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/facedamon/golang-ethereum-homework/task2/store"
)

// ContractAddr 合约地址
const ContractAddr = "0xe92c1e7a0e5bd8e8cdb82fcabcaa3d638d189d02"

func main() {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatalln(err)
	}
	//加载私钥
	privateKey, err := crypto.HexToECDSA("1df99d2849e087f2c3b55f265a149df3c0789c95aad678dec7d287e2fe979213")
	if err != nil {
		log.Fatalln(err)
	}
	publicKey := privateKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKey)
	//0x3827772Ac3ef7Bb875576a45B02F2a0fBdDed97b
	fmt.Println("publicAddress:", fromAddress)

	//加载合约
	storeContract, err := store.NewStore(common.HexToAddress(ContractAddr), client)
	if err != nil {
		log.Fatalln(err)
	}
	//准备调用合约前的数据
	key := common.LeftPadBytes([]byte("test key"), 32)
	value := common.LeftPadBytes([]byte("test value"), 32)
	//签名
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatalln(err)
	}
	//调用合约方法
	tx, err := storeContract.SetItem(opt, [32]byte(key), [32]byte(value))
	if err != nil {
		log.Fatalln(err)
	}
	//0xb9797d39c3c1bb82de3e62ab2c20dd8ec373a9f13d901b534878bec3e6d7dd0f
	fmt.Println("tx hash=", tx.Hash().Hex())

	fmt.Println("等待交易确认...")
	time.Sleep(20 * time.Second)
	fmt.Println("等待结束...")

	//查询结果
	callOpt := &bind.CallOpts{Context: context.Background()}
	valueInContract, err := storeContract.Items(callOpt, [32]byte(key))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("is value saving in contract equals to origin value:", valueInContract == [32]byte(value))
	fmt.Println(string(valueInContract[:]))

}
