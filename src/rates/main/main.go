package main

import (
	"flag"
	"log"
	"rates/scrapper"
	"time"

	"os"

	"rates/api"
	"rates/storage"

	"github.com/pkg/errors"
)

type Conf struct {
	ScrapUrl      string
	ScrapInterval string
	ListenAddr    string
	ApiPath       string
	Storage       string
	Dsn           string
}

var (
	conf   Conf
	logger *log.Logger
)

func init() {
	flag.StringVar(&conf.ScrapUrl, `scrap-url`, `https://wex.nz/api/3`, `Url to scrap`)
	flag.StringVar(&conf.ScrapInterval, `scrap-intervall`, `2s`, `Interval to scrap in seconds`)
	flag.StringVar(&conf.ListenAddr, `listen-addr`, `:8080`, `Address to listen`)
	flag.StringVar(&conf.ApiPath, `api-path`, `/rates/v0`, `API path to handle rate requests`)
	flag.StringVar(&conf.Storage, `storage`, `postgres`, `Storage type`)
	flag.StringVar(&conf.Dsn, `dsn`, ``, `Storage DSN`)
	logger = log.New(os.Stdout, ``, log.Flags())
}

func main() {
	flag.Parse()
	stor, err := storage.NewStorage(conf.Storage, conf.Dsn, logger)
	if err != nil {
		panic(errors.Wrap(err, `Could not create storage!`))
	}
	defer func() {
		stor.Stop()
	}()

	scrapInterval, err := time.ParseDuration(conf.ScrapInterval)
	if err != nil {
		panic(errors.Wrap(err, `Could not parse scrap interval!`))
	}
	scrap := scrapper.NewScrapper(scrapInterval, conf.ScrapUrl, stor, logger)
	defer func() {
		scrap.Stop()
	}()
	scrap.Start()

	apiInstance, err := api.NewApi(conf.ListenAddr, conf.ApiPath, stor, logger)
	if err != nil {
		panic(errors.Wrap(err, `Could not create API!`))
	}
	defer func() {
		apiInstance.Stop()
	}()
	err = apiInstance.Start()
	if err != nil {
		panic(errors.Wrap(err, `Could not start API!`))
	}
}
