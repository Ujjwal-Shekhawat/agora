package db

import (
	"fmt"
	"reflect"
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
	cluster.Keyspace = "guild_messages"                                      // this can be friends as well
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

	cluster.Logger = gocql.Logger

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

func ExecQueryWithResponse(query string, responseStructure interface{}, queryParams ...interface{}) error {
	// if err := session.session.Query(query, queryParams...).MapScan(r); err != nil {
	// 	return nil, err
	// }

	refl := reflect.ValueOf(responseStructure)
	if refl.Kind() != reflect.Ptr || refl.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("responmse struct should be a reference to a slice only")
	}

	iter := session.session.Query(query, queryParams...).Iter()
	defer iter.Close()

	for {
		rStructType := refl.Elem().Type().Elem()
		rStruct := reflect.New(rStructType).Elem()

		scanVars := make([]interface{}, rStruct.NumField())
		for i := 0; i < rStruct.NumField(); i++ {
			scanVars[i] = rStruct.Field(i).Addr().Interface()
		}

		if !iter.Scan(scanVars...) {
			if err := iter.Close(); err != nil {
				return err
			}
			break
		}

		refl.Elem().Set(reflect.Append(refl.Elem(), rStruct))

	}

	if err := iter.Close(); err != nil {
		return err
	}

	return nil
}
