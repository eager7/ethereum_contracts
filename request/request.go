package request

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"strings"
	"time"
)

type Requester struct {
	Client *ethclient.Client
}

func Initialize(url string) (*Requester, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("initialize eth client err:%v", err))
	}

	return &Requester{Client: client}, nil
}

func (r *Requester) RequestContract(url string) (contract, abi, code string, err error) {
	res := gorequest.New()
	ret, body, errs := res.Timeout(time.Second * 5).Get(url).End()
	if errs != nil || ret.StatusCode != http.StatusOK {
		req, err := res.MakeRequest()
		if err == nil {
			fmt.Printf("request body:%+v\n", req)
		}
		var errStr string
		for _, e := range errs {
			errStr += e.Error()
		}
		return "", "", "", errors.New(errStr)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", "", "", err
	}
	codeSelector := doc.Find("#dividcode")
	contract = codeSelector.Find("pre[class='js-sourcecopyarea editor']").Text()
	abi = codeSelector.Find("pre[class='wordwrap js-copytextarea2']").Text()
	code = codeSelector.Find("#verifiedbytecode2").Text()

	return "", "", "", nil
}
