Google Cloud Logs via PubSub to Loggly or Logz.io
=================================================

[![Circle CI](https://circleci.com/gh/rainchasers/logz.svg?style=svg)](https://circleci.com/gh/rainchasers/logz)

A daemon that consumes from a gcloud pubsub topic, batches the results and POSTs the results to a bulk upload endpoint (e.g. Loggly or Logz.io). Useful to transport logs from Google Stackdriver to an external service.
