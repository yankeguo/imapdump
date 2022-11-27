# imapdump

![workflow badge](https://github.com/guoyk93/imapdump/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/guoyk93/imapdump.svg)](https://pkg.go.dev/github.com/guoyk93/imapdump)

A tool for dumping emails to local in batch.

## 中文使用说明

* [imapdump - 批量备份邮件到本地](https://mp.weixin.qq.com/s/q-BAjuInDjc6hkpotHtRPg)

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

Check [GitHub Packages](https://github.com/guoyk93/imapdump/pkgs/container/imapdump) for available container images

Check [Dockerfile](Dockerfile) for details

By default, container image will execute `imapdump` every `6 hours`

All you need to do is to mount `/data` for data persistence, and put a `config.yaml` at `/data/config.yaml`

## Notification

Execution result will be delivered to environment variable `$NOTIFY_URL`, if given, by HTTP `POST`.

```
{"text": "MESSAGE..."}
```

## Upstream

<https://git.guoyk.net/go-guoyk/imapdump>

Due to various reasons, codebase is detached from upstream.

## Donation

View <https://guoyk.xyz/donation>

## Credits

Guo Y.K., MIT License
