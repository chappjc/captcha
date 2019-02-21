// Copyright 2019 Jonathan Chappelow. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Benchmark_captchaHandler_ServeHTTP(b *testing.B) {
	id := NewLen(6)
	h := &captchaHandler{
		imgWidth:  280,
		imgHeight: 120,
		opts:      nil,
	}

	r := httptest.NewRequest("GET", "http://example.com/foo/"+id+".png", nil)
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
	}
}

func Test_captchaHandler_ServeHTTP(t *testing.T) {
	id := NewLen(6)

	type fields struct {
		imgWidth  int
		imgHeight int
		opts      *DistortionOpts
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type expected struct {
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    expected
	}{
		{"no", fields{280, 120, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo", nil),
			},
			expected{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{"bad ID", fields{280, 120, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo.png", nil),
			},
			expected{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{"no extension", fields{280, 120, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo/x.", nil),
			},
			expected{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{"no ID", fields{280, 120, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo/.png", nil),
			},
			expected{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{"yes", fields{280, 120, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo/"+id+".png", nil),
			},
			expected{
				code:        http.StatusOK,
				contentType: "image/png",
			},
		},
		{"yes download", fields{280, 120, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo/download/"+id+".png", nil),
			},
			expected{
				code:        http.StatusOK,
				contentType: "application/octet-stream",
			},
		},
		{"yes opts", fields{280, 120,
			&DistortionOpts{
				CircleCount: 40, MaxSkew: 0.8, StrikeCount: 2,
				CanvasWarp: defaultCanvasWarp, StrikeWarp: defaultStrikeWarp,
			}},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo/"+id+".png", nil),
			},
			expected{
				code:        http.StatusOK,
				contentType: "image/png",
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &captchaHandler{
				imgWidth:  tt.fields.imgWidth,
				imgHeight: tt.fields.imgHeight,
				opts:      tt.fields.opts,
			}
			h.ServeHTTP(tt.args.w, tt.args.r)
			rr := tt.args.w.(*httptest.ResponseRecorder)
			if tt.exp.code != rr.Code {
				t.Errorf("Bad code. Got %d, expected %d",
					rr.Code, tt.exp.code)
			}
			if rr.Code == http.StatusOK {
				fsave := fmt.Sprintf("blah%d.png", i)
				ioutil.WriteFile(fsave, rr.Body.Bytes(), os.ModePerm)
			}
			ctype := rr.Header().Get("Content-Type")
			if ctype != tt.exp.contentType {
				t.Errorf("Bad Content-Type. Got %s, expected %s",
					ctype, tt.exp.contentType)
			}
		})
	}
}
