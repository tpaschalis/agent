---
aliases:
- /docs/agent/latest/flow/reference/components/loki.source.syslog
title: loki.source.syslog
---

# loki.source.syslog

`loki.source.syslog` listens for syslog messages over TCP or UDP connections
and forwards them to other `loki.*` components. The messages must be compliant
with the [RFC5424](https://www.rfc-editor.org/rfc/rfc5424) format.

The component starts a new syslog listener for each of the given `config`
blocks and fans out incoming entries to the list of receivers in `forward_to`.

Multiple `loki.source.syslog` components can be specified by giving them
different labels.

## Usage

```river
loki.source.syslog "LABEL" {
  listener {
    address = "LISTEN_ADDRESS"
  }
  ...

  forward_to = RECEIVER_LIST
}
```

## Arguments

`loki.source.syslog` supports the following arguments:

Name         | Type                   | Description          | Default | Required
------------ | ---------------------- | -------------------- | ------- | --------
`forward_to` | `list(LogsReceiver)`   | List of receivers to send log entries to. | | yes

## Blocks

The following blocks are supported inside the definition of
`loki.source.syslog`:

Hierarchy | Name | Description | Required
--------- | ---- | ----------- | --------
listener | [config][] | Configures a listener for IETF Syslog (RFC5424) messages. | no
listener > tls_config | [tls_config][] | Configures TLS settings for connecting to the endpoint for TCP connections. | no

The `>` symbol indicates deeper levels of nesting. For example, `config > tls_config` 
refers to a `tls_config` block defined inside a `config` block.

[listener]: #listener-block
[tls_config]: #tls_config-block

### listener block

The `listener` block defines the listen address and protocol where the listener
will expect syslog messages to, as well as its behavior when receiving
messages.

The following arguments can be used to configure a `listener`. Only the
`listen_address` field is required and any omitted fields take their default
values.

Name                     | Type          | Description | Default | Required
------------------------ | ------------- | ----------- | ------- | --------
`address`                | `string`      | The `<host:port>` address to listen on for syslog messages. | | yes
`protocol`               | `string`      | The protocol to listen on for syslog messages. Must be either `tcp` or `udp`. | `tcp` | no
`idle_timeout`           | `duration`    | The idle timeout for tcp connections. | `"120s"` | no
`label_structured_data`  | `bool`        | Whether to translate syslog structured data to loki labels. | `false` | no
`labels`                 | `map(string)` | The labels to associate with each received syslog record. | `{}` | no
`use_incoming_timestamp` | `bool`        | Whether to set the timestamp to the incoming syslog record timestamp. | `false` | no
`use_rfc5424_message`    | `bool`        | Whether the full RFC5424-formatted syslog message should be forwarded. | `false` | no
`max_message_length`     | `int`         | The maximum limit to the length of syslog messages. | `8192` | no

By default, the component will assign the log entry timestamp as the time it
was processed.

The `labels` map will be applied to every message that the component reads.

All header fields from the parsed RFC5424 messages will be brought in as
internal labels, prefixed with `__syslog_`.

If `label_structured_data` is set, structured data in the syslog header will
also be translated to internal labels in the form of
`__syslog_message_sd_<ID>_<KEY>`. A structured data entry of `[example@99999
test="yes"]` would become the label `__syslog_message_sd_example_99999_test`
with the value `"yes"`.

### tls_config block

{{< docs/shared lookup="flow/reference/components/tls-config-block.md" source="agent" >}}

## Exported fields

`loki.source.syslog` does not export any fields.

## Component health

`loki.source.syslog` is only reported as unhealthy if given an invalid
configuration.

## Debug information

`loki.source.syslog` exposes some debug information per syslog listener:
* Whether the listener is currently running.
* The listen address.
* The labels that the listener applies to incoming log entries.

## Debug metrics
* `loki_source_syslog_entries_total` (counter): Total number of successful entries sent to the syslog component.
* `loki_source_syslog_parsing_errors_total` (counter): Total number of parsing errors while receiving syslog messages.
* `loki_source_syslog_empty_messages_total` (counter): Total number of empty messages received from the syslog component.

## Example

This example collects listens for Syslog messages in valid RFC5424 format over
TCP and UDP in the specified ports and forwards them to a `loki.write`
component.

```river
loki.source.syslog "local" {
  listener {
    address  = "127.0.0.1:51893"
    labels   = { component = "loki.source.syslog", protocol = "tcp" } 
  }

  listener {
    address  = "127.0.0.1:51898"
    protocol = "udp"
    labels   = { component = "loki.source.syslog", protocol = "udp"} 
  }

  forward_to = [loki.write.local.receiver]
}

loki.write "local" {
  endpoint {
    url = "loki:3100/api/v1/push"
  }
}
```
