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

	"github.com/gorilla/mux"
	"github.com/disintegration/imaging"
)

type ImageJson struct {
	PayloadBase64 string `json:"payloadBase64"`
}

func readImage(r *http.Request) (image.Image, string, error) {
	var imgPayload ImageJson
	var err error

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&imgPayload); err != nil {
		return nil, "", err
	}

    indexComma := strings.Index(imgPayload.PayloadBase64, ",")
	rawImage := imgPayload.PayloadBase64[indexComma + 1:]
	var imgDecoded []byte
	if imgDecoded, err = base64.StdEncoding.DecodeString(rawImage); err != nil {
		return nil, "", err
	}

	var img image.Image
	var imgType string
	switch strings.TrimSuffix(imgPayload.PayloadBase64[5:indexComma], ";base64") {
	case "image/png":
		if img, err = png.Decode(bytes.NewReader(imgDecoded)); err != nil {
			return nil, "", err
		}
		imgType = "png"
	case "image/jpeg":
		if img, err = jpeg.Decode(bytes.NewReader(imgDecoded)); err != nil {
			return nil, "", err
		}
		imgType = "jpeg"
	default:
		return nil, "", errors.New("image format not supported")
	}

		return img, imgType, nil
}

func writeImage(w http.ResponseWriter, img image.Image, imgType string) error {
	// write image
	var imgBuffer bytes.Buffer
	if imgType == "png" {
		png.Encode(&imgBuffer, img)
	} else {
		jpeg.Encode(&imgBuffer, img, nil)
	}

	imgEncoded := base64.StdEncoding.EncodeToString(imgBuffer.Bytes())

	var retPayload ImageJson
	retPayload.PayloadBase64 = "data:image/" + imgType + ";base64," + imgEncoded;

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(retPayload); err != nil {
		return err
	}

	return nil
}

func changeBrightness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vals := r.URL.Query()
	valIntensity, okIntensity := vals["intensity"];
	if !okIntensity {
		log.Printf("Error getting intenstiy")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	strIntensity := valIntensity[0]
	if strIntensity == "" {
		log.Printf("var intensity is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	var intensity float64
	if intensity, err = strconv.ParseFloat(strIntensity, 64); err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var img image.Image
	var imgType string
	if img, imgType, err = readImage(r); err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	img = imaging.AdjustBrightness(img, intensity);

	if err = writeImage(w, img, imgType); err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	router := mux.NewRouter()
    router.HandleFunc("/brightness", changeBrightness).Methods("POST")
	log.Printf("Listening on :8080...\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}
