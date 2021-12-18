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

func GetContact(id uint64) *Contact {
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contact
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
