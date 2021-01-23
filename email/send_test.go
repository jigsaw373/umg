package email

import (
	"fmt"
	"testing"
)

const (
	emailAddress = "bbahrami@edgecomenergy.ca"
)

func TestSendWelcome(t *testing.T) {
	err := SendWelcome("Behdad", emailAddress)
	if err != nil {
		t.Errorf("Error while sending welcom email: %v", err)
	} else {
		fmt.Println("Welcome email sent")
	}
}
