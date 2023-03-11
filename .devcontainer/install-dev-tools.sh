# Install oapi-codegen
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4

# Install flyway
wget -qO- https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/9.8.1/flyway-commandline-9.8.1-linux-x64.tar.gz \
    | tar xvz -C ~
sudo ln -s ~/flyway-9.8.1/flyway /usr/local/bin 

# Install k6
wget -qO- https://github.com/grafana/k6/releases/download/v0.43.1/k6-v0.43.1-linux-amd64.tar.gz \
    | sudo tar xvz -C /usr/local/bin --strip-components=1 k6-v0.43.1-linux-amd64/k6 
