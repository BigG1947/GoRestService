package main

import (
	"GoRestService/db"
)

func refreshCash() (map[int64]db.User, error) {
	cash = make(map[int64]db.User)

	users, err := db.GetAllUsers(connection)
	if err != nil {
		return nil, err
	}

	for i := range users {
		cash[users[i].Id] = users[i]
	}
	return cash, nil
}
