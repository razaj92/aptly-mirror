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

## Docker

You can run the mirrors from a container too. To do so you need to provide GPG Keys for publishing or it will generate its own dummy ones. If you want to use your own you can provide them, along with the config files on the following paths:

```
docker run --rm \
  -v mykey.pub:/gpg_pub.gpg \
  -v mykey.key:/gpg_key.gpg  \
  -v aptly_mirrors.yml:/etc/aptly_mirrors.yml \
  -v aptly.conf:/etc/aptly.conf \
  razaj92/aptly-mirror run
```

You can also choose to mount a current aptly database as a volume to `/var/lib/aptly` if you want the data to persist across runs.
if you dont wish to provide exported gpg keys, you can also mount your .gnupg folder with the keyrings and pass through the location as the `GNUPGHOME` env var.
