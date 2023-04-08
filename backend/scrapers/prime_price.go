package scrapers

import (
	"bytes"
	"fmt"
	"html"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrr/cookiejar"
	"github.com/valyala/fasthttp"

	"github.com/hashicorp/go-retryablehttp"
	"primeprice.com/dal"
	"primeprice.com/pkg/fasthtml"
	"primeprice.com/pkg/logger"
	"primeprice.com/pkg/proxy"
)

const (
	maxRetries = 2
	waitMin    = time.Second
	waitMax    = 10 * time.Second // original = 30;
)

var priceCutter = regexp.MustCompile(`[-]?\d[\d,]*[\.,\]?[\d{2}]*`)

var cj = cookiejar.AcquireCookieJar()

var cookieContainer = make(map[string]string)

var strPost = []byte("POST")
var flgDispatch = true /* set filter "dispatches from amazon" */
var flgSoldBy = true   /* set filter "sold by Amazon" */
const Max = 100000000

func processProduct(product string, baseURL string, isWareHouse bool, proxies []dal.Proxy) *dal.Product {

	logger.Println("processProduct...")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	fp := float64(0)
	priceField := ""
	prime := float64(Max)
	qid := ""
	sr := ""

	asin := fasthtml.GetAttr(product, "data-asin")
	if asin == "" {
		return nil
	}

	// check if prime fields exist  and parse it , get params list...
	linkUrl := checkPrimeExists(product)
	if len(linkUrl) < 10 {
		return nil
	}

	u, err := url.Parse(linkUrl)
	if err != nil {
		logger.Println(err)
		return nil
	}
	if u.RawQuery != "" {
		m, err := url.ParseQuery(u.RawQuery)
		if err == nil {
			for k, v := range m {
				if k == "qid" {
					qid = strings.Join(v, "")
				} else if k == "sr" {
					sr = strings.Join(v, "")
				}

			}
		}
	}

	/* -------create prime filter url and get Content------- */
	filterUrl := ""
	filterUrl = "https://" + baseURL + "/gp/aod/ajax/ref=aod_f_primeEligible?"
	filterUrl += "asin=" + asin
	filterUrl += "&qid=" + qid
	filterUrl += "&sr=" + sr
	if isWareHouse {
		filterUrl += "&pc=dp&pageno=1&filters={\"all\":true,\"primeEligible\":true}&isonlyrenderofferlist=true"
	} else {
		filterUrl += "&pc=dp&pageno=1&filters={\"all\":true,\"primeEligible\":true}&isonlyrenderofferlist=true"
	}

	filterContent, err := retryLoadDocument(filterUrl, proxies)
	if err != nil {
		logger.Println("get FilterContent...", err)
		return nil
	}

	/* pase filterContent to each prime list and get Suitable and cheap prime */
	items := fasthtml.FindNodes(filterContent, "id=\"aod-offer\"")

	if len(items) == 0 {
		return nil
	}

	if isWareHouse {
		flgDispatch = true
		flgSoldBy = true
	} else {
		flgDispatch = false
		flgSoldBy = false
	}

	for _, it := range items {
		//set dispatch and sold by  flag ..
		isDispatch := true
		isSold := true
		// check dispaches from amazon......
		if flgDispatch {
			dispatchHtml := fasthtml.GetTagWithParams(it, "id", "aod-offer-shipsFrom")
			dispatchString := fasthtml.GetInner(fasthtml.GetTagWithParams(dispatchHtml, "class", "a-size-small a-color-base"))
			isDispatch = strings.Contains(dispatchString, "Amazon")
		}

		if flgSoldBy {
			soldHtml := fasthtml.GetTagWithParams(it, "id", "aod-offer-soldBy")
			soldString := fasthtml.GetInner(fasthtml.GetTagWithParams(soldHtml, "class", "a-size-small a-link-normal"))
			isSold = strings.Contains(soldString, "Amazon Warehouse")
		}

		if isDispatch && isSold {
			_primeHtml := fasthtml.GetInner(fasthtml.GetTagWithParams(it, "class", "a-offscreen"))
			_prime := getPriceValue(_primeHtml)
			if _prime < prime {
				prime = _prime
			}
		}

	}

	/*---prime filter end ---------------------------------*/

	title := ""
	titleNode := fasthtml.GetTagWithParams(product, "class", "a-size-medium a-color-base a-text-normal")
	if len(titleNode) < 3 {
		titleNode = fasthtml.GetTagWithParams(product, "class", "a-size-base-plus a-color-base a-text-normal")
	}
	if len(titleNode) > 0 {
		titleTxt := fasthtml.GetInner(titleNode)
		title = html.UnescapeString(titleTxt)
		if title == "" {
			logger.Println("return with title nil")
			return nil
		}
	} else {
		return nil
	}

	// get price...
	if isWareHouse {
		priceField = fasthtml.GetInner(fasthtml.GetTagWithParams(product, "class", "a-color-base"))
	} else {
		htmlPrice := fasthtml.GetTagWithParams(product, "class", "a-offscreen")
		if len(htmlPrice) > 5 {
			priceField = fasthtml.GetInner(htmlPrice)
		}
	}

	submatchall := priceCutter.FindAllString(priceField, -1)
	fp = parsePrice(strings.Join(submatchall, ""))
	// href := product.Find("a.a-link-normal.a-text-normal").First().Attr("href")
	href := fasthtml.GetAttr(fasthtml.GetTagWithParams(product, "class", "a-link-normal a-text-normal"), "href")

	pd := &dal.Product{
		Asin:  asin,
		Title: title,
		Price: getRealValue(prime),
		Prime: getRealValue(prime),
		URL:   baseURL + href,
	}

	logger.Println("\n--------------------------------------- Amazon Product with PRIME ------------------------------------------------------------------------------------------------")
	logger.Println("----- ", title[:7], "----- price:", fp, "----- ", pd.Asin, "----- ", href[:10], "----- PRIME:", getRealValue(prime))
	logger.Println("-----------------------------------------------------------------------------------------------------------------------------------------------\n")
	return pd
}

