package scrapper

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"rates/common"
	"rates/scrapper/mock_scrapper"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

var expectedPairs string = `bch_btc-bch_dsh-bch_eth-bch_eur-bch_ltc-bch_rur-bch_usd-bch_zec-bchet_bch-btc_eur-btc_rur-btc_usd-btcet_btc-dsh_btc-dsh_eth-dsh_eur-dsh_ltc-dsh_rur-dsh_usd-dsh_zec-dshet_dsh-eth_btc-eth_eur-eth_ltc-eth_rur-eth_usd-eth_zec-ethet_eth-eur_rur-eur_usd-euret_eur-ltc_btc-ltc_eur-ltc_rur-ltc_usd-ltcet_ltc-nmc_btc-nmc_usd-nmcet_nmc-nvc_btc-nvc_usd-nvcet_nvc-ppc_btc-ppc_usd-ppcet_ppc-ruret_rur-usd_rur-usdet_usd-usdt_usd-zec_btc-zec_ltc-zec_usd`

func TestScrapper_RefreshPairs(t *testing.T) {
	var pairsResponseBody = `{"server_time":1524387235,"pairs":{"btc_usd":{"decimal_places":3,"min_price":0.1,"max_price":500000,"min_amount":0.001,"hidden":0,"fee":0.2},"btc_rur":{"decimal_places":5,"min_price":1,"max_price":30000000,"min_amount":0.001,"hidden":0,"fee":0.2},"btc_eur":{"decimal_places":5,"min_price":0.1,"max_price":500000,"min_amount":0.001,"hidden":0,"fee":0.2},"ltc_btc":{"decimal_places":5,"min_price":0.0001,"max_price":10,"min_amount":0.01,"hidden":0,"fee":0.2},"ltc_usd":{"decimal_places":6,"min_price":0.0001,"max_price":10000,"min_amount":0.01,"hidden":0,"fee":0.2},"ltc_rur":{"decimal_places":5,"min_price":0.01,"max_price":1000000,"min_amount":0.01,"hidden":0,"fee":0.2},"ltc_eur":{"decimal_places":3,"min_price":0.0001,"max_price":10000,"min_amount":0.01,"hidden":0,"fee":0.2},"nmc_btc":{"decimal_places":5,"min_price":0.0001,"max_price":10,"min_amount":0.1,"hidden":0,"fee":0.2},"nmc_usd":{"decimal_places":3,"min_price":0.001,"max_price":1000,"min_amount":0.1,"hidden":0,"fee":0.2},"nvc_btc":{"decimal_places":5,"min_price":0.0001,"max_price":10,"min_amount":0.1,"hidden":0,"fee":0.2},"nvc_usd":{"decimal_places":3,"min_price":0.001,"max_price":1000,"min_amount":0.1,"hidden":0,"fee":0.2},"usd_rur":{"decimal_places":5,"min_price":25,"max_price":150,"min_amount":0.1,"hidden":0,"fee":0.2},"eur_usd":{"decimal_places":5,"min_price":0.5,"max_price":2,"min_amount":0.1,"hidden":0,"fee":0.2},"eur_rur":{"decimal_places":5,"min_price":30,"max_price":200,"min_amount":0.1,"hidden":0,"fee":0.2},"ppc_btc":{"decimal_places":5,"min_price":0.0001,"max_price":10,"min_amount":0.1,"hidden":0,"fee":0.2},"ppc_usd":{"decimal_places":3,"min_price":0.001,"max_price":1000,"min_amount":0.1,"hidden":0,"fee":0.2},"dsh_btc":{"decimal_places":5,"min_price":0.0001,"max_price":10,"min_amount":0.01,"hidden":0,"fee":0.2},"dsh_usd":{"decimal_places":5,"min_price":0.1,"max_price":50000,"min_amount":0.01,"hidden":0,"fee":0.2},"dsh_rur":{"decimal_places":3,"min_price":1,"max_price":3000000,"min_amount":0.01,"hidden":0,"fee":0.2},"dsh_eur":{"decimal_places":3,"min_price":0.1,"max_price":50000,"min_amount":0.01,"hidden":0,"fee":0.2},"dsh_ltc":{"decimal_places":3,"min_price":0.1,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"dsh_eth":{"decimal_places":3,"min_price":0.1,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"dsh_zec":{"decimal_places":3,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"eth_btc":{"decimal_places":5,"min_price":0.0001,"max_price":10,"min_amount":0.01,"hidden":0,"fee":0.2},"eth_usd":{"decimal_places":5,"min_price":0.0001,"max_price":50000,"min_amount":0.01,"hidden":0,"fee":0.2},"eth_eur":{"decimal_places":5,"min_price":0.0001,"max_price":50000,"min_amount":0.01,"hidden":0,"fee":0.2},"eth_ltc":{"decimal_places":5,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"eth_rur":{"decimal_places":5,"min_price":0.0001,"max_price":3000000,"min_amount":0.01,"hidden":0,"fee":0.2},"eth_zec":{"decimal_places":3,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"bch_usd":{"decimal_places":3,"min_price":0.1,"max_price":100000,"min_amount":0.001,"hidden":0,"fee":0.2},"bch_btc":{"decimal_places":4,"min_price":0.001,"max_price":100,"min_amount":0.001,"hidden":0,"fee":0.2},"bch_rur":{"decimal_places":3,"min_price":0.0001,"max_price":6000000,"min_amount":0.01,"hidden":0,"fee":0.2},"bch_eur":{"decimal_places":3,"min_price":0.0001,"max_price":100000,"min_amount":0.01,"hidden":0,"fee":0.2},"bch_ltc":{"decimal_places":3,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"bch_eth":{"decimal_places":3,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"bch_dsh":{"decimal_places":3,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"bch_zec":{"decimal_places":3,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"zec_btc":{"decimal_places":4,"min_price":0.0001,"max_price":10,"min_amount":0.01,"hidden":0,"fee":0.2},"zec_usd":{"decimal_places":3,"min_price":0.0001,"max_price":30000,"min_amount":0.01,"hidden":0,"fee":0.2},"zec_ltc":{"decimal_places":3,"min_price":0.0001,"max_price":1000,"min_amount":0.01,"hidden":0,"fee":0.2},"usdet_usd":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"ruret_rur":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"euret_eur":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"btcet_btc":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"ltcet_ltc":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"ethet_eth":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"nmcet_nmc":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"nvcet_nvc":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"ppcet_ppc":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"dshet_dsh":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"bchet_bch":{"decimal_places":3,"min_price":0.0001,"max_price":1,"min_amount":0.01,"hidden":0,"fee":0.2},"usdt_usd":{"decimal_places":3,"min_price":0.1,"max_price":2,"min_amount":0.01,"hidden":0,"fee":0.2}}}`

	var expectedLog string = ``

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set(`Content-Length`, strconv.Itoa(len(pairsResponseBody)))
		fmt.Fprintln(w, pairsResponseBody)
	}))
	defer srv.Close()

	lb := bytes.NewBuffer(make([]byte, 0))

	s := &Scrapper{
		urlPairs: srv.URL,
		logger:   log.New(lb, ``, log.Flags()),
	}

	s.RefreshPairs()

	if s.pairs != expectedPairs {
		t.Errorf("Wrong pairs fetched!\nexpected:\n%#v\n\nGot:\n%#v\n", expectedPairs, s.pairs)
	}

	lc := lb.Bytes()
	if string(lc) != `` {
		t.Errorf("Wrong logging!\nexpected:\n%#v\n\nGot:\n%#v\n", expectedLog, string(lc))
	}
}

