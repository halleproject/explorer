package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/client"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/db"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/errors"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/models"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/schema"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/utils"

	"github.com/gorilla/mux"
)

// Account is a account handler
type Account struct {
	l      *log.Logger
	client *client.Client
	db     *db.Database
}

// NewAccount creates a new account handler with the given params
func NewAccount(l *log.Logger, client *client.Client, db *db.Database) *Account {
	return &Account{l, client, db}
}

// GetAccount returns account information
func (a *Account) GetAccount(rw http.ResponseWriter, r *http.Request) {
	a.l.Printf("================= GetAccount ===============\n")

	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		errors.ErrRequiredParam(rw, http.StatusBadRequest, "address is required")
		return
	}

	addrLen := len(address)
	if addrLen != 44 && addrLen != 45 { //44 halle,45 cosmos, bnb 42
		//if len(address) != 42 {
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "address is invalid")
		return
	}

	account, err := a.client.Account(address)
	if err != nil {
		a.l.Printf("failed to request account information: %s\n", err)
	}

	utils.Respond(rw, account)
	return
}

// GetAccountTxs returns transactions associated with an account
func (a *Account) GetAccountTxs(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("============= GetAccountTxs ===================")
	vars := mux.Vars(r)
	address := vars["address"]

	if address == "" {
		errors.ErrRequiredParam(rw, http.StatusBadRequest, "address is required")
		return
	}

	addrLen := len(address)
	if addrLen != 44 && addrLen != 45 { //44 halle,45 cosmos, bnb 42
		//if len(address) != 42 {
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "address is invalid")
		return
	}

	page := int(1)
	rows := int(10)

	if len(r.URL.Query()["page"]) > 0 {
		page, _ = strconv.Atoi(r.URL.Query()["page"][0])
	}

	if len(r.URL.Query()["rows"]) > 0 {
		rows, _ = strconv.Atoi(r.URL.Query()["rows"][0])
	}

	if rows < 1 {
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "'rows' cannot be less than")
		return
	}

	if rows > 50 {
		errors.ErrInvalidParam(rw, http.StatusBadRequest, "'rows' cannot be greater than 50")
		return
	}

	//acctTxs, err := a.client.AccountTxs(address, page, rows)

	//if err != nil {
	//	a.l.Printf("failed to get account txs: %s\n", err)
	//}
	//
	//txArray := make([]models.AccountTxArray, 0)
	//
	//for _, tx := range acctTxs {
	//	var toAddr string
	//	if tx.ToAddress != "" {
	//		toAddr = tx.ToAddress
	//	}
	//
	//	tempTxArray := &models.AccountTxArray{
	//		BlockHeight: tx.Height,
	//		TxHash:      tx.TxHash,
	//		Code:        int64(tx.Code),
	//		//TxType:        tx.Messages,
	//		//TxAsset:       tx.TxAsset,
	//		//TxQuoteAsset:  tx.TxQuoteAsset,
	//		//Value:         tx.Value,
	//		//TxFee:         tx.TxFee,
	//		//TxAge:         tx.TxAge,
	//		FromAddr: tx.FromAddress,
	//		ToAddr:   toAddr,
	//		//Log:           tx.Log,
	//		//ConfirmBlocks: tx.ConfirmBlocks,
	//		Memo: tx.Memo,
	//		//Source:        tx.Source,
	//		Timestamp: tx.Timestamp.Unix(),
	//	}
	//
	//	//txType TRANSFER shouldn't throw message data
	//	//var data models.AccountTxData
	//	//if tx.Memo != "" {
	//	//	err = json.Unmarshal([]byte(tx.Data), &data)
	//	//	if err != nil {
	//	//		a.l.Printf("failed to unmarshal AssetTxData: %s", err)
	//	//	}
	//	//
	//	//	tempTxArray.Message = &data
	//	//}
	//
	//	txArray = append(txArray, *tempTxArray)
	//}
	//
	//result := &models.ResultAccountTxs{
	//	TxNums:  len(acctTxs),
	//	TxArray: txArray,
	//}

	txs, err := a.db.QueryTxsByAddress(address, 0, (page-1)*rows, rows)
	if err != nil {
		a.l.Printf("failed to query txs: %s\n", err)
	}

	if len(txs) <= 0 {
		utils.Respond(rw, models.ResultTxs{})
		return
	}

	result, err := a.setTxs(txs)
	if err != nil {
		a.l.Printf("failed to set txs: %s\n", err)
	}

	totalTxsNum, err := a.db.CountTotalTxsNum()
	if err != nil {
		a.l.Printf("failed to query total number of txs: %s\n", err)
	}

	after := (page - 1) * rows
	// Handling before and after since their ordering data is different
	if after >= 0 {
		result.Paging.Total = totalTxsNum
		result.Paging.Before = result.Data[0].ID
		result.Paging.After = result.Data[len(result.Data)-1].ID
	} else {
		result.Paging.Total = totalTxsNum
		result.Paging.Before = result.Data[len(result.Data)-1].ID
		result.Paging.After = result.Data[0].ID
	}

	utils.Respond(rw, result)
	return
}

func (a *Account) setTxs(txs []schema.Transaction) (*models.ResultTxs, error) {
	data := make([]models.TxData, 0)

	for _, tx := range txs {
		msgs := make([]models.Message, 0)
		err := json.Unmarshal([]byte(tx.Messages), &msgs)
		if err != nil {
			return &models.ResultTxs{}, fmt.Errorf("failed to unmarshal msgs: %s", err)
		}

		//sigs := make([]models.Signature, 0)
		//err = json.Unmarshal([]byte(tx.Signatures), &sigs)
		//if err != nil {
		//	return &models.ResultTxs{}, fmt.Errorf("failed to unmarshal sigs: %s", err)
		//}
		sigs := make([]models.Signature, 0)

		txResult := true
		if tx.Code != 0 {
			txResult = false
		}

		tempData := &models.TxData{
			ID:          tx.ID,
			Height:      tx.Height,
			Result:      txResult,
			TxHash:      tx.TxHash,
			FromAddress: tx.FromAddress,
			ToAddress:   tx.ToAddress,
			Messages:    msgs,
			Signatures:  sigs,
			Memo:        tx.Memo,
			Code:        tx.Code,
			Timestamp:   tx.Timestamp,
		}

		data = append(data, *tempData)
	}

	result := &models.ResultTxs{
		Data: data,
	}

	return result, nil
}
