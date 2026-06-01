//go:build !windows

package container

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"go.uber.org/dig"

	memoryRepo "github.com/Tencent/WeKnora/internal/application/repository/memory/neo4j"
	neo4jRepo "github.com/Tencent/WeKnora/internal/application/repository/retriever/neo4j"
	"github.com/Tencent/WeKnora/internal/logger"
)

// registerNeo4jProviders registers neo4j-related providers in the DI container.
// On non-Windows platforms the full neo4j driver is available.
func registerNeo4jProviders(c *dig.Container) {
	must(c.Provide(initNeo4jClient))
	must(c.Provide(neo4jRepo.NewNeo4jRepository))
	must(c.Provide(memoryRepo.NewMemoryRepository))
	must(c.Provide(provideGraphEngine))
}

// provideGraphEngine returns the graph engine name for the system handler.
func provideGraphEngine(driver neo4j.Driver) string {
	if driver == nil {
		return ""
	}
	return "Neo4j"
}

func initNeo4jClient() (neo4j.Driver, error) {
	ctx := context.Background()
	if strings.ToLower(os.Getenv("NEO4J_ENABLE")) != "true" {
		logger.Debugf(ctx, "NOT SUPPORT RETRIEVE GRAPH")
		return nil, nil
	}
	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	// Retry configuration
	maxRetries := 30                 // Max retry attempts
	retryInterval := 2 * time.Second // Wait between retries

	var driver neo4j.Driver
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
		if err != nil {
			logger.Warnf(ctx, "Failed to create Neo4j driver (attempt %d/%d): %v", attempt, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		err = driver.VerifyAuthentication(ctx, nil)
		if err == nil {
			if attempt > 1 {
				logger.Infof(ctx, "Successfully connected to Neo4j after %d attempts", attempt)
			}
			return driver, nil
		}

		logger.Warnf(ctx, "Failed to verify Neo4j authentication (attempt %d/%d): %v", attempt, maxRetries, err)
		driver.Close(ctx)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to Neo4j after %d attempts: %w", maxRetries, err)
}
