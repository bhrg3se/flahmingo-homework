package store

type User struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	IsVerified  bool   `db:"is_verified"`
	PhoneNumber string `db:"phone_number"`
}