func getRealValue(value float64) float64 {
	if value < Max {
		return value
	}
	return 0
}

func getPriceValue(content string) float64 {
	submatchall := priceCutter.FindAllString(content, -1)
	value := parsePrice(strings.Join(submatchall, ""))
	return value
}

func checkPrimeExists(content string) string {
	isPrime := fasthtml.GetTagWithParams(content, "data-action", "s-show-all-offers-display")
	if len(isPrime) > 1 {
		linkUrl := fasthtml.GetTagWithParams(isPrime, "class", "a-link-normal")
		return fasthtml.GetAttr(linkUrl, "href")
	} else {
		return ""
	}
}

func parsePrice(s string) float64 {
	price := strings.ReplaceAll(s, ".", "")
	price = strings.ReplaceAll(price, ",", "")
	fpr, _ := strconv.ParseFloat(price, 64)
	return fpr / 100
}

func GetStringInBetweenTwoString(str string, startS string, endS string) (result string) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result
	}
	result = newS[:e]
	return result
}

func ChangeZipCode(purl string, proxies []dal.Proxy, zipCode string) {
	cookieContainer = make(map[string]string)
	u, err := url.Parse(purl)
	if err != nil {
		logger.Println(err)
	} else {
		logger.Println(u)
	}
	formData := "locationType=COUNTRY&district=ES&countryCode=ES&storeContext=generic&deviceType=web&pageType=Gateway&actionSource=glow&almBrandId=undefined"
	domainFix := "com"
	searchTxt := "\"Tu dirección de envío actual es:\""
	if strings.Contains(purl, "www.amazon.es") {
		domainFix = "es"
		searchTxt = "\"Tu dirección de envío actual es:\""
		formData = "locationType=LOCATION_INPUT&zipCode=" + zipCode + "&storeContext=generic&deviceType=web&pageType=Gateway&actionSource=glow&almBrandId=undefined"
	} else if strings.Contains(purl, "www.amazon.de") {
		domainFix = "de"
		searchTxt = "\"Sie kaufen gerade ein für:\""
	} else if strings.Contains(purl, "www.amazon.fr") {
		domainFix = "fr"
		searchTxt = "\"Votre lieu de livraison est désormais:\""
	} else if strings.Contains(purl, "www.amazon.co.uk") {
		domainFix = "co.uk"
		searchTxt = "\"You're now shopping for delivery to:\""
	} else if strings.Contains(purl, "www.amazon.it") {
		domainFix = "it"
		searchTxt = "\"L'indirizzo di consegna selezionato è:\""
	} else {
		domainFix = "com"
		searchTxt = "\"You're now shopping for delivery to:\""
	}
	homePage := ""
	homePage, err = getRequest("https://www.amazon."+domainFix, proxies)

	//uidCookie := strings.TrimLeft(strings.TrimRight(homePage, "\" })</script>"), "/ah/ajax/counter?ctr=desktop_ajax_atf")
	uidCookie := GetStringInBetweenTwoString(homePage, "/ah/ajax/counter?ctr=desktop_ajax_atf", "\" })</script>")

	uidCookieUrl := "https://www.amazon." + domainFix + "/ah/ajax/counter?ctr=desktop_ajax_atf" + uidCookie
	postRequest(uidCookieUrl, proxies, "")
	tokenPage := ""
	tokenPage, err = getRequest("https://www.amazon."+domainFix+"/gp/glow/get-address-selections.html?deviceType=desktop&pageType=Gateway&storeContext=NoStoreName", proxies)
	//crosToken := strings.TrimLeft(strings.TrimRight(tokenPage, "\", IDs:{\"ADDRESS_LIST\":\"GLUXAddressList\""), "\"You're now shopping for delivery to:\", CSRF_TOKEN : \"")
	crosToken := GetStringInBetweenTwoString(tokenPage, searchTxt+", CSRF_TOKEN : \"", "\", IDs:{\"ADDRESS_LIST\":\"GLUXAddressList\"")
	changeZipCodePostRequest("https://www.amazon."+domainFix+"/gp/delivery/ajax/address-change.html", proxies, formData, crosToken)
	postRequest("https://www.amazon."+domainFix+"/gp/glow/get-location-label.html", proxies, "storeContext=hpc&pageType=Landing")

}

