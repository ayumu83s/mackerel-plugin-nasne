# mackerel-plugin-nasne
nasne custom metrics plugin for mackerel.io agent.

monitor the following items

- capacity of HDD
- recorded count
- recording failures count

## Build this program
```shell
glide install
go build
```

## Synopsis

```shell
mackerel-plugin-nasne 192.168.11.3
```

## Example of mackerel-agent.conf

```shell
[plugin.metrics.nasne]
command = "/path/to/mackerel-plugin-nasne 192.168.11.3"
```

## Screenshot
![graphs-screenshot](https://user-images.githubusercontent.com/1732800/28264577-83b35940-6b26-11e7-8ac2-f11642672800.png)
