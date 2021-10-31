package clickhouse

import (
	"github.com/jmoiron/sqlx"
)

func NewClickhouseConf() (*sqlx.DB, error) {
	clickhouseConn, err := sqlx.Open("clickhouse", "tcp://127.0.0.1:9000")
	if err != nil {
		return nil, err
	}

	if _, err = clickhouseConn.Exec(`
    CREATE TABLE IF NOT EXISTS users_log (
      info         String,
      action_time  DateTime
    ) engine=TinyLog
  `); err != nil {
		return nil, err
	}

	return clickhouseConn, nil
}
