package user

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/niklod/highload-social-network/internal/user/city"
	"github.com/niklod/highload-social-network/internal/user/interest"
)

type UserCreateRequest struct {
	Login     string `form:"inputLogin" validate:"required,min=5,max=20"`
	Password  string `form:"inputPassword" validate:"required,min=6,max=40"`
	FirstName string `form:"inputName" validate:"required,max=50"`
	LastName  string `form:"inputLastName" validate:"required,max=50"`
	Age       int    `form:"inputAge" validate:"gte=0,lte=120"`
	Sex       string `form:"inputSex" validate:""`
	City      string `form:"inputCity" validate:""`
	Interests string `form:"inputInterests" validate:""`
}

func (u *UserCreateRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserCreateRequest) ConverIntoUser() *User {
	interests := []interest.Interest{}

	splitedFormInterests := strings.Split(strings.ReplaceAll(u.Interests, " ", ""), ",")

	for _, name := range splitedFormInterests {
		interests = append(interests, interest.Interest{Name: strings.ToLower(name)})
	}

	return &User{
		FirstName: u.FirstName,
		Lastname:  u.LastName,
		Age:       u.Age,
		Sex:       u.Sex,
		City:      city.City{Name: u.City},
		Login:     u.Login,
		Password:  u.Password,
		Interests: interests,
	}
}

type UserLoginRequest struct {
	Login    string `form:"inputLogin" validate:"required"`
	Password string `form:"inputPassword" validate:"required"`
}

func (u *UserLoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type fieldError struct {
	err validator.FieldError
}

func (q fieldError) String() string {
	var sb strings.Builder

	sb.WriteString("Не пройдена валидация поля '" + q.err.Field() + "'")

	return sb.String()
}

type UserSearchRequest struct {
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
}
