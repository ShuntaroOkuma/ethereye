package favorites

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type AddressStorage struct {
	filename  string
	addresses map[string][]string
}

func NewAddressStorage(filename string) *AddressStorage {
	return &AddressStorage{filename: filename, addresses: make(map[string][]string)}
}

func (s *AddressStorage) Load() error {
	data, err := ioutil.ReadFile(s.filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &s.addresses)
}

func (s *AddressStorage) Save() error {
	data, err := json.Marshal(s.addresses)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.filename, data, 0644)
}

func (s *AddressStorage) AddFavoriteAddress(addressType, address string) {
	s.addresses[addressType] = append(s.addresses[addressType], address)
}

func (s *AddressStorage) GetFavoriteAddresses(addressType string) []string {
	return s.addresses[addressType]
}

func FavoriteAddressHandler(s *AddressStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		addressType := r.URL.Query().Get("type")
		if addressType != "wallet" && addressType != "token" {
			http.Error(w, "Invalid address type", http.StatusBadRequest)
			return
		}

		switch r.Method {
		// GET /favorites?type={wallet|token}
		case http.MethodGet:
			addresses := s.GetFavoriteAddresses(addressType)
			response, err := json.Marshal(addresses)
			if err != nil {
				http.Error(w, "Failed to get favorite addresses", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(response)

		// POST /favorites
		case http.MethodPost:
			var requestBody struct {
				Address string `json:"address"`
			}
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			s.AddFavoriteAddress(addressType, requestBody.Address)
			if err := s.Save(); err != nil {
				http.Error(w, "Failed to save favorite address", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
