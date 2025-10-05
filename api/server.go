package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	// Wallet and gateway connection placeholders.
	// Before running, ensure you have:
	// - connection profile YAML at ./gateway/connection-org1.yaml
	// - wallet with identity "appUser" at ./wallet
	ccpPath := "./gateway/connection-org1.yaml"

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		log.Println("Wallet identity appUser does not exist. Please populate wallet with user identity before running.")
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(gateway.ConfigFromPath(ccpPath)),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("asset")

	r := mux.NewRouter()

	r.HandleFunc("/asset/create", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			DealerID    string `json:"dealerId"`
			MSISDN      string `json:"msisdn"`
			MPIN        string `json:"mpin"`
			Balance     string `json:"balance"`
			Status      string `json:"status"`
			TransAmount string `json:"transAmount"`
			TransType   string `json:"transType"`
			Remarks     string `json:"remarks"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		_, err := contract.SubmitTransaction("CreateAsset", req.DealerID, req.MSISDN, req.MPIN, req.Balance, req.Status, req.TransAmount, req.TransType, req.Remarks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Error submitting transaction: %v", err)))
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Asset created"))
	}).Methods("POST")

	r.HandleFunc("/asset/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		evaluateResult, err := contract.EvaluateTransaction("ReadAsset", id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf("Asset not found: %v", err)))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(evaluateResult)
	}).Methods("GET")

	r.HandleFunc("/asset/{id}/history", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		histBytes, err := contract.EvaluateTransaction("GetHistory", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Error getting history: %v", err)))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(histBytes)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("REST API listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
