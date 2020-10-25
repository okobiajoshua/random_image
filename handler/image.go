package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"random_image/database"
	"random_image/model"
	"strconv"

	"github.com/gorilla/mux"
)

// ImageHandler struct
type ImageHandler struct {
	db database.Database
}

// NewImageHandler method
func NewImageHandler(db database.Database) *ImageHandler {
	return &ImageHandler{db: db}
}

// GetRandomImage method
func (ih *ImageHandler) GetRandomImage(w http.ResponseWriter, r *http.Request) {
	img, err := getRandomImageDetail()
	if err != nil {
		log.Print(err)
		http.Error(w, "Error!", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(img)
}

// GetImageByID method
func (ih *ImageHandler) GetImageByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID must be an integer", http.StatusInternalServerError)
	}
	img, err := ih.db.GetImageByID(id)
	if err != nil {
		img, err = getRandomImageDetail()
		if err != nil {
			http.Error(w, "Error!", http.StatusInternalServerError)
		}
		ih.db.AddNewImage(id, img)
	}
	json.NewEncoder(w).Encode(img)
}

// GetStoredImages method
func (ih *ImageHandler) GetStoredImages(w http.ResponseWriter, r *http.Request) {
	img, err := ih.db.GetAll()
	if err != nil {
		json.NewEncoder(w).Encode("Error!")
	}
	json.NewEncoder(w).Encode(img)
}

func genRandomNumber() int {
	return rand.Intn(1000)
}

func getImageDetail(id int) (*model.Image, error) {
	url := fmt.Sprintf("https://picsum.photos/id/%d/info", id)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var img model.Image
	err = json.NewDecoder(res.Body).Decode(&img)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func getRandomImageDetail() (*model.Image, error) {
	n := genRandomNumber()
	return getImageDetail(n)
}
