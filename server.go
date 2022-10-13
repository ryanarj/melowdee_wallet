package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	url = "https://goerli.infura.io/v3/a197fe7fa0684b3b8ad84bf01fd2da89"
)

func genWalletHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if r.URL.Path != "/genWallet" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	password := r.FormValue("password")

	path := generate(password)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	key, err := keystore.DecryptKey(b, password)

	if err != nil {
		log.Fatal(err)
	}

	pData := crypto.FromECDSA(key.PrivateKey)
	fmt.Println("Private ", hexutil.Encode(pData))

	pubData := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	fmt.Println("Public ", hexutil.Encode(pubData))

	fmt.Println("Add ", crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex())

	data := make(map[string]string)
	data["private"] = hexutil.Encode(pData)
	data["public"] = hexutil.Encode(pubData)
	data["address"] = crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex()

	if err != nil {
		// handle error
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func generate(password string) string {
	key := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	passord := password
	account, err := key.NewAccount(passord)
	if err != nil {
		log.Fatal(err)
	}
	return account.URL.Path
}

func checkBalance(address string) *big.Int {

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	a1 := common.HexToAddress(address)

	b1, err := client.BalanceAt(context.Background(), a1, nil)
	if err != nil {
		log.Fatal(err)
	}

	return b1
}

func main() {

	http.HandleFunc("/genWallet", genWalletHandler) // Update this line of code

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
