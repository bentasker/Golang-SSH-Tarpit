# Use the official golang image as a base
FROM golang

# Upload the tarpit to the image
ADD . /go/

# Compile and install
RUN go build /go/ssh_tarpit.go

# Image should run the tarpit
ENTRYPOINT /go/ssh_tarpit

# We bind to 2222
EXPOSE 2222
