# mackerel-plugin-apt-check

This is a custom metrics plugin for [mackerel.io](https://mackerel.io/) agent, which runs a helper script included in `update-notifier-common` package and reports the number of available package updates.

## Synopsis

```shell
mackerel-plugin-apt-check [-script=<path-to-apt-check-script>]
                          [-metric-key-prefix=<prefix>]
```

## Example of mackerel-agent.conf

```toml
[plugin.metrics.apt-check]
command = "/path/to/mackerel-plugin-apt-check"
```

## License

MIT License.
