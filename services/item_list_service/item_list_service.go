package item_list_service

import (
	"database/sql"

	"SplendifeList-Server-Go/models"
)

func GetAllLists(dbConnection *sql.DB) ([]models.ItemList, error) {
	rows, err := dbConnection.Query(
		"SELECT id, name, crossed_out, user FROM item_lists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemLists []models.ItemList
	for rows.Next() {
		var itemList models.ItemList
		err = rows.Scan(&itemList.Id, &itemList.Name, &itemList.CrossedOut,
			&itemList.User)
		if err != nil {
			return nil, err
		}

		itemLists = append(itemLists, itemList)
	}

	return itemLists, nil
}

func CreateList(dbConnection *sql.DB, newList models.ItemList) (int, error) {
	result, err := dbConnection.Exec(
		"INSERT INTO item_lists (name, user, created_at, updated_at)"+
			" VALUES (?, ?, NOW(), NOW())",
		newList.Name,
		newList.User,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
