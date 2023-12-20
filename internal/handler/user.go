package user
import(
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"os"
	// "rest1/internal/repository"
)
var db *pgx.Conn

func init() {
	// Establish a connection to the PostgreSQL database.
	connString := "user=your_username password=your_password host=localhost port=5432 dbname=your_database sslmode=disable"
	var err error
	db, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Error connecting to PostgreSQL:", err)
	}

	// Close the database connection when the application exits.
	defer db.Close()
}

func main() {
	// Initialize a user repository with the PostgreSQL connection.
	userRepo := repository.NewUserRepo(db)

	// Handle HTTP requests.
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepo.GetAll()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting users: %s", err), http.StatusInternalServerError)
			return
		}

		// Print the users to the response.
		for _, user := range users {
			fmt.Fprintf(w, "ID: %s, Name: %s, AccountNo: %s, Password: %s\n", user.ID, user.Name, user.AccountNo, user.Password)
		}
	})

	// Start the HTTP server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}