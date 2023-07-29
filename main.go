package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}
	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s, body:\n%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second,
		body)
	io.WriteString(w, "This is my website!\n")

	log.Printf("%s: Got the request. first(%t)=%s, second(%t)=%s\n", ctx.Value(keyServerAddr), hasFirst, first,
		hasSecond, second)
	io.WriteString(w, "website root")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Printf("%s: Got the request", ctx.Value(keyServerAddr))
	myName := r.PostFormValue("myName")
	if myName == "" {
		myName = "HTTP"
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()
}
