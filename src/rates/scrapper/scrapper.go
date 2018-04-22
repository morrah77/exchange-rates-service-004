// go generate src/rates/scrapper/scrapper.go -> src/rates/scrapper/mock_scrapper/mock_scrapper.go
package scrapper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"rates/common"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const pairCurrencyDelimiter = `_`
const defaultInterval = time.Second * 2

type Storage interface {
	Save(interface{}) error
}

type Scrapper struct {
	fetchInterval time.Duration
	urlPairs      string
	urlTicker     string
	pairs         string
	stor          Storage
	logger        *log.Logger
	stopChannel   chan struct{}
}

func NewScrapper(fetchInterval time.Duration, scrapUrl string, stor Storage, logger *log.Logger) *Scrapper {
	if fetchInterval.Nanoseconds() == 0 {
		fetchInterval = defaultInterval
	}
	return &Scrapper{
		fetchInterval: fetchInterval,
		urlPairs:      scrapUrl + `/info`,
		urlTicker:     scrapUrl + `/ticker/`,
		stor:          stor,
		logger:        logger,
		stopChannel:   make(chan struct{}),
	}
}

func (s *Scrapper) Start() {
	go func() {
		ticker := time.NewTicker(s.fetchInterval)
		for {
			select {
			case _ = <-ticker.C:
				s.RefreshPairs()
				s.Scrap()
			case _ = <-s.stopChannel:
				break
			}
		}
	}()
}

func (s *Scrapper) Stop() {
	s.stopChannel <- struct{}{}
}

type Info struct {
	ServerTime int64                  `json:"server_time"`
	Pairs      map[string]interface{} `json:"pairs"`
}

func (s *Scrapper) RefreshPairs() {
	resp, err := http.Get(s.urlPairs)
	if err != nil {
		s.logger.Println(errors.Wrap(err, `scrapper: Could not fetch from `+s.urlPairs))
	}

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Println(errors.Wrap(err, `scrapper: Could not read response`))
		return
	}
	info := &Info{}
	err = json.Unmarshal(b, info)
	if err != nil {
		s.logger.Println(errors.Wrap(err, `scrapper: Could not unmarshal response`).Error())
	}
	pairs := make([]string, 0)
	for p := range info.Pairs {
		pairs = append(pairs, p)
	}
	sort.Strings(pairs)
	s.pairs = strings.Join(pairs, `-`)
}

func (s *Scrapper) Scrap() {
	urlString := s.urlTicker + s.pairs + `?ignore_invalid=1`
	resp, err := http.Get(urlString)
	if err != nil {
		s.logger.Println(errors.Wrap(err, `scrapper: Could not fetch from `+urlString))
		return
	}

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Println(errors.Wrap(err, `scrapper: Could not read response`))
	}
	var t interface{}
	err = json.Unmarshal(b, &t)
	if err != nil {
		s.logger.Println(errors.Wrap(err, `scrapper: Could not unmarshal response`).Error())
		return
	}
	tickers, ok := t.(map[string]interface{})
	if !ok {
		s.logger.Println(`scrapper: Could not assert response`)
	}
	var records []common.Record
	for t := range tickers {
		ticker, ok := tickers[t].(map[string]interface{})
		if !ok {
			s.logger.Println(`scrapper: Could not assert ticker`)
			continue
		}
		dp := strings.Split(t, pairCurrencyDelimiter)
		if len(dp) != 2 {
			s.logger.Println(`scrapper: Wrong pair key`, t)
			return
		}
		r := common.Record{
			CurrencyFrom: dp[0],
			CurrencyTo:   dp[1],
			Rate:         ticker[`last`].(float64),
			Time:         ticker[`updated`].(float64),
		}
		records = append(records, r)
	}
	sort.Slice(records, func(i, j int) bool {
		if records[i].CurrencyFrom < records[j].CurrencyFrom {
			return true
		}
		if records[i].CurrencyFrom == records[j].CurrencyFrom {
			return records[i].CurrencyTo < records[j].CurrencyTo
		}
		return false
	})
	err = s.stor.Save(records)
	if err != nil {
		s.logger.Println(errors.Wrap(err, `scrapper: Could not save tickers`).Error())
		return
	}
}
