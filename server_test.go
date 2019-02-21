package captcha

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_captchaHandler_ServeHTTP(t *testing.T) {
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
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exp    expected
	}{
		{"ok", fields{128, 256, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo", nil),
			},
			expected{
				code: http.StatusNotFound,
			},
		},
		{"ok", fields{128, 256, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo.png", nil),
			},
			expected{
				code: http.StatusNotFound,
			},
		},
		{"ok", fields{128, 256, nil},
			args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://example.com/foo/download.", nil),
			},
			expected{
				code: http.StatusNotFound,
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
			t.Log(rr.Code)
			if rr.Code == http.StatusOK {
				fsave := fmt.Sprintf("blah%d.png", i)
				ioutil.WriteFile(fsave, rr.Body.Bytes(), os.ModePerm)
			}
		})
	}
}
