package main

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	CreateUser(*User) (int, error)
	GetUserByID(int) (*User, error)
	GetUserByUsername(string) (*User, error)
	UpdateUser(*User) (*User, error)
}

type JSONFileStorage struct {
	path	string
}

type JSONSchema struct {
	IDCount	int		`json:"idCount"`
	Users	[]*User	`json:"users"`
}

func NewJSONFileStorage(path string) (*JSONFileStorage, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return Init(path)
		}
		return nil, err
	}
	return &JSONFileStorage{path: path}, nil
}

func Init(path string) (*JSONFileStorage, error) {
	data, err := json.Marshal(JSONSchema{})
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return nil, err
	}
	return &JSONFileStorage{path: path}, nil
}

func (fs *JSONFileStorage) GetUserByID(id int) (*User, error) {
	data, err := fs.fetchData()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch data")
	}
	for _, usr := range data.Users {
		if usr.ID == id {
			return usr, nil
		}
	}
	return nil, fmt.Errorf("user id %d does not exist", id)
}

func (fs *JSONFileStorage) GetUserByUsername(username string) (*User, error) {
	data, err := fs.fetchData()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch data")
	}
	for _, usr := range data.Users {
		if usr.Username == username {
			return usr, nil
		}
	}
	return nil, fmt.Errorf("username %s does not exist", username)
}

func (fs *JSONFileStorage) CreateUser(u *User) (int, error) {
	data, err := fs.fetchData()
	if err != nil {
		return 0, fmt.Errorf("cannot fetch data")
	}
	data.IDCount++
	u.ID = data.IDCount
	enc, err := bcrypt.GenerateFromPassword([]byte(u.EncPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u.EncPassword = string(enc)
	data.Users = append(data.Users, u)
	if err = fs.saveData(data); err != nil {
		return 0, err
	}
	return data.IDCount, nil
}

func (fs *JSONFileStorage) UpdateUser(u *User) (*User, error) {
	data, err := fs.fetchData()
	if err != nil {
		return nil, fmt.Errorf("cannot fetch data")
	}
	for _, usr := range data.Users {
		if usr.ID == u.ID {
			enc, err := bcrypt.GenerateFromPassword([]byte(u.EncPassword), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			usr.EncPassword = string(enc)
			usr.UserProfile = u.UserProfile
			if err = fs.saveData(data); err != nil {
				return nil, err
			}
			return usr, nil
		}
	}
	return nil, fmt.Errorf("user id %d does not exist", u.ID)
}

func (fs *JSONFileStorage) fetchData() (*JSONSchema, error) {
	f, err := os.Open(fs.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data JSONSchema
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (fs *JSONFileStorage) saveData(d *JSONSchema) error {
	f, err := os.OpenFile(fs.path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = f.Truncate(0); err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err = enc.Encode(d); err != nil {
		return err
	}
	return nil
}