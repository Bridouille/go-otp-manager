package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func auth(session sessions.Session, r render.Render) {
	if session.Get("user") == nil {
		r.JSON(http.StatusUnauthorized, map[string]string{"msg": "You need to be logged to acess this part"})
	}
}

func getIndex(session sessions.Session, w http.ResponseWriter, req *http.Request, r render.Render) {
	if session.Get("user") == nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	r.HTML(http.StatusOK, "otp", nil)
}

func getLogin(session sessions.Session, r render.Render) {
	r.HTML(http.StatusOK, "login", nil)
}

func postLogin(session sessions.Session, w http.ResponseWriter, req *http.Request, r render.Render) {
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")

	if username == viper.GetString("config.credentials.username") && password == viper.GetString("config.credentials.password") {
		session.Set("user", username)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	r.HTML(http.StatusOK, "login", nil)
}

func postLogout(session sessions.Session, w http.ResponseWriter, req *http.Request) {
	session.Delete("user")
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func applyMigrations(db *sql.DB) (n int, err error) {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	if n, err = migrate.Exec(db, "sqlite3", migrations, migrate.Up); err != nil {
		return
	} else if n > 0 {
		log.Printf("Applied %d migrations!\n", n)
	}
	return
}

func dbConn() martini.Handler {
	db, err := sql.Open("sqlite3", viper.GetString("dbfile"))
	if err != nil {
		log.Fatal(err.Error())
	}

	if _, err := applyMigrations(db); err != nil {
		log.Fatalf("Can not apply migrations : %v", err)
	}

	return func(c martini.Context) {
		c.Map(db)
		c.Next()
	}
}

func loadConfig() error {
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}

	if os.Getenv("OTP_MANAGER_ENV") == "production" {
		viper.Set("config", viper.Get("production"))
	} else {
		viper.Set("config", viper.Get("development"))
	}

	return nil
}

func main() {
	if err := loadConfig(); err != nil {
		log.Fatalf("%v", err.Error())
	}

	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("jesuisunabricot"))
	store.Options(sessions.Options{MaxAge: viper.GetInt("config.sessiontime")})

	m.Get("/", getIndex)
	m.Get("/login", getLogin)
	m.Post("/login", postLogin)
	m.Post("/logout", postLogout)

	// API part
	m.Get("/otp", auth, getOtps)
	m.Get("/otp/(?P<id>[0-9]+)", auth, getOtp)
	m.Post("/otp", auth, postOtp)
	m.Delete("/otp/(?P<id>[0-9]+)", auth, deleteOtp)

	m.Use(dbConn())                                                             // inject our db to handlers
	m.Use(render.Renderer(render.Options{Delims: render.Delims{"{[%", "%]}"}})) // inject the renderer to handlers
	m.Use(martini.Static("assets"))
	m.Use(sessions.Sessions("otpToken", store)) // inject the sessions to handlers

	m.Run()
}
