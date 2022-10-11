package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

func main() {

	http.HandleFunc("/genWallet", genWalletHandler) // Update this line of code

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}