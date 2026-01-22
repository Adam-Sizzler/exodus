package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"

	_ "github.com/mattn/go-sqlite3"
)

// OpenAndInitDB opens and initializes a SQLite database.
func OpenAndInitDB(dbPath string, dbType string, cfg *config.BackendConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		cfg.Logger.Error("Failed to open database", "dbType", dbType, "path", dbPath, "error", err)
		return nil, fmt.Errorf("failed to open %s database: %v", dbType, err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	var tableCount int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='node_name'").Scan(&tableCount)
	if err != nil {
		cfg.Logger.Error("Failed to check table existence", "dbType", dbType, "error", err)
		db.Close()
		return nil, fmt.Errorf("failed to check table existence for %s database: %v", dbType, err)
	}
	if tableCount > 0 {
		cfg.Logger.Debug("Tables already exist", "dbType", dbType)
		return db, nil
	}

	sqlStmt := `
		PRAGMA foreign_keys = ON;
        PRAGMA cache_size = 2000;
        PRAGMA journal_mode = WAL;
        PRAGMA synchronous = NORMAL;
        PRAGMA temp_store = MEMORY;
        PRAGMA busy_timeout = 5000;

		-- Таблица для нод
		CREATE TABLE IF NOT EXISTS nodes (
			node_name TEXT PRIMARY KEY,
			address TEXT NOT NULL,
			port TEXT NOT NULL,
			CONSTRAINT unique_address_port UNIQUE (address, port)
		);

		-- Таблица для статистики трафика по нодам
		CREATE TABLE IF NOT EXISTS bound_traffic (
			node_name TEXT,
			source TEXT,
			rate INTEGER DEFAULT 0,
			uplink INTEGER DEFAULT 0,
			downlink INTEGER DEFAULT 0,
			sess_uplink INTEGER DEFAULT 0,
			sess_downlink INTEGER DEFAULT 0,
			PRIMARY KEY (node_name, source),
			FOREIGN KEY (node_name) REFERENCES nodes(node_name) ON DELETE CASCADE
		);

		-- Таблица для данных пользователей
		CREATE TABLE IF NOT EXISTS user_data (
			user TEXT PRIMARY KEY,
			traffic_cap INTEGER DEFAULT 0,
			sub_end INTEGER DEFAULT 0,
			renew INTEGER DEFAULT 0,
			lim_ip INTEGER DEFAULT 0,
			ips TEXT DEFAULT ''
		);

		-- Таблица для статистики пользователей по нодам
		CREATE TABLE IF NOT EXISTS user_traffic (
			node_name TEXT,
			user TEXT,
			last_seen INTEGER DEFAULT 0,
			rate INTEGER DEFAULT 0,
			uplink INTEGER DEFAULT 0,
			downlink INTEGER DEFAULT 0,
			sess_uplink INTEGER DEFAULT 0,
			sess_downlink INTEGER DEFAULT 0,
			created INTEGER DEFAULT 0,
			enabled TEXT DEFAULT 'true',
			PRIMARY KEY (node_name, user),
			FOREIGN KEY (node_name) REFERENCES nodes(node_name) ON DELETE CASCADE,
			FOREIGN KEY (user) REFERENCES user_data(user) ON DELETE CASCADE
		);

		-- Таблица для ID пользователей
		CREATE TABLE IF NOT EXISTS user_ids (
			node_name TEXT,
			user TEXT,
			id TEXT,
			inbound_tag TEXT,
			PRIMARY KEY (node_name, user, id, inbound_tag),
			FOREIGN KEY (node_name, user) REFERENCES user_traffic(node_name, user) ON DELETE CASCADE,
			FOREIGN KEY (node_name) REFERENCES nodes(node_name) ON DELETE CASCADE
		);

		-- Таблица для DNS-статистики
		CREATE TABLE IF NOT EXISTS user_dns (
			node_name TEXT,
			user TEXT,
			count INTEGER DEFAULT 1,
			domain TEXT,
			PRIMARY KEY (node_name, user, domain),
			FOREIGN KEY (node_name) REFERENCES nodes(node_name) ON DELETE CASCADE,
			FOREIGN KEY (user) REFERENCES user_data(user) ON DELETE CASCADE
		);

		-- Триггер для добавления пользователя в user_data при вставке в user_traffic
		CREATE TRIGGER IF NOT EXISTS insert_user_traffic_trigger
		AFTER INSERT ON user_traffic
		BEGIN
			INSERT OR IGNORE INTO user_data (user) VALUES (NEW.user);
		END;

		-- Триггер для удаления пользователя из user_data, если он больше не существует в user_traffic
		CREATE TRIGGER IF NOT EXISTS delete_user_traffic_trigger
		AFTER DELETE ON user_traffic
		BEGIN
			DELETE FROM user_data
			WHERE user = OLD.user
			AND NOT EXISTS (
				SELECT 1 FROM user_traffic WHERE user = OLD.user
			);
		END;

		-- Индексы для оптимизации запросов
		CREATE INDEX IF NOT EXISTS idx_user_traffic_user ON user_traffic(user);
		CREATE INDEX IF NOT EXISTS idx_user_traffic_rate ON user_traffic(rate);
		CREATE INDEX IF NOT EXISTS idx_user_traffic_last_seen ON user_traffic(last_seen);
		CREATE INDEX IF NOT EXISTS idx_user_traffic_sess_uplink ON user_traffic(sess_uplink);
		CREATE INDEX IF NOT EXISTS idx_user_traffic_sess_downlink ON user_traffic(sess_downlink);
		CREATE INDEX IF NOT EXISTS idx_user_traffic_uplink ON user_traffic(uplink);
		CREATE INDEX IF NOT EXISTS idx_user_traffic_downlink ON user_traffic(downlink);
		CREATE INDEX IF NOT EXISTS idx_user_dns_domain ON user_dns(domain);

		-- Таблица для головных заданий
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			operation TEXT NOT NULL,
			payload TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			created_at INTEGER NOT NULL,
			completed_at INTEGER,
			timeout_at INTEGER NOT NULL
		);

		-- Таблица для дочерних заданий
		CREATE TABLE IF NOT EXISTS task_nodes (
			id TEXT PRIMARY KEY,
			task_id TEXT NOT NULL,
			node_name TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			error_message TEXT,
			sent_at INTEGER,
			completed_at INTEGER,
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
			FOREIGN KEY (node_name) REFERENCES nodes(node_name) ON DELETE CASCADE
		);

		-- Индексы для таблиц заданий
		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
		CREATE INDEX IF NOT EXISTS idx_tasks_timeout ON tasks(timeout_at);
		CREATE INDEX IF NOT EXISTS idx_task_nodes_task_id ON task_nodes(task_id);
		CREATE INDEX IF NOT EXISTS idx_task_nodes_status ON task_nodes(status);
    `
	if _, err = db.Exec(sqlStmt); err != nil {
		cfg.Logger.Error("Failed to execute SQL script", "dbType", dbType, "error", err)
		db.Close()
		return nil, fmt.Errorf("failed to execute SQL script for %s database: %v", dbType, err)
	}

	for _, node := range cfg.Nodes {
		_, err = db.Exec("INSERT OR IGNORE INTO nodes (node_name, address, port) VALUES (?, ?, ?)", node.NodeName, node.Address, node.Port)
		if err != nil {
			cfg.Logger.Error("Failed to insert node", "node_name", node.NodeName, "address", node.Address, "port", node.Port, "error", err)
		}
	}

	cfg.Logger.Debug("Database initialized", "dbType", dbType)
	return db, nil
}

// InitDatabase initializes in-memory and file databases.
func InitDatabase(cfg *config.BackendConfig) (memDB, fileDB *sql.DB, err error) {
	memDB, err = OpenAndInitDB(":memory:", "in-memory", cfg)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			memDB.Close()
		}
	}()

	fileExists := true
	if _, err := os.Stat(cfg.Paths.Database); os.IsNotExist(err) {
		cfg.Logger.Warn("File database does not exist, will create new", "path", cfg.Paths.Database)
		fileExists = false
	} else if err != nil {
		cfg.Logger.Error("Failed to check file database", "path", cfg.Paths.Database, "error", err)
		return nil, nil, fmt.Errorf("error checking file database: %v", err)
	}

	fileDB, err = OpenAndInitDB(cfg.Paths.Database, "file", cfg)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			fileDB.Close()
		}
	}()

	if fileExists {
		var tableCount int
		err = fileDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='user_traffic'").Scan(&tableCount)
		if err == nil && tableCount > 0 {
			tempManager, err := manager.NewDatabaseManager(fileDB, context.Background(), 1, 300, 500, cfg)
			if err != nil {
				cfg.Logger.Error("Failed to create temporary DatabaseManager", "error", err)
				return nil, nil, fmt.Errorf("failed to create temporary DatabaseManager: %v", err)
			}
			syncCtx, syncCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer syncCancel()
			if err = tempManager.SyncDBWithContext(syncCtx, memDB, "file to memory"); err != nil {
				cfg.Logger.Error("Failed to synchronize database (file to memory)", "error", err)
			}
			tempManager.Close()
		}
	}

	cfg.Logger.Debug("Database initialization completed", "in-memory", true, "file", true)
	return memDB, fileDB, nil
}

