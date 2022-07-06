package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"iota_dex/config"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var baseUrl = "http://localhost:" + strconv.Itoa(config.HttpPort)

func RunTest() {
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	signature, err := crypto.Sign(crypto.Keccak256Hash([]byte(ts)).Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	sign := hexutil.Encode(signature)

	params := "ts=" + ts + "&sign=" + sign

	testPairs()
	testPrice()
	testBalance()
	testOrderSwap(params)
	testGetPendingSwapOrder(params)
	testCancelSwapOrder(params)
	testGetSwapOrderList(params)
	testOrderCollectCoin(params)
	testGetPendingCoinOrder(params)
	testCancelCoinOrder(params)
	testOrderRetrieveCoin(params)
	testGetCoinOrderList(params)
	testOrderAddLiquidity(params)
	testGetPendingLiquidityOrder(params)
	testCancelLiquidityOrder(params)
	testOrderRemoveLiquidity(params)
	testGetLiquidityOrderList(params)

	//
	testOrderSwap(params)
	testOrderCollectCoin(params)
	testOrderAddLiquidity(params)
}

func testPairs() {
	url := baseUrl + "/public/pairs"
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testPrice() {
	url := baseUrl + "/public/price?coin1=smr&coin2=iota"
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testBalance() {
	url := baseUrl + "/public/balance?account=0x96216849c49358B10257cb55b28eA603c874b05E"
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testOrderSwap(sign string) {
	params := fmt.Sprintf("source=%s&target=%s&to=%s&amount=%s&min_amount=%s", "IOTA", "SMR", "0x96216849c49358B10257cb55b28eA603c874b05E", "100", "10")
	url := baseUrl + "/order/swap?" + params + "&" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testGetPendingSwapOrder(sign string) {
	url := baseUrl + "/order/swap/pending?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testCancelSwapOrder(sign string) {
	url := baseUrl + "/order/swap/cancel?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testGetSwapOrderList(sign string) {
	url := baseUrl + "/order/swap/list?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testOrderCollectCoin(sign string) {
	params := fmt.Sprintf("coin=%s&amount=%s&account=%s", "SMR", "1000000000000000000", "0x96216849c49358B10257cb55b28eA603c874b05E")
	url := baseUrl + "/order/coin/collect?" + params + "&" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testGetPendingCoinOrder(sign string) {
	url := baseUrl + "/order/coin/pending?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testCancelCoinOrder(sign string) {
	url := baseUrl + "/order/coin/cancel?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testOrderRetrieveCoin(sign string) {
	params := fmt.Sprintf("coin=%s&amount=%s&to=%s", "SMR", "1000000000000000000", "0x96216849c49358B10257cb55b28eA603c874b05E")
	url := baseUrl + "/order/coin/retrieve?" + params + "&" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testGetCoinOrderList(sign string) {
	url := baseUrl + "/order/coin/list?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testOrderAddLiquidity(sign string) {
	params := fmt.Sprintf("coin1=%s&amount1=%s&coin2=%s", "IOTA", "1000000000000000000", "SMR")
	url := baseUrl + "/order/liquidity/add?" + params + "&" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testGetPendingLiquidityOrder(sign string) {
	url := baseUrl + "/order/liquidity/pending?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testCancelLiquidityOrder(sign string) {
	url := baseUrl + "/order/liquidity/cancel?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testOrderRemoveLiquidity(sign string) {
	params := fmt.Sprintf("coin1=%s&coin2=%s&lp=%s", "IOTA", "SMR", "1000000000000000000")
	url := baseUrl + "/order/liquidity/remove?" + params + "&" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func testGetLiquidityOrderList(sign string) {
	url := baseUrl + "/order/liquidity/list?" + sign
	fmt.Println(HttpRequest(url, "GET", "", nil))
}

func HttpRequest(url string, method string, postParams string, headers map[string]string) (int, string) {
	httpClient := &http.Client{}

	var reader io.Reader
	if len(postParams) > 0 {
		reader = strings.NewReader(postParams)
		if headers == nil {
			headers = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
		}
	} else {
		reader = nil
	}

	//构建request
	request, err := http.NewRequest(method, url, reader)
	if nil != err {
		return -100, err.Error()
	}

	//添加header
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	// 发出请求
	response, err := httpClient.Do(request)
	if nil != err {
		return -200, err.Error()
	}

	defer response.Body.Close()

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return -300, err.Error()
	}

	return 0, string(body)
}
