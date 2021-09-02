package store

import (
	"context"
	"crypto/rsa"
	"github.com/bhrg3se/flahmingo-homework/utils"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetConfig() utils.Config {
	args := m.Called()
	return args.Get(0).(utils.Config)
}

func (m *MockStore) CreateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockStore) GetUser(phoneNumber string) (*User, error) {
	args := m.Called(phoneNumber)
	r0, r1 := args.Get(0), args.Error(1)
	if r0 == nil {
		return nil, r1
	}
	return r0.(*User), r1
}

func (m *MockStore) SendOTP(ctx context.Context, otp, phoneNumber string) {
	m.Called(ctx, otp, phoneNumber)
}

func (m *MockStore) SaveOTP(otp, phoneNumber string) error {
	args := m.Called(otp, phoneNumber)
	return args.Error(0)
}

func (m *MockStore) GetOTP(phoneNumber string) (string, error) {
	args := m.Called(phoneNumber)
	return args.String(0), args.Error(1)
}

func (m *MockStore) VerifyUser(phoneNumber string) error {
	args := m.Called(phoneNumber)
	return args.Error(0)
}

func (m *MockStore) GetJWTPublicKey() *rsa.PublicKey {
	args := m.Called()
	return args.Get(0).(*rsa.PublicKey)
}

func (m *MockStore) GetJWTPrivateKey() *rsa.PrivateKey {
	args := m.Called()
	return args.Get(0).(*rsa.PrivateKey)
}
