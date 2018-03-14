# aptly-mirror

A simple tool to mirror a list of repos using `aptly` and publishing them to S3

## Dependencies

You need `aptly` along with your chosen imported `gpg` keys, you can `export GPUPGHOME='/dir/gpg'` if the keyfile is stored in a non default location. Your Aptly config must also inclue the [S3PublishEndpoints](https://www.aptly.info/doc/configuration/) if you wish to publish your mirrors to S3.

## Install

`go get github.com/razaj92/aptly-mirror`

## Usage

After setting up the config file as descibed in `aptly_mirrors.yml` you can run `aptly-mirror run` to create a local mirror of the Repositories. `aptly-mirror run -p` will mirror the repo and publish to S3

```
Usage:
  aptly_mirror [command]

Available Commands:
  help        Help about any command
  run         Mirrors repos from config file

Flags:
      --aptly-path string    Path to Aptly Binary (default "aptly")
      --config-path string   Path to ./aptly_mirrors.yaml file containing config (default "/etc")
  -d, --debug                Debug mode: prints command outputs..
  -h, --help                 help for aptly_mirror
```

## Configuration

The config yaml file can takes the following attributes

```
aptly:
  endpoint: s3:[aptly s3 publish endpoint name]

gpg:
  key: [your local gpg key id]
  servers:
    - [list of gpg key servers you want to use]

repos:
  - name: [repo name]
    url: [repo url eg. https://download.app.com/linux/apt]
    release: [release name eg. xenial, trusty etc]
    components: [component name eg. main]
    arch: [supported architectures]
    gpgkeys: [repo public gpg id]
```
