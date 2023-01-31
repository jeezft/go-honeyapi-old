package wb

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/corpix/uarand"
	"github.com/go-resty/resty/v2"
)

type BrandInfo struct {
	State   int    `json:"state"`
	Version int    `json:"version"`
	Params  Params `json:"params"`
	Data    Data   `json:"data"`
}

type Params struct {
	Curr    string `json:"curr"`
	Spp     int    `json:"spp"`
	Version int    `json:"version"`
}

type Colors struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Sizes struct {
	Name     string `json:"name"`
	OrigName string `json:"origName"`
	Rank     int    `json:"rank"`
	OptionID int    `json:"optionId"`
	Wh       int    `json:"wh"`
	Sign     string `json:"sign"`
}

type Products struct {
	Sort            int      `json:"__sort"`
	Ksort           int      `json:"ksort"`
	Time1           int      `json:"time1"`
	Time2           int      `json:"time2"`
	Dist            int      `json:"dist"`
	ID              int      `json:"id"`
	Root            int      `json:"root"`
	KindID          int      `json:"kindId"`
	SubjectID       int      `json:"subjectId"`
	SubjectParentID int      `json:"subjectParentId"`
	Name            string   `json:"name"`
	Brand           string   `json:"brand"`
	BrandID         int      `json:"brandId"`
	SiteBrandID     int      `json:"siteBrandId"`
	SupplierID      int      `json:"supplierId"`
	Sale            int      `json:"sale"`
	PriceU          int      `json:"priceU"`
	SalePriceU      int      `json:"salePriceU"`
	LogisticsCost   int      `json:"logisticsCost"`
	Pics            int      `json:"pics"`
	Rating          int      `json:"rating"`
	Feedbacks       int      `json:"feedbacks"`
	Volume          int      `json:"volume"`
	Colors          []Colors `json:"colors"`
	Sizes           []Sizes  `json:"sizes"`
	DiffPrice       bool     `json:"diffPrice"`
}

type Data struct {
	Products []Products `json:"products"`
}

func (a *Api) GetBrandList(letter string) (*BrandlistO, error) {
	c := resty.New()

	resp, err := c.R().SetHeader("User-Agent", uarand.GetRandom()).Post("https://www.wildberries.ru/webapi/wildberries/brandlist/data?letter=" + letter)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		fmt.Println(resp.StatusCode())
		return nil, ErrNotfound
	}

	var o BrandlistO

	if err = json.Unmarshal(resp.Body(), &o); err != nil {
		return nil, err
	}

	return &o, nil
}

func (a *Api) GetBrand(id int) (*BrandInfo, error) {
	c := resty.New()

	resp, err := c.R().SetHeader("User-Agent", uarand.GetRandom()).Get("https://catalog.wb.ru/brands/%D0%B0/catalog?appType=1&brand=" + strconv.Itoa(id) + "&couponsGeo=2,12,7,3,6,18,21&curr=rub&dest=-1029256,-85617,-543140,-1586361&emp=0&lang=ru&locale=ru&pricemarginCoeff=1.0&reg=0&regions=80,64,83,4,38,33,70,68,69,86,30,40,48,1,66,31,22&sort=popular&spp=0")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, ErrNotfound
	}

	var brand BrandInfo

	if err = json.Unmarshal(resp.Body(), &brand); err != nil {
		return nil, err
	}

	return &brand, err
}
