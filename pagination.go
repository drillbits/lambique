// Copyright 2019 drillbits
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lambique

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type PagingLink struct {
	Rel string
	URL *url.URL
}

func (link *PagingLink) String() string {
	return fmt.Sprintf(`<%s>; rel="%s"`, link.URL.String(), link.Rel)
}

type PagingLinks []*PagingLink

func (ls PagingLinks) String() string {
	var links []string
	for _, link := range ls {
		links = append(links, link.String())
	}
	return strings.Join(links, ",")
}

type Pagination interface {
	ParseURL(u *url.URL) error
	FirstPagingLink(u *url.URL) (*PagingLink, error)
	PrevPagingLink(u *url.URL) (*PagingLink, error)
	NextPagingLink(u *url.URL) (*PagingLink, error)
	LastPagingLink(u *url.URL) (*PagingLink, error)
	PagingLinks(u *url.URL) (PagingLinks, error)
}

type PageNumberPagination struct {
	PageNumberKey string
	PageNumber    int
	parsed        bool
}

func NewPageNumberPagination() *PageNumberPagination {
	return &PageNumberPagination{
		PageNumberKey: "page",
		PageNumber:    1,
	}
}

func (p *PageNumberPagination) ParseURL(u *url.URL) error {
	vs := u.Query()
	p.PageNumber = vsGetInt(vs, p.PageNumberKey, 1)
	p.parsed = true
	return nil
}

func (p *PageNumberPagination) FirstPagingLink(u *url.URL) (*PagingLink, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	vs := u.Query()
	vs.Set(p.PageNumberKey, "1")
	u = copyURL(u)
	u.RawQuery = vs.Encode()

	return &PagingLink{
		Rel: "first",
		URL: u,
	}, nil
}

func (p *PageNumberPagination) PrevPagingLink(u *url.URL) (*PagingLink, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	if p.PageNumber < 2 {
		return nil, nil
	}

	vs := u.Query()
	vs.Set(p.PageNumberKey, strconv.Itoa(p.PageNumber-1))
	u = copyURL(u)
	u.RawQuery = vs.Encode()

	return &PagingLink{
		Rel: "prev",
		URL: u,
	}, nil
}

func (p *PageNumberPagination) NextPagingLink(u *url.URL) (*PagingLink, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	vs := u.Query()
	vs.Set(p.PageNumberKey, strconv.Itoa(p.PageNumber+1))
	u = copyURL(u)
	u.RawQuery = vs.Encode()

	return &PagingLink{
		Rel: "next",
		URL: u,
	}, nil
}

func (p *PageNumberPagination) LastPagingLink(u *url.URL) (*PagingLink, error) {
	// TODO: implement
	return nil, nil
}

func (p *PageNumberPagination) PagingLinks(u *url.URL) (PagingLinks, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	var links PagingLinks

	firstPage, err := p.FirstPagingLink(u)
	if err != nil {
		return nil, err
	} else if firstPage != nil {
		links = append(links, firstPage)
	}

	prevPage, err := p.PrevPagingLink(u)
	if err != nil {
		return nil, err
	} else if prevPage != nil {
		links = append(links, prevPage)
	}

	nextPage, err := p.NextPagingLink(u)
	if err != nil {
		return nil, err
	} else if nextPage != nil {
		links = append(links, nextPage)
	}

	// TODO: lastPage

	return links, nil
}

type OffsetLimitPagination struct {
	OffsetKey string
	LimitKey  string
	PageSize  int
	Offset    int
	Limit     int
	parsed    bool
}

func NewOffsetLimitPagination() *OffsetLimitPagination {
	pageSize := 50
	return &OffsetLimitPagination{
		OffsetKey: "offset",
		LimitKey:  "limit",
		PageSize:  pageSize,
		Offset:    0,
		Limit:     pageSize,
	}
}

func (p *OffsetLimitPagination) ParseURL(u *url.URL) error {
	vs := u.Query()
	p.Offset = vsGetInt(vs, p.OffsetKey, 0)
	p.Limit = vsGetInt(vs, p.LimitKey, p.PageSize)
	p.parsed = true
	return nil
}

func (p *OffsetLimitPagination) FirstPagingLink(u *url.URL) (*PagingLink, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	vs := u.Query()
	vs.Set(p.OffsetKey, "0")
	vs.Set(p.LimitKey, strconv.Itoa(p.PageSize))
	u = copyURL(u)
	u.RawQuery = vs.Encode()

	return &PagingLink{
		Rel: "first",
		URL: u,
	}, nil
}

func (p *OffsetLimitPagination) PrevPagingLink(u *url.URL) (*PagingLink, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	if p.Offset < p.PageSize {
		return nil, nil
	}

	vs := u.Query()
	vs.Set(p.OffsetKey, strconv.Itoa(p.Offset-p.PageSize))
	vs.Set(p.LimitKey, strconv.Itoa(p.PageSize))
	u = copyURL(u)
	u.RawQuery = vs.Encode()

	return &PagingLink{
		Rel: "prev",
		URL: u,
	}, nil
}

func (p *OffsetLimitPagination) NextPagingLink(u *url.URL) (*PagingLink, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	vs := u.Query()
	vs.Set(p.OffsetKey, strconv.Itoa(p.Offset+p.PageSize))
	vs.Set(p.LimitKey, strconv.Itoa(p.PageSize))
	u = copyURL(u)
	u.RawQuery = vs.Encode()

	return &PagingLink{
		Rel: "next",
		URL: u,
	}, nil
}

func (p *OffsetLimitPagination) LastPagingLink(u *url.URL) (*PagingLink, error) {
	// TODO: implement
	return nil, nil
}

func (p *OffsetLimitPagination) PagingLinks(u *url.URL) (PagingLinks, error) {
	if !p.parsed {
		err := p.ParseURL(u)
		if err != nil {
			return nil, err
		}
	}

	var links PagingLinks

	firstPage, err := p.FirstPagingLink(u)
	if err != nil {
		return nil, err
	} else if firstPage != nil {
		links = append(links, firstPage)
	}

	prevPage, err := p.PrevPagingLink(u)
	if err != nil {
		return nil, err
	} else if prevPage != nil {
		links = append(links, prevPage)
	}

	nextPage, err := p.NextPagingLink(u)
	if err != nil {
		return nil, err
	} else if nextPage != nil {
		links = append(links, nextPage)
	}

	// TODO: lastPage

	return links, nil
}

func vsGetInt(vs url.Values, key string, defaultValue int) int {
	s := vs.Get(key)
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return i
}

func copyURL(u *url.URL) *url.URL {
	copy := *u
	return &copy
}
