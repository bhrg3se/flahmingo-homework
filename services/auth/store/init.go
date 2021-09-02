package store

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/bhrg3se/flahmingo-homework/utils"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type GenericStore interface {
	GetConfig() utils.Config
	CreateUser(user *User) error
	GetUser(phoneNumber string) (*User, error)
	PublishOTP(ctx context.Context, otp, phoneNumber string)
	SaveOTP(otp, phoneNumber string) error
	GetOTP(phoneNumber string) (string, error)
	VerifyUser(phoneNumber string) error
	GetJWTPublicKey() *rsa.PublicKey
	GetJWTPrivateKey() *rsa.PrivateKey
}

type Store struct {
	db     *sql.DB
	config utils.Config
	pubsub *pubsub.Client
	jwtKey struct {
		public  *rsa.PublicKey
		private *rsa.PrivateKey
	}
}

// NewStore creates a new store with all dependencies like database, pubsub client etc
func NewStore(config utils.Config) Store {

	//initialise pub sub client
	absPath, err := filepath.Abs("/etc/flahmingo/key.json")
	if err != nil {
		logrus.Errorf("google key not found: %v", err)
	}
	opt := option.WithCredentialsFile(absPath)
	psClient, err := pubsub.NewClient(context.Background(), config.GoogleCloud.ProjectID, opt)
	if err != nil {
		panic(err)
	}

	privateKey := initJWTKeys()

	//create database connection
	db := createDBPool(config)
	return Store{
		db:     db,
		config: config,
		jwtKey: struct {
			public  *rsa.PublicKey
			private *rsa.PrivateKey
		}{public: &privateKey.PublicKey, private: privateKey},
		pubsub: psClient,
	}
}

func initJWTKeys() *rsa.PrivateKey {

	f, err := os.Open("/etc/flahmingo/jwt.key")
	if err != nil {
		if os.IsNotExist(err) {
			key, errGen := rsa.GenerateKey(rand.Reader, 2048)
			if errGen != nil {
				logrus.Fatalf("could not generate private key file: %v", errGen)
			}
			keyBytes := x509.MarshalPKCS1PrivateKey(key)

			err = ioutil.WriteFile("/etc/flahmingo/jwt.key", keyBytes, os.ModeType)
			if err != nil {
				logrus.Fatalf("could not write private key file: %v", err)
			}
			return key
		}
		logrus.Fatalf("could not open private key file: %v", err)
	}
	defer f.Close()
	keyBytes, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.Fatalf("could not read private key file: %v", err)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		logrus.Fatalf("could not parse private key file: %v", err)
	}

	return privateKey
}

// createDBPool creates the connection to postgres database
func createDBPool(config utils.Config) *sql.DB {
	var str string

	if config.Database.SSL {

		caCert, _ := filepath.Abs(config.Database.CaCertPath)
		userCert, _ := filepath.Abs(config.Database.UserCertPath)
		userKey, _ := filepath.Abs(config.Database.UserKeyPath)

		str = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=%s&sslcert=%s&sslkey=%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			caCert,
			userCert,
			userKey,
		)

	} else {
		str = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
	}

	db, err := sql.Open("postgres", str)
	if err != nil {
		panic(err.Error())
	}

	//Check if the connection is successful by establishing a connection.
	//Retry upto 10 times if connection is not successful
	for retryCount := 0; retryCount < 10; retryCount++ {
		err = db.Ping()
		if err == nil {
			logrus.Info("database connection successful")
			return db
		}

		logrus.Error(err)
		logrus.Info("could not connect to database: retrying...")
		time.Sleep(time.Second)
	}

	panic("could not connect to database")

}

// GetConfig returns config
func (s Store) GetConfig() utils.Config {
	return s.config
}

// GetJWTPrivateKey gets the private key used for generating JWT tokens
func (s Store) GetJWTPrivateKey() *rsa.PrivateKey {
	return s.jwtKey.private
}

// GetJWTPublicKey gets the private key used to verify JWT tokens
func (s Store) GetJWTPublicKey() *rsa.PublicKey {
	return s.jwtKey.public
}
