[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/DoctorOgg/sensu-zcard)
![goreleaser](https://github.com/DoctorOgg/sensu-zcard/workflows/goreleaser/badge.svg)

# Sensu Redis Zcard Metrics Plugin

## Overview

This plugin will scan redis for keys matching a pattern and return the cardinality (zcard) of the set in graphite format.

## Table of Contents

- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)

## Files

- sensu-zcard

## Usage examples

```bash
./sensu-zcard
fubar:farts:43 3 1673374962
fubar:farts:222 6 1673374962
```

Help:

```bash
$ ./sensu-zcard -h
Sensu Check to return zcard metrics

Usage:
  sensu-zcard [flags]
  sensu-zcard [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -d, --database int      The Redis db to connect to
  -h, --help              help for sensu-zcard
  -i, --host string       The Redis host to connect to (default "localhost")
  -k, --key string        The Redis key to report on (default "fubar:farts*")
  -w, --password string   The Redis password
  -p, --port int          The Redis port to connect to (default 6379)

Use "sensu-zcard [command] --help" for more information about a command.

```

## Configuration

### Asset registration

Sensu Assets are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add DoctorOgg/sensu-zcard
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/DoctorOgg/sensu-zcard].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-zcard
  namespace: default
spec:
  command: sensu-zcard -h localhost -p 6379 -k "sensu:zcard" 
  subscriptions:
  - system
  runtime_assets:
  - DoctorOgg/sensu-zcard
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-zcard repository:

```
go build
```
