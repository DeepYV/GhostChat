package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var wsChannel = make(chan WsPayload)

var clients = make(map[WebsocketConnection]string)

var views = jet.NewSet(jet.NewOSFileSystemLoader("./html"), jet.InDevelopmentMode())

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketConnection struct {
	*websocket.Conn
}
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"Message_type"`
}

type WsPayload struct {
	Action  string              `json:"Action"`
	Message string              `json:"Message"`
	Conn    WebsocketConnection `json:"-"`
}

func WsEndPoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	var response WsJsonResponse
	conn := WebsocketConnection{
		Conn: ws,
	}

	clients[conn] = ""
	response.Message = `<small> connected to server </small>`
	err = ws.WriteJSON(response)
	if err != nil {
		return
	}

	go ListenForWs(&conn)
}

func ListToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChannel

		response.Action = "got here"
		response.Message = fmt.Sprintf("some message and action %s", e.Action)
		broadcastToAll(response)
	}

}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.ReadJSON(response)
		if err != nil {
			_ = client.Close()
			delete(clients, client)
		}

	}
}

func ListenForWs(conn *WebsocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Print("error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			return
		} else {
			payload.Conn = *conn
			wsChannel <- payload

		}
	}
}
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		return
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		return err
	}
	return nil
}
