# Containerize your app
Before we can deploy anything to kubernetes we need to containerize the different microservices. We will be using Docker as the container platform. If you have not [installed Docker](https://docs.docker.com/install/) on your system you will need to do so now to continue this lab. [[Install Docker]](https://docs.docker.com/install/)

To maintain CI/CD principles we will use Docker to create a clean room for building your application and to create the image which k8s will run. All of this will happen in a single Dockerfile.

In your `api` directory create a file named `Dockerfile`. In this file add the following lines.

```dockerfile
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
docker build -t myapi:v1 .
```

This command will create a docker image with the name `myapi`. To create a running container from the image, run the following.

```sh
docker run -it -p 8082:8082 myapi:v1
```

This will run a container locally on your machine. You should see output similar to the following.

```logs
[GIN-debug] POST   /social-post              --> main.main.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8082
```

Let's hit the service:
```sh
curl \
 -H "Accept: application/json" \
 -H "Content-type: application/json" \
 -X POST \
 -d '{"title":"Some Post", "body":"Some Content", "imageUrl": "https://qvxfxxo9ak-flywheel.netdna-ssl.com/wp-content/uploads/2018/03/Jasper-canoe-tour-at-Pyramid-Lake.jpg"}' \
 http://localhost:8082/social-post 
```

You will get a 500 error, but you should see a new log line reporting the error. This means that your service is up and running.

The next exercise is to take what you learned here and create a Dockerfile for the `analyzer` service.

As this part of the tutorial is focused on running services in Kubernetes, for simplicity, we are going to be running creating the Dockerfile in the analyar directory, not the analyzer-tf directory.

# Write deployment configuration
You should now have **2 Dockerfiles** that create valid images that can be run locally. If you do not have both Docker files go back to the last exercise and create them now.

The latest Docker install allows you to enable k8s locally. Open the docker preferences and make sure that k8s is enabled. [[Mac Instructions]](https://docs.docker.com/docker-for-mac/#kubernetes) [[Windows Instructions]](https://docs.docker.com/docker-for-windows/#kubernetes)

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
        image: myapi:v1                # Using the 'latest' tag is a bad practice but easy for demos
        ports:
        - containerPort: 8082          # This container port will be exposed
        env:                           # These key/value pairs will be available in the containers environment
        - name: PORT
          value: "8082"
        - name: API_ANALYZER_ENDPOINT     
          value: http://analyzer-service:8080/

```

Run the following command to deploy your app to k8s.
```sh
kubectl apply -f deployment.yaml
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

Take a few minutes and build out a deployment.yaml file the `analyzer` app and then deploy it to k8s.

If you make any mistakes you can just fix the deployment.yaml file and run
```sh
kubectl apply -f deployment.yaml
```
This will update the configuration in k8s. K8s will then adjust your deployment to match the new configuration.

# Write service configuration
Your applications are now running in k8s. Each container is exposed with its own IP address on the ports you have specified. There are guarantees on how many containers will stay running, but there are no guarantees that the IP address will stay the same if something restarts. To build out a guarantee you will create a service policy. This will give your containers a virtual name and IP address that can always be used to access whichever pods are available. The set of pods targeted by a service should be replicas of the same container. To implement a k8s service, create a file called service.yaml in the api directory. Add the following configuration to the file.
```yaml
apiVersion: v1            # The version of the k8s api
kind: Service             # Specify what we are configuring
metadata:
  name: api-service       # The name of this service
spec:
  ports:
  - port: 80              # This is the port the service will expose
    targetPort: 8082      # This is the port the container exposes
    protocol: TCP
    name: http
  selector:
    app: api-pod          # This service will front pods that match the selector `app: api-pod`
```

To deploy this configuration run the following.

```sh
kubectl apply -f service.yaml
```

Your service is now configured. To see the service run:
```sh
kubectl get services
```

This should return something like the following.
```
NAME            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
api-service     ClusterIP   10.104.154.176   <none>        80/TCP         2d
kubernetes      ClusterIP   10.96.0.1        <none>        443/TCP        81d
```

Notice the `PORT` definition for the `api-service`. The service is now available on port 80, but it is only available inside the cluster. This is great because in a microservice world most microservices are not customer facing. All microservices can communicate with each other as though they were public, but they have the security of only being available on the cluster. 

Notice also the `NAME` and `CLUSTER-IP` values for your services. The `NAME` value can be thought of as a hostname or cname for the service. The `NAME` value will resolve to the `CLUSTER-IP` address. Remember in the api deployment that we set the environment variable `API_ANALYZER_ENDPOINT` to `http://analyzer-service:8080/`? This was so the api would know where the analyzer microservice could be found. If the pods restart and move to different IP address there will be no need to make a configuration change. The service will maintain a static name and address. 

Now create a second `service.yaml` file in the analyzer directory for that service. Expose port 8080 to the cluster for this service and give it the name `analyzer-service`.

# Write ingress configuration
Generally whatever cluster you are using will have a couple of ingress controllers configured for you, f.e. AKS uses ALBs and Ethos uses Contour. In our case the local docker instance of k8s does not have an ingress controller configured. Our first step will be to install and nginx ingress controller. (If you are using a production grade cluster you will not need this step.)

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml
```

The good people who wrote the ingress-nginx controller also supplied the config needed to install it locally. Running the above two commands installed and configured k8s for ingress. If you run `kubectl get all --namespace ingress-nginx` you will see that the above commands configured a Pod, a Deployment, a ReplicaSet, and a Service in it's own namespace. This started an nginx container and configured k8s to accept external requests (from localhost).

Now that an ingress controller is in place, we can tie the ingress service and our api service together allowing our internal api service to accept requests from external clients. 

Create a new file in the api directory called `ingress.yaml` and add the following configuration to it.

```yaml
apiVersion: extensions/v1beta1         # The version of the k8s api
kind: Ingress                          # Specify what we are configuring
metadata:
  name: api-ingress                    # The name of this ingress object
spec:
  rules:
  - host: localhost                    # The HOST header value used for routing
    http:
      paths:
      - path: /                        # The PATH header value to match for routing
        backend:
          serviceName: api-service     # The name of the service to route to
          servicePort: 80              # The port of the service to route to
```

Send the above configuration to the k8s cluster in the usual way.

```sh
kubectl apply -f ingress.yaml
```

At this point you should be able to hit your service from your browser at `http://localhost/recognizer/identification?url=%E2%80%9Csomefakeurl.jpeg`

# Retrieving logs
You can directly get log output from your containers to help you debug issues. The best practice for logging in a containerized environment is to log to standard out and let the orchestrator handle the logs for you. This is by default what is happening in our setup. If you run `kubectl get pods` and select the name of one of the pods and then run `kubectl logs -f <podname>` you will be able to follow the logs for the specified container. While debugging it is often convenient to decrease the number of running containers to one.
