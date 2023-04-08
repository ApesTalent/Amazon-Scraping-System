package controller

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"primeprice.com/dal"
	"primeprice.com/pkg/logger"
	"primeprice.com/scrapers"
)

var proxies []dal.Proxy

// RunScraping scrape products with search link
func RunScraping(searchID primitive.ObjectID) error {
	// check app configuration for scraping
	cfg, err := dal.GetAppConfig()
	if err != nil {
		return err
	}

	if cfg.RunScrape == 0 {
		return errors.New("scraping not allowed")
	}

	// set waiting time
	wt := cfg.Interval
	if wt < 10 {
		wt = 10
	}

	// check scraping and its status of scraping
	s, err := dal.GetSearchByID(searchID)
	if err != nil {
		return err
	}
	if s == nil {
		return errors.New("search doesn't exist")
	}

	if s.Status == 0 {
		time.Sleep(time.Second * time.Duration(2*wt))
		return nil
	}

	// start scraping products
	pds, err := scrapers.GetAmazonPrimeProducts(s.URL, proxies, s.Postal)
	if err != nil {
		logger.Println(err)
		return nil
	}
	for _, v := range pds {
		err = processNewProduct(v, s.ID, cfg.Discount)
		if err != nil {
			logger.Println(err)
		}
	}
	t := time.Now()
	s.LastSearch = &t
	err = dal.UpdateSearch(*s)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * time.Duration(wt))
	return nil
}

func processNewProduct(p *dal.Product, searchID primitive.ObjectID, discount float64) error {
	if p.Price == 0 || p.Asin == "" {
		return nil
	}
	pd, err := dal.CreateProduct(p, searchID)
	if err != nil {
		return nil
	}

	now := time.Now()
	nsd := dal.Standard{
		Asin:  p.Asin,
		Price: p.Price,
		Date:  &now,
	}

	sd, err := dal.GetStandardByAsin(pd.Asin)
	if err != nil {
		// if standard with that asin doesn't exist, create new one
		err = dal.CreateStandard(nsd)
		if err != nil {
			return err
		}
	} else {
		// if time difference is bigger than two weeks, update the date, price
		tdiff := now.Sub(*sd.Date)
		if tdiff.Hours() > 336 {
			sd.Price = nsd.Price
			sd.Date = nsd.Date
			err = dal.UpdateStandard(*sd)
			if err != nil {
				return err
			}
		}

		// if new price is less than discount, send telegram notification
		if nsd.Price < sd.Price*(1-discount/100) {
			msg := buildTgMsg(sd.Price, nsd.Price, p.Title, p.URL, p.Asin)
			err = SendTgNotification(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func buildTgMsg(op, np float64, title, purl, asin string) string {
	flag := ""
	currency := "€"
	if len(purl) > 7 && purl[0:5] != "https" {
		purl = "https://" + purl
	}
	u, _ := url.Parse(purl)
	host := u.Host
	sd := ""
	if len(host) > 3 {
		sd = host[len(host)-2:]
	}
	switch sd {
	case "uk":
		flag = "🇬🇧"
		currency = "£"
		break
	case "fr":
		flag = "🇫🇷"
		currency = "€"
		break
	case "it":
		flag = "🇮🇹"
		currency = "€"
		break
	case "de":
		flag = "🇩🇪"
		currency = "€"
		break
	case "es":
		flag = "🇪🇸"
		currency = "€"
		break
	default:
		break
	}
	link := host + "/dp/" + asin
	msgText := fmt.Sprintf("%s \n\n <b>Saved Price: %.2f%s</b> \n\n <b>Price: %.2f%s</b> \n\n <b>link:</b> %s", flag+title, op, currency, np, currency, link)
	return msgText
}

func StartRecursiveScraping(searchID primitive.ObjectID) {
	logger.Println("Start RecursiveScraping............................................")
	for {
		err := RunScraping(searchID)
		if err != nil {
			logger.Println(err)
			break
		}
	}
}

// StartsAllScraping starts all scraping for each search
func StartsAllScraping() {
	ps, err := dal.ListProxy()
	if err != nil {
		logger.Println(err)
		return
	}
	proxies = ps

	ss, err := dal.ListSearch()
	if err != nil {
		logger.Println(err)
		return
	}

	for _, v := range ss {
		go func(search dal.Search) {
			StartRecursiveScraping(search.ID)
		}(v)
	}
}
