package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/rogerwesterbo/familytree/internal/repositories/arangorepository"
	"github.com/rogerwesterbo/familytree/internal/services/v1personservice"
	"github.com/rogerwesterbo/familytree/internal/services/v1relationshipservice"
	"github.com/rogerwesterbo/familytree/pkg/clients/arangodbclient"
	"github.com/rogerwesterbo/familytree/pkg/consts"
	"github.com/spf13/viper"
)

var (
	ArangoClient        *arangodbclient.Client
	PersonService       *v1personservice.PersonService
	RelationshipService *v1relationshipservice.RelationshipService
)

// Init initializes all clients, repositories, and services
func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize ArangoDB client
	host := viper.GetString(consts.ARANGODB_HOST)
	port := viper.GetString(consts.ARANGODB_PORT)
	endpoint := fmt.Sprintf("http://%s:%s", host, port)

	client, err := arangodbclient.NewClient(ctx,
		arangodbclient.WithEndpoints([]string{endpoint}),
		arangodbclient.WithUsername(viper.GetString(consts.ARANGODB_USERNAME)),
		arangodbclient.WithPassword(viper.GetString(consts.ARANGODB_PASSWORD)),
		arangodbclient.WithDatabase(viper.GetString(consts.ARANGODB_DATABASE_NAME)),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize ArangoDB client: %w", err)
	}

	// Test the connection
	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping ArangoDB: %w", err)
	}

	ArangoClient = client

	// Initialize repositories
	personsCollection, err := client.GetCollection(ctx, "persons")
	if err != nil {
		return fmt.Errorf("failed to get persons collection: %w", err)
	}

	relationshipsCollection, err := client.GetCollection(ctx, "relationships")
	if err != nil {
		return fmt.Errorf("failed to get relationships collection: %w", err)
	}

	personRepo := arangorepository.NewPersonRepository(client.GetDatabase(), personsCollection)
	relationshipRepo := arangorepository.NewRelationshipRepository(client.GetDatabase(), relationshipsCollection)

	// Initialize services
	PersonService = v1personservice.NewPersonService(personRepo)
	RelationshipService = v1relationshipservice.NewRelationshipService(relationshipRepo)

	return nil
}
