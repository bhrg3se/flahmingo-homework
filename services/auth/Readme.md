# Auth Service

This service will handle authentication and user profile.

## Build and Run
- Go to services/auth folder
- Run `go build`
- Run `./auth` (you may need sudo)

## Endpoints

#### SignupWithPhoneNumber
Takes phone number and user's name as argument, creates user profile and sends OTP to verify phone number.

#### VerifyPhoneNumber
Takes OTP as argument and marks users as verified if OTP is correct

#### LoginWithPhoneNumber
Takes phone number as argument, sends OTP to login.

#### ValidatePhoneNumberLogin
Takes OTP as argument and returns a auth token if OTP is correct

#### GetProfile
Takes auth token and return user profile  based on that auth token if the token is valid
