// Package client provides the client application logic.
package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/Pklerik/gophKeep/internal/client"
	"github.com/Pklerik/gophKeep/internal/client/httpclient"
	config "github.com/Pklerik/gophKeep/internal/config/client"
	"github.com/Pklerik/gophKeep/internal/logger"
	"github.com/Pklerik/gophKeep/internal/models"
)

// App represents the client application.
type App struct {
	client client.IClient
	config config.Config
}

// NewApp creates a new client application.
func NewApp(cfg config.Config) *App {
	return &App{
		client: httpclient.NewHTTPClient(cfg.ServerURL, cfg.Timeout),
		config: cfg,
	}
}

// StartApp starts the client application.
func StartApp(cfg config.Config) {
	app := NewApp(cfg)
	app.Run()
}

// Run starts the interactive CLI loop.
func (a *App) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== GophKeeper ===")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. List Secrets")
		fmt.Println("4. Add Secret")
		fmt.Println("5. Get Secret")
		fmt.Println("6. Update Secret")
		fmt.Println("7. Delete Secret")
		fmt.Println("8. Exit")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			a.handleRegister(reader)
		case "2":
			a.handleLogin(reader)
		case "3":
			a.handleListSecrets()
		case "4":
			a.handleAddSecret(reader)
		case "5":
			a.handleGetSecret(reader)
		case "6":
			a.handleUpdateSecret(reader)
		case "7":
			a.handleDeleteSecret(reader)
		case "8":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

// handleRegister handles user registration.
func (a *App) handleRegister(reader *bufio.Reader) {
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	password := passwordPrompt("Enter password:")
	fmt.Println() // New line after password input
	password = strings.TrimSpace(password)

	resp, err := a.client.Register(username, password)
	if err != nil {
		logger.Sugar.Errorf("Registration failed: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Registration successful!\nToken: %s\n", resp.Token)
}

// handleLogin handles user login.
func (a *App) handleLogin(reader *bufio.Reader) {
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	password := passwordPrompt("Enter password:")
	fmt.Println() // New line after password input
	password = strings.TrimSpace(password)

	resp, err := a.client.Login(username, password)
	if err != nil {
		logger.Sugar.Errorf("Login failed: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Login successful!\nToken: %s\n", resp.Token)
}

// handleListSecrets lists all user secrets.
func (a *App) handleListSecrets() {
	secrets, err := a.client.ListSecrets()
	if err != nil {
		logger.Sugar.Errorf("Failed to list secrets: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found")
		return
	}

	fmt.Println("\n=== Your Secrets ===")
	for _, secret := range secrets {
		fmt.Printf("ID: %s\nType: %s\nTitle: %s\nMetadata: %s\n---\n",
			secret.ID, secret.Type, secret.Title, secret.Metadata)
	}
}

// handleAddSecret adds a new secret.
func (a *App) handleAddSecret(reader *bufio.Reader) {
	fmt.Println("\t1. credentials")
	fmt.Println("\t2. text")
	fmt.Println("\t3. binary")
	fmt.Println("\t4. card")
	fmt.Print("Choose secret type (credentials/text/binary/card): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var secretType models.SecretType
	switch choice {
	case "1":
		secretType = models.SecretTypeCredentials
	case "2":
		secretType = models.SecretTypeText
	case "3":
		secretType = models.SecretTypeBinary
	case "4":
		secretType = models.SecretTypeCard
	default:
		fmt.Println("Invalid secret type")
		return
	}

	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Data: ")
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)

	fmt.Print("Metadata (optional): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	secret, err := a.client.CreateSecret(models.SecretType(secretType), title, data, metadata)
	if err != nil {
		logger.Sugar.Errorf("Failed to create secret: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Secret created successfully!\nID: %s\n", secret.ID)
}

// handleGetSecret retrieves a specific secret.
func (a *App) handleGetSecret(reader *bufio.Reader) {
	fmt.Print("Secret ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	secret, err := a.client.GetSecret(id)
	if err != nil {
		logger.Sugar.Errorf("Failed to get secret: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("ID: %s\nType: %s\nTitle: %s\nData: %s\nMetadata: %s\n",
		secret.ID, secret.Type, secret.Title, string(secret.Data), secret.Metadata)
}

// handleUpdateSecret updates an existing secret.
func (a *App) handleUpdateSecret(reader *bufio.Reader) {
	fmt.Print("Secret ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("New Type (credentials/text/binary/card): ")
	secretType, _ := reader.ReadString('\n')
	secretType = strings.TrimSpace(secretType)

	fmt.Print("New Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("New Data: ")
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)

	fmt.Print("New Metadata (optional): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	secret, err := a.client.UpdateSecret(id, models.SecretType(secretType), title, data, metadata)
	if err != nil {
		logger.Sugar.Errorf("Failed to update secret: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Secret updated successfully!\nID: %s\n", secret.ID)
}

// handleDeleteSecret deletes a secret.
func (a *App) handleDeleteSecret(reader *bufio.Reader) {
	fmt.Print("Secret ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	err := a.client.DeleteSecret(id)
	if err != nil {
		logger.Sugar.Errorf("Failed to delete secret: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Secret deleted successfully!")
}

// passwordPrompt displays a label and securely reads a password from the terminal.
func passwordPrompt(label string) string {
	fmt.Fprint(os.Stderr, label+" ") // Print label to stderr to keep stdout clean

	// ReadPassword disables terminal echo and reads input until a newline
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))

	if err != nil {
		// Handle errors as appropriate for your application
		fmt.Printf("\nError reading password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println() // Print a newline after input
	return string(bytePassword)
}
