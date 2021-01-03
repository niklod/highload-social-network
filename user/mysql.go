package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/niklod/highload-social-network/user/city"
)

type mysql struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return &mysql{db: db}
}

func (m *mysql) Create(user *User) (*User, error) {
	query := queryMap[createQuery]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	res, err := m.db.ExecContext(ctx, query.SQL,
		sql.Named("first_name", user.FirstName),
		sql.Named("last_name", user.Lastname),
		sql.Named("age", user.Age),
		sql.Named("sex_id", user.Sex.ID),
		sql.Named("city_id", user.City.ID),
		sql.Named("password", user.Password),
		sql.Named("login", user.Login),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new user: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("gettings last insert id: %v", err)
	}

	user.ID = int(id)

	return user, nil
}

func (m *mysql) List() ([]User, error) {
	query := queryMap[createQuery]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query.SQL)
	if err != nil {
		return nil, fmt.Errorf("list users: %v", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		var cityName sql.NullString
		var cityID sql.NullInt64

		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.Lastname,
			&u.Age,
			&u.Sex.ID,
			&u.Sex.Name,
			&u.Login,
			&cityID,
			&cityName,
		)
		if err != nil {
			log.Printf("scanning users list row: %v", err)
			continue
		}

		u.City = city.City{}

		if cityName.Valid && cityID.Valid {
			u.City.Name = cityName.String
			u.City.ID = int(cityID.Int64)
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating through rows: %v", err)
	}

	return users, nil
}

func (m *mysql) GetByID(id int) (*User, error) {
	query := queryMap[getByID]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	var user User
	var cityName sql.NullString
	var cityID sql.NullInt64

	row := m.db.QueryRowContext(ctx, query.SQL, sql.Named("id", id))
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.Lastname,
		&user.Age,
		&user.Sex.ID,
		&user.Sex.Name,
		&user.Login,
		&cityID,
		&cityName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user id %d not found: %v", id, err)
		}
		return nil, fmt.Errorf("get user by id: scanning user sql row: %v", err)
	}

	user.City = city.City{}

	if cityName.Valid && cityID.Valid {
		user.City.Name = cityName.String
		user.City.ID = int(cityID.Int64)
	}

	return &user, nil
}

func (m *mysql) GetByLogin(login string) (*User, error) {
	query := queryMap[getByLogin]

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	var user User
	var cityName sql.NullString
	var cityID sql.NullInt64

	row := m.db.QueryRowContext(ctx, query.SQL, sql.Named("login", login))
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.Lastname,
		&user.Age,
		&user.Sex.ID,
		&user.Sex.Name,
		&user.Login,
		&cityID,
		&cityName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by login: scanning user sql row: %v", err)
	}

	user.City = city.City{}

	if cityName.Valid && cityID.Valid {
		user.City.Name = cityName.String
		user.City.ID = int(cityID.Int64)
	}

	return &user, nil
}