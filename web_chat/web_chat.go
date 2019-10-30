package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

const (
	httpListenAddr = "localhost:4000"
)

func main() {
	http.HandleFunc("/", rootPage)
	http.Handle("/socket", websocket.Handler(socketHandler))
	http.ListenAndServe(httpListenAddr, nil)
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, httpListenAddr)
}

type socket struct {
	*websocket.Conn
	done chan bool
}

func (s socket) Read(b []byte) (int, error)  { return s.Conn.Read(b) }
func (s socket) Write(b []byte) (int, error) { return s.Conn.Write(b) }

func (s socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(ws *websocket.Conn) {
	s := socket{Conn: ws, done: make(chan bool)}
	go match(s)
	<-s.done
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprintln(c, "Waiting for a parnter...")
	select {
	case partner <- c:
	// now handled by other goroutine
	case p := <-partner:
		chat(p, c)
	}
}
func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

func cp(w, r io.ReadWriteCloser, errc chan error) {
	_, err := io.Copy(w, r)
	errc <- err
}
