# [WIP] Docker Compose Example

This is an example of the instructions to run the gateway server on a remote debian
server using docker-compose from scratch.

Locally

```bash

# SSH into your server
ssh user@your-remote-server

# Install the necessary dependencies
sudo apt update
sudo apt install -y postgresql postgresql-contrib git docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# Make a workspace directory
mkdir workspace
cd workspace

# Add the following to ~/.profile
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Install go
wget https://go.dev/dl/go1.21.12.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.12.linux-amd64.tar.gz

# Clone the gateway server repository
git clone https://github.com/pokt-network/gateway-server.git
cd gateway-server/


# Uncomment the postgres related lines in the docker-compose.yml file
cp docker-compose.yml.sample docker-compose.yml

# Update the .env accordingly
cp .env.sample .env

# Start the gateway server
docker-compose up -d

# OPTIONAL: Connect to the postgres container directly
PGPASSWORD=mypassword psql -h localhost -U myuser -d postgres -p 5433
```
