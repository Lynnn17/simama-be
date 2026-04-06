package infras

import (
	"fmt"

	_ "github.com/lib/pq"

	"lms-be/configs"
	"lms-be/shared/failure"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxIdleConnection = 10
	maxOpenConnection = 10
)

// Block contains a transaction block
type Block func(db *sqlx.Tx, c chan error)

// PostgresqlConn wraps a pair of read/write Postgresql connections
type PostgresqlConn struct {
	Read      *sqlx.DB
	Write     *sqlx.DB
	GormRead  *gorm.DB
	GormWrite *gorm.DB
	Conn      *sqlx.DB
}

func ProvidePostgreSQLConn(config *configs.Config) *PostgresqlConn {
	return &PostgresqlConn{
		Read:      CreatePostgreSQLReadConn(*config),
		Write:     CreatePostgreSQLWriteConn(*config),
		GormRead:  CreateGormReadConn(*config),
		GormWrite: CreateGormWriteConn(*config),
	}
}

func CreatePostgreSQLWriteConn(config configs.Config) *sqlx.DB {
	return CreateDBConnection("write",
		config.DB.PostgreSQL.Write.Host,
		config.DB.PostgreSQL.Write.Port,
		config.DB.PostgreSQL.Write.Name,
		config.DB.PostgreSQL.Write.Username,
		config.DB.PostgreSQL.Write.Password,
		config.DB.PostgreSQL.Write.Timezone)
}

func CreatePostgreSQLReadConn(config configs.Config) *sqlx.DB {
	return CreateDBConnection("read",
		config.DB.PostgreSQL.Read.Host,
		config.DB.PostgreSQL.Read.Port,
		config.DB.PostgreSQL.Read.Name,
		config.DB.PostgreSQL.Read.Username,
		config.DB.PostgreSQL.Read.Password,
		config.DB.PostgreSQL.Read.Timezone)
}

func CreateGormWriteConn(config configs.Config) *gorm.DB {
	return CreateGormConn("write",
		config.DB.PostgreSQL.Write.Host,
		config.DB.PostgreSQL.Write.Port,
		config.DB.PostgreSQL.Write.Name,
		config.DB.PostgreSQL.Write.Username,
		config.DB.PostgreSQL.Write.Password,
		config.DB.PostgreSQL.Write.Timezone)
}

func CreateGormReadConn(config configs.Config) *gorm.DB {
	return CreateGormConn("read",
		config.DB.PostgreSQL.Read.Host,
		config.DB.PostgreSQL.Read.Port,
		config.DB.PostgreSQL.Read.Name,
		config.DB.PostgreSQL.Read.Username,
		config.DB.PostgreSQL.Read.Password,
		config.DB.PostgreSQL.Read.Timezone)
}

func CreateDBConnection(name, host, port, dbName, username, password, timeZone string) *sqlx.DB {
	descriptor := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		dbName)
	log.Info().Msg(descriptor)

	db, err := sqlx.Connect("postgres", descriptor)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Str("name", name).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Failed connecting to database")
	} else {
		log.
			Info().
			Str("name", name).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Connected to database")
	}
	db.SetMaxIdleConns(maxIdleConnection)
	db.SetMaxOpenConns(maxOpenConnection)

	return db
}

func CreateGormConn(name, host, port, dbName, username, password, timeZone string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.
			Fatal().
			Err(err).
			Str("Gorm->name", name).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Failed connecting to database")
	} else {
		log.
			Info().
			Str("Gorm->name", name).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Connected to database")
	}
	// Auto Migrate
	// db.AutoMigrate(&item.Item{})
	// db.Table("m_item").Migrator().CreateTable(&item.Item{})
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(maxIdleConnection)
	sqlDB.SetMaxOpenConns(maxOpenConnection)
	return db
}

// WithTransaction performs queries with transaction
func (m *PostgresqlConn) WithTransaction(block Block) (err error) {
	e := make(chan error)
	tx, err := m.Write.Beginx()
	if err != nil {
		return
	}
	go block(tx, e)
	err = <-e
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			err = failure.InternalError(errTx)
		}
		return
	}
	err = tx.Commit()
	return
}
