package main

import (
	"dbconnection/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err!= nil {
		respondWithError(w, 400, fmt.Sprint("Error parsing JSON: %v", err))
	}

	feedFollows, err  := apiCfg.DB.CreateFeedFollows(r.Context(),database.CreateFeedFollowsParams {
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,


	})
	if err!= nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJson(w, 200, databaseFeedFollowToFeedFollow(feedFollows))
}


func (apiCfg *apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){

	feedFollows, err  := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err!= nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJson(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}


func (apiCfg *apiConfig) handlerRemoveFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err!=nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feedFollowID: %v", err))
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(),database.DeleteFeedFollowsParams { 
		ID: feedFollowID,
		UserID: user.ID,
	})
	if err!= nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't unfollowfeed: %v", err))
		return
	}

	respondWithJson(w, 200, struct{}{})
}

