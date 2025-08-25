package postgres

import (
	"testing"

	"github.com/latimeri-compute/go-exam-exchanger/gw-currency-wallet/internal/storages"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("пропуск интеграционных тестов")
	}

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "неповторяющийся",
			email:    "newuser@new.com",
			password: "password",
		},
		{
			name:     "повторяющийся",
			email:    "newuser@new.com",
			password: "password",
			wantErr:  storages.ErrRecordExists,
		},
	}

	db := newTestDB(t)
	setupDB(t, db)
	defer teardownDB(t, db)

	model := NewUserModel(db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bytes, err := bcrypt.GenerateFromPassword([]byte(test.password), bcrypt.DefaultCost)
			if err != nil {
				t.Fatal(err)
			}

			user := &storages.User{
				Email:        test.email,
				PasswordHash: bytes,
			}
			err = model.CreateUser(user)
			assert.NotEqual(t, 0, user.ID)
			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}

func TestFindUser(t *testing.T) {
	if testing.Short() {
		t.Skip("пропуск интеграционных тестов")
	}

	tests := []struct {
		name    string
		user    storages.User
		wantErr error
	}{
		{
			name: "существующий пользователь",
			user: storages.User{
				Email: "hello@hello.com",
			},
			wantErr: nil,
		},
		{
			name: "несуществующий пользователь",
			user: storages.User{
				Email: "noooo",
			},
			wantErr: storages.ErrRecordNotFound,
		},
		{
			name:    "пустой пользователь",
			wantErr: storages.ErrRecordNotFound,
		},
	}
	db := newTestDB(t)
	setupDB(t, db)
	defer teardownDB(t, db)

	model := NewUserModel(db)

	ph, err := bcrypt.GenerateFromPassword([]byte("asdasd"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	create := &storages.User{
		Email:        tests[0].user.Email,
		PasswordHash: ph,
	}
	err = model.CreateUser(create)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			find := storages.User{
				Email: test.user.Email,
			}
			err := model.FindUser(&find)
			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}
