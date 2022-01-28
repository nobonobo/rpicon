package backend

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/gorilla/websocket"

	"github.com/nobonobo/rpicon/procon"
)

var upgrader = websocket.Upgrader{} // use default options

type response struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response{
		Status: "ok",
		Result: result,
	}); err != nil {
		http.Error(w, Error(err), http.StatusInternalServerError)
	}
}

func Error(err error) string {
	_, fn, ln, _ := runtime.Caller(1)
	log.Printf("%s:%d %s", filepath.Base(fn), ln, err)
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode(map[string]interface{}{
		"status": "error",
		"error":  err.Error(),
	})
	return b.String()
}

type API struct {
	http.Handler
	Script string
	client *procon.Client
}

func New(scriptPath string) *API {
	mux := http.NewServeMux()
	api := &API{Handler: mux, Script: scriptPath}
	mux.HandleFunc("/api/health", api.health)
	mux.HandleFunc("/api/ws", api.ws)
	api.client = procon.New()
	return api
}

func (api *API) health(w http.ResponseWriter, r *http.Request) {
	Success(w, nil)
}

func (api *API) ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	log.Print("start procon")
	if err := api.client.Start(api.Script); err != nil {
		log.Print(err)
		return
	}
	defer func() {
		log.Print("stop procon")
		api.client.Stop()
	}()
	if err := api.client.Connect(); err != nil {
		log.Print(err)
		return
	}
	defer api.client.Disconnect()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		v := procon.Input{}
		if err := json.Unmarshal([]byte(message), &v); err != nil {
			log.Println("unmarshal:", err)
			break
		}
		if err := api.client.Input(v); err != nil {
			e := err
			s, err := api.client.State()
			if err != nil {
				log.Print(e, err)
			} else {
				log.Print(e, s)
			}
			if err := api.client.Connect(); err != nil {
				log.Print(err)
				return
			}
		}
		//log.Printf("recv: %v", v)
	}
}
