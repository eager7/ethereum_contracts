package contracts

import (
	"context"
	"fmt"
	"github.com/eager7/ethereum_contracts/config"
	"github.com/eager7/ethereum_contracts/database"
	"github.com/eager7/ethereum_contracts/request"
	"testing"
)

func TestWriteFile(t *testing.T) {
	if err := WriteFile("/tmp/test_dir/", "code.bin", "test"); err != nil {
		t.Fatal(err)
	}
}

func TestParseContracts(t *testing.T) {
	cfg, err := config.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	db, err := database.Initialize(cfg.DbOpt.Address, cfg.DbOpt.User, cfg.DbOpt.Password, cfg.DbOpt.DbName, cfg.DbOpt.MaxOpenConn, cfg.DbOpt.MaxIdleConn)
	if err != nil {
		t.Fatal(err)
	}
	requester, err := request.Initialize(cfg.EthOpt.Address)
	if err != nil {
		t.Fatal(err)
	}

	c := Contract{hash: "ff24fbc66d5f41cfdb397e735e7bb4b74480738871564d1e75dd5b183f9892a8", count: 92}
	if err := c.RequestContractCode(context.Background(), db.LogMode(true), requester, cfg.EthOpt.ApiAddress, cfg.Path); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", c)
}
