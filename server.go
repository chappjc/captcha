// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Copyright 2019 Jonathan Chappelow. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"image/png"
	"net/http"
	"path"
	"strings"
)

type captchaHandler struct {
	imgWidth  int
	imgHeight int
	opts      *DistortionOpts
}

// Server returns a handler that serves HTTP requests with captcha images.
//
// To serve a captcha as a downloadable file, the URL must be constructed in
// such a way as if the file to serve is in the "download" subdirectory:
// "/download/LBm5vMjHDtdUfaWYXiQX.png".
//
// To reload captcha (get a different solution for the same captcha id), append
// "?reload=x" to URL, where x may be anything (for example, current time or a
// random number to make browsers refetch an image instead of loading it from
// cache).
func Server(imgWidth, imgHeight int, opts *DistortionOpts) http.Handler {
	return &captchaHandler{imgWidth, imgHeight, opts}
}

func (h *captchaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := strings.TrimSuffix(file, ext)
	if ext != ".png" || id == "" {
		http.NotFound(w, r)
		return
	}

	if r.FormValue("reload") != "" && !Reload(id) {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	digits := Digits(id)
	if len(digits) == 0 {
		http.NotFound(w, r)
		return
	}

	img := NewImage(id, digits, h.imgWidth, h.imgHeight, h.opts)
	if img == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	switch path.Base(dir) {
	case "download":
		w.Header().Set("Content-Type", "application/octet-stream")
	default:
		w.Header().Set("Content-Type", "image/png")
	}

	enc := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	if err := enc.Encode(w, img.Paletted); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
