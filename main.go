package main

import (
	"context"
	"fmt"
	"github.com/eager7/ethereum_contracts/config"
	"github.com/eager7/ethereum_contracts/contracts"
	"github.com/eager7/ethereum_contracts/database"
	"github.com/eager7/ethereum_contracts/request"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Initialize()
	if err != nil {
		panic(err)
	}
	db, err := database.Initialize(cfg.DbOpt.Address, cfg.DbOpt.User, cfg.DbOpt.Password, cfg.DbOpt.DbName, cfg.DbOpt.MaxOpenConn, cfg.DbOpt.MaxIdleConn)
	if err != nil {
		panic(err)
	}
	requester, err := request.Initialize(cfg.EthOpt.Address)
	if err != nil {
		panic(err)
	}
	cons, err := database.SearchContracts(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("search contracts len:", len(cons))
	ctx, cancel := context.WithCancel(context.Background())
	if err := contracts.ParseContracts(ctx, db, requester, cons, cfg.EthOpt.ApiAddress, cfg.Path); err != nil {
		panic(err)
	}

	pause()
	cancel()
}

func pause() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(interrupt)
	<-interrupt
}
