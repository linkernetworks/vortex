# Monitoring

## Prometheus
Prometheus is a systems and service monitoring system. It collects metrics from configured targets at given intervals, evaluates rule expressions, displays the results, and can trigger alerts if some condition is observed to be true.

### Create a monitoring namespace

``` bash
kubectl create -f monitoring-namespace.yaml
```

### ConfigMap for Prometheus
It collects metrics from ```kubernetes-apiservers```, ```kubernetes-cadvisor```, ```kubernetes-nodes```, ```kubernetes-pods```(every node-exporter), ```kubernetes-service-endpoints``` by default.

``` bash
kubectl create -f prometheus-config.yaml
```

### Create the deployment and service

``` bash
kubectl create -f prometheus-deployment.yaml
kubectl create -f prometheus-svc.yaml
```

### Use the Prometheus
Open the browser and connect to ```hostIP:30003```, you can see whether the Prometheus collect the metrics from targets correctly in ```Status->Targets```.

You also can enter a PromQL expression to query the data and show them in graph.

``` SQL
sum by(pod_name)(container_memory_usage_bytes{namespace="kube-system"})
```

Or pass the expression through url and get the json.
``` bash
http://hostIP:30003/api/v1/query?query=sum%20by(pod_name)(container_memory_usage_bytes%7Bnamespace%3D%22kube-system%22%7D)
```


## Node Exporter

Node exporter for machine metrics, written in Go with pluggable metric collectors.

### Create the node-exporter daemonset

``` bash
kubectl create -f node-exporter.yaml
```

If you want to check whether the daemonset is created successfully, you can create a service to expose the node-exporter.

``` bash
kubectl create -f node-exporter-svc.yaml
``` 

Open the browser and connect to ```hostIP:30910/metrics```, you can see the metrics of the host node (service will route the request to the pod in round robin order).