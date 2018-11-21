package storagehandler

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
)

type users struct {
	Users []User `json:"users"`
}

type User struct {
	UID      int    `json:"uid"`
	Name     string `json:"name"`
	Password []byte `json:"password"`
	Urlaub   bool   `json:"urlaub"`
}

const saltSize = 16

func saltedHash(secret []byte) []byte {
	buf := make([]byte, saltSize, saltSize+sha256.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Errorf("random read failed: %v", err))
	}
	h := sha256.New()
	h.Write(buf)
	h.Write(secret)
	return h.Sum(buf)
}

func match(data, secret []byte) bool {
	if len(data) != saltSize+sha256.Size {
		fmt.Println("wrong length of data")
	}
	h := sha256.New()
	h.Write(data[:saltSize])
	h.Write(secret)
	return bytes.Equal(h.Sum(nil), data[saltSize:])
}

func verifyUser(username string, password string) bool {
	var user = readSpecificUserFromMemory(username)
	return match(user.Password, []byte(password))
}

func readAllUsersFromMemory() users {
	var byteValue = readJSONFromFile("../../../storage/users.json")
	var users users
	json.Unmarshal(byteValue, &users)
	return users
}

func readSpecificUserFromMemory(userName string) User {
	var specUser User
	var users = readAllUsersFromMemory()
	for _, user := range users.Users {
		if user.Name == userName {
			specUser = user
			break
		}
	}
	return specUser
}

func isUserAvailable(userName string) bool {
	var isUserAvailable = false
	var users = readAllUsersFromMemory()
	for i := 0; i < len(users.Users); i++ {
		if users.Users[i].Name == userName {
			return true
		}
	}
	return isUserAvailable
}

func addUser(userName string, password string) bool {
	if isUserAvailable(userName) {
		fmt.Println("user existiert bereits")
		return false
	}

	var byteValue = readJSONFromFile("../../../storage/users.json")
	var users users
	json.Unmarshal(byteValue, &users)
	var hashedPwd = saltedHash([]byte(password))
	users.Users = append(users.Users, User{0, userName, hashedPwd, false})
	result, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error while add user")
	}
	return writeJSONToFile("../../../storage/users.json", result)
}

func deleteUser(userName string) bool {

	if isUserAvailable(userName) == false {
		fmt.Println("user does not exists")
		return false
	}

	var byteValue = readJSONFromFile("../../../storage/users.json")
	var oldUsers, newUsers users
	json.Unmarshal(byteValue, &oldUsers)
	for i := 0; i < len(oldUsers.Users); i++ {
		if oldUsers.Users[i].Name != userName {
			newUsers.Users = append(newUsers.Users, oldUsers.Users[i])
		}
	}
	result, err := json.Marshal(newUsers)
	if err != nil {
		fmt.Println("Error while delete user")
	}
	return writeJSONToFile("../../../storage/users.json", result)
}
