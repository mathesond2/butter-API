package models

import (
	"github.com/jinzhu/gorm"
)

// type Transaction struct {
// 	gorm.Model
// 	Status               string `json:"status"`
// 	MonitorId            string `json:"monitorId"`
// 	MonitorVersion       string `json:"monitorVersion"`
// 	TimePending          string `json:"timePending"`
// 	BlocksPending        string `json:"blocksPending"`
// 	PendingTimeStamp     string `json:"pendingTimeStamp"`
// 	PendingBlockNumber   string `json:"pendingBlockNumber"`
// 	Hash                 string `json:"hash"`
// 	From                 string `json:"from"`
// 	To                   string `json:"to"`
// 	Value                string `json:"value"`
// 	Gas                  string `json:"gas"`
// 	Nonce                string `json:"nonce"`
// 	BlockHash            string `json:"blockHash"`
// 	BlockNumber          string `json:"blockNumber"`
// 	V                    string `json:"v"`
// 	R                    string `json:"r"`
// 	S                    string `json:"s"`
// 	Input                string `json:"input"`
// 	GasPrice             string `json:"gasPrice"`
// 	GasPriceGwei         string `json:"gasPriceGwei"`
// 	GasUsed              string `json:"gasUsed"`
// 	Type                 string `json:"type"`
// 	MaxFeePerGas         string `json:"maxFeePerGas"`
// 	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
// 	BaseFeePerGas        string `json:"baseFeePerGas"`
// 	TransactionIndex     string `json:"transactionIndex"`
// 	Asset                string `json:"asset"`
// 	BlockTimeStamp       string `json:"blockTimeStamp"`
// 	WatchedAddress       string `json:"watchedAddress"`
// 	Direction            string `json:"direction"`
// 	Counterparty         string `json:"counterparty"`
// 	ServerVersion        string `json:"serverVersion"`
// 	EventCode            string `json:"eventCode"`
// 	TimeStamp            string `json:"timeStamp"`
// 	SispatchTimestamp    string `json:"dispatchTimestamp"`
// 	System               string `json:"system"`
// 	Network              string `json:"network"`
// }

type Transaction struct {
	gorm.Model
	Asset          string `json:"asset"`
	Status         string `json:"status"`
	From           string `json:"from"`
	To             string `json:"to"`
	WatchedAddress string `json:"watchedAddress"`
	Direction      string `json:"direction"`
}

// func (invoice *Transaction) ValidateTransaction() (map[string]interface{}, bool) {
// 	if invoice.Sender_Address == "" {
// 		return u.Message(false, "sender address should be on the payload"), false
// 	}

// 	if invoice.Token_Address == "" {
// 		return u.Message(false, "token address should be on the payload"), false
// 	}

// 	if invoice.Amount <= 0 {
// 		return u.Message(false, "amount should be on the payload"), false
// 	}

// 	if invoice.Recipient_Address == "" {
// 		return u.Message(false, "recipient address should be on the payload"), false
// 	}

// 	if invoice.UserId <= 0 {
// 		return u.Message(false, "User is not recognized"), false
// 	}

// 	return u.Message(true, "success"), true
// }

// func (invoice *Invoice) CreateInvoice() map[string]interface{} {
// 	if resp, ok := invoice.ValidateInvoice(); !ok {
// 		return resp
// 	}

// 	GetDB().Create(invoice)

// 	resp := u.Message(true, "success")
// 	resp["invoice"] = invoice
// 	return resp
// }

// func UpdateInvoice(id uint64, reqinvoice *Invoice) map[string]interface{} {
// 	if resp, ok := reqinvoice.ValidateInvoice(); !ok {
// 		return resp
// 	}

// 	invoice := &Invoice{}
// 	err := GetDB().Table("invoices").Where("id = ?", id).First(invoice).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	invoice.Name = reqinvoice.Name
// 	invoice.Description = reqinvoice.Description
// 	invoice.Sender_Address = reqinvoice.Sender_Address
// 	invoice.Token_Address = reqinvoice.Token_Address
// 	invoice.Amount = reqinvoice.Amount
// 	invoice.To = reqinvoice.To
// 	invoice.Recipient_Address = reqinvoice.Recipient_Address
// 	invoice.Status = reqinvoice.Status

// 	updatedErr := GetDB().Table("invoices").Where("id = ?", id).Save(invoice).Error
// 	if updatedErr != nil {
// 		fmt.Println(updatedErr)
// 		return nil
// 	}

// 	resp := u.Message(true, "success")
// 	resp["invoice"] = invoice
// 	return resp
// }

// func DeleteInvoice(id uint64) *Invoice {
// 	invoice := &Invoice{}

// 	invoiceBeforeDeletion := &Invoice{}
// 	err := GetDB().Table("invoices").Where("id = ?", id).First(invoiceBeforeDeletion).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	deleteErr := GetDB().Table("invoices").Where("id = ?", id).Delete(invoice).Error
// 	if deleteErr != nil {
// 		fmt.Println(deleteErr)
// 		return nil
// 	}
// 	//todo: fix this..returning null
// 	invoiceBeforeDeletion.DeletedAt = invoice.DeletedAt
// 	return invoiceBeforeDeletion
// }

// func GetInvoice(id uint64) *Invoice {
// 	invoice := &Invoice{}
// 	err := GetDB().Table("invoices").Where("id = ?", id).First(invoice).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}
// 	return invoice
// }

// func GetInvoices(user uint) []*Invoice {
// 	invoices := make([]*Invoice, 0)
// 	err := GetDB().Table("invoices").Where("user_id = ?", user).Find(&invoices).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	return invoices
// }
