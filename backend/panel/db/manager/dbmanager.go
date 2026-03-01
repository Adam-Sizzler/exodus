package manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"v2ray-stat/backend/panel/config"
	"v2ray-stat/backend/panel/dbutil"
)

const (
	defaultMaxOpenConns = 10
	defaultMaxIdleConns = 5
)

// DBExecutor abstracts SQL execution with rebinding support.
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (TxExecutor, error)
}

// TxExecutor abstracts SQL transactions with rebinding support.
type TxExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Commit() error
	Rollback() error
}

// RebindDB wraps *sql.DB to translate '?' placeholders into PostgreSQL-style '$1'.
type RebindDB struct {
	db *sql.DB
}

// NewRebindDB creates a rebind wrapper for a database handle.
func NewRebindDB(db *sql.DB) *RebindDB {
	return &RebindDB{db: db}
}

// BeginTx begins a transaction with rebinding.
func (r *RebindDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (TxExecutor, error) {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &RebindTx{tx: tx}, nil
}

func (r *RebindDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return r.db.ExecContext(ctx, dbutil.Rebind(query), args...)
}

func (r *RebindDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, dbutil.Rebind(query), args...)
}

func (r *RebindDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return r.db.QueryRowContext(ctx, dbutil.Rebind(query), args...)
}

func (r *RebindDB) Exec(query string, args ...any) (sql.Result, error) {
	return r.db.Exec(dbutil.Rebind(query), args...)
}

func (r *RebindDB) Query(query string, args ...any) (*sql.Rows, error) {
	return r.db.Query(dbutil.Rebind(query), args...)
}

func (r *RebindDB) QueryRow(query string, args ...any) *sql.Row {
	return r.db.QueryRow(dbutil.Rebind(query), args...)
}

// RebindTx wraps *sql.Tx to translate '?' placeholders into PostgreSQL-style '$1'.
type RebindTx struct {
	tx *sql.Tx
}

func (r *RebindTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return r.tx.ExecContext(ctx, dbutil.Rebind(query), args...)
}

func (r *RebindTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return r.tx.QueryContext(ctx, dbutil.Rebind(query), args...)
}

func (r *RebindTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return r.tx.QueryRowContext(ctx, dbutil.Rebind(query), args...)
}

func (r *RebindTx) Exec(query string, args ...any) (sql.Result, error) {
	return r.tx.Exec(dbutil.Rebind(query), args...)
}

func (r *RebindTx) Query(query string, args ...any) (*sql.Rows, error) {
	return r.tx.Query(dbutil.Rebind(query), args...)
}

func (r *RebindTx) QueryRow(query string, args ...any) *sql.Row {
	return r.tx.QueryRow(dbutil.Rebind(query), args...)
}

func (r *RebindTx) Commit() error {
	return r.tx.Commit()
}

func (r *RebindTx) Rollback() error {
	return r.tx.Rollback()
}

// DatabaseManager manages sequential database access through prioritized request channels.
type DatabaseManager struct {
	raw               *sql.DB                     // Raw database connection
	db                DBExecutor                  // Private field for database connection
	cfg               *config.BackendConfig       // Configuration for logging
	highPriority      chan func(DBExecutor) error // Channel for high-priority requests
	lowPriority       chan func(DBExecutor) error // Channel for low-priority requests
	ctx               context.Context             // Context for cancellation
	cancel            context.CancelFunc          // Cancel function for context
	workerPool        chan struct{}               // Worker pool for concurrent request processing
	isClosed          bool                        // Flag indicating if the manager is closed
	closedMu          sync.Mutex                  // Mutex for isClosed
	highPriorityCount uint64                      // Counter for high-priority requests
	lowPriorityCount  uint64                      // Counter for low-priority requests
}

