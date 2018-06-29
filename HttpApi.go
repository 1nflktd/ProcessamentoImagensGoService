package main

import (
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"strings"
	"encoding/base64"
	"image"
	"image/png"
	"image/jpeg"
	"errors"
	"strconv"
)

type JsonResponse struct {
	PayloadBase64 string `json:"payloadBase64"`
}

type HttpApi struct {
	Writer http.ResponseWriter
	Request *http.Request
	Image *Image
	Parameters *Parameters
}

type Parameters struct {
	Intensity float64
}

func HttpApiNew(w http.ResponseWriter, r *http.Request) *HttpApi {
	return &HttpApi{
		Writer: w,
		Request: r,
	}
}

func (api *HttpApi) Init() error {
	api.Writer.Header().Set("Content-Type", "application/json")

	var err error
	if err = api.getParameters(); err != nil {
		log.Printf(err.Error())
		api.Writer.WriteHeader(http.StatusBadRequest)
		return err
	}

	if err = api.readImage(); err != nil {
		log.Printf(err.Error())
		api.Writer.WriteHeader(http.StatusBadRequest)
		return err
	}

	return nil
}

func (api *HttpApi) getParameters() (error) {
	vals := api.Request.URL.Query()
	valIntensity, okIntensity := vals["intensity"];
	if !okIntensity {
		return errors.New("Error getting intenstiy")
	}
	strIntensity := valIntensity[0]
	if strIntensity == "" {
		return errors.New("var intensity is empty")
	}

	var err error
	var intensity float64
	if intensity, err = strconv.ParseFloat(strIntensity, 64); err != nil {
		return err
	}

	api.Parameters = &Parameters{
		Intensity: intensity,
	}

	return nil
}

func (api *HttpApi) readImage() (error) {
	var imgPayload JsonResponse
	var err error

	decoder := json.NewDecoder(api.Request.Body)
	defer api.Request.Body.Close()
	if err = decoder.Decode(&imgPayload); err != nil {
		return err
	}

    indexComma := strings.Index(imgPayload.PayloadBase64, ",")
	rawImage := imgPayload.PayloadBase64[indexComma + 1:]
	var imgDecoded []byte
	if imgDecoded, err = base64.StdEncoding.DecodeString(rawImage); err != nil {
		return err
	}

	var img image.Image
	var imageType string
	switch strings.TrimSuffix(imgPayload.PayloadBase64[5:indexComma], ";base64") {
	case "image/png":
		if img, err = png.Decode(bytes.NewReader(imgDecoded)); err != nil {
			return err
		}
		imageType = "png"
	case "image/jpeg":
		if img, err = jpeg.Decode(bytes.NewReader(imgDecoded)); err != nil {
			return err
		}
		imageType = "jpeg"
	default:
		return errors.New("image format not supported")
	}

	api.Image = ImageNew(img, imageType)
	return nil
}

func (api *HttpApi) writeImage() error {
	var imgBuffer bytes.Buffer
	if api.Image.Type == "png" {
		png.Encode(&imgBuffer, api.Image.Image)
	} else {
		jpeg.Encode(&imgBuffer, api.Image.Image, nil)
	}

	imgEncoded := base64.StdEncoding.EncodeToString(imgBuffer.Bytes())

	var retPayload JsonResponse
	retPayload.PayloadBase64 = "data:image/" + api.Image.Type + ";base64," + imgEncoded;

	api.Writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(api.Writer).Encode(retPayload); err != nil {
		return err
	}

	return nil
}
