package repos

import (
	"database/sql"
	"fmt"
	"friendfy-api/src/models"
)

type users struct {
	db *sql.DB
}

func NewRepositoryOfUsers(db *sql.DB) *users {
	return &users{db: db}
}

// rep -> repository of users
func (rep users) Create(user models.User) (uint64, error) {
	statement, err := rep.db.Prepare("INSERT INTO users (name, nick, email, password, created_at) VALUES (?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (rep users) Get(nameOrNick string) ([]models.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick) // LIKE %nameOrNick%

	rows, err := rep.db.Query("SELECT id, name, nick, email, created_at FROM users WHERE name LIKE ? OR nick LIKE ?", nameOrNick, nameOrNick)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)

		if len(users) == 0 {
			return nil, sql.ErrNoRows
		}
	}

	return users, nil
}

func (rep users) GetByID(id uint64) (models.User, error) {
	rows, err := rep.db.Query("SELECT id, name, nick, email, created_at FROM users WHERE id = ?", id)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	var user models.User

	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt)
		if err != nil {
			return models.User{}, err
		}
	} else {
		return models.User{}, sql.ErrNoRows
	}

	return user, nil
}

func (rep users) Update(id uint64, user models.User) error {
	statement, err := rep.db.Prepare("UPDATE users SET name = ?, nick = ?, email = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Nick, user.Email, id); err != nil {
		return err
	}

	return nil
}

func (rep users) Delete(id uint64) error {
	statement, err := rep.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return err
	}

	return nil
}

func (rep users) SearchByEmail(email string) (models.User, error) {
	row, err := rep.db.Query("SELECT id, password FROM users WHERE email = ?", email)
	if err != nil {
		return models.User{}, err
	}
	defer row.Close()

	var user models.User
	if row.Next() {
		err := row.Scan(&user.ID, &user.Password)
		if err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (rep users) Follow(userID uint64, followerID uint64) error {
	statement, err := rep.db.Prepare("INSERT IGNORE INTO followers (user_id, follower_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return err
	}

	return nil
}

func (rep users) Unfollow(userID uint64, followerID uint64) error {
	statement, err := rep.db.Prepare("DELETE FROM followers WHERE user_id = ? AND follower_id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return err
	}

	return nil
}

func (rep users) SearchFollowers(userID uint64) ([]models.User, error) {

	// Here down below is how we get the info of all the followers of a specific user
	rows, err := rep.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.created_at 
		FROM users u inner join followers f on u.id = f.follower_id where f.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)

		if len(users) == 0 {
			return nil, sql.ErrNoRows
		}
	}

	return users, nil
}

func (rep users) GetFollowingUsers(userID uint64) ([]models.User, error) {

	rows, err := rep.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.created_at 
    FROM users u inner join followers f on u.id = f.user_id where f.follower_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)

		if len(users) == 0 {
			return nil, sql.ErrNoRows
		}
	}

	return users, nil
}

func (rep users) SearchPassword(followerID uint64) (string, error) {
	row, err := rep.db.Query("SELECT password FROM users WHERE id = ?", followerID)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var user models.User
	if row.Next() {
		err := row.Scan(&user.Password)
		if err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (rep users) UpdatePassword(userID uint64, password string) error {
	statement, err := rep.db.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(password, userID); err != nil {
		return err
	}

	return nil
}
