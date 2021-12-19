package models

import (
	"fmt"
	u "go-contacts/utils"

	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserId uint   `json:"user_id"` //The user that this contact belongs to
}

type Invoice struct {
	gorm.Model
	Name              string `json:"name"`
	Description       string `json:"description"`
	UserId            uint   `json:"user_id"` //The user that this invoice belongs to
	Sender_Address    string `json:"sender_address"`
	Token_Address     string `json:"token_address"`
	Amount            uint   `json:"amount"`
	To                string `json:"to"`
	Recipient_Address string `json:"recipient_address"`
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

/*
 This struct function validate the required parameters sent through the http request body

returns message and true if the requirement is met
*/
func (contact *Contact) Validate() (map[string]interface{}, bool) {
	if contact.Name == "" {
		return u.Message(false, "Contact name should be on the payload"), false
	}

	if contact.Phone == "" {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	if contact.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	//All the required parameters are present
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

func (contact *Contact) Create() map[string]interface{} {
	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := u.Message(true, "success")
	resp["contact"] = contact
	return resp
}

func Update(id uint64, reqcontact *Contact) map[string]interface{} {
	if resp, ok := reqcontact.Validate(); !ok {
		return resp
	}

	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	contact.Name = reqcontact.Name
	contact.Phone = reqcontact.Phone

	updatedErr := GetDB().Table("contacts").Where("id = ?", id).Save(contact).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		return nil
	}

	resp := u.Message(true, "success")
	resp["contact"] = contact
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

	updatedErr := GetDB().Table("invoices").Where("id = ?", id).Save(invoice).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		return nil
	}

	resp := u.Message(true, "success")
	resp["invoice"] = invoice
	return resp
}

func DeleteContact(id uint64) *Contact {
	contact := &Contact{}

	contactBeforeDeletion := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contactBeforeDeletion).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	deleteErr := GetDB().Table("contacts").Where("id = ?", id).Delete(contact).Error
	if deleteErr != nil {
		fmt.Println(deleteErr)
		return nil
	}
	//todo: fix this..returning null
	contactBeforeDeletion.DeletedAt = contact.DeletedAt
	return contactBeforeDeletion
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

func GetContact(id uint64) *Contact {
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contact
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

func GetContacts(user uint) []*Contact {
	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return contacts
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
