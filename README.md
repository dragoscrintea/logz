Google Cloud Log Service via PubSub to Loggly
=============================================

[![Circle CI](https://circleci.com/gh/rainchasers/gcloud-pubsub-to-loggly.svg?style=svg)](https://circleci.com/gh/rainchasers/gcloud-pubsub-to-loggly)

[![Docker Repository on Quay](https://quay.io/repository/rainchasers/pubsub2loggly/status "Docker Repository on Quay")](https://quay.io/repository/rainchasers/pubsub2loggly)

A daemon that consumes from a gcloud pubsub topic, batches the results and POSTs the results to a Loggly bulk upload endpoint. Useful to transport logs from Google Cloud Log service, via PubSub to Loggly.

Kubernetes
----------

For rainchasers, the container executes within a Kubernetes cluster. A sample copy of the replication controller spec is included in the file `replication-controller.yml` which you'll need to amend with your own env var values.

    kubectl create -f replication-controller.yml
    kubectl get rc
    kubectl get pods