// GetAmazonPrimeProducts get all prime products from search url
func GetAmazonPrimeProducts(purl string, proxies []dal.Proxy, postal string) ([]*dal.Product, error) {

	// purl = "https://www.amazon.co.uk/s?k=laptops&rh=n%3A340831031%2Cp_76%3A419159031%2Cp_89%3AASUS&dc&qid=1617938829&rnid=1632651031&ref=sr_nr_p_89_1"
	//purl = "https://www.amazon.fr/s?i=computers&bbn=3581943031&rh=n%3A3581943031%2Cn%3A340858031%2Cn%3A429879031%2Cp_36%3A80000-&dc&pf_rd_i=3581943031&pf_rd_m=A1X6FK5RDHNB96&pf_rd_p=42630c51-0c6e-47d7-be65-54b20c5deb6e&pf_rd_r=BWY1QSRN565BZDA8C7CS&pf_rd_s=merchandised-search-3&pf_rd_t=101&qid=1613071424&rnid=9733298031&ref=sr_nr_p_36_7"

	var ps []*dal.Product
	curPage := 1
	baseURL := ""
	isWareHouse := false
	u, err := url.Parse(purl)

	if err != nil {
		logger.Println(err)
	} else {
		baseURL = u.Host
	}

	ChangeZipCode(purl, proxies, postal)

	for {

		tt := time.Now()
		surl := fmt.Sprintf("%s&page=%d", purl, curPage)

		content, err := retryLoadDocument(surl, proxies)
		if strings.Contains(content, "To discuss automated access to Amazon data") {
			logger.Println("Scrapping Detected.Please change cookie")
			ChangeZipCode(purl, proxies, postal)
		}
		if err != nil {
			logger.Println("page:", curPage, err)
			break
		}

		logger.Println("request done", time.Since(tt).Seconds(), "sec")
		tt = time.Now()

		if strings.Contains(fasthtml.GetTagWithParams(content, "selected", "selected"), "Amazon Warehouse") {
			isWareHouse = true
		}

		// Cut upper js part for collision safety
		/* in case of other design...
		if isWareHouse {
			content = fasthtml.GetTagWithParams(content, "class", "s-main-slot s-result-list s-search-results")
			if len(content) > 0 {
				logger.Println("warehoues search div found")
			} else {
				logger.Println("warehoues search div not found")
			}
		}
		if len(content) > 10 {
			logger.Println("search id div found")
		} else {
			logger.Println("search id empty")
		}
		*/

		searchIdx := strings.Index(content, "<div id=\"search\">")
		if searchIdx == -1 {
			logger.Println("page:", curPage, "search block not found")
			break
		}
		content = content[searchIdx:]

		items := fasthtml.FindNodes(content, "data-component-type=\"s-search-result\"")
		logger.Println("find done", isWareHouse, len(items), time.Since(tt).Seconds(), "sec... items...", len(items))
		tt = time.Now()

		if len(items) == 0 {
			break
		}

		for _, it := range items {
			pd := processProduct(it, baseURL, isWareHouse, proxies)
			if pd != nil {
				ps = append(ps, pd)
			}
		}

		logger.Println("process items done len ", time.Since(tt).Seconds(), "sec... len", len(ps))

		if len(ps) > 200 {
			break
		}
		/* get next page url...
		pagination := fasthtml.GetTagWithParams(content, "class", "a-pagination")
		next := fasthtml.GetTagWithParams(pagination, "class", "a-last")
		if len(next) < 10 {
			break
		}
		surl = fasthtml.GetAttr(next, "href") */
		curPage++
	}

	return ps, nil
}

