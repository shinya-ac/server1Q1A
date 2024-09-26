package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"

	config "github.com/shinya-ac/server1Q1A/configs"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
)

const maxRetries = 5
const delay = 5 * time.Second

var (
	once  sync.Once
	dbcon *sql.DB
)

func SetDB(d *sql.DB) {
	dbcon = d
}

func GetDB() *sql.DB {
	return dbcon
}

func NewMainDB(cnf config.ConfigList) {
	once.Do(func() {
		var err error
		logging.Logger.Info("DBHost:", "", cnf.DBHost)
		dbcon, err := connect(
			cnf.DBUser,
			cnf.DBPassword,
			cnf.DBHost,
			cnf.DBPort,
			cnf.DBName,
			cnf.CACertPath,
		)
		if err != nil {
			logging.Logger.Error("DBの初期化に失敗", "error", err)
			panic(err)
		}
		SetDB(dbcon)
	})
}

func connect(user string, password string, host string, port string, name string, caCertPath string) (*sql.DB, error) {
	// DBをTiDBで動かす場合はtlsModeをtidbにする
	tlsMode := config.Config.DBTLSMode
	// CA証明書の読み込み
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(caCertPath)
	if err != nil {
		logging.Logger.Error("CA証明書の読み込みに失敗", "error", err)
		return nil, fmt.Errorf("CA証明書の読み込みに失敗しました: %w", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		logging.Logger.Error("CA証明書の追加に失敗")
		return nil, fmt.Errorf("CA証明書の追加に失敗しました")
	}

	// TLS設定の作成
	err = mysql.RegisterTLSConfig("tidb", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		logging.Logger.Error("TLS設定の登録に失敗", "error", err)
		return nil, fmt.Errorf("TLS設定の登録に失敗しました: %w", err)
	}

	for i := 0; i < maxRetries; i++ {
		var connect string
		if tlsMode == "tidb" {
			connect = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=tidb", user, password, host, port, name)
		} else {
			connect = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)
		}

		db, err := sql.Open("mysql", connect)
		if err != nil {
			logging.Logger.Error("MySQLの接続に失敗", "error", err)
			return nil, fmt.Errorf("DBに接続できません。: %w", err)
		}

		err = db.Ping()
		if err == nil {
			logging.Logger.Info("DB接続が確立")
			return db, nil
		}

		logging.Logger.Warn("DBへの接続に失敗。再試行中...", "attempt", i+1, "error", err)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("DBへの接続に %d 回失敗しました", maxRetries)
}
