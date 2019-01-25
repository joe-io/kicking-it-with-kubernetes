# Containerize your app
Before we can deploy anything to kubernetes we need to containerize the different microservices. We will be using Docker as the container platform. If you have not [installed Docker](https://docs.docker.com/install/) on your system you will need to do so now to continue this lab. [[Install Docker]](https://docs.docker.com/install/)

To maintain CI/CD principles we will use Docker to create a clean room for building your application and to create the image which k8s will run. All of this will happen in a single Dockerfile.

In your `api` directory create a file named `Dockerfile`. In this file add the following lines.

```docker
FROM golang:1.11-alpine as build
RUN apk add --no-cache --update alpine-sdk git gcc
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=build /build/api /
CMD ["/api"]
```

Let's look at each line in detail. `FROM golang:1.11-alpine as build` states that at the base level our image will be an alpine linux image with golang 1.11 preinstalled. You can browse around [Docker Hub](https://hub.docker.com/) and find images of all shapes and sizes. The `as build` part just gives a name to this image. The next line installs a few alpine packages that are required to build our go application.

`COPY . /build` copies all of the contents of the api directory into the /build directory inside the container. We then change the current working directory inside the image to /build.

The last line of the first section actually builds your application inside the `build` container. This isolates your build from your own machine making the build portable. You could build using just `go build` but that would dynamically link your app to a specific set of underlying libraries. The more complex build line specified in the Dockerfile statically links all required libraries so there is a single stand-alone binary.

The second section starts with `FROM scratch`. This starts the creation of a second image. The `build` image will ultimately be thrown away after this second image is complete. `scratch` means that the image will be a base linux vm with no OS. We can do this because we built a static binary with no dependencies. We then copy the go application that we built from the `build` image to this second image and place it at the root of the file system.

The `CMD ["/api"]` statement is the command that will be run inside the container. At this point the `build` image is cleaned up and we have a single image with our application available to run.

We have a configuration file that defines what image we would like created. We need to have docker actually build the image now. Run the following from the `api` directory where the `Dockerfile` is.

```sh
docker build -t myapi .
```

This command will create a docker image with the name `myapi`. To create a running container from the image, run the following.

```sh
docker run -it -p 8082:8082 myapi
```

This will run a container locally on your machine. You should see output similar to the following.

```logs
[GIN-debug] POST   /recognizer/training-image --> main.trainImage (3 handlers)
[GIN-debug] GET    /recognizer/identification --> main.identifyImage (3 handlers)
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8082
```

If you hit http://localhost:8082/recognizer/identification from a browser you will get a 500 error, but you should see a new log line reporting the error. This means that your service is up and running.

The next exercise is to take what you learned here and create a Dockerfile for the `model` service.

# Write deployment configuration
You should now have **2 Dockerfiles** that create a valid image that can be run locally. If you do not have both Docker files go back to the last exercise and create them now.

The latest Docker installs allow you to enable k8s locally. Open the docker preferences and make sure that k8s is enabled. [[Mac Instructions]](https://docs.docker.com/docker-for-mac/#kubernetes) [[Windows Instructions]](https://docs.docker.com/docker-for-windows/#kubernetes)

We will now define a k8s node and replica set in context of a deployment controller. Create a new file called deployment.yaml in the `api` directory and add the following configuration.

```yaml
apiVersion: apps/v1                    # The version of the k8s api
kind: Deployment                       # Specify what we are configuring
metadata:
  name: api-deployment                 # A label for this config
spec:
  replicas: 2                          # This defines the replica set
  selector:
    matchLabels:
      app: api-pod                     # This value must match the pod label below
  template:
    metadata:
      labels:
        app: api-pod                   # All pods have the key/value pair app:myapi added to their list of labels
    spec:
      containers:
      - name: myapi                    # The name assigned to the container in the docker daemon
        image: myapi:latest            # Using the 'latest' tag is a bad practice but easy for demos
        ports:
        - containerPort: 8082          # This container port will be exposed
        env:
        - name: API_MODEL_ENDPOINT     # This key/value pair will be available in the containers environment
          value: http://model-service/

```

Run the following command to deploy your app to k8s.
```sh
kubectl create -f deployment.yaml
```

To see if the deployment was successful run
```sh
kubectl get all
```

You shoud see output similar to the following.
```
NAME                                  READY   STATUS    RESTARTS   AGE
pod/api-deployment-57449d868c-4rcmg   1/1     Running   0          40s
pod/api-deployment-57449d868c-mlg2z   1/1     Running   0          38s

NAME                             DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/api-deployment   2         2         2            2           1m

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/api-deployment-57449d868c   2         2         2       41s
```
This shows the deployment config, the replica set config, and the 2 pods. These are all of the things requested in the deployment.yaml configuration.

Take a few minutes and build out a deployment.yaml file the `model` app and then deploy it to k8s.

If you make any mistakes you can just fix the deployment.yaml file and run
```sh
kubectl apply -f deployment.yaml
```
This will update the configuration in k8s. K8s will then adjust your deployment to match the new configuration.

# Write service configuration

# Write ingress configuration

