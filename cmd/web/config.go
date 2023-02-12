package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
)

// Config struct veri yapısını programımızda çeşitli yerlerde receiver olarak kullanacağız.
// Böylelikle application configuration farklı modüller arasında kolaylıkla paylaşılabilecek.
type Config struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
}
