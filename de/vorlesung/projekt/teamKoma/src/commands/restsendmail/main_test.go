//Matrikelnummern:
//9188103
//1798794
//4717960
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrabMailsToSend(t *testing.T) {
	assert.Nil(t, grabMailsToSend("localhost", 8443, "Werner", "password"))
}

func TestSetSentFlag(t *testing.T) {
	mails2send := grabMailsToSend("localhost", 8443, "Werner", "password")
	assert.Nil(t, setSentFlag("localhost", 8443, "Werner", "password", mails2send))
}
