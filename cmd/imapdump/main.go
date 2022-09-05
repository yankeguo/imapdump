package main

import (
	"context"
	"github.com/guoyk93/grace"
	"github.com/guoyk93/grace/graceconf"
	"github.com/guoyk93/grace/gracemain"
	"github.com/guoyk93/grace/gracenotify"
	"github.com/guoyk93/grace/gracetrack"
	"github.com/guoyk93/imapdump"
	"log"
	"path/filepath"
)

type Options struct {
	Dir      string `yaml:"dir" validate:"required"`
	Accounts []struct {
		Name     string   `yaml:"name" validate:"required"`
		Host     string   `yaml:"host" validate:"required"`
		Username string   `yaml:"username" validate:"required"`
		Password string   `yaml:"password" validate:"required"`
		Prefixes []string `yaml:"prefixes" validate:"required"`
	} `yaml:"accounts"`
}

func main() {
	var (
		err error

		ctx, _ = gracemain.WithSignalCancel(
			gracetrack.Init(context.Background()),
		)
	)

	defer gracemain.Exit(&err)
	defer gracenotify.Notify("[IMAPDUMP]", &ctx, &err)
	defer grace.Guard(&err)

	opts := grace.Must(graceconf.LoadYAMLFlagConf[Options]())

	_ = gracemain.WriteLastRun(opts.Dir)

	for _, account := range opts.Accounts {
		log.Println("started:", account.Name)
		if err = imapdump.DumpAccount(ctx, imapdump.DumpAccountOptions{
			DisplayName: account.Name,
			Dir:         filepath.Join(opts.Dir, account.Name),
			Host:        account.Host,
			Username:    account.Username,
			Password:    account.Password,
			Prefixes:    account.Prefixes,
		}); err != nil {
			log.Println("failed to dump account:", account.Name, ":", err.Error())
			err = nil
			continue
		}
	}
}
