# imapdump

![workflow badge](https://github.com/yankeguo/imapdump/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/yankeguo/imapdump.svg)](https://pkg.go.dev/github.com/yankeguo/imapdump)

一个用来批量下载邮件到本地的工具

## 用法

**命令**

```
./imapdump -conf config.yaml
```

**配置文件**

```yaml
dir: dump
accounts:
  # 任务名称，用于生成子目录
  - name: username@mydomain.com
    # IMAP 服务器地址
    host: imap.mydomain.com:993
    # 登录用户名
    username: username@mydomain.com
    # 登录密码
    password: xxxxxxxxxxxxxxxxxxx
    # 要备份的文件夹前缀
    prefixes:
      - Archive # 可以匹配 'Archive', 'Archives', 'Archived' 和 'Archives/2022' 等
      - 存档
```

## 容器镜像

查看 [GitHub Packages](https://github.com/yankeguo/imapdump/pkgs/container/imapdump)

容器构建细节，请检查 [Dockerfile](Dockerfile)

默认情况下，容器镜像会每隔六个小时执行一次 `imapdump`，你只需要挂载 `/data` 以保持数据持久化，并在 `/data/config.yaml` 中放置配置文件即可

## 许可证

GUO YANKE, MIT License
