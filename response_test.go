package humblebundle

import (
	"testing"
)

var tests = []struct {
	RequestID int
	PageSize  int
	Page      int
	Sort      string
	Platform  string
	Drm       string
	Search    string
}{
	{
		RequestID: 0,
		PageSize:  20,
		Page:      0,
		Sort:      SortBestselling,
	},
	{
		RequestID: 1,
		PageSize:  5,
		Page:      0,
		Sort:      SortBestselling,
	},
	{
		RequestID: 2,
		PageSize:  20,
		Page:      1,
		Sort:      SortBestselling,
	},
	{
		RequestID: 3,
		PageSize:  20,
		Page:      0,
		Sort:      SortBestselling,
		Platform:  PlatformWindows,
	},
	{
		RequestID: 4,
		PageSize:  20,
		Page:      0,
		Sort:      SortBestselling,
		Platform:  PlatformWindows,
		Drm:       DrmFree,
	},
	{
		RequestID: 5,
		PageSize:  20,
		Page:      0,
		Sort:      SortNewest,
		Platform:  PlatformWindows,
		Drm:       DrmFree,
	},
	{
		RequestID: 6,
		PageSize:  20,
		Page:      0,
		Sort:      SortDiscount,
		Platform:  PlatformWindows,
		Drm:       DrmFree,
	},
	{
		RequestID: 7,
		PageSize:  20,
		Page:      0,
		Sort:      SortAlphabetical,
		Platform:  PlatformWindows,
		Drm:       DrmFree,
	},
	{
		RequestID: 8,
		PageSize:  20,
		Page:      0,
		Sort:      SortNewest,
		Platform:  PlatformWindows,
		Drm:       DrmSteam,
	},
	{
		RequestID: 9,
		PageSize:  20,
		Page:      0,
		Sort:      SortNewest,
		Platform:  PlatformLinux,
		Drm:       DrmFree,
	},
	{
		RequestID: 10,
		PageSize:  20,
		Page:      0,
		Sort:      SortNewest,
		Platform:  PlatformWindows,
		Drm:       DrmUplay,
	},
}

func TestRequest(t *testing.T) {
	var (
		resp *Response
		err  error
	)
	for _, req := range tests {
		resp, err = Request(req.RequestID, req.PageSize, req.Page, req.Sort, req.Platform, req.Drm, req.Search)
		if err != nil {
			t.Errorf("Request(%v): %s", req, err)
		}
		if resp.RequestID != req.RequestID {
			t.Errorf("RequestID %d different to ResponseID %d", req.RequestID, resp.RequestID)
		}
		if resp.NumResults <= 0 {
			t.Errorf("Number of results are %d", resp.NumResults)
		}
		if len(resp.Results) == 0 {
			t.Errorf("Response.Results is empty: %v", resp.Results)
		}
	}
}

func TestRequestAll(t *testing.T) {
	var (
		pageSize = 20 // The maximum of 07.12.2014
		pages    = 26 // as of 09.12.2014
		sort     = SortAlphabetical
		platform = PlatformAll
		drm      = DrmAll
		search   string
		err      error
		resp     *Response
	)
	resp, err = Request(0, pageSize, 0, sort, platform, drm, search)
	if err != nil {
		t.Errorf("Request(%d): %s", 0, err)
	}

	pages = resp.NumResults / pageSize

	for i := 1; i <= pages; i++ {
		resp, err = Request(i, pageSize, i, sort, platform, drm, search)
		if err != nil {
			t.Errorf("Request(%d): %s", i, err)
		}
		if resp.RequestID != i {
			t.Errorf("RequestID %d different to ResponseID %d", i, resp.RequestID)
		}
	}
}
