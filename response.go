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
	NumResults int        `json:"num_results"`
	RequestID  int        `json:"request"`
	Results    []*Product `json:"results"`
}

func (p0 *Product) Eq(p1 *Product) bool {
	if p0.MachineName != p1.MachineName {
		return false
	}
	for key, value := range p0.IconDict {
		v, ok := p1.IconDict[key]
		if !ok {
			return false
		}
		if value.Eq(v) == false {
			return false
		}
	}
	// TODO AlertMessages and the rest
	return true
}

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	if strings.HasSuffix(str, ".0") {
		str = str[:len(data)-2]
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(n, 0))
	return nil
}

type Product struct {
	MachineName                  string              `json:"machine_name"`
	IconDict                     map[string]*Icon    `json:"icon_dict"`
	AlertMessages                []map[string]string `json:"alert_messages"`
	StoreFrontFeaturedImageSmall string              `json:"storefront_featured_image_small"`
	YoutubeLink                  string              `json:"youtube_link"`
	Platforms                    []string            `json:"platforms"`
	PromotionalMessage           interface{}         `json:"promotional_message"` // TODO Lookup
	UskRating                    string              `json:"usk_rating"`
	ForcePopup                   bool                `json:"force_popup"`
	RatingDetails                interface{}         `json:"rating_details"` // TODO Lookup
	EsrbRating                   string              `json:"esrb_rating"`
	Developers                   []*Developer        `json:"developers"`
	Publishers                   interface{}         `json:"publishers"` // TODO Lookup
	DeliveryMethods              []string            `json:"delivery_methods"`
	StoreFrontIcon               string              `json:"storefront_icon"`
	Description                  string              `json:"description"`
	AllowedTerritories           interface{}         `json:"allowed_territories"` // TODO Lookup
	MinimumAge                   interface{}         `json:"minimum_age"`         // TODO Lookup
	SystemRequirements           string              `json:"system_requirements"`
	PegiRating                   string              `json:"pegi_rating"`
	StoreFrontFeaturedImageLarge string              `json:"storefront_featured_image_large"`
	ContentTypes                 []string            `json:"content_types"`
	StoreFrontPreviewImage       interface{}         `json:"storefront_preview_image"` // TODO Lookup
	HumanName                    string              `json:"human_name"`
	CurrentPrice                 *Currency           `json:"current_price"` // value float, currency string
	SaleEnd                      Time                `json:"sale_end,number"`
	SaleType                     string              `json:"sale_type"`
	FullPrice                    *Currency           `json:"full_price"` // value float, currency string
}

type Icon struct {
	Available   []string `json:"available"`
	Unavailable []string `json:"unavailable"`
}

func (i0 *Icon) Eq(i1 *Icon) bool {
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
	Name string `json:"developer-name"`
	URL  string `json:"developer-url"`
}

type Currency []interface{}

func (c Currency) Value() float64 {
	return c[1].(float64)
}

func (c Currency) Name() string {
	return c[0].(string)
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
