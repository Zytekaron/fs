package config

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Server struct {
	// Addr is the address to host on, such as :80 or 127.0.0.1:1337
	Addr string `yaml:"addr"`
	// Tokens is the set of pre-defined authentication tokens
	Tokens *ServerTokens `yaml:"tokens"`
	// RateLimit is the ratelimiting config for the file server
	RateLimit *ServerRateLimit `yaml:"ratelimit"`
}

func (s *Server) Validate() error {
	err := checkBindAddress(s.Addr)
	if err != nil {
		return err
	}
	err = s.Tokens.Validate()
	if err != nil {
		return err
	}
	return s.RateLimit.Validate()
}

type ServerTokens struct {
	// NoRatelimit is the token used to disable ratelimiting
	NoRatelimit string `yaml:"no_ratelimit"`
	// Admin is the token used to allow all access
	Admin string `yaml:"admin"`
}

func (s *ServerTokens) Validate() error {
	return nil
}

type ServerRateLimit struct {
	Limit  int64         `yaml:"limit"`
	Burst  int64         `yaml:"burst"`
	Refill time.Duration `yaml:"refill"`
}

func (s *ServerRateLimit) Validate() error {
	if s.Limit < 0 {
		return errors.New("server.ratelimit.limit cannot be negative")
	}
	if s.Burst < 0 {
		return errors.New("server.ratelimit.burst cannot be negative")
	}
	if s.Refill < 0 {
		return errors.New("server.ratelimit.refill cannot be negative")
	}
	return nil
}

func (s *ServerRateLimit) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data struct {
		Limit  int64  `yaml:"limit"`
		Burst  int64  `yaml:"burst"`
		Refill string `yaml:"refill"`
	}

	err := unmarshal(&data)
	if err != nil {
		return err
	}

	refill, err := time.ParseDuration(data.Refill)
	if err != nil {
		return fmt.Errorf("error parsing duration '%s': %w", data.Refill, err)
	}

	s.Limit = data.Limit
	s.Burst = data.Burst
	s.Refill = refill
	return nil
}

func checkBindAddress(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("server.addr is not a valid bind address: %w", err)
	}
	if host != "" && net.ParseIP(host) == nil {
		return errors.New("server.addr is not a valid bind address")
	}
	_, err = net.LookupPort("tcp", port)
	if err != nil {
		return fmt.Errorf("server.addr is not a valid bind address: %w", err)
	}
	return nil
}
