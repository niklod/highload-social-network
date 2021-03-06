package post

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_mysql_PostsByUserId_OneRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	userId := 22

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "body", "first_name", "last_name", "login", "id"})
	rows.AddRow(1, time.Now(), time.Now(), "Test", "TestFirst", "TestLast", "Testlogin", 1)

	mock.ExpectQuery("SELECT p.id").WithArgs(userId).WillReturnRows(rows)

	res, err := repo.PostsByUserId(userId)

	assert.Equal(t, 1, len(res))
	assert.Nil(t, err)
}

func Test_mysql_PostsByUserId_TwoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	userId := 22

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "body", "first_name", "last_name", "login", "id"})
	rows.AddRow(1, time.Now(), time.Now(), "Test", "TestFirst", "TestLast", "TestLogin", 1)
	rows.AddRow(1, time.Now(), time.Now(), "Test", "TestFirst", "TestLast", "TestLogin", 1)

	mock.ExpectQuery("SELECT p.id").WithArgs(userId).WillReturnRows(rows)

	res, err := repo.PostsByUserId(userId)

	assert.Equal(t, 2, len(res))
	assert.Nil(t, err)
}

func Test_mysql_PostsByUserId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	sqlError := fmt.Errorf("test sql error")

	mock.ExpectQuery("SELECT p.id").WillReturnError(sqlError)

	res, err := repo.PostsByUserId(2)

	assert.Nil(t, res)
	assert.Contains(t, err.Error(), sqlError.Error())
}

func Test_mysql_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	post := &Post{
		ID:        0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Body:      "Test",
	}
	userId := 22

	mock.ExpectExec("INSERT INTO posts").WithArgs(userId, post.Body).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add(post, userId)

	assert.Nil(t, err)
}

func Test_mysql_UserFeed_OneRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	userId := 22

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "body", "first_name", "last_name", "login", "id"})
	rows.AddRow(1, time.Now(), time.Now(), "Test", "TestFirst", "TestLast", "Testlogin", 1)

	mock.ExpectQuery("SELECT p.id").WithArgs(userId).WillReturnRows(rows)

	res, err := repo.UserFeed(userId)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func Test_mysql_UserFeed_TwoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	userId := 22

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "body", "first_name", "last_name", "login", "id"})
	rows.AddRow(1, time.Now(), time.Now(), "Test", "TestFirst", "TestLast", "Testlogin", 1)
	rows.AddRow(1, time.Now(), time.Now(), "Test1", "TestFirst1", "TestLast1", "Testlogin1", 1)

	mock.ExpectQuery("SELECT p.id").WithArgs(userId).WillReturnRows(rows)

	res, err := repo.UserFeed(userId)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
}

func Test_mysql_UserFeed_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepository(db)

	userId := 22
	testErr := fmt.Errorf("test user feed error")

	mock.ExpectQuery("SELECT p.id").WithArgs(userId).WillReturnError(testErr)

	res, err := repo.UserFeed(userId)

	assert.Nil(t, res)
	assert.Contains(t, err.Error(), testErr.Error())
}