func MonitorSubscriptionsAndSync(ctx context.Context, manager *manager.DatabaseManager, fileDB *sql.DB, cfg *config.BackendConfig, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		cfg.Logger.Debug("Starting subscription and sync monitoring")

		for {
			select {
			case <-ticker.C:
				fileDB = performSyncCycle(ctx, manager, fileDB, cfg)

			case <-ctx.Done():
				cfg.Logger.Debug("Stopped subscription and sync monitoring")
				return
			}
		}
	}()
}

func performSyncCycle(ctx context.Context, mgr *manager.DatabaseManager, currentFileDB *sql.DB, cfg *config.BackendConfig) *sql.DB {
	if _, err := os.Stat(cfg.Paths.Database); os.IsNotExist(err) {
		cfg.Logger.Warn("File database does not exist, recreating", "path", cfg.Paths.Database)

		newDB, err := OpenAndInitDB(cfg.Paths.Database, "file", cfg)
		if err != nil {
			cfg.Logger.Error("Failed to recreate file database", "error", err)
			return currentFileDB
		}

		if currentFileDB != nil {
			currentFileDB.Close()
		}
		currentFileDB = newDB
	}

	syncCtx, syncCancel := context.WithTimeout(ctx, 15*time.Second)
	defer syncCancel()

	if err := mgr.SyncDBWithContext(syncCtx, currentFileDB, "memory to file"); err != nil {
		cfg.Logger.Error("Failed to synchronize database", "error", err)
	} else {
		cfg.Logger.Info("Database synchronized successfully")
	}

	return currentFileDB
}
