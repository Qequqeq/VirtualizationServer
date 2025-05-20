package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

// Data structures remain unchanged for backward compatibility
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VMInfo struct {
	Addr   string `json:"address"`
	Username string `json:"username"`
	Port int    `json:"port"`
	Message string `json:"message"`
	Password string `json:"password"`
}

var db *sql.DB
var publicIP string

func main() {
	// Database initialization remains the same
	var err error
	db, err = sql.Open("sqlite3", "/root/VirtualizationServer/database/vm.db")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Existing table schema preserved
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		vm_port INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Table creation error:", err)
	}


	// Simplified routing with explicit method checks
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/vm", authMiddleware(vmHandler))
	http.HandleFunc("/wtf", wtfHandler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func wtfHandler(w http.ResponseWriter, r *http.Request) {
    examples := map[string]interface{}{
        "register_example": map[string]string{
            "method":  "POST",
            "command": `curl -X POST -H "Content-Type: application/json" -d '{"username":"example","password":"example"}' https://<yourdomain>.ru.tuna.am/register`,
            "comment": "Registration",
        },
        "login_example": map[string]string{
            "method":  "POST",
            "command": `curl -X POST -H "Content-Type: application/json" -d '{"username":"example","password":"example"}' https://<yourdomain>.ru.tuna.am/login`,
            "comment": "Authorization",
        },
        "get_vm_example": map[string]string{
            "method":  "GET",
            "command": `curl -H "Authorization: Bearer 5b907122-af7a-4105-a95c-86efcfb8cbf6" "https://<yourdomain>.ru.tuna.am/vm?username=example"`,
            "comment": "Get VM data (need to authorization)",
        },
        "ssh_connect": map[string]string{
            "command": `ssh username@address -p port`,
            "comment": "For connect to VM",
        },
        "notes": []string{
            "Change example to real username",
            "Use token from login in  Authorization",
            "Port is unique for all users",
        },
    }

    w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    if err := enc.Encode(examples); err != nil {
        http.Error(w, "Error generating examples", http.StatusInternalServerError)
    }
}
// Existing registration logic preserved
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		user.Username,
		string(hashedPassword),
	)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s registered successfully", user.Username)
}

// Existing login logic preserved
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	var storedHash string
	err := db.QueryRow(
		"SELECT password_hash FROM users WHERE username = ?",
		user.Username,
	).Scan(&storedHash)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Simplified token generation
	token := uuid.New().String()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":  token,
		"status": "authenticated",
	})
}

// Simplified auth middleware (only checks for token presence)
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// Updated VM handler with error logging
func vmHandler(w http.ResponseWriter, r *http.Request) {

    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "Username required", http.StatusBadRequest)
        return
    }
    port, err := createVM(username)
    if err != nil {
	    log.Printf("VM creation failed: %v", err)
	    http.Error(w, "You are here", http.StatusBadRequest)
	    return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(VMInfo{
        Addr:   "ru.tuna.am",
	Username: "root",
        Port: port,
	Message: "please wait like 2 mins",
	Password: "toor",
    })
}

// CreateVM remains similar but with better error handling
func createVM(username string) (int, error) {
    scriptPath := "/root/VirtualizationServer/vm-api/scripts/vmctl.sh"
    cmd := exec.Command("/bin/bash", scriptPath, "create", username)
    output, err := cmd.CombinedOutput()
    if err != nil {
	    log.Printf("SMTH went wrong")
	    return 0, nil
    }
    log.Printf("output: %s", output)

    // tuna-port
    time.Sleep(3 * time.Second)
    tunaPort, err := getTunaPort("/root/VirtualizationServer/database/tuna_ports")
    if err != nil {
        log.Printf("Tuna port error: %v | VM Script output: %s", err, string(tunaPort))
        return 0, fmt.Errorf("tuna port error: %v", err)
    }

    return tunaPort, nil
}

//get tuna port
func getTunaPort(filePath string) (int, error) {
	cmd := exec.Command("python3", "/root/VirtualizationServer/vm-api/parser.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("bad parsing, return %s", output)
	}
	portStr := string(output)
	port, err := strconv.Atoi(strings.TrimSpace(portStr))
	if err != nil {
		return 0, fmt.Errorf("convertation error")
	}
	return port, nil
}


