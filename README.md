# Kube Scraper Backend

This is a small project and is part of
[Kube Scraper](https://github.com/SunSince90/kube-scraper), It
connects to a backend and stores useful information, such as *chats*
information through *gRPC*.

This is mostly done as a personal project, to improve my skills and learn some
new programming patterns and technologies.

If you'd like to include this in your project please be aware that I will no be
giving no warranty about its functionalities.
You are more than welcome to open issues, discussions, ask for help and make
pull requests if you'd like to.

## Kubernetes nature

The program is meant to run in Kubernetes, so that other pods - *scrapers* -
can connect to this program running as a pod and get the desired information.

### K3s & Raspberry Pi

As stated above, this project serves as a good exercise for me and it is
currently running in a cluster of *Raspberry Pis* running
[k3s](https://k3s.io/).

## Backend flexibility

Different backends can be used, but as of now, only firestore is implemented.
In future, when I will be able to test other backends I will add other
implementations.

### Why a backend pod

Instead of using the backend api directly from each scraper pod, this backend
pod is used. This is made to reduce call to the backend API as much as
possible, so that free quotas - i.e. Firestore quotas - are not exceeded.
This will make this project cheap to run on your *Raspberry Pi*.

## Get it

```bash
git clone https://github.com/SunSince90/kube-scraper-backend.git
cd kube-scraper-backend
```

## Build it

```bash
make build
```

## Example

```bash
./backend firestore \
--chats-collection chats
--project-id my-project-name
--service-account-path ./credentials/service-account.json
--address 10.23.55.77
--port 8787
```

Some description about the command above:

* `--chats-collection` is the name of the firestore collection that contains
all chats documents
* `--project-id` is the id of the firebase project
* `--service-account-path` is the path of the *JSON* service account
* `--address` specifies the address where we will serve requests
* `--port` specifies the port where we will serve requests

## Deploy on Kubernetes

### Build and push the image

Please note that the image that is going to be built will run on *ARM*, as it
will run on a Raspberry Pi. Make sure to edit the `Dockerfile` in case you want
to build for another architecture.

Login to your repository and

```bash
make docker-build docker-push IMG=<image>
```

### Create the namespace

```bash
kubectl create namespace kube-scraper
```

### Create the project id secret

Get the `project id` from your firebase console and run:

```bash
kubectl create secret generic firebase-project-id --from-literal=project-id=<your-project-id> -n kube-scraper
```

### Create the service account secret

Get the service account from your firebase console (or from gcp) and run:

```bash
kubectl create secret generic gcp-service-account --from-file=service-account.json=<path-to-your-service-account> -n kube-scraper
```

### Create the firestore chats collection

```bash
kubectl create configmap chats-config --from-literal=firestore.chats-collection=<chats-collection> -n kube-scraper
```

### Create the service

```bash
kubectl create deploy/service.yaml
```

### Create the deployment

```bash
kubectl create deploy/deployment.yaml
```
