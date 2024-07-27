package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"slices"
	"time"

	"log"
)

type Server struct {
	dir string
}

func NewServer(datadir string) *Server {
	return &Server{
		dir: datadir,
	}
}

type songRequest struct {
	Emoji     int    `json:"emoji"`
	Song      string `json:"song"`
	Requester string `json:"requester"`
}

func (r songRequest) String() string {
	return fmt.Sprintf("[song: %s, requester: %s]å", r.Song, r.Requester)
}

// song is the way a song request is internally represented and stored
type song struct {
	Emoji     int       `json:"emoji"`
	Name      string    `json:"name"`
	Requester string    `json:"requester"`
	Created   time.Time `json:"created"`
	Id        int       `json:"id"`
}

func (s Server) saveSongRequest(name string, emoji int, requester string) error {
	id := rand.Int()
	created := time.Now()
	newSong := song{
		Emoji:     emoji,
		Id:        id,
		Created:   created,
		Name:      name,
		Requester: requester,
	}
	bytes, err := json.Marshal(newSong)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(s.dir, fmt.Sprintf("%d.json", id)), bytes, 0755)
}

func (s Server) SongRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("request was not a POST"))
		return
	}
	slurp, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to parse input: %s", err.Error())))
		return
	}
	log.Println("Received /request with ", string(slurp))
	// log.Printf("Received /request with body %s", req)

	var req songRequest
	err = json.Unmarshal(slurp, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to unmarshal json: %s", err.Error())))
		return
	}

	if req.Song == "" || req.Requester == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("both song and requester need to be non-empty"))
		return
	}

	s.saveSongRequest(req.Song, req.Emoji, req.Requester)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("got request"))
}

func (srv Server) getQueue() ([]song, error) {
	files, err := os.ReadDir(srv.dir)
	if err != nil {
		return nil, err
	}
	songs := make([]song, len(files))
	for i, file := range files {
		bytes, err := os.ReadFile(path.Join(srv.dir, file.Name()))
		if err != nil {
			log.Println("Unable to get file at ", srv.dir+file.Name())
			continue
		}
		var s song
		err = json.Unmarshal(bytes, &s)
		if err != nil {
			log.Println("Unable to unmarshal file at ", srv.dir+file.Name())
			continue
		}
		songs[i] = s
	}
	slices.SortFunc(songs, func(a song, b song) int {
		return -1 * a.Created.Compare(b.Created) // descending order
	})
	return songs, nil
}

type queueResponse struct {
	Songs []song `json:"songs"`
}

func (s Server) Queue(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /queue with ")
	queue, err := s.getQueue()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to get song queue: %s", err.Error())))
		return
	}
	q := queueResponse{
		Songs: queue,
	}
	bytes, err := json.Marshal(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unable to marshal JSON: %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
