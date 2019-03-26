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
	"reflect"
	"testing"
)

var (
	_ Pagination = (*OffsetLimitPagination)(nil)
	_ Pagination = (*PageNumberPagination)(nil)
)

func TestPagingLink_String(t *testing.T) {
	tests := []struct {
		name   string
		fields PagingLink
		want   string
	}{
		{
			name: "format",
			fields: PagingLink{
				Rel: "prev",
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=1"),
			},
			want: "<https://www.example.com/foo?bar+baz%3Aqux&page=1>; rel=\"prev\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := &PagingLink{
				Rel: tt.fields.Rel,
				URL: tt.fields.URL,
			}
			if got := link.String(); got != tt.want {
				t.Errorf("PagingLink.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagingLinks_String(t *testing.T) {
	tests := []struct {
		name string
		ls   PagingLinks
		want string
	}{
		{
			name: "join",
			ls: PagingLinks{
				&PagingLink{Rel: "prev", URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=1")},
				&PagingLink{Rel: "next", URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3")},
			},
			want: "<https://www.example.com/foo?bar+baz%3Aqux&page=1>; rel=\"prev\",<https://www.example.com/foo?bar+baz%3Aqux&page=3>; rel=\"next\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ls.String(); got != tt.want {
				t.Errorf("PagingLinks.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPageNumberPagination(t *testing.T) {
	tests := []struct {
		name string
		want *PageNumberPagination
	}{
		{
			name: "new",
			want: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    1,
				parsed:        false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPageNumberPagination(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPageNumberPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageNumberPagination_ParseURL(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *PageNumberPagination
		args    args
		want    *PageNumberPagination
		wantErr bool
	}{
		{
			name: "parse",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    1,
				parsed:        false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=2"),
			},
			want: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    2,
				parsed:        true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PageNumberPagination{
				PageNumberKey: tt.fields.PageNumberKey,
				PageNumber:    tt.fields.PageNumber,
				parsed:        tt.fields.parsed,
			}
			err := p.ParseURL(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageNumberPagination.ParseURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("PageNumberPagination.ParseURL() = %#v, want %#v", p, tt.want)
			}
		})
	}
}

func TestPageNumberPagination_FirstPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *PageNumberPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    1,
				parsed:        false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: &PagingLink{
				Rel: "first",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=1"),
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    3,
				parsed:        true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: &PagingLink{
				Rel: "first",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=1"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PageNumberPagination{
				PageNumberKey: tt.fields.PageNumberKey,
				PageNumber:    tt.fields.PageNumber,
				parsed:        tt.fields.parsed,
			}
			got, err := p.FirstPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageNumberPagination.FirstPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PageNumberPagination.FirstPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageNumberPagination_PrevPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *PageNumberPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    1,
				parsed:        false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: &PagingLink{
				Rel: "prev",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=2"),
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    3,
				parsed:        true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: &PagingLink{
				Rel: "prev",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=2"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PageNumberPagination{
				PageNumberKey: tt.fields.PageNumberKey,
				PageNumber:    tt.fields.PageNumber,
				parsed:        tt.fields.parsed,
			}
			got, err := p.PrevPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageNumberPagination.PrevPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PageNumberPagination.PrevPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageNumberPagination_NextPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *PageNumberPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    1,
				parsed:        false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: &PagingLink{
				Rel: "next",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=4"),
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    3,
				parsed:        true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: &PagingLink{
				Rel: "next",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=4"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PageNumberPagination{
				PageNumberKey: tt.fields.PageNumberKey,
				PageNumber:    tt.fields.PageNumber,
				parsed:        tt.fields.parsed,
			}
			got, err := p.NextPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageNumberPagination.NextPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PageNumberPagination.NextPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageNumberPagination_LastPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *PageNumberPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    1,
				parsed:        false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			// TODO: implement
			want:    nil,
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    3,
				parsed:        true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			// TODO: implement
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PageNumberPagination{
				PageNumberKey: tt.fields.PageNumberKey,
				PageNumber:    tt.fields.PageNumber,
				parsed:        tt.fields.parsed,
			}
			got, err := p.LastPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageNumberPagination.LastPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PageNumberPagination.LastPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageNumberPagination_PagingLinks(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *PageNumberPagination
		args    args
		want    PagingLinks
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    1,
				parsed:        false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: PagingLinks{
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				&PagingLink{
					Rel: "first",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=1"),
				},
				&PagingLink{
					Rel: "prev",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=2"),
				},
				&PagingLink{
					Rel: "next",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=4"),
				},
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &PageNumberPagination{
				PageNumberKey: "page",
				PageNumber:    3,
				parsed:        true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&page=3"),
			},
			want: PagingLinks{
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				&PagingLink{
					Rel: "first",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=1"),
				},
				&PagingLink{
					Rel: "prev",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=2"),
				},
				&PagingLink{
					Rel: "next",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&page=4"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PageNumberPagination{
				PageNumberKey: tt.fields.PageNumberKey,
				PageNumber:    tt.fields.PageNumber,
				parsed:        tt.fields.parsed,
			}
			got, err := p.PagingLinks(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("PageNumberPagination.PagingLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PageNumberPagination.PagingLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOffsetLimitPagination(t *testing.T) {
	tests := []struct {
		name string
		want *OffsetLimitPagination
	}{
		{
			name: "new",
			want: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    0,
				Limit:     50,
				parsed:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOffsetLimitPagination(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOffsetLimitPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetLimitPagination_ParseURL(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *OffsetLimitPagination
		args    args
		want    *OffsetLimitPagination
		wantErr bool
	}{
		{
			name: "parse",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    0,
				Limit:     50,
				parsed:    false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&offset=100&limit=50"),
			},
			want: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    100,
				Limit:     50,
				parsed:    true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OffsetLimitPagination{
				OffsetKey: tt.fields.OffsetKey,
				LimitKey:  tt.fields.LimitKey,
				PageSize:  tt.fields.PageSize,
				Offset:    tt.fields.Offset,
				Limit:     tt.fields.Limit,
				parsed:    tt.fields.parsed,
			}
			err := p.ParseURL(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("OffsetLimitPagination.ParseURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("OffsetLimitPagination.ParseURL() = %#v, want %#v", p, tt.want)
			}
		})
	}
}

func TestOffsetLimitPagination_FirstPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *OffsetLimitPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    0,
				Limit:     50,
				parsed:    false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: &PagingLink{
				Rel: "first",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=0"),
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    100,
				Limit:     50,
				parsed:    true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: &PagingLink{
				Rel: "first",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=0"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OffsetLimitPagination{
				OffsetKey: tt.fields.OffsetKey,
				LimitKey:  tt.fields.LimitKey,
				PageSize:  tt.fields.PageSize,
				Offset:    tt.fields.Offset,
				Limit:     tt.fields.Limit,
				parsed:    tt.fields.parsed,
			}
			got, err := p.FirstPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("OffsetLimitPagination.FirstPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OffsetLimitPagination.FirstPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetLimitPagination_PrevPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *OffsetLimitPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    0,
				Limit:     50,
				parsed:    false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: &PagingLink{
				Rel: "prev",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=50"),
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    100,
				Limit:     50,
				parsed:    true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: &PagingLink{
				Rel: "prev",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=50"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OffsetLimitPagination{
				OffsetKey: tt.fields.OffsetKey,
				LimitKey:  tt.fields.LimitKey,
				PageSize:  tt.fields.PageSize,
				Offset:    tt.fields.Offset,
				Limit:     tt.fields.Limit,
				parsed:    tt.fields.parsed,
			}
			got, err := p.PrevPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("OffsetLimitPagination.PrevPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OffsetLimitPagination.PrevPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetLimitPagination_NextPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *OffsetLimitPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    0,
				Limit:     50,
				parsed:    false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: &PagingLink{
				Rel: "next",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=150"),
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    100,
				Limit:     50,
				parsed:    true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: &PagingLink{
				Rel: "next",
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=150"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OffsetLimitPagination{
				OffsetKey: tt.fields.OffsetKey,
				LimitKey:  tt.fields.LimitKey,
				PageSize:  tt.fields.PageSize,
				Offset:    tt.fields.Offset,
				Limit:     tt.fields.Limit,
				parsed:    tt.fields.parsed,
			}
			got, err := p.NextPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("OffsetLimitPagination.NextPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OffsetLimitPagination.NextPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetLimitPagination_LastPagingLink(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *OffsetLimitPagination
		args    args
		want    *PagingLink
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    0,
				Limit:     50,
				parsed:    false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			// TODO: implement
			want:    nil,
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    100,
				Limit:     50,
				parsed:    true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			// TODO: implement
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OffsetLimitPagination{
				OffsetKey: tt.fields.OffsetKey,
				LimitKey:  tt.fields.LimitKey,
				PageSize:  tt.fields.PageSize,
				Offset:    tt.fields.Offset,
				Limit:     tt.fields.Limit,
				parsed:    tt.fields.parsed,
			}
			got, err := p.LastPagingLink(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("OffsetLimitPagination.LastPagingLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OffsetLimitPagination.LastPagingLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetLimitPagination_PagingLinks(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		fields  *OffsetLimitPagination
		args    args
		want    PagingLinks
		wantErr bool
	}{
		{
			name: "not parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    0,
				Limit:     50,
				parsed:    false,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: PagingLinks{
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				&PagingLink{
					Rel: "first",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=0"),
				},
				&PagingLink{
					Rel: "prev",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=50"),
				},
				&PagingLink{
					Rel: "next",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=150"),
				},
			},
			wantErr: false,
		},
		{
			name: "parsed",
			fields: &OffsetLimitPagination{
				OffsetKey: "offset",
				LimitKey:  "limit",
				PageSize:  50,
				Offset:    100,
				Limit:     50,
				parsed:    true,
			},
			args: args{
				mustParseURL("https://www.example.com/foo?bar+baz%3Aqux&limit=50&offset=100"),
			},
			want: PagingLinks{
				// net/url: url.RawQuery = url.Query().Encode() can change the URL
				// `?foo&bar=baz` -> `?foo=&bar=baz`
				// https://github.com/golang/go/issues/16460
				&PagingLink{
					Rel: "first",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=0"),
				},
				&PagingLink{
					Rel: "prev",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=50"),
				},
				&PagingLink{
					Rel: "next",
					URL: mustParseURL("https://www.example.com/foo?bar+baz%3Aqux=&limit=50&offset=150"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &OffsetLimitPagination{
				OffsetKey: tt.fields.OffsetKey,
				LimitKey:  tt.fields.LimitKey,
				PageSize:  tt.fields.PageSize,
				Offset:    tt.fields.Offset,
				Limit:     tt.fields.Limit,
				parsed:    tt.fields.parsed,
			}
			got, err := p.PagingLinks(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("OffsetLimitPagination.PagingLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OffsetLimitPagination.PagingLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_vsGetInt(t *testing.T) {
	type args struct {
		vs           url.Values
		key          string
		defaultValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "int value",
			args: args{
				vs: url.Values{
					"foo": []string{"13"},
				},
				key:          "foo",
				defaultValue: 1,
			},
			want: 13,
		},
		{
			name: "string value",
			args: args{
				vs: url.Values{
					"foo": []string{"str"},
				},
				key:          "foo",
				defaultValue: 1,
			},
			want: 1,
		},
		{
			name: "nil value",
			args: args{
				vs: url.Values{
					"foo": nil,
				},
				key:          "foo",
				defaultValue: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vsGetInt(tt.args.vs, tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("vsGetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyURL(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name string
		args args
		want *url.URL
	}{
		{
			name: "copied",
			args: args{
				mustParseURL("https://www.example.com/"),
			},
			want: mustParseURL("https://www.example.com/"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := copyURL(tt.args.u)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("copyURL() = %v, want %v", got, tt.want)
			}
			if fmt.Sprintf("%p", got) == fmt.Sprintf("%p", tt.args.u) {
				t.Errorf("copyURL() pointer = %p, want %p", got, tt.args.u)
			}
		})
	}
}

func mustParseURL(rawurl string) *url.URL {
	u, _ := url.Parse(rawurl)
	return u
}
