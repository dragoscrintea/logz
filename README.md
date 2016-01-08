Google Cloud Log Service via PubSub to Loggly
=============================================

[![Circle CI](https://circleci.com/gh/rainchasers/gcloud-pubsub-to-loggly.svg?style=svg)](https://circleci.com/gh/rainchasers/gcloud-pubsub-to-loggly)

[![Docker Repository on Quay](https://quay.io/repository/rainchasers/pubsub2loggly/status "Docker Repository on Quay")](https://quay.io/repository/rainchasers/pubsub2loggly)

A daemon that consumes from a gcloud pubsub topic, batches the results and POSTs the results to a Loggly bulk upload endpoint. Useful to transport logs from Google Cloud Log service, via PubSub to Loggly.

Log Format from PubSub
----------------------

Unstructured log:

```
{"insertId":"2016-01-07|13:36:58.968679-08|10.194.237.35|1677983050","log":"fluentd-cloud-logging","metadata":{"labels":{"compute.googleapis.com/resource_id":"7206755018967447229","compute.googleapis.com/resource_name":"fluentd-cloud-logging-gke-eud-b527d921-node-in2m","compute.googleapis.com/resource_type":"instance","container.googleapis.com/cluster_name":"eud","container.googleapis.com/container_name":"fluentd-cloud-logging","container.googleapis.com/instance_id":"7206755018967447229","container.googleapis.com/namespace_name":"kube-system","container.googleapis.com/pod_name":"fluentd-cloud-logging-gke-eud-b527d921-node-in2m","container.googleapis.com/stream":"stdout"},"projectId":"tuleyprod","serviceName":"container.googleapis.com","timestamp":"2016-01-07T21:36:56Z","zone":"europe-west1-d"},"textPayload":"2016-01-07 21:36:56 +0000 [warn]: emit transaction failed: error_class=Fluent::BufferQueueLimitError error=\"queue size exceeds limit\" tag=\"fluent.warn\"\n"}
```

Structured log:

```
{"insertId":"2016-01-07|19:52:54.873350-08|10.194.235.73|838957687","log":"com-textcaptcha","metadata":{"labels":{"compute.googleapis.com/resource_id":"9903089256222766859","compute.googleapis.com/resource_name":"fluentd-cloud-logging-gke-eud-b527d921-node-opq3","compute.googleapis.com/resource_type":"instance","container.googleapis.com/cluster_name":"eud","container.googleapis.com/container_name":"com-textcaptcha","container.googleapis.com/instance_id":"9903089256222766859","container.googleapis.com/namespace_name":"default","container.googleapis.com/pod_name":"textcaptcha-ot0qk","container.googleapis.com/stream":"stdout"},"projectId":"tuleyprod","serviceName":"container.googleapis.com","timestamp":"2016-01-08T03:52:49Z","zone":"europe-west1-d"},"structPayload":{"event":"runtime.gc","gc_per_second":0.033333350237786351,"pause_max_ms":2.202556,"pause_ms_per_second":0.073418533333333327,"pause_ms_total":164955.525249,"service":"textcaptcha","timestamp":"2016-01-08T03:52:49.661831825Z","type":"info"}}
```

Kubernetes
----------

For rainchasers, the container executes within a Kubernetes cluster. A sample copy of the replication controller spec is included in the file `replication-controller.yml` which you'll need to amend with your own env var values.

    kubectl create -f replication-controller.yml
    kubectl get rc
    kubectl get pods