package model

type User struct {
	Id           string `json:"id,omitempty" db:"id"`
	Email        string `json:"email,omitempty" db:"email"`
	PasswordHash string `json:"password_hash,omitempty" db:"password_hash"`
	Name         string `json:"name,omitempty" db:"name"`
	Surname      string `json:"surname,omitempty" db:"surname"`
	BirthdayDate string `json:"birthday_date,omitempty" db:"birthday_date"`
}
