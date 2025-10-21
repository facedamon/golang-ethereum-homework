package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	if err != nil {
		log.Fatalln(err)
	}

	blockNumber := big.NewInt(9452401)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalln(err)
	}
	//0xdea2020841cad6438d2e102f202d3e117ed843595202c93957078ed829ffcdab
	fmt.Println("block.hash.hex=", block.Hash().Hex())
	//1760965656
	fmt.Println("block.time=", block.Time())
	//96
	fmt.Println("block.Transactions=", block.Transactions().Len())

	//加载私钥
	privateKey, err := crypto.HexToECDSA("1df99d2849e087f2c3b55f265a149df3c0789c95aad678dec7d287e2fe979213")
	if err != nil {
		log.Fatalln(err)
	}
	publicKey := privateKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKey)
	//0x3827772Ac3ef7Bb875576a45B02F2a0fBdDed97b
	fmt.Println("publicAddress:", fromAddress)
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("balance wei:", balance)
	//对方账户
	toAddress := common.HexToAddress("0x21D2c46662Bde5850109b8d60c863BC332235BAa")
	//准备交易的参数：nonce、gasPrice、gasLimit、value
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalln(err)
	}
	value := big.NewInt(4000000000000000)
	var data []byte
	data = append(data, common.LeftPadBytes(value.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(toAddress.Bytes(), 32)...)
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	//开始生成事务准备交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	//发送人对事务签名
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		log.Fatalln(err)
	}
	//开始广播事务
	fmt.Println("开始广播事务")
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalln(err)
	}
	//0x0240f4a4989d2f74ceacf5122651bb68b48c70d862453c57afaa1af0028e1c90
	fmt.Println("tx hash:", signedTx.Hash().Hex())
}
