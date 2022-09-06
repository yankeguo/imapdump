# imapdump

![workflow badge](https://github.com/guoyk93/imapdump/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/guoyk93/imapdump.svg)](https://pkg.go.dev/github.com/guoyk93/imapdump)

A tool for dumping emails to local in batch.

## 中文使用说明

* [imapdump - 批量备份邮件到本地](https://mp.weixin.qq.com/s?__biz=Mzg2ODIyNzg2Ng==&mid=2247483664&idx=1&sn=1748de50e7acff3738f6c03971b77b1e&chksm=ceaecf65f9d946733c3d06d28d43461469f61c303b06eff75c68d7d720ddd49363dc85b4a495#rd)

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

## Notification

Execution result will be delivered to environment variable `$NOTIFY_URL`, if given, by HTTP `POST`.

```
{"text": "MESSAGE..."}
```

## Upstream

https://git.guoyk.net/go-guoyk/imapdump

Due to various reasons, codebase is detached from upstream.

## Credits

Guo Y.K., MIT License
