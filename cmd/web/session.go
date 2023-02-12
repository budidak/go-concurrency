package main

import (
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

// initSession initializes the session and configures some settings
func initSession() *scs.SessionManager {
	// setup session
	session := scs.New()
	session.Store = redisstore.New(initRedis()) // store session informations inside redis
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // we want to be persist between visits to the website
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true // true for production phase
	return session
}

// initRedis initializes redis pool which will store the session informations
func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS")) // gets env.var. from Makefile
		},
	}
	return redisPool
}
