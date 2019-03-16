package orm

import (
	"fmt"
	"net/url"
	"sync"
)

// GatewayPool represents a gateway pool
type GatewayPool struct {
	// URL is the connection string
	URL string

	// Migrations to be executed on get
	Migrations FileSystem

	// Routines to be loaded on get
	Routines FileSystem

	// Isolated for each gateway instance creates a new schema and set the
	// search path to this schema
	Isolated bool

	m  map[string]*Gateway
	mu sync.RWMutex
}

// Get returns a gateway for given key
func (p *GatewayPool) Get(name string) (*Gateway, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.m == nil {
		p.m = make(map[string]*Gateway)
	}

	gateway, ok := p.m[name]
	if ok {
		return gateway, nil
	}

	addr, err := p.url(name)
	if err != nil {
		return nil, p.error(name, err)
	}

	if gateway, err = Connect(addr); err != nil {
		return nil, p.error(name, err)
	}

	if err = p.schema(gateway, name); err != nil {
		return nil, p.error(name, err)
	}

	if fileSystem := p.Migrations; fileSystem != nil {
		if err = gateway.Migrate(fileSystem); err != nil {
			return nil, p.error(name, err)
		}
	}

	if fileSystem := p.Routines; fileSystem != nil {
		if err = gateway.ReadDir(fileSystem); err != nil {
			return nil, p.error(name, err)
		}
	}

	p.m[name] = gateway

	return gateway, nil
}

// Close closes all gateways
func (p *GatewayPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var errs ErrorSlice

	for key, gateway := range p.m {
		if err := gateway.Close(); err != nil {
			errs = append(errs, p.error(key, err))
		}

		delete(p.m, key)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (p *GatewayPool) url(name string) (string, error) {
	if !p.Isolated {
		return p.URL, nil
	}

	addr, err := url.Parse(p.URL)
	if err != nil {
		return "", err
	}

	if addr.Scheme == "postgres" {
		values := addr.Query()
		values.Set("application_name", name)
		values.Set("search_path", name)

		addr.RawQuery = values.Encode()

		return addr.String(), nil
	}

	return "", fmt.Errorf("not supported driver %q", addr.Scheme)
}

func (p *GatewayPool) schema(gateway *Gateway, name string) error {
	if !p.Isolated {
		return nil
	}

	param := Map{
		"schema": name,
	}

	query := SQL("CREATE SCHEMA IF NOT EXISTS {{schema}};", param)

	_, err := gateway.Exec(query)
	return err
}

func (p *GatewayPool) error(name string, err error) error {
	return fmt.Errorf("name: %v error: %v", name, err)
}
