package contracts

import (
	"context"
	"fmt"
	"github.com/eager7/ethereum_contracts/database/tables"
	"github.com/eager7/ethereum_contracts/request"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"os"
)

type Contract struct {
	hash     string
	count    int64
	code     string
	abi      string
	contract string
}

func ParseContracts(ctx context.Context, db *gorm.DB, requester *request.Requester, contracts []*tables.ContractCountInfo, api, dir string) error {
	for _, contract := range contracts {
		c := Contract{hash: contract.Hash, count: contract.Count}
		if err := c.RequestContractCode(ctx, db, requester, api, dir); err != nil {
			return err
		}
	}

	return nil
}

func (c *Contract) RequestContractCode(ctx context.Context, db *gorm.DB, requester *request.Requester, api, dir string) error {
	var contractInfo tables.TableContractInfo
	if err := db.Where("hash=?", c.hash).Find(&contractInfo).Error; err != nil {
		return err
	}
	tx, _, err := requester.Client.TransactionByHash(ctx, common.HexToHash(contractInfo.Hash))
	if err != nil {
		return err
	}
	c.code = common.Bytes2Hex(tx.Data())
	c.contract, c.abi, _, err = requester.RequestContract(api + "/" + common.HexToAddress(c.hash).Hex())
	if err != nil {
		return err
	}
	return nil
}

func (c *Contract) StoreContract(dir string) error {
	dir = fmt.Sprintf("%s/%v-%s/", dir, c.count, c.hash)
	if err := WriteFile(dir, "code.bin", c.code); err != nil {
		return err
	}
	if err := WriteFile(dir, "code.abi", c.abi); err != nil {
		return err
	}
	if err := WriteFile(dir, "code.txt", c.contract); err != nil {
		return err
	}
	return nil
}

func WriteFile(dir, name, data string) error {
	if err := os.Mkdir(dir, 0755); err != nil {
		return err
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
