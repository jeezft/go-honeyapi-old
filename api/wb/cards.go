package wb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/corpix/uarand"
	"github.com/go-resty/resty/v2"
)

type Card struct {
	Name       string
	ID         int
	SupplierID int
	Price      int
	CPM        int
	Rating     int
}

type Cpm struct {
	ID  int
	CPM int
}

type CpmQuery struct {
	Pages []struct {
		Positions []int `json:"positions"`
		Page      int   `json:"page"`
		Count     int   `json:"count"`
	} `json:"pages"`
	PrioritySubjects []int     `json:"prioritySubjects"`
	Adverts          []Adverts `json:"adverts"`
	MinCPM           int       `json:"minCPM"`
}

type Adverts struct {
	Code     string `json:"code"`
	AdvertID int    `json:"advertId"`
	ID       int    `json:"id"`
	Cpm      int    `json:"cpm"`
	Subject  int    `json:"subject"`
}

type SearchRes struct {
	Metadata struct {
		Name         string `json:"name"`
		CatalogType  string `json:"catalog_type"`
		CatalogValue string `json:"catalog_value"`
	} `json:"metadata"`
	State   int `json:"state"`
	Version int `json:"version"`
	Params  struct {
		Curr    string `json:"curr"`
		Spp     int    `json:"spp"`
		Version int    `json:"version"`
	} `json:"params"`
	Data struct {
		Products []struct {
			Sort            int    `json:"__sort"`
			Ksort           int    `json:"ksort"`
			Time1           int    `json:"time1"`
			Time2           int    `json:"time2"`
			Dist            int    `json:"dist"`
			ID              int    `json:"id"`
			Root            int    `json:"root"`
			KindID          int    `json:"kindId"`
			SubjectID       int    `json:"subjectId"`
			SubjectParentID int    `json:"subjectParentId"`
			Name            string `json:"name"`
			Brand           string `json:"brand"`
			BrandID         int    `json:"brandId"`
			SiteBrandID     int    `json:"siteBrandId"`
			SupplierID      int    `json:"supplierId"`
			Sale            int    `json:"sale"`
			PriceU          int    `json:"priceU"`
			SalePriceU      int    `json:"salePriceU"`
			Pics            int    `json:"pics"`
			Rating          int    `json:"rating"`
			Feedbacks       int    `json:"feedbacks"`
			Volume          int    `json:"volume"`
			Colors          []struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			} `json:"colors"`
			Sizes []struct {
				Name     string `json:"name"`
				OrigName string `json:"origName"`
				Rank     int    `json:"rank"`
				OptionID int    `json:"optionId"`
				Wh       int    `json:"wh"`
				Sign     string `json:"sign"`
			} `json:"sizes"`
			DiffPrice    bool   `json:"diffPrice"`
			PanelPromoID int    `json:"panelPromoId,omitempty"`
			PromoTextCat string `json:"promoTextCat,omitempty"`
			IsNew        bool   `json:"isNew,omitempty"`
		} `json:"products"`
	} `json:"data"`
}

type BrandlistO struct {
	ResultState int   `json:"resultState"`
	Value       Value `json:"value"`
}
type Brand struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	LogoPath string `json:"logoPath"`
}
type Value struct {
	BrandsCount string  `json:"brandsCount"`
	BrandsList  []Brand `json:"brandsList"`
}

func (a *Api) Search(srpx string, page int) ([]*Card, error) {
	c := resty.New()
	req, err := c.R().Get("https://search.wb.ru/exactmatch/ru/common/v4/search?appType=1&couponsGeo=2,12,7,3,6,18,21&curr=rub&dest=-1029256,-85617,-543140,-1586361&emp=0&lang=ru&locale=ru&pricemarginCoeff=1.0&reg=0&regions=80,64,83,4,38,33,70,68,69,86,30,40,48,1,66,31,22&resultset=catalog&sort=popular&spp=0&suppressSpellcheck=false&page=" + strconv.Itoa(page) + "&query=" + url.QueryEscape(srpx))
	if err != nil {
		return nil, err
	}

	if req.StatusCode() != 200 {
		fmt.Println(req.StatusCode())
		return nil, errors.New(req.Status())
	}

	var rsp SearchRes

	if err = json.Unmarshal(req.Body(), &rsp); err != nil {
		return nil, err
	}

	var products []*Card

	cpm, err := a.GetCPM(srpx)
	// if err != nil {
	// 	return nil, errors.New("cannot retrieve CPMs")
	// }

	for _, itm := range rsp.Data.Products {
		card := &Card{
			Name:       itm.Name,
			ID:         itm.ID,
			SupplierID: itm.SupplierID,
			Price:      itm.SalePriceU,
			Rating:     itm.Rating,
		}
		for _, c := range cpm {
			if c.ID == itm.ID {
				card.CPM = c.CPM
			}
		}
		products = append(products, card)
	}

	return products, nil
}

func (a *Api) GetCPM(srpx string) ([]*Cpm, error) {
	c := resty.New()

	resp, err := c.R().SetHeader("User-Agent", uarand.GetRandom()).Get("https://catalog-ads.wildberries.ru/api/v5/search?keyword=" + url.QueryEscape(srpx))
	if err != nil {
		return nil, err
	}

	cpm := &CpmQuery{}

	if err = json.Unmarshal(resp.Body(), &cpm); err != nil {
		fmt.Println(string(resp.Body()))
		return nil, err
	}

	adverts := []*Cpm{}

	for _, el := range cpm.Adverts {
		adverts = append(adverts, &Cpm{
			ID:  el.ID,
			CPM: el.Cpm,
		})
	}
	return adverts, nil
}

// func (a *Api) GetCpm(keyword string) (*CpmQuery, error) {
// 	c := resty.New()
// 	resp, err := c.R().SetHeader("User-Agent", uarand.GetRandom()).Get("https://catalog-ads.wildberries.ru/api/v5/search?keyword=" + url.QueryEscape(keyword))
// 	if err != nil {
// 		return nil, err
// 	}

// 	cpm := &CpmQuery{}
// 	if err = json.Unmarshal(resp.Body(), &cpm); err != nil {
// 		return nil, err
// 	}

// 	return cpm, nil
// }

// func (a *Api)

func volhostv2(vol int) string {
	if vol <= 143 {
		return "basket-01.wb.ru"
	} else if vol <= 287 {
		return "basket-02.wb.ru"
	} else if vol <= 431 {
		return "basket-03.wb.ru"
	} else if vol <= 719 {
		return "basket-04.wb.ru"
	} else if vol <= 1007 {
		return "basket-05.wb.ru"
	} else if vol <= 1061 {
		return "basket-06.wb.ru"
	} else if vol <= 1115 {
		return "basket-07.wb.ru"
	} else if vol <= 1169 {
		return "basket-08.wb.ru"
	} else if vol <= 1313 {
		return "basket-09.wb.ru"
	} else if vol <= 1601 {
		return "basket-10.wb.ru"
	} else {
		return "basket-11.wb.ru"
	}
}

func getbasketstatic(nmID int) string {
	vol := (nmID / 100000)
	part := (nmID / 1000)
	host := volhostv2(vol)

	return fmt.Sprintf("https://%s/vol%d/part%d/%d/", host, vol, part, nmID)
}
