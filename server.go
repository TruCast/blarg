package main

import (
	"encoding/gob"
	"os"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/garyburd/redigo/redis"
	"github.com/pjebs/tokbox"
	"gopkg.in/boj/redistore.v1"
)

const (
	Development = iota
	Production
)

// Environment Variables:

// TOKBOX_KEY
// TOKBOX_SECRET_KEY
// BLARG_ENV
// TWITTER_KEY
// TWITTER_SECRET_KEY
// REDIS_PORT_6379_TCP_ADDR
// REDIS_PORT_6379_TCP_POST
// REDIS_PASSWORD

type Server struct {
	RediStore     *redistore.RediStore
	RedisPool     *redis.Pool
	OAuthClient   oauth.Client
	State         State
	ClientBaseUrl string
	ForceHTTPS    bool
	TokBox        *tokbox.Tokbox
	TokBoxKey     string
}

func newPool(addr, password string) *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 100,

		// Maximum number of connections allocated by the pool at a given time.
		// When zero, there is no limit on the number of connections in the pool.
		MaxActive: 100,

		// Close connections after remaining idle for this duration. If the value
		// is zero, then idle connections are not closed. Applications should set
		// the timeout to a value less than the server's timeout.
		IdleTimeout: 240 * time.Second,

		// Dial is an application supplied function for creating and configuring a
		// connection
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
	}
}

func NewServer() Server {
	// Get environment
	environment := Development
	if os.Getenv("BLARG_ENV") == "production" {
		environment = Production
	}

	// Redis
	redisAddr := os.Getenv("REDIS_PORT_6379_TCP_ADDR") +
		":" + os.Getenv("REDIS_PORT_6379_TCP_PORT")
	redisPool := newPool(redisAddr, os.Getenv("REDIS_PASSWORD"))

	// Redis store
	rediStore, err := redistore.NewRediStoreWithPool(redisPool, []byte("secret string "))
	if err != nil {
		panic(err)
	}

	// TokBox
	tokBoxKey := os.Getenv("TOKBOX_KEY")
	tb := tokbox.New(
		tokBoxKey,
		os.Getenv("TOKBOX_SECRET_KEY"),
	)

	// OAuth
	gob.Register(&oauth.Credentials{})
	oauthClient := oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}
	oauthClient.Credentials.Token = os.Getenv("TWITTER_KEY")
	oauthClient.Credentials.Secret = os.Getenv("TWITTER_SECRET_KEY")

	// State
	state := NewLocalState()

	// Config
	clientBaseUrl := "http://localhost:8000"
	forceHTTPS := false
	if environment == Production {
		clientBaseUrl = "https://blarg-im.herokuapp.com"
		forceHTTPS = true
	}

	return Server{
		RedisPool:     redisPool,
		RediStore:     rediStore,
		OAuthClient:   oauthClient,
		State:         &state,
		ClientBaseUrl: clientBaseUrl,
		ForceHTTPS:    forceHTTPS,
		TokBox:        tb,
		TokBoxKey:     tokBoxKey,
	}
}
