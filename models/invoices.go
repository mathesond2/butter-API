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
	Webhook_Name      string  `json:"webhook_name"`
}

type UserInfo struct {
	UserId uint `json:"user_id"`
	ID     uint `json:"ID"`
}

func (i *Invoice) ValidateInvoice() (map[string]interface{}, bool) {
	if i.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	if i.Sender_Address == "" {
		return u.Message(false, "sender address should be on the payload"), false
	}

	if i.Token_Address == "" {
		return u.Message(false, "token address should be on the payload"), false
	}

	if i.Amount <= 0 {
		return u.Message(false, "amount should be on the payload"), false
	}

	if i.Recipient_Address == "" {
		return u.Message(false, "recipient address should be on the payload"), false
	}

	if i.Status != "paid" && i.Status != "unpaid" {
		return u.Message(false, "invoice status must be either 'paid' or 'unpaid'"), false
	}

	return u.Message(true, "success"), true
}

func (i *Invoice) CreateInvoice() map[string]interface{} {
	if resp, ok := i.ValidateInvoice(); !ok {
		return resp
	}

	GetDB().Create(i)
	resp := u.Message(true, "success")
	resp["data"] = i
	return resp
}

func UpdateInvoice(reqinvoice *Invoice) map[string]interface{} {
	if resp, ok := reqinvoice.ValidateInvoice(); !ok {
		return resp
	}

	invoice := &Invoice{}
	err := GetDB().Table("invoices").Where("id = ? AND user_id = ?", reqinvoice.ID, reqinvoice.UserId).First(invoice).Error
	if err != nil {
		fmt.Println(err)
		errStr := err.Error()
		return u.Message(false, errStr)
	}

	invoice.Name = reqinvoice.Name
	invoice.Description = reqinvoice.Description
	invoice.Sender_Address = reqinvoice.Sender_Address
	invoice.Token_Address = reqinvoice.Token_Address
	invoice.Amount = reqinvoice.Amount
	invoice.To = reqinvoice.To
	invoice.Recipient_Address = reqinvoice.Recipient_Address
	invoice.Status = reqinvoice.Status
	invoice.Webhook_Name = reqinvoice.Webhook_Name

	updatedErr := GetDB().Table("invoices").Where("id = ?", reqinvoice.ID).Save(invoice).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		errStr := err.Error()
		return u.Message(false, errStr)
	}

	resp := u.Message(true, "success")
	resp["data"] = invoice
	return resp
}

func DeleteInvoice(userInfo *UserInfo) (string, error) {
	invoice := &Invoice{}
	err := GetDB().Table("invoices").Where("id = ? AND user_id = ?", userInfo.ID, userInfo.UserId).First(invoice).Error
	if err != nil {
		fmt.Println("DeleteInvoice find error: ", err)
		return "", err
	}

	deleteErr := GetDB().Table("invoices").Where("id = ? AND user_id = ?", userInfo.ID, userInfo.UserId).Delete(invoice).Error
	if deleteErr != nil {
		fmt.Println("DeleteInvoice delete error: ", deleteErr)
		return "", deleteErr
	}

	return "invoice successfully deleted", nil
}

func GetInvoice(id uint) *Invoice {
	i := &Invoice{}

	err := GetDB().Table("invoices").Where("id = ?", id).First(i).Error
	if err != nil {
		fmt.Println("GetInvoice error: ", err)
		return nil
	}

	return i
}

func UpdateInvoiceStatusAsPending(id uint) *Invoice {
	invoice := &Invoice{}

	err := GetDB().Table("invoices").Where("id = ?", id).First(invoice).Error
	if err != nil {
		fmt.Println("UpdateInvoiceStatusAsPending error: ", err)
		return nil
	}

	invoice.Status = "pending"

	updatedErr := GetDB().Table("invoices").Where("id = ?", invoice.ID).Save(invoice).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		return nil
	}

	return invoice
}

func GetInvoices(user uint) []*Invoice {
	invoices := make([]*Invoice, 0)

	err := GetDB().Table("invoices").Where("user_id = ?", user).Find(&invoices).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return invoices
}
