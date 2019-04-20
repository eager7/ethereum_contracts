package contracts

import (
	"context"
	"fmt"
	"github.com/eager7/ethereum_contracts/database/tables"
	"github.com/eager7/ethereum_contracts/request"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

type Contract struct {
	hash     string
	count    int64
	full     string
	contract string
	abi      string
	code     string
}

func ParseContracts(ctx context.Context, db *gorm.DB, requester *request.Requester, contracts []*tables.ContractCountInfo, api, dir string) error {
	for _, contract := range contracts {
		c := Contract{hash: contract.Hash, count: contract.Count}
	retry:
		if err := c.RequestContractCode(ctx, db, requester, api, dir); err != nil {
			fmt.Println("get code err:", err)
			time.Sleep(time.Second * 1)
			goto retry
		}
	}
	return nil
}

func (c *Contract) RequestContractCode(ctx context.Context, db *gorm.DB, requester *request.Requester, api, dir string) error {
	var contractsInfo []*tables.TableContractInfo
	if err := db.Where("hash=?", c.hash).Find(&contractsInfo).Error; err != nil {
		return err
	}
	for _, ct := range contractsInfo {
		fmt.Println("handle contract:", ct.Address)
		tx, _, err := requester.Client.TransactionByHash(ctx, common.HexToHash(ct.Transaction))
		if err != nil {
			return err
		}
		c.code = common.Bytes2Hex(tx.Data())
		c.full, c.contract, c.abi, _, err = requester.RequestContract(api + common.HexToAddress(ct.Address).Hex())
		if err != nil {
			return err
		}
		if err := c.StoreContract(fmt.Sprintf("%s/%d_%s/%s/", dir, c.count, c.hash, common.HexToAddress(ct.Address).Hex())); err != nil {
			return err
		}
	}
	return nil
}

func (c *Contract) StoreContract(dir string) error {
	if err := WriteFile(dir, "full.txt", c.full); err != nil {
		return err
	}
	if err := WriteFile(dir, "abi.bin", c.abi); err != nil {
		return err
	}
	if err := WriteFile(dir, "code.bin", c.code); err != nil {
		return err
	}
	if err := WriteFile(dir, "code.txt", c.contract); err != nil {
		return err
	}
	return nil
}

func WriteFile(dir, name, data string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	file, err := os.Create(dir + name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}
