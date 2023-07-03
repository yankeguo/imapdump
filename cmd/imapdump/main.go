package main

import (
	"context"
	"flag"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/guoyk93/imapdump"
	"github.com/guoyk93/rg"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()
	defer rg.Guard(&err)

	var ctx context.Context
	{
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(context.Background())
		go func() {
			chSig := make(chan os.Signal, 1)
			signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)
			sig := <-chSig
			log.Println("received signal:", sig.String())
			cancel()
		}()
	}

	var (
		optConf string
	)
	flag.StringVar(&optConf, "conf", "config.yaml", "config file")
	flag.Parse()

	buf := rg.Must(os.ReadFile(optConf))

	var opts Options
	rg.Must0(yaml.Unmarshal(buf, &opts))
	rg.Must0(defaults.Set(&opts))
	rg.Must0(validator.New().Struct(&opts))

	for _, account := range opts.Accounts {
		log.Println("account:", account.Name)
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
