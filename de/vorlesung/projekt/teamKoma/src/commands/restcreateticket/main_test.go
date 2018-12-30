//Matrikelnummern:
//9188103
//1798794
//4717960

//Auf weitere Tests wurde aufgrund der fehlenden Konnektivität zu einem Ticketsystem verzichtet!
package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestGenJSONData(t *testing.T) {
	assert.NotNil(t, genJSONData("test@test.com", "test", "test", "test", "test"))
}
