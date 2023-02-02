# Docker Compose

### In this project, Docker Compose was used to run an application (HTML, CSS and JS) in an Apache container, and another multistage image (A binary of a GO application, combined with an Alpine Linux container). Also available are the Dockerfiles used in the two images, and the required application files. For Docker Compose the images were made available via Docker Hub.

.
.
.

## **Why was it made?**

.
.

To master Docker, it's good to practice the commands needed to make a dockerfile and a docker compose. The multistage image is important for knowing how to include a binary in a docker image. So it's a practical example that can be changed by changing applications or linux containers. It is very good for study purposes.

## **Where to start?**

The Dockerfiles are the first step. To start you can see the Valida Dockerfile.

### Valida > Dockerfile

...

**FROM debian** - To set our Linux Image

**RUN apt-get update && apt-get install -y apache2 && apt-get clean** - Docker recommends using apt-get over apt as it says it brings more information. For good practices it is recommended to put everything in just one RUN. The clean command is used to delete the remains of files that were used in the installation.

**ENV APACHE_LOCK_DIR="var/lock"** - To avoid having more than one execution of apache in the same container.

**ENV APACHE_PID_FILE="var/run/apache2.pid"** - Need to specify where this PID type file will be.

**ENV APACHE_RUN_USER="www-data"** - www-data will be the user that will run apache, we can put any other user. It is not advisable to use the root.

**ENV APACHE_RUN_GROUP="www-data"**

**ENV APACHE_LOG_DIR="/var/log/apache2"** - Specify log directory.

**ADD valida.tar /var/www/html** - Default folder, where the file will be copied and already unzipped.

**LABEL description = "Apache webserver 1.0"** - To specify the description.

**VOLUME /var/www/html** To specify where data will be saved.

**EXPOSE 80** To expose port 80.

**ENTRYPOINT ["/usr/sbin/apachectl"]** To specify the run file.

**CMD ["-D", "FOREGROUND"]** To specify that the run needs to be in the foreground.
.
.

After the Dockerfile is ready, we can build the image:

**docker image build -t imagename:1.0 .** Where 1.0 is the image version tag. And don't forget the "." to use the dockerfile from the current directory.

**docker run -dti -p 80:80 --name my-apache imagename:1.0** Where we set port 80, as we had put in the image. So we put this apache container to run with the site files inside it.

.
.
.

## **Generating a Multistage Image**

Let's check out an example where we'll make an application in GO, and we'll place its binary next to the Linux Alpine image. From the Dockerfile we will use the GO image to generate the application binary:

### Go > Dockerfile

...

**FROM golang as exec** Golang is the name of the image, and since we are going to use the result of the binary that will be inside the FROM of this image, we will call it executable to be able to import it into the next stage.

**COPY main.go /go/src/app/** To copy the app file to the indicated address.

**ENV GO111MODULE=auto** Informs that we will be able to generate the executable from any place inside the container. Without this command, an error will appear saying that go.mod was not found.

**WORKDIR /go/src/app/** To inform the directory we will work with.

**RUN go build -o app.go .** Run will check app.go to generate the executable file, this change in the image needs to come from Run, because it executes an instruction and sends the return of the instruction to the first layer of the Dockerfile, which is where we can make changes (such as generating the binary).

--- The top part was the first stage, now comes the second stage inside the Dockerfile ---

**FROM alpine**

**WORKDIR /appexec** We create a working directory, it could be another name.

**COPY --from=exec /go/src/app /appexec** Here we copy the exec from the first directory /go/src/app to the second stage directory /appexec.

**RUN chmod -R 755 /appexec** As we are going to make changes to the file, in this case changes to the directory permissions, so we need to use Run again.

**EXPOSE 8080**

**ENTRYPOINT ./app.go**
.
.

We close the dockerfile.

**docker image build -t app-go:1.0 .** To generate the image.

**docker run -ti --name myappOk app-go:1.0** To run the image.
