FROM larjim/kademlialab:latest

# Copy the current directory contents into the container at /app
COPY . /home/go/src/D7024E

#RUN echo "hello world"

# Make port 80 available to the world outside this container
#EXPOSE 

# Run app.py when the container launches
#CMD ["python", "app.py"]


#För att bygga in ny image kör man:
#
# docker build -t testimage:latest .
#
# docker swarm init
#
# docker stack deploy -c docker-compose-lab.yml
#

