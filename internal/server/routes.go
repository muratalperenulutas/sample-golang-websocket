package server

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"
	"time"
	"websocket/internal/file"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type SensorData struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
	Light       float64 `json:"light"`
}
type Data struct {
	ClientTime int64      `json:"clientTime"`
	ServerTime int64      `json:"serverTime"`
	ClientId   string     `json:"clientId"`
	Data       SensorData `json:"data"`
}

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.Use(s.corsMiddleware)

	r.HandleFunc("/", s.ServerRunning)

	r.HandleFunc("/websocket", s.websocketHandler)

	return r
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) ServerRunning(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Server is running"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

var upgrader = websocket.Upgrader{}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	clientId := r.Header.Get("client-type")
	isIotDevice := clientId == "iot-device"

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func() {
		conn.Close()
		s.clients.RemoveClient(conn)
	}()

	s.clients.AddClient(conn, isIotDevice)

	if isIotDevice {
		for {
			var msg Data
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			msg.ServerTime = time.Now().UnixNano()

			log.Printf("Received data from non-IoT device: %v", msg)

			file.SaveTextFile(fmt.Sprintf("%+v\n", msg))

			for clientConn, isIot := range s.clients.GetAllConnections() {
				if !isIot {
					err := clientConn.WriteJSON(msg)
					if err != nil {
						log.Printf("Error sending message to client: %v", err)
					}
				}
			}
		}
	} else {
		for {
			payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
			err := conn.WriteMessage(websocket.TextMessage, []byte(payload))
			if err != nil {
				break
			}
			time.Sleep(time.Second * 5)
		}
	}

}
