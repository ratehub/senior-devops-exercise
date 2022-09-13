package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var kClientset *kubernetes.Clientset

var (
	version string = "dev"
	commit  string = ""
	date    string = ""
)

// https://stackoverflow.com/a/51270134
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func setup() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	kClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func verifyToken(clientId string) (bool, error) {
	ctx := context.TODO()
	audienceName := os.Getenv("AUDIENCE_NAME")
	if len(audienceName) == 0 {
		panic("AUDIENCE_NAME expected")
	}

	tr := authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token:     clientId,
			Audiences: []string{audienceName},
		},
	}
	result, err := kClientset.AuthenticationV1().TokenReviews().Create(ctx, &tr, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	log.Printf("%s\n", prettyPrint(result.Status))

	if result.Status.Authenticated {
		return true, nil
	}
	return false, nil

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Read the value of the client identifier from the request header
	clientId := r.Header.Get("X-Client-Id")
	if len(clientId) == 0 {
		http.Error(w, "X-Client-Id not supplied", http.StatusUnauthorized)
		return
	}
	authenticated, err := verifyToken(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !authenticated {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}
	io.WriteString(w, "Hello from data store. You have been authenticated")
}

func main() {
	log.Printf("arena-config-generator version: %s; commit: %s; build date: %s; go build: %s", version, commit, date, runtime.Version())
	setup()

	listenAddrPortString := os.Getenv("LISTEN_ADDR_PORT")
	if len(listenAddrPortString) == 0 {
		panic("LISTEN_ADDR_PORT expected")
	}
	var listenAddr = fmt.Sprintf(":%s", listenAddrPortString)

	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(listenAddr, nil)
}
