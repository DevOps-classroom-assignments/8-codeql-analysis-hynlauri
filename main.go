package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	cleanFile := filepath.Clean(filename)
	if cleanFile == "." || strings.HasPrefix(cleanFile, "..") || filepath.IsAbs(cleanFile) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	path := filepath.Join(allowedDir, cleanFile)
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Write(data)
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	cmdParam := r.URL.Query().Get("cmd")
	if cmdParam != "ls" {
		http.Error(w, "Command is not allowed", http.StatusForbidden)
		return
	}

	out, err := exec.Command("ls").Output()
	if err != nil {
		http.Error(w, "Command failed", http.StatusInternalServerError)
		return
	}

	w.Write(out)
}