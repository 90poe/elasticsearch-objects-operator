package elasticsearch

import (
	"fmt"

	"github.com/olivere/elastic/v7"
)

// Option is a type of options for Executor
type Option func(*Client) error

// Client is structure of ES client for XO
type Client struct {
	esURL string
	es    *elastic.Client
}

// URL is option function to set ES URL for Client
func URL(esURL string) Option {
	return func(c *Client) error {
		if len(esURL) == 0 {
			return fmt.Errorf("url for ES can't be empty")
		}
		c.esURL = esURL
		return nil
	}
}

// ESclient is option function to set ES cluster object - mainly for mocking
func ESclient(es *elastic.Client) Option {
	return func(c *Client) error {
		if es == nil {
			return fmt.Errorf("es cluster must not be nil")
		}
		c.es = es
		return nil
	}
}

// New would create ES Client
func New(options ...Option) (*Client, error) {
	c := Client{}
	var err error
	for _, option := range options {
		err = option(&c)
		if err != nil {
			return nil, fmt.Errorf("can't make new ES Client: %w", err)
		}
	}
	if c.es != nil {
		return &c, nil
	}
	// ES cluster client not provided - create one
	c.es, err = elastic.NewClient(
		elastic.SetURL(c.esURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, fmt.Errorf("can't make new ES Client: %w", err)
	}
	return &c, nil
}
