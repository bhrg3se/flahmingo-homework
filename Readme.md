# Flahmingo Homework

## Setup and running

### 1. Setup Google Cloud PubSub
- Create a project in google cloud.  
- Create a topic in PubSub named "verification".  
- Create a service account which has permission of publishing and subscribing pub sub.
- Download the key file for that account.
   
### 2. Setup Twilio Account
Sign up for twilio and get account SID, auth token and phone number
### 3. Setup Microservices

#### 3.1 Using `docker-compose`

- Go to setup/ directory
- Edit the config file, fill up googleCloud and twillio fields
- Copy the google cloud key file to setup/key.json
  > You can change the directory by changing the volume source in docker-compose.yml file.
- Run `docker-compose up`

#### 3.2 Installing natively

- Make sure you have go:1.15+ installed
- Create a config file in /etc/flahmingo/config.toml
- Use setup/config.toml as a starting point
- Copy the google cloud key file to /etc/flahmingo/key.json
- Start a postgres database and configure the host,name,user and password in /etc/flahmingo/config.toml
- Run setup/init.sql in postgres
- Go to services/auth. Run `go build && ./auth`
- Go to services/otp. Run `go build && ./otp`
> You may run into permission issues because it will try to create private key file and log file.
> You may just run the binary with sudo



## File Structure

- utils (utility functions)
- setup (docker-compose, config samples and initial sql file)
- services
  - auth
    - proto (protobuf definitions)
    - pb (protobuf generated files)
    - server (gRPC server and APIS)
      - apis.go (gRPC API handlers)
      - jwt.go (auth token generation and verification)
      - server.go 
    - store (database and other dependencies)
      - init.go (initialization of database and other dependencies)
      - db.go (database functions)
      - pubsub.go (pubsub functions)
      - mock.go (mock store for testing)
  - otp


   