func retryLoadDocument(surl string, proxies []dal.Proxy) (string, error) {
	var lastErr error
	for n := 0; n < maxRetries; n++ {
		document, err, shouldRetry := loadDocument(surl, proxies)
		if !shouldRetry {
			return document, err
		}

		lastErr = err
		backoff := retryablehttp.DefaultBackoff(waitMin, waitMax, n, nil)
		logger.Println(err, "Retrying in", backoff)
		time.Sleep(backoff)
	}
	return "", lastErr
}

func getRequest(surl string, proxies []dal.Proxy) (string, error) {
	logger.Println("processing", surl)
	var client fasthttp.Client
	if len(proxies) > 0 {
		px := getRandomProxy(proxies)
		client = fasthttp.Client{
			Dial: proxy.FastHTTPProxyDialer(px),
		}

		logger.Println("with proxy", px)
	}

	defer client.CloseIdleConnections()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	// Acquire cookie jar
	u, errUrl := url.Parse(surl)
	if errUrl == nil {
		cj = cookiejar.AcquireCookieJar()
		for key, value := range cookieContainer {
			if strings.Contains(key, u.Host) {
				key = strings.Replace(key, u.Host, "", -1)
				valueArry := strings.Split(value, "=")
				value = strings.Split(valueArry[1], ";")[0]
				cj.Set(key, value)
			}
		}
	}
	cj.FillRequest(req)

	req.SetRequestURI(surl)

	req.Header.Set("Content-Type", "text/html;charset=UTF-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("user-agent", getRandomUserAgent())
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Connection", "keep-alive")
	err := client.DoTimeout(req, resp, 30*time.Second)
	if err != nil {
		return "", err
	}

	resp.Header.VisitAllCookie(func(key, value []byte) {
		c := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(c)

		c.ParseBytes(value)
		var emptyContent = string(key) + "=-;"
		if !strings.Contains(string(value), emptyContent) {
			var middle = strings.Replace(string(value), "Domain=.amazon", "domain=.www.amazon", -1)
			middle = strings.Replace(middle, "domain=.amazon", "domain=.www.amazon", -1)
			cookieContainer[string(key)+u.Host] = middle
		}
	})
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		fmt.Println("Unzipping...")
		body, _ = resp.BodyGunzip()
	} else {
		body = resp.Body()
	}
	content := string(body)
	return content, nil
}

func changeZipCodePostRequest(surl string, proxies []dal.Proxy, formData string, token string) (string, error, bool) {
	logger.Println("processing", surl)
	var client fasthttp.Client
	if len(proxies) > 0 {
		px := getRandomProxy(proxies)
		client = fasthttp.Client{
			Dial: proxy.FastHTTPProxyDialer(px),
		}

		logger.Println("with proxy", px)
	}

	defer client.CloseIdleConnections()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	// Acquire cookie jar
	u, errUrl := url.Parse(surl)
	if errUrl == nil {
		cj = cookiejar.AcquireCookieJar()
		for key, value := range cookieContainer {
			if strings.Contains(key, u.Host) {
				key = strings.Replace(key, u.Host, "", -1)
				valueArry := strings.Split(value, "=")
				value = strings.Split(valueArry[1], ";")[0]
				cj.Set(key, value)
			}
		}
	}
	cj.FillRequest(req)

	req.SetRequestURI(surl)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("anti-csrftoken-a2z", token)
	req.Header.Set("Connection", "keep-alive")
	req.Header.SetMethodBytes(strPost)
	req.SetBodyString(formData)
	err := client.DoTimeout(req, resp, 30*time.Second)
	if err != nil {
		return "", err, true
	}
	resp.Header.VisitAllCookie(func(key, value []byte) {
		c := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(c)

		c.ParseBytes(value)
		var emptyContent = string(key) + "=-;"
		if !strings.Contains(string(value), emptyContent) {
			var middle = strings.Replace(string(value), "Domain=.amazon", "domain=.www.amazon", -1)
			middle = strings.Replace(middle, "domain=.amazon", "domain=.www.amazon", -1)
			cookieContainer[string(key)+u.Host] = middle
		}
	})
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		fmt.Println("Unzipping...")
		body, _ = resp.BodyGunzip()
	} else {
		body = resp.Body()
	}
	content := string(body)
	return content, nil, false
}

