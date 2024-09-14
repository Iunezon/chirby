package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorId, _ := strconv.Atoi(r.URL.Query().Get("author_id"))
	sortType := r.URL.Query().Get("sort")

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	if authorId == 0 {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:       dbChirp.ID,
				Body:     dbChirp.Body,
				AuthorID: dbChirp.AuthorID,
			})
		}
	} else {
		for _, dbChirp := range dbChirps {
			if dbChirp.AuthorID == authorId {
				chirps = append(chirps, Chirp{
					ID:       dbChirp.ID,
					Body:     dbChirp.Body,
					AuthorID: dbChirp.AuthorID,
				})
			}
		}
	}

	switch sortType {
	case "desc": // Sort by descending order
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	default: // Default to ascending order
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
