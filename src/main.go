package main
import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"github.com/gocql/gocql"
	"os"
	"time"
	"github.com/Sirupsen/logrus"
)

func main() {
	time.Sleep(20 * time.Second)

	address := os.Getenv("QUOTES_CASSANDRA_1_PORT_9042_TCP_ADDR")
	cluster := gocql.NewCluster(address)
	session, _ := cluster.CreateSession()

	app := Application{session:session}
	err := app.createKeyspaceAndTables()
	if err != nil {
		panic(err)
	}


	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/quote", app.AddQuote)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World !")
}

func (app *Application) createKeyspaceAndTables() error {
	err := app.session.Query("CREATE KEYSPACE IF NOT EXISTS quotes WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };").Exec()
	if err != nil {
		logrus.Error(err)
	}

	err = app.session.Query("CREATE TABLE IF NOT EXISTS quotes.quote (quote_id uuid, quote text, author text, PRIMARY_KEY(quote_id)").Exec()
	if err != nil {
		logrus.Error(err)
	}

	if err != nil { return err } else { return nil }
}

func (app *Application) AddQuote(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := app.session.Query("INSERT INTO quote (quote_id, quote, author) VALUES (uuid(), ?, ?)", "ReactJS C'est comme Swing !", "Cédric Hauber").Exec()
	if err != nil {
		logrus.Error(err)
	}
}

func (app *Application) Quotes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := app.session.Query("SELECT * FROM quote").Exec()
	if err != nil {
		logrus.Error(err)
	}
}

type Application struct {
	session *gocql.Session
}