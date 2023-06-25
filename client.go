/*
	Random Word Generator
	==========================
	This module is a powerful tool for generating random words by scraping and wrapping the functionality of RandomWordGenerator.com.

	It allows easy access to the website's interface, enabling efficient random word generation for various applications.

	Author: @z3ntl3
	License: GNU

*/

package randomwordgenerator

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	Version string = "v1.0.0" // exports the current wrapper version
	uri string = "https://randomwordgenerator.com/json/words.php?"
	maxQuery int = 50
)

type (
	
	Client struct {
		quantity int // determines how many words to generate
		proxy url.URL // http/https proxy URI
		ctx context.Context // context 
		useProxy bool

		initialized bool
		httpClient http.Client
	}

	words = []string // slice object holding all the words
)

/*
Sets the quantity which determines how many words will be generated at once at a time
“quantity“ can not be more than 50.

*/ 
func (c *Client) SetQuantity(quantity int) error {
	
	if quantity > maxQuery || quantity <= 0 {
		return errors.New(
			fmt.Sprintf("Quantity needs to be equal or less than %d and greater than 0", maxQuery),
		)
	}

	c.quantity = quantity
	return nil
}

/*
Gets the current set quantity
*/
func (c *Client) GetQuantity() (quantity int) {
	quantity = c.quantity
	return
}


/*
Sets the proxy to use while connecting and gathering data from the script utility hosted on randomwordgenerator.com

This method needs to be called before the initialize method
*/
func(c *Client) SetProxy(proxy string) error {
	proxy = strings.ToLower(proxy)
	uri, err := url.Parse(proxy); if err != nil {
		return err
	}
	
	if uri.Scheme != "http" && uri.Scheme != "https" {
		return errors.New("Only http/https proxies are supported")
	}
	c.proxy = *uri
	c.useProxy = true

	return nil
}


/*
Sets the given context to the client, it can be used to modify but that is not recommended
*/
func (c *Client) SetContext(ctx context.Context) {
	c.ctx = ctx
}

/*
This method is meant to check if proxy use is active
*/
func (c *Client) UseProxy() bool {
	return c.useProxy
}

/*
Modifies the proxy, if this is called after Initialize it has no effect
*/
func (c *Client) ModifyProxy(proxy string) error {
	err := c.SetProxy(proxy); if err != nil {
		return err
	}

	return nil
}

/*
This function can be used to initialize a new client with a given context
*/
func WithContext(ctx context.Context) *Client {
	client := &Client{
		ctx: ctx,
	}
	return client
}

/*
Initializes a new client with zero values and a dummy context

Use ``WithContext`` instead to use your own context
*/
func NewClient() *Client {
	client := &Client{
		ctx: context.Background(),
	}

	return client
}

// Initialize, must be called before starting fetching words
func (c *Client) Initialize(){
	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	if c.useProxy {
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(&c.proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	c.initialized = true
	c.httpClient = httpClient // copies into c.httpClient
}	

/*
This is the main scraper function
*/
func(c *Client) fetchWords() (words , error){
	if !c.initialized {
		return make([]string,0), errors.New("You should call the initializing method before executing this one")
	}

	var readCloser io.ReadCloser
	var data interface{}

	req, err := http.NewRequest(
		"GET", 
		fmt.Sprintf(
			"%sqty=%d&category=es&first_letter=&last_letter=&word_size_by=length&operator=equals&length=",
			uri,
			c.quantity,
		),
		readCloser,
	); if err != nil {
		return make([]string,0), errors.New(
			fmt.Sprintf("Could not build request, failed --> %s", 
			err.Error()),
		)
	}

	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := c.httpClient.Do(req); if err != nil {
		return make([]string,0), errors.New(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		dataB, err := io.ReadAll(resp.Body); if err != nil {
			return make([]string,0), errors.New(err.Error())
		}
		data = strings.Split(
			strings.Trim(strings.ReplaceAll(
				strings.ReplaceAll(
					strings.Replace(
						strings.Replace(
							string(dataB), "[", "",1) , "]", "",1),"\"",""),
							"\r",""),"\r\n"),
							 ",",
		)
		return data.(words), nil
	}

	return make([]string,0), errors.New(
		fmt.Sprintf("Server responded with bad code %d", resp.StatusCode),
	)
}

/*
This method can be used to generate random words


If a context with a timeout or cancellation is provided, when the cancellation or timeout occurs, it will immediately cancel and return.
*/
func (c *Client) GenerateWords() (words, error) {

	completed := make(chan int)
	var results words
	var Error error

	go func(done chan <-int){
		defer func() {
			done <- 1
		}()

		words, err := c.fetchWords(); 

		results = words
		Error = err

		return 
	}(completed)

	select {
		case <-c.ctx.Done():
			return make([]string,0), c.ctx.Err()

		case <-completed:
			if Error != nil {
				return make(words, 0), Error
			}

			return results, nil
	}
}