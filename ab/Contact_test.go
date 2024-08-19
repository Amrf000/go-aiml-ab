package ab

import (
	"fmt"
	"testing"
)

func TestContact(t *testing.T) {
	// Example usage:
	contact := NewContact("John Doe", "Home", "123-456-7890", "Work", "john.doe@example.com", "2000-01-01")
	fmt.Println("Contact ID:", ContactId("John"))
	fmt.Println("Display Name:", DisplayName(contact.ContactId))
	fmt.Println("Dial Number (Home):", DialNumber("Home", contact.ContactId))
	fmt.Println("Email Address (Work):", EmailAddress("Work", contact.ContactId))
	fmt.Println("Birthday:", Birthday(contact.ContactId))
	fmt.Println("Multiple IDs (John):", MultipleIds("John"))
}
