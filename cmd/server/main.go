package main

import (
	"os"

	"github.com/revernelly/mcp-server-shop/config"
	"github.com/revernelly/mcp-server-shop/internal/client"
	"github.com/revernelly/mcp-server-shop/internal/mcp"
	"github.com/revernelly/mcp-server-shop/internal/tools/products"
	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()
	logger.SetOutput(os.Stderr) // Use stderr for logs (stdout is for JSON-RPC)
	logger.SetFormatter(&logrus.JSONFormatter{})

	config, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logger.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)

	logger.WithFields(logrus.Fields{
		"api_url":               config.APIURL,
		"log_level":             config.LogLevel,
		"auth_token_configured": config.AuthToken != "",
	}).Info("Starting MCP E-commerce Server")

	restClient := client.NewRestClient(config.APIURL, config.AuthToken, logger)

	toolRegistry := mcp.NewRegistry(logger)

	products.NewProductToolset(toolRegistry, restClient, logger)

	logger.WithField("tool_count", len(toolRegistry.ListTools())).Info("Registered tools")

	mcpServer := mcp.NewServer(toolRegistry, logger)

	logger.Info("Starting server on stdio transport")

	if err := mcpServer.Start(); err != nil {
		logger.WithError(err).Fatal("Server error")
	}
}
