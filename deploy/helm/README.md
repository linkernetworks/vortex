This folder contains the charts that you can deploy them by [helm](https://helm.sh).


### Install helm

- ubunut
```
curl -L https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz > helm-v2.9.1-linux-amd64.tar.gz
tar -zxvf helm-v2.9.1-linux-amd64.tar.gz
chmod +x linux-amd64/helm
sudo mv linux-amd64/helm /usr/local/bin/helm
```
- macOS
```
brew install kubernetes-helm
```

### Initial helm
Helm has two parts: a client (helm) and a server (tiller), Tiller runs inside of your Kubernetes cluster as the deployment, and manages releases (installations) of your charts. Helm will figure out where to install Tiller by reading your Kubernetes configuration file (usually $HOME/.kube/config). This is the same file that kubectl uses.
```
make apps.init-helm
```

### Using helm deploy apps
This will deploy mongodb and prometheus in your cluster
```
make apps.launch
```
If you wnat to deploy certain chart, you can type
```
helm install --debug --wait --set global.environment=<environmtneSetting> <chart path>

#example
helm install --debug --wait --set global.environment=local deploy/helm/apps/prometheus/charts/cadvisor
```

### Delete all release
```
make apps.teardown
``` 
