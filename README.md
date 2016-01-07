Google Cloud Log Service via PubSub to Loggly
=============================================

[![Circle CI](https://circleci.com/gh/rainchasers/gcloud-pubsub-to-loggly.svg?style=svg)](https://circleci.com/gh/rainchasers/gcloud-pubsub-to-loggly)

[![Docker Repository on Quay](https://quay.io/repository/rainchasers/pubsub2loggly/status "Docker Repository on Quay")](https://quay.io/repository/rainchasers/pubsub2loggly)

A daemon that consumes from a gcloud pubsub topic, batches the results and POSTs the results to a Loggly bulk upload endpoint. Useful to transport logs from Google Cloud Log service, via PubSub to Loggly.