//Matrikelnummern:
//9188103
//1798794
//4717960
package storagehandler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
)

// User struct defines the user information of processors
type User struct {
	UID        int    `json:"uid"`
	Name       string `json:"name"`
	Password   []byte `json:"password"`
	HasHoliday bool   `json:"hasHoliday"`
}

const saltSize = 16

func saltedHash(secret []byte) []byte {
	buf := make([]byte, saltSize, saltSize+sha256.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		log.Fatal("Error gen salt")
	}
	h := sha256.New()
	h.Write(buf)
	h.Write(secret)
	return h.Sum(buf)
}

func match(data, secret []byte) bool {
	if len(data) != saltSize+sha256.Size {
		log.Fatal("wrong length of data")
	}
	h := sha256.New()
	h.Write(data[:saltSize])
	h.Write(secret)
	return bytes.Equal(h.Sum(nil), data[saltSize:])
}

func (handler *StorageHandler) verifyUser(username string, password string) bool {
	var user = handler.GetUserByUserName(username)
	return match(user.Password, []byte(password))
}

func (handler *StorageHandler) loadUserFromMemory() []User {
	var byteValue = readJSONFromFile(handler.userStoreFile)
	json.Unmarshal(byteValue, &handler.users)
	return handler.users
}

// GetUserByUserName return the user by the given name
func (handler *StorageHandler) GetUserByUserName(userName string) User {
	var specUser User
	var users = handler.GetUsers()
	for _, user := range *users {
		if user.Name == userName {
			specUser = user
			break
		}
	}
	return specUser
}

func (handler *StorageHandler) toggleHoliday(userName string) bool {
	var newHolidayState = false
	var users = handler.GetUsers()
	for i := 0; i < len(*users); i++ {
		if (*users)[i].Name == userName {
			if (*users)[i].HasHoliday == false {
				(*users)[i].HasHoliday = true
				newHolidayState = true
			} else {
				(*users)[i].HasHoliday = false
			}
		}
	}
	return newHolidayState
}

func (handler *StorageHandler) isUserAvailable(userName string) bool {
	var isUserAvailable = false
	var users = handler.GetUsers()
	for i := 0; i < len(*users); i++ {
		if (*users)[i].Name == userName {
			return true
		}
	}
	return isUserAvailable
}

func (handler *StorageHandler) addUser(userName string, password string) bool {
	if handler.isUserAvailable(userName) {
		return false
	}
	var hashedPwd = saltedHash([]byte(password))
	handler.users = append((handler).users, User{0, userName, hashedPwd, false})
	result, err := json.Marshal(handler.users)
	if err != nil {
		log.Fatal("Error while add user")
	}
	return writeJSONToFile(handler.userStoreFile, result)
}

func (handler *StorageHandler) deleteUser(userName string) bool {

	if handler.isUserAvailable(userName) == false {
		return false
	}
	var i int
	for i = 0; i < len(*handler.GetUsers()); i++ {
		if (*handler.GetUsers())[i].Name == userName {
			break
		}
	}
	handler.users[i] = handler.users[len(handler.users)-1]
	handler.users[len(handler.users)-1] = User{}
	handler.users = handler.users[:len(handler.users)-1]

	result, err := json.Marshal(*handler.GetUsers())
	if err != nil {
		log.Fatal("Error while delete user")
	}
	return writeJSONToFile(handler.userStoreFile, result)
}
