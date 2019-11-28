package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"net/http"
)

func (p *Plugin) loadImageFromURL(panelURL string) (image.Image, error) {
	resp, err := http.Get(panelURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