func postRequest(surl string, proxies []dal.Proxy, formData string) (string, error, bool) {
	logger.Println("processing", surl)
	var client fasthttp.Client
	if len(proxies) > 0 {
		px := getRandomProxy(proxies)
		client = fasthttp.Client{
			Dial: proxy.FastHTTPProxyDialer(px),
		}

		logger.Println("with proxy", px)
	}

	defer client.CloseIdleConnections()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	// Acquire cookie jar
	u, errUrl := url.Parse(surl)
	if errUrl == nil {
		cj = cookiejar.AcquireCookieJar()
		for key, value := range cookieContainer {
			if strings.Contains(key, u.Host) {
				key = strings.Replace(key, u.Host, "", -1)
				valueArry := strings.Split(value, "=")
				value = strings.Split(valueArry[1], ";")[0]
				cj.Set(key, value)
			}
		}
	}
	cj.FillRequest(req)

	req.SetRequestURI(surl)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.SetMethodBytes(strPost)
	req.SetBodyString(formData)
	err := client.DoTimeout(req, resp, 30*time.Second)
	if err != nil {
		return "", err, true
	}
	resp.Header.VisitAllCookie(func(key, value []byte) {
		c := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(c)

		c.ParseBytes(value)
		var emptyContent = string(key) + "=-;"
		if !strings.Contains(string(value), emptyContent) {
			var middle = strings.Replace(string(value), "Domain=.amazon", "domain=.www.amazon", -1)
			middle = strings.Replace(middle, "domain=.amazon", "domain=.www.amazon", -1)
			cookieContainer[string(key)+u.Host] = middle
		}
	})
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		fmt.Println("Unzipping...")
		body, _ = resp.BodyGunzip()
	} else {
		body = resp.Body()
	}
	content := string(body)
	return content, nil, false
}

func loadDocument(surl string, proxies []dal.Proxy) (string, error, bool) {
	logger.Println("processing", surl)
	var client fasthttp.Client
	if len(proxies) > 0 {
		px := getRandomProxy(proxies)
		client = fasthttp.Client{
			Dial: proxy.FastHTTPProxyDialer(px),
		}

		logger.Println("with proxy", px)
	}
	defer client.CloseIdleConnections()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	cj.FillRequest(req)

	req.SetRequestURI(surl)
	req.Header.Set("Content-Type", "text/html;charset=UTF-8")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("user-agent", getRandomUserAgent())
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Connection", "keep-alive")

	err := client.DoTimeout(req, resp, 30*time.Second)
	if err != nil {
		return "", err, true
	}
	u, errUrl := url.Parse(surl)
	if errUrl == nil {
		cj = cookiejar.AcquireCookieJar()
		for key, value := range cookieContainer {
			if strings.Contains(key, u.Host) {
				key = strings.Replace(key, u.Host, "", -1)
				valueArry := strings.Split(value, "=")
				value = strings.Split(valueArry[1], ";")[0]
				cj.Set(key, value)
			}
		}
	}
	resp.Header.VisitAllCookie(func(key, value []byte) {
		c := fasthttp.AcquireCookie()
		defer fasthttp.ReleaseCookie(c)

		c.ParseBytes(value)
		var emptyContent = string(key) + "=-;"
		if !strings.Contains(string(value), emptyContent) {
			var middle = strings.Replace(string(value), "Domain=.amazon", "domain=.www.amazon", -1)
			middle = strings.Replace(middle, "domain=.amazon", "domain=.www.amazon", -1)
			cookieContainer[string(key)+u.Host] = middle
		}
	})
	content := string(resp.Body())
	return content, nil, false
}

func getRandomProxy(ps []dal.Proxy) string {
	if len(ps) == 0 {
		return ""
	}
	i := rand.Intn(len(ps))
	return ps[i].Proxy
}
