package data

import (
	"database/sql"
	"time"
)

// Bu dosyada, test dosyalarımızda kullanabilmek için fake veri yarattık.
// models.go içerisindeki kısmı interface olarak düzenlediğimiz için,
// aynı metotları buradaki veriye de implement edip
// hiçbir database bağlantısı olmadan testlerimizi gerçekleştirebiliriz.

func TestNew(dbPool *sql.DB) Models {
	db = dbPool
	return Models{}
}

type UserTest struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    int
	IsAdmin   int
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      *Plan
}

func (u *UserTest) GetAll() ([]*User, error) {
	var users []*User
	user := User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "abc",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	users = append(users, &user)
	return users, nil
}

func (u *UserTest) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "abc",
		Active:    1,
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &user, nil
}

func (u *UserTest) GetOne(id int) (*User, error) {
	// body GetByEmail ile aynı olduğu için direkt onu return ettik, kodu kısaltmak için.
	return u.GetByEmail("")
}

func (u *UserTest) Update() error {
	return nil
}

func (u *UserTest) Delete() error {
	return nil
}

func (u *UserTest) DeleteByID(id int) error {
	return nil
}

func (u *UserTest) Insert(user User) (int, error) {
	return 2, nil
}

func (u *UserTest) ResetPassword(password string) error {
	return nil
}

func (u *UserTest) PasswordMatches(plainText string) (bool, error) {
	return true, nil
}