// NewDatabaseManager creates a new DatabaseManager and starts processing requests.
func NewDatabaseManager(db *sql.DB, ctx context.Context, workerCount, highPriorityBuffer, lowPriorityBuffer int, cfg *config.BackendConfig) (*DatabaseManager, error) {
	if workerCount < 1 || highPriorityBuffer < 0 || lowPriorityBuffer < 0 {
		cfg.Logger.Fatal("Invalid parameters", "workerCount", workerCount, "highPriorityBuffer", highPriorityBuffer, "lowPriorityBuffer", lowPriorityBuffer)
		return nil, fmt.Errorf("invalid parameters: workerCount=%d, highPriorityBuffer=%d, lowPriorityBuffer=%d", workerCount, highPriorityBuffer, lowPriorityBuffer)
	}

	// Configure connection pool
	maxOpen := defaultMaxOpenConns
	maxIdle := defaultMaxIdleConns
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(0)

	mCtx, cancel := context.WithCancel(ctx)
	manager := &DatabaseManager{
		raw:          db,
		db:           NewRebindDB(db),
		cfg:          cfg,
		highPriority: make(chan func(DBExecutor) error, highPriorityBuffer),
		lowPriority:  make(chan func(DBExecutor) error, lowPriorityBuffer),
		ctx:          mCtx,
		cancel:       cancel,
		workerPool:   make(chan struct{}, workerCount),
	}

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		go manager.processRequests(i)
	}
	cfg.Logger.Debug("DatabaseManager created", "workerCount", workerCount, "highPriorityBuffer", highPriorityBuffer, "lowPriorityBuffer", lowPriorityBuffer)
	return manager, nil
}

// processRequests handles requests from prioritized channels with strict prioritization.
func (m *DatabaseManager) processRequests(workerID int) {
	for {
		select {
		case <-m.ctx.Done():
			m.cfg.Logger.Debug("Worker stopped", "workerID", workerID)
			return
		case req, ok := <-m.highPriority:
			if !ok {
				m.cfg.Logger.Debug("High-priority channel closed", "workerID", workerID)
				return
			}
			atomic.AddUint64(&m.highPriorityCount, 1)
			m.processRequest(req, workerID, "highPriority")
		default:
			select {
			case req, ok := <-m.lowPriority:
				if !ok {
					m.cfg.Logger.Debug("Low-priority channel closed", "workerID", workerID)
					return
				}
				atomic.AddUint64(&m.lowPriorityCount, 1)
				m.processRequest(req, workerID, "lowPriority")
			case <-m.ctx.Done():
				m.cfg.Logger.Debug("Worker stopped", "workerID", workerID)
				return
			default:
				time.Sleep(10 * time.Millisecond) // Avoid busy-waiting
			}
		}
	}
}

// processRequest executes a single database request.
func (m *DatabaseManager) processRequest(req func(DBExecutor) error, workerID int, priority string) {
	m.workerPool <- struct{}{}
	defer func() { <-m.workerPool }()
	m.cfg.Logger.Trace("Processing request", "workerID", workerID, "priority", priority)
	start := time.Now()
	if err := req(m.db); err != nil {
		m.cfg.Logger.Error("Failed to execute request", "workerID", workerID, "priority", priority, "error", err)
	}
	m.cfg.Logger.Debug("Processed request", "workerID", workerID, "priority", priority, "duration", time.Since(start))
	m.cfg.Logger.Debug("Processed request counts", "highPriority", atomic.LoadUint64(&m.highPriorityCount), "lowPriority", atomic.LoadUint64(&m.lowPriorityCount))
}

