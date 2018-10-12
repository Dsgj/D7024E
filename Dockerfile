#Copy the image from Lara
FROM larjim/kademlialab:latest

# Copy the current directory contents into the container at /home/go/src/D7024E
COPY . /home/go/src/D7024E

#Copy shit into the container (go get doesnt work properly)
COPY /import/protobuf-master/proto /home/go/src/github.com/golang/protobuf/proto
COPY /import/urfave /home/go/src/github.com/urfave
COPY /import/takama /home/go/src/github.com/takama
COPY /import/go-chi /home/go/src/github.com/go-chi
COPY /import/resty.v1 /home/go/src/gopkg.in/resty.v1

ENV GOPATH=/home/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN

WORKDIR /home/go/src/D7024E/server
RUN /usr/local/go/bin/go build server.go
RUN /usr/local/go/bin/go install
WORKDIR /home/go/src/D7024E/client
RUN /usr/local/go/bin/go build client.go
RUN /usr/local/go/bin/go install

#Starts the server and passes the ip of the container as an argument
WORKDIR /home/go/bin
CMD server $(hostname -i)
