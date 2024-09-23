package db

import (
	"time"

	"github.com/gocql/gocql"
)

func DatabaseSession() (*gocql.Session, error) {
	cluster := gocql.NewCluster("localhost:9042", "localhost:9043", "localhost:9044") // Fetch from config later
	cluster.Keyspace = "users"                                                        // this can be friends as well
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
