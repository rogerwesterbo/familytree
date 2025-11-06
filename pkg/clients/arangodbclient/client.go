package arangodbclient

import (
	"context"
	"fmt"

	"strings"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"
	"github.com/spf13/viper"
)

// Client wraps the ArangoDB client with additional functionality
type Client struct {
	conn    arangodb.Client
	db      arangodb.Database
	options *Options
}

// NewClient creates a new ArangoDB client with the given options
func NewClient(ctx context.Context, optFuncs ...OptionFunc) (*Client, error) {
	opts := ApplyOptions(DefaultOptions(), optFuncs...)

	var conn connection.Connection
	var err error

	if viper.GetBool("ARANGODB_HTTP2_ENABLED") {
		conn = connection.NewHttp2Connection(connection.Http2Configuration{
			Endpoint:       connection.NewRoundRobinEndpoints(opts.Endpoints),
			Authentication: connection.NewBasicAuth(opts.Username, opts.Password),
			ContentType:    connection.ApplicationJSON,
		})

	} else {
		conn = connection.NewHttpConnection(connection.HttpConfiguration{

			Authentication: connection.NewBasicAuth(opts.Username, opts.Password),
			ContentType:    connection.ApplicationJSON,
			Endpoint:       connection.NewRoundRobinEndpoints(opts.Endpoints),
		})
	}

	if conn == nil {
		return nil, fmt.Errorf("failed to create ArangoDB connection")
	}

	// Create HTTP connection with authentication
	// conn := connection.NewHttp2Connection(connection.Http2Configuration{
	// 	Endpoint:       connection.NewRoundRobinEndpoints(opts.Endpoints),
	// 	Authentication: connection.NewBasicAuth(opts.Username, opts.Password),
	// 	ContentType:    connection.ApplicationJSON,
	// })

	// conn, err := connection.NewConnection(connection.Config{
	// 	Endpoints:    opts.Endpoints,
	// 	Authentication: connection.NewBasicAuth(opts.Username, opts.Password),
	// 	ContentType: connection.ApplicationJSON,
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create ArangoDB connection: %w", err)
	// }

	// Create client
	client := arangodb.NewClient(conn)

	// Get or create database
	db, err := ensureDatabase(ctx, client, opts.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure database: %w", err)
	}

	c := &Client{
		conn:    client,
		db:      db,
		options: opts,
	}

	// Initialize collections
	if err := c.initializeCollections(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize collections: %w", err)
	}

	return c, nil
}

// ensureDatabase ensures the database exists, creating it if necessary
func ensureDatabase(ctx context.Context, client arangodb.Client, dbName string) (arangodb.Database, error) {
	// Try to get the database - if it exists, return it
	db, err := client.GetDatabase(ctx, dbName, nil)
	if err == nil {
		return db, nil
	}

	// If error is not "not found", return the error
	if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "404") {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// Database doesn't exist, create it
	db, err = client.CreateDatabase(ctx, dbName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return db, nil
}

// initializeCollections creates the necessary collections if they don't exist
func (c *Client) initializeCollections(ctx context.Context) error {
	// Create persons collection (document collection)
	if err := c.ensureCollection(ctx, "persons", false); err != nil {
		return fmt.Errorf("failed to create persons collection: %w", err)
	}

	// Create relationships collection (edge collection)
	if err := c.ensureCollection(ctx, "relationships", true); err != nil {
		return fmt.Errorf("failed to create relationships collection: %w", err)
	}

	return nil
}

// ensureCollection ensures a collection exists, creating it if necessary
func (c *Client) ensureCollection(ctx context.Context, name string, isEdge bool) error {
	exists, err := c.db.CollectionExists(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if exists {
		return nil
	}

	props := &arangodb.CreateCollectionPropertiesV2{}
	if isEdge {
		collType := arangodb.CollectionTypeEdge
		props.Type = &collType
	}

	_, err = c.db.CreateCollectionV2(ctx, name, props)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

// GetDatabase returns the database instance
func (c *Client) GetDatabase() arangodb.Database {
	return c.db
}

// GetCollection returns a collection by name
func (c *Client) GetCollection(ctx context.Context, name string) (arangodb.Collection, error) {
	return c.db.GetCollection(ctx, name, nil)
}

// Close closes the client connection
func (c *Client) Close() error {
	// The v2 driver doesn't have an explicit close method
	// Connection management is handled internally
	return nil
}

// Ping checks if the connection to ArangoDB is alive
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.conn.Version(ctx)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
