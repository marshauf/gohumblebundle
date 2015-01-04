package humblebundle

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SortBestselling  = "bestselling"
	SortNewest       = "newest"
	SortDiscount     = "discount"
	SortAlphabetical = "alphabetical"

	PlatformAll     = ""
	PlatformAndroid = "android"
	PlatformLinux   = "linux"
	PlatformMac     = "mac"
	PlatformWindows = "windows"

	DrmAll   = ""
	DrmFree  = "download"
	DrmSteam = "steam"
	DrmUplay = "uplay"
)

type Response struct {
	NumResults int        `bson:"num_results" json:"num_results"`
	RequestID  int        `bson:"request" json:"request"`
	Results    []*Product `bson:"results" json:"results"`
}

func (p0 *Product) Equal(p1 *Product) bool {
	if p0.MachineName != p1.MachineName {
		return false
	}
	for key, value := range p0.IconDict {
		v, ok := p1.IconDict[key]
		if !ok {
			return false
		}
		if value.Equal(v) == false {
			return false
		}
	}

	// TODO AlertMessages to StoreFrontPreviewImage

	if p0.HumanName != p1.HumanName {
		return false
	}
	if p0.CurrentPrice.Equal(p1.CurrentPrice) == false {
		return false
	}
	if p0.SaleEnd != nil {
		if p1.SaleEnd != nil {
			if p0.SaleEnd.Equal(p1.SaleEnd) == false {
				return false
			}
		} else {
			return false
		}
	} else if p1.SaleEnd != nil {
		return false
	}
	if p0.SaleType != p1.SaleType {
		return false
	}
	if p0.FullPrice != nil {
		if p1.FullPrice != nil {
			if p0.FullPrice.Equal(p1.FullPrice) == false {
				return false
			}
		} else {
			return false
		}
	} else if p1.FullPrice != nil {
		return false
	}
	return true
}

type Time struct {
	time.Time
}