// executeOnce executes a database request once with timeout handling.
func (m *DatabaseManager) executeOnce(fn func(DBExecutor) error, priority bool, sendTimeout, waitTimeout time.Duration) error {
	errChan := make(chan error, 1)
	m.closedMu.Lock()
	if m.isClosed {
		m.closedMu.Unlock()
		m.cfg.Logger.Error("DatabaseManager is closed")
		return fmt.Errorf("DatabaseManager is closed")
	}
	m.closedMu.Unlock()

	// Select channel based on priority
	requestChan, priorityStr := m.lowPriority, "lowPriority"
	if priority {
		requestChan, priorityStr = m.highPriority, "highPriority"
	}

	// Log channel status before sending request
	m.cfg.Logger.Debug("Channel status before sending request",
		"priority", priorityStr,
		"tasks", len(requestChan),
		"capacity", cap(requestChan))

	// Send request to channel
	select {
	case requestChan <- func(db DBExecutor) error {
		err := fn(db)
		if err == sql.ErrNoRows {
			errChan <- nil
			return nil
		}
		errChan <- err
		return err
	}:
		m.cfg.Logger.Debug("Request sent to channel", "priority", priorityStr)
	case <-m.ctx.Done():
		m.cfg.Logger.Warn("Context canceled while sending request", "priority", priorityStr)
		return m.ctx.Err()
	case <-time.After(sendTimeout):
		m.cfg.Logger.Error("Request send timeout", "priority", priorityStr, "timeout", sendTimeout)
		return fmt.Errorf("request send timeout (%s, %v)", priorityStr, sendTimeout)
	}

	// Wait for request execution
	select {
	case err := <-errChan:
		m.cfg.Logger.Trace("Received response from request", "priority", priorityStr, "error", err)
		return err
	case <-m.ctx.Done():
		m.cfg.Logger.Warn("Context canceled while waiting for response", "priority", priorityStr)
		return m.ctx.Err()
	case <-time.After(waitTimeout):
		m.cfg.Logger.Error("Response wait timeout", "priority", priorityStr, "timeout", waitTimeout)
		return fmt.Errorf("response wait timeout (%s, %v)", priorityStr, waitTimeout)
	}
}

// isRetryableError checks if an error is retryable.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "database is locked") ||
		strings.Contains(err.Error(), "i/o timeout") ||
		strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no such file or directory")
}

// ExecuteWithTimeout executes a database request with retries and timeout handling.
func (m *DatabaseManager) ExecuteWithTimeout(fn func(DBExecutor) error, priority bool, sendTimeout, waitTimeout time.Duration) error {
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := m.executeOnce(fn, priority, sendTimeout, waitTimeout); err == nil {
			m.cfg.Logger.Debug("Request executed successfully", "attempt", attempt, "priority", priority)
			return nil
		} else if isRetryableError(err) {
			m.cfg.Logger.Warn("Retryable error, attempting retry", "attempt", attempt, "maxRetries", maxRetries, "error", err)
			time.Sleep(time.Duration(attempt*100) * time.Millisecond)
			continue
		} else {
			m.cfg.Logger.Error("Failed to execute request", "attempt", attempt, "error", err)
			return err
		}
	}
	err := fmt.Errorf("failed to execute request after %d retries", maxRetries)
	m.cfg.Logger.Error("Failed to execute request after all retries", "maxRetries", maxRetries)
	return err
}

// ExecuteHighPriority executes a high-priority database request.
func (m *DatabaseManager) ExecuteHighPriority(fn func(DBExecutor) error) error {
	return m.ExecuteWithTimeout(fn, true, 5*time.Second, 10*time.Second)
}

// ExecuteLowPriority executes a low-priority database request.
func (m *DatabaseManager) ExecuteLowPriority(fn func(DBExecutor) error) error {
	return m.ExecuteWithTimeout(fn, false, 5*time.Second, 15*time.Second)
}

// Close shuts down the DatabaseManager and completes request processing.
func (m *DatabaseManager) Close() {
	m.closedMu.Lock()
	if m.isClosed {
		m.closedMu.Unlock()
		m.cfg.Logger.Warn("DatabaseManager already closed")
		return
	}
	m.isClosed = true
	m.closedMu.Unlock()

	m.cfg.Logger.Debug("Shutting down DatabaseManager")
	m.cancel() // Stop accepting new requests

	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			m.cfg.Logger.Warn("Timeout during shutdown, forcing channel closure")
			close(m.highPriority)
			close(m.lowPriority)
			return
		case <-ticker.C:
			if len(m.highPriority) == 0 && len(m.lowPriority) == 0 {
				m.cfg.Logger.Debug("All requests processed, closing channels")
				close(m.highPriority)
				close(m.lowPriority)
				return
			}
			m.cfg.Logger.Debug("Waiting for requests", "highPriority", len(m.highPriority), "lowPriority", len(m.lowPriority))
		}
	}
}

// DB returns the sql.DB pointer (use only for initialization or testing; prefer ExecuteLowPriority/ExecuteHighPriority for database operations).
func (m *DatabaseManager) DB() *sql.DB {
	m.cfg.Logger.Warn("Direct access to DB() should be avoided; use ExecuteLowPriority or ExecuteHighPriority")
	return m.raw
}
