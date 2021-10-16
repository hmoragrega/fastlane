package ws

import (
	"github.com/gorilla/websocket"
	"github.com/hmoragrega/fastlane"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

func Ws(svc *fastlane.Syncer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var mx sync.Mutex

		changes := svc.Subscribe()

		defer svc.Unsubscribe(changes)

		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("ws upgrade failed:", err)
			return
		}
		defer func() {
			_ = socket.Close()
		}()

		//_ = socket.WriteJSON(fastlane.Event{Name: fastlane.ReviewsEventName, Data: svc.OpenReviews()})

		go func() {
			for evt := range changes{
				writeToSocket(&mx, socket, evt)
			}
		}()

		for {
			var evt fastlane.Event
			err := socket.ReadJSON(&evt)
			//_, message, err := socket.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Println("ws read error:", err)
				}
				break
			}

			log.Printf("recv: %+v", evt)
			for _, res := range svc.Handle(r.Context(), evt) {
				log.Printf("response: %+v", res)
				writeToSocket(&mx, socket, res)
			}
		}
	}
}

func writeToSocket(mx *sync.Mutex, socket *websocket.Conn, evt fastlane.Event) {
	mx.Lock()
	defer mx.Unlock()

	_ = socket.WriteJSON(evt)
}