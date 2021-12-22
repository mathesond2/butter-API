package models

import (
	"fmt"
	u "go-invoices/utils"

	"github.com/jinzhu/gorm"
)

type Invoice struct {
	gorm.Model
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	UserId            uint    `json:"user_id"`
	Sender_Address    string  `json:"sender_address"`
	Token_Address     string  `json:"token_address"`
	Amount            float64 `json:"amount"`
	To                string  `json:"to"`
	Recipient_Address string  `json:"recipient_address"`
	Status            string  `json:"status"`
}

func (invoice *Invoice) ValidateInvoice() (map[string]interface{}, bool) {
	if invoice.Sender_Address == "" {
		return u.Message(false, "sender address should be on the payload"), false
	}

	if invoice.Token_Address == "" {
		return u.Message(false, "token address should be on the payload"), false
	}

	if invoice.Amount <= 0 {
		return u.Message(false, "amount should be on the payload"), false
	}

	if invoice.Recipient_Address == "" {
		return u.Message(false, "recipient address should be on the payload"), false
	}

	if invoice.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	return u.Message(true, "success"), true
}

func (invoice *Invoice) CreateInvoice() map[string]interface{} {
	if resp, ok := invoice.ValidateInvoice(); !ok {
		return resp
	}

	GetDB().Create(invoice)

	resp := u.Message(true, "success")
	resp["invoice"] = invoice
	return resp
}

func UpdateInvoice(id uint64, reqinvoice *Invoice) map[string]interface{} {
	if resp, ok := reqinvoice.ValidateInvoice(); !ok {
		return resp
	}

	invoice := &Invoice{}
	err := GetDB().Table("invoices").Where("id = ?", id).First(invoice).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	invoice.Name = reqinvoice.Name
	invoice.Description = reqinvoice.Description
	invoice.Sender_Address = reqinvoice.Sender_Address
	invoice.Token_Address = reqinvoice.Token_Address
	invoice.Amount = reqinvoice.Amount
	invoice.To = reqinvoice.To
	invoice.Recipient_Address = reqinvoice.Recipient_Address
	invoice.Status = reqinvoice.Status

	updatedErr := GetDB().Table("invoices").Where("id = ?", id).Save(invoice).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		return nil
	}

	resp := u.Message(true, "success")
	resp["invoice"] = invoice
	return resp
}

func DeleteInvoice(id uint64) *Invoice {
	invoice := &Invoice{}

	invoiceBeforeDeletion := &Invoice{}
	err := GetDB().Table("invoices").Where("id = ?", id).First(invoiceBeforeDeletion).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	deleteErr := GetDB().Table("invoices").Where("id = ?", id).Delete(invoice).Error
	if deleteErr != nil {
		fmt.Println(deleteErr)
		return nil
	}
	//todo: fix this..returning null
	invoiceBeforeDeletion.DeletedAt = invoice.DeletedAt
	return invoiceBeforeDeletion
}

func GetInvoice(id uint64) *Invoice {
	invoice := &Invoice{}
	err := GetDB().Table("invoices").Where("id = ?", id).First(invoice).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return invoice
}

func GetInvoices(user uint64) []*Invoice {
	invoices := make([]*Invoice, 0)
	err := GetDB().Table("invoices").Where("user_id = ?", user).Find(&invoices).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return invoices
}

func GetAssociatedTxn(txn *ParsedTransaction) *Invoice {
	invoice := &Invoice{}

	err := GetDB().Table("invoices").Where(&Invoice{
		Sender_Address:    txn.WatchedAddress,
		Recipient_Address: txn.To,
		Amount:            txn.Value,
	}).First(&invoice).Error
	if err != nil {
		fmt.Println(err, "GetAssociatedTxn")
		return nil
	}

	var status string
	if txn.Status == "confirmed" {
		status = "confirmed"
	} else {
		status = "pending"
	}
	invoice.Status = status

	updatedErr := GetDB().Table("invoices").Where("id = ?", invoice.ID).Save(invoice).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		return nil
	}

	return invoice
}
