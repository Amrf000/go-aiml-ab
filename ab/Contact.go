package ab

import (
	"fmt"
	"regexp"
	"strings"
)

var contactCount = 0
var idContactMap = make(map[string]*Contact)
var nameIdMap = make(map[string]string)

type Contact struct {
	ContactId   string
	DisplayName string
	Birthday    string
	Phones      map[string]string
	Emails      map[string]string
}

func MultipleIds(contactName string) string {
	patternString := " (" + strings.ToUpper(contactName) + ") "
	patternString = strings.ReplaceAll(patternString, " ", "(.*)")
	pattern := regexp.MustCompile(patternString)

	var result strings.Builder
	idCount := 0

	for key, value := range nameIdMap {
		matches := pattern.FindStringSubmatch(key)
		if len(matches) > 0 {
			result.WriteString(value + " ")
			idCount++
		}
	}

	if idCount <= 1 {
		return "false"
	}
	return strings.TrimSpace(result.String())
}

func ContactId(contactName string) string {
	patternString := " " + strings.ToUpper(contactName) + " "
	patternString = strings.ReplaceAll(patternString, " ", ".*")
	pattern := regexp.MustCompile(patternString)

	result := "unknown"
	for key, value := range nameIdMap {
		if pattern.MatchString(key) {
			result = value + " "
		}
	}
	return strings.TrimSpace(result)
}

func DisplayName(id string) string {
	c, ok := idContactMap[strings.ToUpper(id)]
	if !ok {
		return "unknown"
	}
	return c.DisplayName
}

func DialNumber(typ, id string) string {
	c, ok := idContactMap[strings.ToUpper(id)]
	if !ok {
		return "unknown"
	}
	dialNumber := c.Phones[strings.ToUpper(typ)]
	if dialNumber == "" {
		return "unknown"
	}
	return dialNumber
}

func EmailAddress(typ, id string) string {
	c, ok := idContactMap[strings.ToUpper(id)]
	if !ok {
		return "unknown"
	}
	emailAddress := c.Emails[strings.ToUpper(typ)]
	if emailAddress == "" {
		return "unknown"
	}
	return emailAddress
}

func Birthday(id string) string {
	c, ok := idContactMap[strings.ToUpper(id)]
	if !ok {
		return "unknown"
	}
	return c.Birthday
}

func NewContact(displayName, phoneType, dialNumber, emailType, emailAddress, birthday string) *Contact {
	contact := &Contact{
		ContactId:   fmt.Sprintf("ID%d", contactCount),
		DisplayName: displayName,
		Birthday:    birthday,
		Phones:      make(map[string]string),
		Emails:      make(map[string]string),
	}
	contactCount++

	contact.IdContactMapAdd()
	contact.AddPhone(phoneType, dialNumber)
	contact.AddEmail(emailType, emailAddress)
	contact.AddName(displayName)
	contact.AddBirthday(birthday)

	return contact
}

func (c *Contact) AddPhone(typ, dialNumber string) {
	c.Phones[strings.ToUpper(typ)] = dialNumber
}

func (c *Contact) AddEmail(typ, emailAddress string) {
	c.Emails[strings.ToUpper(typ)] = emailAddress
}

func (c *Contact) AddName(name string) {
	c.DisplayName = name
	nameIdMap[strings.ToUpper(name)] = c.ContactId
}

func (c *Contact) AddBirthday(birthday string) {
	c.Birthday = birthday
}

func (c *Contact) IdContactMapAdd() {
	idContactMap[strings.ToUpper(c.ContactId)] = c
}
