# mackerel-plugin-nasne
nasne custom metrics plugin for mackerel.io agent.

monitor the following items

- capacity of HDD
- recorded count
- recording failures count

## How to install a plugin
It is easy to use `mkr plugin install`

```shell
mkr plugin install ayumu83s/mackerel-plugin-nasne
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
![graphs-screenshot](https://user-images.githubusercontent.com/1732800/28486860-0ff281f8-6ec2-11e7-9b6f-c5854908adc2.png)
