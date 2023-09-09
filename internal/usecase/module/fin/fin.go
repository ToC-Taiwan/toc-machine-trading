// Package fin package fin
package fin

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Fin struct {
	c     *resty.Client
	token string
}

func NewFin() (*Fin, error) {
	client := resty.New()
	res, err := client.R().
		SetFormData(map[string]string{
			"user_id":  "chindada",
			"password": "ryxPof-sinqoj-bipse1",
		}).
		Post(loginURL)
	if err != nil {
		return nil, err
	}

	loginRes := make(map[string]interface{})
	err = json.Unmarshal(res.Body(), &loginRes)
	if err != nil {
		return nil, err
	}

	return &Fin{
		c:     client,
		token: loginRes["token"].(string),
	}, nil
}

func (f *Fin) GetTaiwanStockMarginPurchaseShortSale() error {
	res, err := f.c.R().
		SetQueryParams(map[string]string{
			"dataset":    "TaiwanStockMarginPurchaseShortSale",
			"data_id":    "2330",
			"start_date": "2023-09-06",
			"token":      f.token,
		}).
		Get(dataURL)
	if err != nil {
		return err
	}

	fmt.Println(string(res.Body()))
	return nil
}
