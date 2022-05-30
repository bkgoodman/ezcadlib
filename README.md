# Stuff

This just makes a simple CGI executable

# Initialize like:

`sqlite3 params.db < schema.sq3`

# Build like:

`CGO_ENABLED=0 go build`

# Test Like:

`REQUEST_METHOD=GET SERVER_PROTOCOL=HTTP/1.1 ./ezcaddb`