func (t *Time) Equal(u *Time) bool {
	return t.Time.Equal(u.Time)
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	if strings.HasSuffix(str, ".0") {
		str = str[:len(data)-2]
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(n, 0)
	return nil
}

type Product struct {
	MachineName                  string              `bson:"machine_name" json:"machine_name"`
	IconDict                     map[string]*Icon    `bson:"icon_dict" json:"icon_dict"`
	AlertMessages                []map[string]string `bson:"alert_messages" json:"alert_messages"`
	StoreFrontFeaturedImageSmall string              `bson:"storefront_featured_image_small" json:"storefront_featured_image_small"`
	YoutubeLink                  string              `bson:"youtube_link" json:"youtube_link"`
	Platforms                    []string            `bson:"platforms" json:"platforms"`
	PromotionalMessage           interface{}         `bson:"promotional_message" json:"promotional_message"` // TODO Lookup
	UskRating                    string              `bson:"usk_rating" json:"usk_rating"`
	ForcePopup                   bool                `bson:"force_popup" json:"force_popup"`
	RatingDetails                interface{}         `bson:"rating_details" json:"rating_details"` // TODO Lookup
	EsrbRating                   string              `bson:"esrb_rating" json:"esrb_rating"`
	Developers                   []*Developer        `bson:"developers" json:"developers"`
	Publishers                   interface{}         `bson:"publishers" json:"publishers"` // TODO Lookup
	DeliveryMethods              []string            `bson:"delivery_methods" json:"delivery_methods"`
	StoreFrontIcon               string              `bson:"storefront_icon" json:"storefront_icon"`
	Description                  string              `bson:"description" json:"description"`
	AllowedTerritories           interface{}         `bson:"allowed_territories" json:"allowed_territories"` // TODO Lookup
	MinimumAge                   interface{}         `bson:"minimum_age" json:"minimum_age"`                 // TODO Lookup
	SystemRequirements           string              `bson:"system_requirements" json:"system_requirements"`
	PegiRating                   string              `bson:"pegi_rating" json:"pegi_rating"`
	StoreFrontFeaturedImageLarge string              `bson:"storefront_featured_image_large" json:"storefront_featured_image_large"`
	ContentTypes                 []string            `bson:"content_types" json:"content_types"`
	StoreFrontPreviewImage       interface{}         `bson:"storefront_preview_image" json:"storefront_preview_image"` // TODO Lookup
	HumanName                    string              `bson:"human_name" json:"human_name"`
	CurrentPrice                 Prices              `bson:"current_price" json:"current_price"` // value float currency string
	SaleEnd                      *Time               `bson:"sale_end" json:"sale_end,number"`
	SaleType                     string              `bson:"sale_type" json:"sale_type"`
	FullPrice                    Prices              `bson:"full_price" json:"full_price"` // value float currency string
}

type Icon struct {
	Available   []string `bson:"available" json:"available"`
	Unavailable []string `bson:"unavailable" json:"unavailable"`
}

func (i0 *Icon) Equal(i1 *Icon) bool {
	if len(i0.Available) != len(i1.Available) {
		return false
	}
	if len(i0.Unavailable) != len(i1.Unavailable) {
		return false
	}
	for i := range i0.Available {
		if i0.Available[i] != i1.Available[i] {
			return false
		}
	}
	for i := range i0.Unavailable {
		if i0.Unavailable[i] != i1.Unavailable[i] {
			return false
		}
	}
	return true
}

type Developer struct {
	Name string `bson:"developer_name" json:"developer-name"`
	URL  string `bson:"developer_url" json:"developer-url"`
}

type Prices []Price

func (p0 Prices) Equal(p1 Prices) bool {
	var found bool
	for i := range p0 {
		found = false
		for j := range p1 {
			if p0[i].Equal(p1[j]) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (p *Prices) UnmarshalJSON(data []byte) error {
	list := make([]interface{}, 0)
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	*p = make([]Price, len(list)/2)
	for i := 0; i < len(list)/2; i++ {
		switch v := list[i*2].(type) {
		case float64:
			(*p)[i].Value = v
		case string:
			(*p)[i].Currency = v
		}
		switch v := list[i*2+1].(type) {
		case float64:
			(*p)[i].Value = v
		case string:
			(*p)[i].Currency = v
		}
	}
	return nil
}

type Price struct {
	Currency string  `bson:"currency"`
	Value    float64 `bson:"value"`
}

func (p0 Price) Equal(p1 Price) bool {
	return p0.Currency == p1.Currency && p0.Value == p1.Value
}

func Request(requestID, pageSize, page int, sort, platform, drm, search string) (*Response, error) {
	/*
		https://www.humblebundle.com/store/api/humblebundle?request=1&page_size=20&sort=bestselling&page=0
		https://www.humblebundle.com/store/api/humblebundle?request=1&page_size=5&sort=bestselling&page=0
		https://www.humblebundle.com/store/api/humblebundle?request=2&page_size=20&sort=bestselling&page=1
		https://www.humblebundle.com/store/api/humblebundle?request=3&page_size=20&sort=bestselling&page=0&platform=windows
		https://www.humblebundle.com/store/api/humblebundle?request=4&page_size=20&sort=bestselling&page=0&platform=windows&drm=download
		https://www.humblebundle.com/store/api/humblebundle?request=5&page_size=20&sort=newest&page=0&platform=windows&drm=download
		https://www.humblebundle.com/store/api/humblebundle?request=6&page_size=20&sort=discount&page=0&platform=windows&drm=download
		https://www.humblebundle.com/store/api/humblebundle?request=7&page_size=20&sort=alphabetical&page=0&platform=windows&drm=download
	*/
	url := "https://www.humblebundle.com/store/api/humblebundle?request=" + strconv.Itoa(requestID)
	if pageSize > 0 {
		url += "&page_size=" + strconv.Itoa(pageSize)
	}
	if page >= 0 {
		url += "&page=" + strconv.Itoa(page)
	}
	if len(sort) > 0 {
		url += "&sort=" + sort
	}
	if len(platform) > 0 {
		url += "&platform=" + platform
	}
	if len(drm) > 0 {
		url += "&drm=" + drm
	}
	if len(search) > 0 {
		url += "&search=" + search
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response Response
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	return &response, err
}
