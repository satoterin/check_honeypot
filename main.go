package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to an Ethereum node
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/61c46185479343c48a3f62e809d7028c")
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress("0xaec8fd4BE5d770a5f0d93bA48cA4D4AdBd4Cb9F4")  //Sample Honeypot contract address
	routerAddress := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D") //Uniswap v2 router address

	// If liquidity is present (Uniswap v2/v3 pair)
	liquidity, err := checkLiquidity(client, routerAddress)
	if err != nil {
		log.Fatal(err)
	}

	if liquidity {
		fmt.Println("Liquidity is Present!")
	} else {
		fmt.Println("Liduidity is Not Present!")
	}

	// Validate buy and sell tax
	buyTax, sellTax, err := validateTax(client, tokenAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Buy Tax: %v\nSell Tax: %v\n", buyTax, sellTax)

	// Check if token contract is a honeypot -  if liquidity is present and both buy and sell taxes are 0, not a honeypot
	if liquidity && buyTax.Cmp(big.NewInt(0)) == 0 && sellTax.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("Valid!")
	} else {
		fmt.Println("Honeypot!")
	}
}

func checkLiquidity(client *ethclient.Client, routerAddress common.Address) (bool, error) {
	// Check the balance of the router contract
	balance, err := client.BalanceAt(context.Background(), routerAddress, nil)
	if err != nil {
		return false, err
	}

	// If balance is over than 0, liquidity is present
	return balance.Cmp(big.NewInt(0)) > 0, nil
}

func validateTax(client *ethclient.Client, tokenAddress common.Address) (buyTax, sellTax *big.Int, err error) {

	// Get buy tax
	buyTax, err = getTokenBuyTax(client, tokenAddress)
	if err != nil {
		return nil, nil, err
	}

	// Get sell tax
	sellTax, err = getTokenSellTax(client, tokenAddress)
	if err != nil {
		return nil, nil, err
	}

	return buyTax, sellTax, nil
}

func getTokenBuyTax(client *ethclient.Client, tokenAddress common.Address) (*big.Int, error) {
	return big.NewInt(10), nil // Example buy tax value
}

func getTokenSellTax(client *ethclient.Client, tokenAddress common.Address) (*big.Int, error) {
	return big.NewInt(5), nil // Example sell tax value
}
