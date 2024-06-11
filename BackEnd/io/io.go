package io

import (
	"BackEnd/config"
	"BackEnd/models"
	"encoding/json"
	"io/ioutil"
)

// 函数功能：从文件中读取数据
func LoadTodos() ([]models.Todo, error) {
	data, err := ioutil.ReadFile(config.TodosFile)
	if err != nil {
		return nil, err
	}

	var todos []models.Todo
	err = json.Unmarshal(data, &todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

// 函数功能：将数据保存至文件中
func SaveTodos(todos []models.Todo) error {
	data, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.TodosFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// 函数功能:读取用户数据
func LoadUsers() ([]models.User, error) {
	data, err := ioutil.ReadFile(config.UsersFile)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// 函数功能:保存用户数据
func SaveUsers(users []models.User) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.UsersFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
