# mike-admission-controller

This is the code to a simple admission controller that looks to see if a container is deployed with the tag "latest". If so, it will reject the deployment.

## Usage

Build the container:

```
docker build . -t michaelbraunbass/mike-admission-controller:main
docker push michaelbraunbass/mike-admission-controller:main
```

Deploy to Kubernetes:

```
kubectl apply -f kubernetes/deploy-mike-admission-controller.yaml
```

Register the webhook:

```
kubectl apply -f kubernetes/register-webhook.yaml
```