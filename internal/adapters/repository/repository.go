package repository

import (
	"database/sql"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (uuid, name, email, age, password, role_id) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, user.UUID, user.Name, user.Email, user.Age, user.Password, user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) ReadUser(id string) (*domain.User, error) {
	query := `SELECT id, uuid, name, email, age, password, role_id FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.Age, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) ReadUsers() ([]*domain.User, error) {
	query := `SELECT id, uuid, name, email, age, password, role_id FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.Age, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) UpdateUser(id, email, password string) error {
	query := `UPDATE users SET email=?, password=? WHERE id=?`
	_, err := r.db.Exec(query, email, password, id)
	return err
}

func (r *userRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id=?`
	_, err := r.db.Exec(query, id)
	return err
}
