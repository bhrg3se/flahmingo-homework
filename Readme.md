# Flahmingo Homework

## Setup and running

### Setup Google Cloud PubSub
- Create a project in google cloud.  
- Create a topic in PubSub named "verification".  
- Create a service account which has permission of publishing and subscribing pub sub.
- Download the Google Cloud key file for that account.
   
### Setup Auth-Service

#### Using `docker-compose`

- Go to docker/ directory
- Copy the google cloud key file as /etc/flahmingo/key.json
  > You can change the directory by changing the volume mapping in docker-compose.yml file.
- Run `docker-compose up`

#### Installing natively

- Create a config file in /etc/flahmingo/config.toml
- Use setup/config.toml as a starting point
- Copy the google cloud key file as /etc/flahmingo/key.json


   

