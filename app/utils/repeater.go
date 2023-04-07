package utils

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Token struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Repeater struct {
	ClientTimeout time.Duration
	Attempts      time.Duration
	URI           string
	Headers       http.Header
	Body          string
	Count         int
}

func (r *Repeater) Get() ([]byte, error) {
	var res []byte
	client := http.Client{
		Timeout: r.ClientTimeout * time.Second,
	}
	request, err := http.NewRequest("GET", r.URI, nil)
	if err != nil {
		log.Printf("[ERROR] cannot create GET request: %#v; URL: %s", err, r.URI)
		return nil, err
	}
	request.Header.Set("X-API-V2-MEDREG", "true")

	response, err := client.Do(request)
	if err != nil {
		log.Printf("[ERROR] can not make Get request: %#v", err)
		if errClose := response.Body.Close(); errClose != nil {
			log.Printf("[ERROR] can not close response body %#v", errClose)
		}
		sumTimeout := r.Attempts * time.Second
		ticker := time.NewTicker(sumTimeout)
		cancel := make(chan struct{})
		go func() {
			defer func() {
				cancel <- struct{}{}
			}()
			time.Sleep(10 * sumTimeout)
		}()
		for {
			select {
			case <-ticker.C:
				response, err = client.Get(r.URI)
				if err != nil {
					log.Printf("[ERROR] can not make Get request: %#v ", err)
					if errClose := response.Body.Close(); errClose != nil {
						log.Printf("[ERROR] can not close response body %#v", errClose)
					}
					continue
				}
				break
			case <-cancel:
				log.Printf("[WARN] completed repeater call. API is not reachible")
				break
			}
		}
	}
	if response == nil && err != nil {
		return nil, err
	}

	res, err = io.ReadAll(response.Body)
	if err != nil {
		log.Printf("[ERROR] can not read response body %#v", err)
		return nil, err
	}
	return res, nil
}
