package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"mahayoga.org/models"

	"github.com/thedevsaddam/govalidator"
)

const (
	PORT = ":5300"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/cards", CreateCard).Methods("POST")
	r.HandleFunc("/api/cards", CreateCard).Methods("PUT")
	r.HandleFunc("/api/cards", DeleteCard).Methods("DELETE")
	r.HandleFunc("/api/cards/a", GetCards).Methods("GET")
	r.HandleFunc("/api/cards/m", GetCardsUnscoped).Methods("GET")

	log.Fatal(http.ListenAndServe(PORT, r))
}

func CreateCard(res http.ResponseWriter, req *http.Request) {
	rules := govalidator.MapData{
		"title":        []string{"required", "between:5,150"},
		"description":  []string{"required", "between:5,500"},
		"language":     []string{"required", "numeric", "between:0,30"},
		"file:cover":   []string{"ext:jpg,png,jpeg", fmt.Sprintf("size:%d", 1048576), "mime:image/jpeg,image/png", "required"},
		"file:content": []string{"ext:pdf", fmt.Sprintf("size:%d", 1048576*50), "mime:application/pdf", "required"},
	}

	messages := govalidator.MapData{
		"file:cover":   []string{"ext:Only jpg/png/jpeg is allowed", "size:The file size could not be greater than 1MB", "required:Cover is required"},
		"file:content": []string{"ext:Only pdf is allowed", "size:The file size could not be greater than 50MB", "required:Content is required"},
	}

	opts := govalidator.Options{
		Request:         req,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) != 0 {
		err := map[string]interface{}{"validationError": e}
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(res).Encode(err)
		return
	}

	uid := uuid.New().String()

	file_cover, _, err := req.FormFile("cover")
	if err != nil {
		log.Print(err)
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(res).Encode(err)
		return
	}

	buffer_cover, err := io.ReadAll(file_cover)
	if err != nil {
		log.Print(err)
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(res).Encode(err)
		return
	}

	file_content, _, err := req.FormFile("content")
	if err != nil {
		log.Print(err)
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(res).Encode(err)
		return
	}

	buffer_content, err := io.ReadAll(file_content)
	if err != nil {
		log.Print(err)
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(res).Encode(err)
		return
	}

	language, err := strconv.Atoi(req.FormValue("language"))
	if err != nil {
		log.Printf("Error parsing language : %v", err)
		http.Error(res, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	Card := &models.Cards{
		UID:         uid,
		Title:       req.FormValue("title"),
		Description: req.FormValue("description"),
		Language:    language,
		Cover:       buffer_cover,
		Content:     buffer_content,
	}

	err = Card.CreateCard()
	if err != nil {
		log.Fatalf("Error creating card: %v", err)
		http.Error(res, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(Card)
}

func UpdateCard(res http.ResponseWriter, req *http.Request) {
	rules := govalidator.MapData{
		"card_id":      []string{"between:5:40", "required"},
		"title":        []string{"between:5,150"},
		"description":  []string{"between:5,500"},
		"language":     []string{"numeric", "between:0,30"},
		"file:cover":   []string{"ext:jpg,png,jpeg", fmt.Sprintf("size:%d", 1048576), "mime:image/jpeg,image/png"},
		"file:content": []string{"ext:pdf", fmt.Sprintf("size:%d", 1048576*50), "mime:application/pdf"},
	}

	messages := govalidator.MapData{
		"file:cover":   []string{"ext:Only jpg/png/jpeg is allowed", "size:The file size could not be greater than 1MB"},
		"file:content": []string{"ext:Only pdf is allowed", "size:The file size could not be greater than 50MB"},
	}
	opts := govalidator.Options{
		Request:         req,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}

	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) != 0 {
		err := map[string]interface{}{"validationError": e}
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(res).Encode(err)
		return
	}

	card := &models.Cards{}

	if req.FormValue("title") != "" {
		card.Title = req.FormValue("title")
	}

	if req.FormValue("description") != "" {
		card.Description = req.FormValue("description")
	}

	cover_file, _, err := req.FormFile("cover")
	if err == nil {
		buffer, err := io.ReadAll(cover_file)
		if err != nil {
			log.Print(err)
			res.Header().Set("Content-type", "application/json")
			res.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(res).Encode(err)
			return
		}

		card.Cover = buffer

	}

	content_file, _, err := req.FormFile("content")
	if err == nil {
		buffer, err := io.ReadAll(content_file)
		if err != nil {
			log.Print(err)
			res.Header().Set("Content-type", "application/json")
			res.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(res).Encode(err)
			return
		}

		card.Content = buffer
	}

	if req.FormValue("language") != "" {
		card.Language, err = strconv.Atoi(req.FormValue("language"))

		if err != nil {
			log.Print(err)
			http.Error(res, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
			return
		}

	}

	err = card.UpdateCard(req.FormValue("card_id"))
	if err != nil {
		log.Fatalf("Error updating card: %v", err)
		http.Error(res, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(card)
}

func DeleteCard(res http.ResponseWriter, req *http.Request) {

	rules := govalidator.MapData{
		"card_id": []string{"between:5,40", "required"},
	}

	opts := govalidator.Options{
		Request:         req,
		Rules:           rules,
		RequiredDefault: false,
	}

	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) != 0 {
		err := map[string]interface{}{"validationError": e}
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(res).Encode(err)
		return
	}

	err := models.DeleteCard(req.URL.Query().Get("card_id"))
	if err != nil {
		log.Fatalf("Error deleting card: %v", err)
		http.Error(res, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode("DELETED")
}

func GetCards(res http.ResponseWriter, req *http.Request) {
	cards, err := models.GetCards()
	if err != nil {
		log.Fatalf("Error getting cards: %v", err)
		http.Error(res, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(cards)
}

func GetCardsUnscoped(res http.ResponseWriter, req *http.Request) {
	cards, err := models.GetCardsUnscoped(req.URL.Query().Get("modified"))
	if err != nil {
		log.Fatalf("Error getting cards: %v", err)
		http.Error(res, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(cards)
}
