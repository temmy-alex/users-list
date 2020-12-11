package userRepository

import (
	"database/sql"
	"log"
	"users-list/models"
)

type UserRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b UserRepository) GetUsers(db *sql.DB, user models.User, users []models.User) ([]models.User, error) {
	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		return []models.User{}, err
	}

	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email)
		users = append(users, user)
	}

	if err != nil {
		return []models.User{}, err
	}

	return users, nil
}

func (b UserRepository) GetUser(db *sql.DB, user models.User, id int) (models.User, error) {
	rows := db.QueryRow("SELECT * FROM users WHERE id=$1", id)
	err := rows.Scan(&user.ID, &user.Name, &user.Email)

	return user, err
}

func (b UserRepository) AddUser(db *sql.DB, user models.User) (int, error) {
	err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id;",
		user.Name, user.Email).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (b UserRepository) UpdateUser(db *sql.DB, user models.User) (int64, error) {
	result, err := db.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING id",
		&user.Name, &user.Email, &user.ID)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, nil
}

func (b UserRepository) RemoveUser(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)

	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}
