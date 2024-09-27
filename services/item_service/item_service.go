package item_service

import (
	"database/sql"

	"SplendifeList-Server-Go/models"
)

func GetItemsByItemList(
	dbConnection *sql.DB,
	itemListId string,
	userId string,
) ([]models.Item, error) {
	rows, err := dbConnection.Query(
		`SELECT id, name, crossed_out
			FROM items
			WHERE item_list_id = ?
			AND user_id = ?`,
		itemListId,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.Id, &item.Name, &item.CrossedOut)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func CreateItem(dbConnection *sql.DB, newItem models.Item) (int, error) {
	result, err := dbConnection.Exec(
		"INSERT INTO items (name, user, created_at, updated_at)"+
			" VALUES (?, ?, NOW(), NOW())",
		newItem.Name,
		newItem.User,
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
