package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/niklod/highload-social-network/config"
	"github.com/niklod/highload-social-network/user"
	"github.com/niklod/highload-social-network/user/city"
	"github.com/niklod/highload-social-network/user/interest"
	"github.com/tarantool/go-tarantool"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := dbConnect("mysql", cfg.DB.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	userRepo := user.NewRepository(db)
	cityRepo := city.NewRepository(db)
	interestRepo := interest.NewRepository(db)

	// Services
	cityService := city.NewService(cityRepo)
	interestService := interest.NewService(interestRepo)
	userService := user.NewService(userRepo, cityService, interestService)

	opts := tarantool.Opts{
		Timeout:       0,
		Reconnect:     0,
		MaxReconnects: 0,
		User:          cfg.Tarantool.Login,
		Pass:          cfg.Tarantool.Password,
	}

	server := fmt.Sprintf("%s:%d", cfg.Tarantool.Host, cfg.Tarantool.Port)

	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Ping()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userService.Users()
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		resp, err := client.Insert("user", []interface{}{
			uint(user.ID),
			user.FirstName,
			user.Lastname,
			user.Login,
			uint(user.Age),
			uint(user.City.ID),
			user.City.Name,
			user.Sex,
			user.Password,
		})
		if err != nil {
			fmt.Printf("%+v\n", user)
			log.Fatal(err)
		}
		fmt.Println(resp.Code)
	}
}

func dbConnect(driver, connectionString string) (*sql.DB, error) {
	var connErr error

	for i := 1; i <= 5; i++ {
		fmt.Printf("trying to connect to DB, try %d\n", i)

		if i != 1 {
			time.Sleep(5 * time.Second)
		}

		db, err := sql.Open(driver, connectionString)
		if err != nil {
			connErr = err
			continue
		}

		if err := db.Ping(); err != nil {
			connErr = err
			continue
		}

		return db, connErr
	}

	return nil, connErr
}

func tarantoolConnect(cfg config.Tarantool) (*tarantool.Connection, error) {
	opts := tarantool.Opts{
		Timeout:       0,
		Reconnect:     0,
		MaxReconnects: 0,
		User:          cfg.Login,
		Pass:          cfg.Password,
	}

	server := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client, err := tarantool.Connect(server, opts)
	if err != nil {
		return nil, err
	}

	_, err = client.Ping()
	if err != nil {
		return nil, err
	}

	return client, nil
}
