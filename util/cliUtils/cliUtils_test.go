package cliUtils

import (
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSubjectLine(t *testing.T) {
	subjectline := createSubjectLine("abcd", "")
	assert.Equal(t, "abcd", subjectline)
	subjectline = createSubjectLine("abcd", "3FM")
	assert.Equal(t, "[Ticket 3FM] abcd", subjectline)
}

func TestCreateMail(t *testing.T) {
	eMailAddr := "john.doe@example.com"
	subj := "Search field broken"
	tID := "12ab3"
	mes := "a message"
	expected := structs.Mail{
		Email:   eMailAddr,
		Subject: "[Ticket " + tID + "] " + subj,
		Message: mes}
	actual := CreateMail(eMailAddr, subj, tID, mes)
	assert.Equal(t, expected, actual)
}
func TestCheckEmailAddress(t *testing.T) {
	simple := "aName@address.com"
	withDots := "first.last@more.complicated.address.org"
	dashesAndUnderscores := "mike_miller@impressive-institute.de"
	failWithSpaces := "notAn Address@example.com"
	failAt := "notAnAddress.com"
	failDomain := "first.last@notADomain"
	assert.True(t, CheckEmailAddress(simple))
	assert.True(t, CheckEmailAddress(withDots))
	assert.True(t, CheckEmailAddress(dashesAndUnderscores))
	assert.False(t, CheckEmailAddress(failWithSpaces))
	assert.False(t, CheckEmailAddress(failAt))
	assert.False(t, CheckEmailAddress(failDomain))
}
