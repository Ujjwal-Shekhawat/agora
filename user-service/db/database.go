package db

import (
	"strings"
	"sync"
	"time"
	"user_service/config"

	"github.com/gocql/gocql"
)

// Make this readonly later by making its methods exported later
type dbSession struct {
	session *gocql.Session
	closed  bool
}

var mu sync.Mutex

var session *dbSession = &dbSession{
	session: nil,
	closed:  true,
}

func InitSession() error {
	dbSession, err := DatabaseSession()
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	session.session = dbSession
	session.closed = false
	return nil
}

func DatabaseSession() (*gocql.Session, error) {
	cfg := config.LoadConfig()
	cluster := gocql.NewCluster(strings.Split(cfg.CassandraCluster, ",")...) // Fetch from config later
	cluster.Keyspace = "users"                                               // this can be friends as well
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	// Connection pooling settings
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.NumConns = 5
	cluster.SocketKeepalive = 3 * time.Second

	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{
		NumRetries: 3,
	}

	cluster.Timeout = 3 * time.Second
	cluster.ConnectTimeout = 3 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func ExecQuery(query string, queryParams ...interface{}) error {
	if err := session.session.Query(query, queryParams...).Exec(); err != nil {
		return err
	}
	return nil
}
