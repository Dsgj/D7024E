#Copy the image from Lara
FROM larjim/kademlialab:latest

# Copy the current directory contents into the container at /app
COPY . /home/go/src/D7024E

#Sets the working directory to where to code is
WORKDIR /home/go/src/D7024E

#Compiles the file "client.go" and names it "clientBin"
#OBS!! This is done when the image it built!!
#RUN /usr/local/go/bin/go build -o clientBin client.go

#Run the composed file
#OBS!! This is done when the container is running!!
CMD /usr/local/go/bin/go run ./app/main.go

# Make port 80 available to the world outside this container
#EXPOSE 80


#För att bygga in ny image kör man:
#
# docker build -t testimage:latest .
#
# docker swarm init
#
# docker stack deploy -c docker-compose-lab.yml
#
