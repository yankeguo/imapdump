# imapdump

![workflow badge](https://github.com/yankeguo/imapdump/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/yankeguo/imapdump.svg)](https://pkg.go.dev/github.com/yankeguo/imapdump)

A tool for dumping emails to local in batch.

[中文文档](README.zh.md)

## Usage

**Command**

```
./imapdump -conf config.yaml
```

**Configuration**

```yaml
dir: dump
accounts:
  # name, name of the subdirectory
  - name: username@mydomain.com
    # host, host and port for IMAP server, must be TLS
    host: imap.mydomain.com:993
    # username, username of login
    username: username@mydomain.com
    # password, password of login
    password: xxxxxxxxxxxxxxxxxxx
    # prefixes, mailbox name prefixes, if you are not sure of maibox names, you can check the log
    prefixes:
      - Archive # this will match 'Archive', 'Archives', 'Archived' and 'Archives/2022' etc
      - 存档
```

## Container Image

Check [GitHub Packages](https://github.com/yankeguo/imapdump/pkgs/container/imapdump) for available container images

Check [Dockerfile](Dockerfile) for details

By default, container image will execute `imapdump` every `6 hours`

All you need to do is to mount `/data` for data persistence, and put a `config.yaml` at `/data/config.yaml`

## Credits

GUO YANKE, MIT License
