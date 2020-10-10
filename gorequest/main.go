package main

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

func main() {
	/*
		https://www.google.com/finance/converter?a=" . $amount . "&from=" . $from_Currency . "&to=" . $to_Currency;
	*/
	url := `http://api.k780.com/?app=finance.rate&scur=USD&tcur=CNY&appkey=10003&sign=b59bc3ef6191eb9f747dd4e83c99f2a4`

	request := gorequest.New()
	//resp, body, errs := request.Get(url).End()

	//// To reuse same client with no_proxy, use empty string:
	resp, body, errs := request.Proxy("").Get(url).End()

	fmt.Println(resp, body, errs)
}
