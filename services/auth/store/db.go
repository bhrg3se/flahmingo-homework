package store

import (
	"errors"
	"time"
)

// CreateUser inserts new user profile into database
func (s Store) CreateUser(user *User) error {
	_, err := s.db.Exec(`INSERT INTO users (id,name,phone_number) VALUES ($1,$2,$3)`, user.ID, user.Name, user.PhoneNumber)
	return err
}

// GetUser returns user profile based on phone number
func (s Store) GetUser(phoneNumber string) (*User, error) {
	var user User
	err := s.db.QueryRow(`SELECT id,phone_number,name,is_verified FROM users WHERE phone_number= $1 `, phoneNumber).
		Scan(&user.ID, &user.PhoneNumber, &user.Name, &user.IsVerified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyUser marks user as verified
func (s Store) VerifyUser(phoneNumber string) error {
	_, err := s.db.Exec(`UPDATE users SET is_verified=true WHERE phone_number=$1`, phoneNumber)
	return err
}

// SaveOTP saves otp in database
func (s Store) SaveOTP(otp, phoneNumber string) error {
	expiry := time.Now().Add(time.Minute * 5)
	//try to update existing row
	res, err := s.db.Exec(`UPDATE otp SET value = $1 , expiry=$2 WHERE phone_number= $3`, otp, expiry, phoneNumber)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// if row does not exist, insert new one
	if rowsAffected < 1 {
		_, err = s.db.Exec(`INSERT INTO otp  (value,expiry,phone_number) VALUES ($1,$2 ,$3)`, otp, expiry, phoneNumber)
	}
	return err
}

// GetOTP returns otp saved in database. It returns error if otp is expired.
func (s Store) GetOTP(phoneNumber string) (string, error) {
	var expiry time.Time
	var otp string
	err := s.db.QueryRow(`SELECT value,expiry FROM otp WHERE phone_number= $1`, phoneNumber).Scan(&otp, &expiry)
	if err != nil {
		return otp, err
	}
	if time.Now().After(expiry) {
		return otp, errors.New("otp expired")
	}
	return otp, nil
}
