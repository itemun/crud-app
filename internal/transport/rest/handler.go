package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/itemun/crud-app/internal/domain"

	"github.com/gorilla/mux"
)

type Books interface {
	Create(ctx context.Context, car domain.Car) error
	GetByID(ctx context.Context, id int64) (domain.Car, error)
	GetAll(ctx context.Context) ([]domain.Car, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateCarInput) error
}

type Handler struct {
	booksService Cars
}

func NewHandler(cars Cars) *Handler {
	return &Handler{
		carsService: cars,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	cars := r.PathPrefix("/cars").Subrouter()
	{
		cars.HandleFunc("", h.createCar).Methods(http.MethodPost)
		cars.HandleFunc("", h.getAllCars).Methods(http.MethodGet)
		cars.HandleFunc("/{id:[0-9]+}", h.getCarByID).Methods(http.MethodGet)
		cars.HandleFunc("/{id:[0-9]+}", h.deleteCar).Methods(http.MethodDelete)
		cars.HandleFunc("/{id:[0-9]+}", h.updateCar).Methods(http.MethodPut)
	}

	return r
}

func (h *Handler) getCarByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("getCarByID() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	car, err := h.carsService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrCarNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("getCarByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(car)
	if err != nil {
		log.Println("getCarByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createCar(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book domain.Car
	if err = json.Unmarshal(reqBytes, &car); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.carsService.Create(context.TODO(), car)
	if err != nil {
		log.Println("createCar() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteCar(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("deleteCar() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.carsService.Delete(context.TODO(), id)
	if err != nil {
		log.Println("deleteCar() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllCars(w http.ResponseWriter, r *http.Request) {
	cars, err := h.carsService.GetAll(context.TODO())
	if err != nil {
		log.Println("getAllCars() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(cars)
	if err != nil {
		log.Println("getAllCars() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) updateCar(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.UpdateCarInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.carsService.Update(context.TODO(), id, inp)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
