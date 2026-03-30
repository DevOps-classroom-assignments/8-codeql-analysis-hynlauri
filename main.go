package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

const allowedDir = "./safe-files"

func main() {
	http.HandleFunc("/readfile", readFileHandler)
	http.HandleFunc("/exec", execHandler)

	fmt.Println("Listening on http://0.0.0.0:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func readFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	path := allowedDir + "/" + filename

	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Write(data)
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		http.Error(w, "Command failed", http.StatusInternalServerError)
		return
	}

	w.Write(out)
}