package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"syscall"

	"getartist/app/auth"
	"getartist/app/models"
	"getartist/config"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"golang.org/x/sys/unix"
)

type JSONError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type HealthCheck struct {
	Status int
	Result string
}

func APIError(w http.ResponseWriter, errMessage string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JSONError{Error: errMessage, Code: code})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonError)
}

// responseJSON JSON形式に変換する
func responseJSON(w http.ResponseWriter, value interface{}) {
	js, err := json.Marshal(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// getID URLのIDを取得する
func getID(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	return strconv.Atoi(vars["id"])
}

// web artist
func searchArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artist_name")
	if artistName == "" {
		APIError(w, "No artist_name param", http.StatusBadRequest)
		return
	}

	client := models.GetClient()
	result, err := client.Search(artistName, spotify.SearchTypeArtist) // artistName
	if err != nil {
		log.Fatalf("couldn't get artists: %v", err)
	}

	// json出力
	responseJSON(w, result)
}

func getArtistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	ID := vars["id"]

	artist := models.GetSpotifyArtist(ID)

	// json出力
	responseJSON(w, artist)
}

func getArtistInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	artistID := vars["id"]

	if artistID == "" {
		APIError(w, "No artist_id param", http.StatusBadRequest)
		return
	}

	artistInfo := models.GetArtistInfoFromArtistID(artistID)

	// json出力
	responseJSON(w, artistInfo)
}

func getArtistInfosHandler(w http.ResponseWriter, r *http.Request) {
	artistInfos := models.GetArtistInfos()

	var artists []*spotify.FullArtist
	for _, artistInfo := range artistInfos {
		client := models.GetClient()
		artist, err := client.GetArtist(spotify.ID(artistInfo.ArtistId))
		if err != nil {
			fmt.Println(err)
		}
		artists = append(artists, artist)
	}

	// json出力
	responseJSON(w, artists)
}

// web article
func getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articles := models.GetArticles()

	w.Header().Set("X-Total-Count", strconv.Itoa(len(articles)))
	responseJSON(w, articles)
}

func getArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	article := models.GetArticle(ID)

	responseJSON(w, article)
}

// admin artist
func getAdminArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artistInfos := models.GetArtistInfos()

	w.Header().Set("X-Total-Count", strconv.Itoa(len(artistInfos)))
	responseJSON(w, artistInfos)
}

func getAdminArtistHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	artistInfo := models.GetArtistInfo(ID)

	responseJSON(w, artistInfo)
}

func createArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistInfo := models.CreateArtistInfo(r)

	responseJSON(w, artistInfo)
}

func updateArtistHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	artistInfo := models.UpdateArtistInfo(r, ID)

	responseJSON(w, artistInfo)
}

func deleteArtistHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	models.DeleteArtistInfo(ID)
}

// admin article
func getAdminArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articles := models.GetArticles()

	w.Header().Set("X-Total-Count", strconv.Itoa(len(articles)))
	responseJSON(w, articles)
}

func getAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	article := models.GetArticle(ID)

	responseJSON(w, article)
}

func createAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	article := models.CreateArticle(r)

	responseJSON(w, article)
}

func updateAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	article := models.UpdateArticle(r, ID)

	responseJSON(w, article)
}

func deleteAdminArticleHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := getID(w, r)

	if err != nil {
		fmt.Println(err)
	}

	models.DeleteArticle(ID)
}

// healthCheckHandler is ALBによるヘルスチェック用
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ping := HealthCheck{http.StatusOK, "ok"}

	res, err := json.Marshal(ping)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func listenCtrl(network string, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(s uintptr) {
		err = unix.SetsockoptInt(int(s), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1) // portをbindできる設定
		if err != nil {
			return
		}
	})
	return err
}

func StartWebServer() error {
	r := mux.NewRouter()
	// web
	// artist
	r.HandleFunc("/api/search", searchArtistHandler).Methods("GET")
	r.HandleFunc("/api/artist/infos", getArtistInfosHandler).Methods("GET")
	r.HandleFunc("/api/artist/info/{id}", getArtistInfoHandler).Methods("GET")
	r.HandleFunc("/api/artist/{id}", getArtistHandler).Methods("GET")

	// article
	r.HandleFunc("/api/articles", getArticlesHandler).Methods("GET")
	r.HandleFunc("/api/articles/{id}", getArticleHandler).Methods("GET")

	// admin
	// artist
	r.HandleFunc("/api/admin/artists", getAdminArtistsHandler).Methods("GET")
	r.HandleFunc("/api/admin/artists/{id}", getAdminArtistHandler).Methods("GET")
	r.HandleFunc("/api/admin/artists", createArtistHandler).Methods("POST")
	r.HandleFunc("/api/admin/artists/{id}", updateArtistHandler).Methods("PUT")
	r.HandleFunc("/api/admin/artists/{id}", deleteArtistHandler).Methods("DELETE")

	// // article
	r.HandleFunc("/api/admin/articles", getAdminArticlesHandler).Methods("GET")
	r.HandleFunc("/api/admin/articles/{id}", getAdminArticleHandler).Methods("GET")
	r.HandleFunc("/api/admin/articles", createAdminArticleHandler).Methods("POST")
	r.HandleFunc("/api/admin/articles/{id}", updateAdminArticleHandler).Methods("PUT")
	r.HandleFunc("/api/admin/articles/{id}", deleteAdminArticleHandler).Methods("DELETE")

	// auth
	r.HandleFunc("/api/admin/login", auth.Login).Methods("POST")

	r.HandleFunc("/health-check/", healthCheckHandler)
	http.Handle("/", r)

	lc := net.ListenConfig{
		Control: listenCtrl, //portのbindを許可する設定を入れる
	}

	listener, err := lc.Listen(context.Background(), "tcp4", fmt.Sprintf(":%d", config.Config.Port))
	if err != nil {
		panic(err)
	}

	return http.Serve(listener, nil)
}
