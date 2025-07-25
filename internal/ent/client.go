// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"auroraride.com/edocseal/internal/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"auroraride.com/edocseal/internal/ent/certification"
	"auroraride.com/edocseal/internal/ent/document"

	stdsql "database/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Certification is the client for interacting with the Certification builders.
	Certification *CertificationClient
	// Document is the client for interacting with the Document builders.
	Document *DocumentClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Certification = NewCertificationClient(c.config)
	c.Document = NewDocumentClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:           ctx,
		config:        cfg,
		Certification: NewCertificationClient(cfg),
		Document:      NewDocumentClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:           ctx,
		config:        cfg,
		Certification: NewCertificationClient(cfg),
		Document:      NewDocumentClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Certification.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Certification.Use(hooks...)
	c.Document.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Certification.Intercept(interceptors...)
	c.Document.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *CertificationMutation:
		return c.Certification.mutate(ctx, m)
	case *DocumentMutation:
		return c.Document.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// CertificationClient is a client for the Certification schema.
type CertificationClient struct {
	config
}

// NewCertificationClient returns a client for the Certification from the given config.
func NewCertificationClient(c config) *CertificationClient {
	return &CertificationClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `certification.Hooks(f(g(h())))`.
func (c *CertificationClient) Use(hooks ...Hook) {
	c.hooks.Certification = append(c.hooks.Certification, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `certification.Intercept(f(g(h())))`.
func (c *CertificationClient) Intercept(interceptors ...Interceptor) {
	c.inters.Certification = append(c.inters.Certification, interceptors...)
}

// Create returns a builder for creating a Certification entity.
func (c *CertificationClient) Create() *CertificationCreate {
	mutation := newCertificationMutation(c.config, OpCreate)
	return &CertificationCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Certification entities.
func (c *CertificationClient) CreateBulk(builders ...*CertificationCreate) *CertificationCreateBulk {
	return &CertificationCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *CertificationClient) MapCreateBulk(slice any, setFunc func(*CertificationCreate, int)) *CertificationCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &CertificationCreateBulk{err: fmt.Errorf("calling to CertificationClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*CertificationCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &CertificationCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Certification.
func (c *CertificationClient) Update() *CertificationUpdate {
	mutation := newCertificationMutation(c.config, OpUpdate)
	return &CertificationUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CertificationClient) UpdateOne(ce *Certification) *CertificationUpdateOne {
	mutation := newCertificationMutation(c.config, OpUpdateOne, withCertification(ce))
	return &CertificationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *CertificationClient) UpdateOneID(id int) *CertificationUpdateOne {
	mutation := newCertificationMutation(c.config, OpUpdateOne, withCertificationID(id))
	return &CertificationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Certification.
func (c *CertificationClient) Delete() *CertificationDelete {
	mutation := newCertificationMutation(c.config, OpDelete)
	return &CertificationDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *CertificationClient) DeleteOne(ce *Certification) *CertificationDeleteOne {
	return c.DeleteOneID(ce.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *CertificationClient) DeleteOneID(id int) *CertificationDeleteOne {
	builder := c.Delete().Where(certification.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CertificationDeleteOne{builder}
}

// Query returns a query builder for Certification.
func (c *CertificationClient) Query() *CertificationQuery {
	return &CertificationQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeCertification},
		inters: c.Interceptors(),
	}
}

// Get returns a Certification entity by its id.
func (c *CertificationClient) Get(ctx context.Context, id int) (*Certification, error) {
	return c.Query().Where(certification.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CertificationClient) GetX(ctx context.Context, id int) *Certification {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *CertificationClient) Hooks() []Hook {
	return c.hooks.Certification
}

// Interceptors returns the client interceptors.
func (c *CertificationClient) Interceptors() []Interceptor {
	return c.inters.Certification
}

func (c *CertificationClient) mutate(ctx context.Context, m *CertificationMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&CertificationCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&CertificationUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&CertificationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&CertificationDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Certification mutation op: %q", m.Op())
	}
}

// DocumentClient is a client for the Document schema.
type DocumentClient struct {
	config
}

// NewDocumentClient returns a client for the Document from the given config.
func NewDocumentClient(c config) *DocumentClient {
	return &DocumentClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `document.Hooks(f(g(h())))`.
func (c *DocumentClient) Use(hooks ...Hook) {
	c.hooks.Document = append(c.hooks.Document, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `document.Intercept(f(g(h())))`.
func (c *DocumentClient) Intercept(interceptors ...Interceptor) {
	c.inters.Document = append(c.inters.Document, interceptors...)
}

// Create returns a builder for creating a Document entity.
func (c *DocumentClient) Create() *DocumentCreate {
	mutation := newDocumentMutation(c.config, OpCreate)
	return &DocumentCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Document entities.
func (c *DocumentClient) CreateBulk(builders ...*DocumentCreate) *DocumentCreateBulk {
	return &DocumentCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *DocumentClient) MapCreateBulk(slice any, setFunc func(*DocumentCreate, int)) *DocumentCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &DocumentCreateBulk{err: fmt.Errorf("calling to DocumentClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*DocumentCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &DocumentCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Document.
func (c *DocumentClient) Update() *DocumentUpdate {
	mutation := newDocumentMutation(c.config, OpUpdate)
	return &DocumentUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *DocumentClient) UpdateOne(d *Document) *DocumentUpdateOne {
	mutation := newDocumentMutation(c.config, OpUpdateOne, withDocument(d))
	return &DocumentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *DocumentClient) UpdateOneID(id string) *DocumentUpdateOne {
	mutation := newDocumentMutation(c.config, OpUpdateOne, withDocumentID(id))
	return &DocumentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Document.
func (c *DocumentClient) Delete() *DocumentDelete {
	mutation := newDocumentMutation(c.config, OpDelete)
	return &DocumentDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *DocumentClient) DeleteOne(d *Document) *DocumentDeleteOne {
	return c.DeleteOneID(d.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *DocumentClient) DeleteOneID(id string) *DocumentDeleteOne {
	builder := c.Delete().Where(document.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &DocumentDeleteOne{builder}
}

// Query returns a query builder for Document.
func (c *DocumentClient) Query() *DocumentQuery {
	return &DocumentQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeDocument},
		inters: c.Interceptors(),
	}
}

// Get returns a Document entity by its id.
func (c *DocumentClient) Get(ctx context.Context, id string) (*Document, error) {
	return c.Query().Where(document.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *DocumentClient) GetX(ctx context.Context, id string) *Document {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *DocumentClient) Hooks() []Hook {
	return c.hooks.Document
}

// Interceptors returns the client interceptors.
func (c *DocumentClient) Interceptors() []Interceptor {
	return c.inters.Document
}

func (c *DocumentClient) mutate(ctx context.Context, m *DocumentMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&DocumentCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&DocumentUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&DocumentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&DocumentDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Document mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Certification, Document []ent.Hook
	}
	inters struct {
		Certification, Document []ent.Interceptor
	}
)

// ExecContext allows calling the underlying ExecContext method of the driver if it is supported by it.
// See, database/sql#DB.ExecContext for more information.
func (c *config) ExecContext(ctx context.Context, query string, args ...any) (stdsql.Result, error) {
	ex, ok := c.driver.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	return ex.ExecContext(ctx, query, args...)
}

// QueryContext allows calling the underlying QueryContext method of the driver if it is supported by it.
// See, database/sql#DB.QueryContext for more information.
func (c *config) QueryContext(ctx context.Context, query string, args ...any) (*stdsql.Rows, error) {
	q, ok := c.driver.(interface {
		QueryContext(context.Context, string, ...any) (*stdsql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	return q.QueryContext(ctx, query, args...)
}