func TestScrapper_Scrap(t *testing.T) {
	var tickerResponseBody = `{"btc_usd":{"high":8984.999,"low":8639.99,"avg":8812.4945,"vol":7477661.56068,"vol_cur":847.87672,"last":8956.999,"buy":8962.999,"sell":8951,"updated":1524393116},"btc_rur":{"high":531107.6722,"low":513618.20813,"avg":522362.940165,"vol":146304861.92641,"vol_cur":279.61045,"last":529992.50912,"buy":530000,"sell":528000,"updated":1524393116},"btc_eur":{"high":7332.43822,"low":7100.10001,"avg":7216.269115,"vol":1813876.2362,"vol_cur":250.3217,"last":7315.37237,"buy":7335,"sell":7295.10001,"updated":1524393116},"ltc_btc":{"high":0.01716,"low":0.01631,"avg":0.016735,"vol":331.12792,"vol_cur":19877.06676,"last":0.01688,"buy":0.0169,"sell":0.01681,"updated":1524393116},"ltc_rur":{"high":8964.721,"low":8445.77466,"avg":8705.24783,"vol":99493753.86681,"vol_cur":11394.50741,"last":8898.9436,"buy":8909.07909,"sell":8872.16577,"updated":1524393116},"ltc_eur":{"high":124,"low":117,"avg":120.5,"vol":1477119.85211,"vol_cur":12225.15028,"last":123.103,"buy":123.388,"sell":123.103,"updated":1524393116}}`
	var expectedLog string = ``
	var expectedArg = []common.Record{
		common.Record{Id: "", CurrencyFrom: "btc", CurrencyTo: "eur", Rate: 7315.37237, Time: 1524393116},
		common.Record{Id: "", CurrencyFrom: "btc", CurrencyTo: "rur", Rate: 529992.50912, Time: 1524393116},
		common.Record{Id: "", CurrencyFrom: "btc", CurrencyTo: "usd", Rate: 8956.999, Time: 1524393116},
		common.Record{Id: "", CurrencyFrom: "ltc", CurrencyTo: "btc", Rate: 0.01688, Time: 1524393116},
		common.Record{Id: "", CurrencyFrom: "ltc", CurrencyTo: "eur", Rate: 123.103, Time: 1524393116},
		common.Record{Id: "", CurrencyFrom: "ltc", CurrencyTo: "rur", Rate: 8898.9436, Time: 1524393116},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	stor := mock_scrapper.NewMockStorage(ctrl)
	stor.EXPECT().Save(gomock.Any()).Do(func(arg []common.Record) {
		if !reflect.DeepEqual(arg, expectedArg) {
			t.Errorf("Wrong param passed to storage.Save()!\nexpected:\n%#v\n\nGot:\n%#v\n", expectedArg, arg)
		}
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, tickerResponseBody)
	}))
	defer srv.Close()

	lb := bytes.NewBuffer(make([]byte, 0))

	s := &Scrapper{
		urlTicker: srv.URL,
		logger:    log.New(lb, ``, log.Flags()),
		stor:      stor,
	}

	s.Scrap()

	lc := lb.Bytes()
	if string(lc) != `` {
		t.Errorf("Wrong logging!\nexpected:\n%#v\n\nGot:\n%#v\n", expectedLog, string(lc))
	}
}
