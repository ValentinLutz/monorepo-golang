# monorepo

## local kubernetes with minikube

````shell
minikube delete
````

````shell
minikube start 
````

````shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
kubectl create namespace prometheus
helm install my-kube-prometheus-stack prometheus-community/kube-prometheus-stack --version 43.1.1 --namespace prometheus
````