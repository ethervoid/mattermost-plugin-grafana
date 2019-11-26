package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"time"
)

func (p *Plugin) loadImageFromURL(panelUrl string) (image.Image, error) {

	urlParsed, err := url.Parse(panelUrl)
	if err != nil {
		return nil, err
	}
	urlQuery := urlParsed.Query()
	now := time.Now()
	duration, _ := time.ParseDuration("1min")
	then := now.Add(-duration)
	urlQuery.Set("to", string(now.Unix()))
	urlQuery.Set("from", string(then.Unix()))
	urlQuery.Set("width", "400")
	urlQuery.Set("height", "200")
	urlParsed.RawQuery = urlQuery.Encode()
	p.API.LogInfo("Calling URL: " + urlParsed.String())
	resp, err := http.Get(urlParsed.String())
	if err != nil {
		return nil, err
	}

	myImage, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return myImage, nil
}

func (p *Plugin) encodeImageToBase64(image image.Image) (string, error) {

	var buffer bytes.Buffer
	err := png.Encode(&buffer, image)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}
