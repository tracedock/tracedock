# TraceDock

**TraceDock** is a highly extensible and minimalistic telemetry processor built to handle observability data based on user-defined rules.

It aims to be an alternative to OpenTelemetry Collector distributions and offers an alternative for teams who need **fine-grained control** through **simple but powerful configurations**, **plugin support**, and it all provided in a **well documented** platform.

## Why TraceDock?

Traditional observability pipelines are powerful but often:

- Overly complex to configure
- Are poor documented
- Hard to customize at runtime
- Limited in extensibility

**TraceDock** was created to solve these pain points with:

- Config-driven trace processing via YAML
- Native plugin support via `.so` files
- Simple CLI interface
- Well documented features

## Core Concepts

| Concept      | Description                                                           |
| ------------ | --------------------------------------------------------------------- |
| **Routers**  | Conditions that defines what **Pipeline** will process a span.        |
| **Pipeline** | Internal engine that applies **Rules** sequentially                   |
| **Rules**    | YAML-defined filters and actions that dictate how spans are processed |

A **Rule** can use a build-in provided feature in TraceDock or a feature provided by an external Plugin.

## Configuration example


```yaml
log:
  level: INFO

plugins:
  folders: [/etc/trackdock/plugins]

performance:
  memory_limiter:
    strategy: disk_dump
    max_consumption: 4096m

pipelines:
- name: main
  rules:

  # discard all the spans with redis' HGETALL operations lasting less than 100ms
  - provider: eraser
    match:
      attributes:
        db.system: redis
        db.statement: ^HGETALL.*
      duration:
        lt: 100ms

  # set attribute http.route with the same value of http.target when the span 
  # have the second one and its kind is server
  - provider: setter
    match:
      attributes:
        http.target: .*
        kind: server
    missing:
      attributes: [http.route]
    set:
      http.route: ${attributes.http.target}

  # set attribute name with the same value as http.target when the previous
  # name is "HTTP [method]"
  - provider: setter
    match:
      name: ^HTTP (GET|POST|PUT|DELETE|PATCH)$
    set:
      name: ${attributes.http.target}

  # send the spans to jaeger via otlp/grpc 
  - provider: export.otlp
    config:
      endpoint: https://jaeger.my.domain:4317
      protocol: grpc
      timeout: 1s
```
