---
canonical: https://grafana.com/docs/alloy/latest/reference/components/prometheus/prometheus.exporter.docker_state/
aliases:
  - ../prometheus.exporter.docker_state/ # /docs/alloy/latest/reference/components/prometheus.exporter.docker_state/
description: Learn about prometheus.exporter.docker_state
labels:
  stage: general-availability
  products:
    - oss
title: prometheus.exporter.docker_state
---

# `prometheus.exporter.docker_state`

The `prometheus.exporter.docker_state` component embeds a Docker state exporter for collecting container state metrics from the local Docker daemon.

{{< docs/shared lookup="reference/components/exporter-clustering-warning.md" source="alloy" version="<ALLOY_VERSION>" >}}

## Usage

```alloy
prometheus.exporter.docker_state "<LABEL>" {
}
```

## Arguments

You can use the following arguments with `prometheus.exporter.docker_state`:

| Name            | Type     | Description                                                   | Default                       | Required |
| --------------- | -------- | ------------------------------------------------------------- | ----------------------------- | -------- |
| `docker_host`   | `string` | The Docker daemon address.                                    | `"unix:///var/run/docker.sock"` | no       |
| `cache_period`  | `int`    | Cache duration for Docker inspect results, in seconds.        | `1`                           | no       |
| `enable_labels` | `bool`   | Whether to expose container labels as Prometheus labels.      | `true`                        | no       |

## Exported fields

{{< docs/shared lookup="reference/components/exporter-component-exports.md" source="alloy" version="<ALLOY_VERSION>" >}}

## Component health

`prometheus.exporter.docker_state` is only reported as unhealthy if given an invalid configuration.
In those cases, exported fields retain their last healthy values.

## Debug information

`prometheus.exporter.docker_state` doesn't expose any component-specific debug information.

## Debug metrics

`prometheus.exporter.docker_state` doesn't expose any component-specific debug metrics.

## Example

This example uses a [`prometheus.scrape`][scrape] component to collect metrics from `prometheus.exporter.docker_state`:

```alloy
prometheus.exporter.docker_state "example" {
  docker_host   = "unix:///var/run/docker.sock"
  cache_period  = 1
  enable_labels = true
}

// Configure a prometheus.scrape component to collect docker_state_exporter metrics.
prometheus.scrape "demo" {
  targets    = prometheus.exporter.docker_state.example.targets
  forward_to = [prometheus.remote_write.demo.receiver]
}

prometheus.remote_write "demo" {
  endpoint {
    url = "<PROMETHEUS_REMOTE_WRITE_URL>"

    basic_auth {
      username = "<USERNAME>"
      password = "<PASSWORD>"
    }
  }
}
```

Replace the following:

* _`<PROMETHEUS_REMOTE_WRITE_URL>`_: The URL of the Prometheus `remote_write` compatible server to send metrics to.
* _`<USERNAME>`_: The username to use for authentication to the `remote_write` API.
* _`<PASSWORD>`_: The password to use for authentication to the `remote_write` API.

[scrape]: ../prometheus.scrape/

<!-- START GENERATED COMPATIBLE COMPONENTS -->

## Compatible components

`prometheus.exporter.docker_state` has exports that can be consumed by the following components:

- Components that consume [Targets](../../../compatibility/#targets-consumers)

{{< admonition type="note" >}}
Connecting some components may not be sensible or components may require further configuration to make the connection work correctly.
Refer to the linked documentation for more details.
{{< /admonition >}}

<!-- END GENERATED COMPATIBLE COMPONENTS -->
