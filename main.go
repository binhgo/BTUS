package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"time"
	"github.com/InclusION/static"
	"github.com/centrifugal/centrifuge"
)

const maxUploadSize = 2 * 1024 // 2 MB
const uploadPath = "./tmp"

func main() {

	log.Printf("Server started. Listening on port %s", static.PORT)
	log.Printf("UTC Time: %s", time.Now().UTC())

	router := mux.NewRouter()
	router.HandleFunc("/User/Register", Register).Methods(static.HTTP_POST)
	router.HandleFunc("/User/Login", Login).Methods(static.HTTP_POST)
	router.HandleFunc("/User/Logout", Logout).Methods(static.HTTP_POST)
	router.HandleFunc("/User/Report", ReportUser).Methods(static.HTTP_POST)
	router.HandleFunc("/User/Profile/Update", UpdateProfile).Methods(static.HTTP_POST)
	router.HandleFunc("/User/Profile/View", LoadProfile).Methods(static.HTTP_GET)
	//router.HandleFunc("/User/UploadImage", ).Methods(static.HTTP_POST)


	router.HandleFunc("/Post/Create", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Post/Update", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Post/View", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Post/Delete", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Post/Comment/Create", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Post/Comment/Update", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Post/Comment/Delete", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Post/Comment/View", testConnection).Methods(static.HTTP_GET)

	router.HandleFunc("/Chat/CreateChannel", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Chat/BlockUser", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Chat/RemoveImage", testConnection).Methods(static.HTTP_GET)
	router.HandleFunc("/Chat/RemoveChat", testConnection).Methods(static.HTTP_GET)


	// fmc push notification
	router.HandleFunc("/Push/AddToken", addToken).Methods(static.HTTP_POST)
	router.HandleFunc("/Push/RemoveToken", removeToken).Methods(static.HTTP_POST)
	router.HandleFunc("/Push/PushToDevice", pushToDevice).Methods(static.HTTP_POST)
	router.HandleFunc("/Push/PushToUser", pushToUser).Methods(static.HTTP_POST)


	// chat gorilla ws
	go handleMessages()
	router.HandleFunc("/ws", handleConnections)

	// chat centrifuge ws
	node := initCentrifuge()
	router.Handle("/connection/websocket", centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{}))


	// handle files html, js for testing
	fs := http.FileServer(http.Dir("./chat"))
	router.PathPrefix("/chat").Handler(http.StripPrefix("/chat", fs))


	// upload file to server
	router.HandleFunc("/upload", uploadFileHandler())
	fs1 := http.FileServer(http.Dir(uploadPath))
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files", fs1))


	// Start HTTP server async
	go startHTTPServer(router)
	// Run program until interrupted.
	waitExitSignal(node)
}

// Start HTTP server.
func startHTTPServer(handler http.Handler) {
	err := http.ListenAndServe(static.PORT, handler)
	if err != nil {
		log.Fatal(err)
	}
}







