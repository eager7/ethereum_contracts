package request

import "testing"

func TestRequestContract(t *testing.T) {
	_, _, _, err := (&Requester{}).RequestContract("https://etherscan.io/address/0xB8c77482e45F1F44dE1745F52C74426C631bDD52")
	if err != nil {
		t.Fatal(err)
	}
}
