// Package oslookup knows how to get environmental variables
package oslookup

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// EnvFunc has the same signature as os.LookupEnv
type EnvFunc func(string) (string, bool)

// GetAllowedCorsOrigins returns the environment variable's value if exists
func (c CORSHelper) GetAllowedCorsOrigins(osLookupEnv EnvFunc, key string) (map[string]bool, error) {
	val, found := osLookupEnv(key)
	if !found {
		return nil, fmt.Errorf("could not find environment variable %s", key)
	}

	allowed := make(map[string]bool)
	for _, cors := range strings.Split(val, ",") {
		allowed[cors] = true
	}

	return allowed, nil
}

// PrintAllowedCorsOrigins prints allowed cors origins in a nice format
func (c CORSHelper) PrintAllowedCorsOrigins(allowedCorsOrigins map[string]bool) {
	c.logger.Info("Allowed CORS origins:")

	for k := range allowedCorsOrigins {
		c.logger.Infof("- %s\n", k)
	}
}

// CORSHelper helps with CORS stuff
type CORSHelper struct {
	logger *zap.SugaredLogger
}

// NewCORSHelper returns a new CORSHelper
func NewCORSHelper(logger *zap.SugaredLogger) CORSHelper {
	return CORSHelper{
		logger: logger,
	}
}
