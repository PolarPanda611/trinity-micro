package dbx

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PolarPanda611/trinity-micro/core/logx"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type ContextKey string

var (
	DBXContext   ContextKey = "dbx-context"
	SQLLiteType             = "sqllite"
	MysqlType               = "mysql"
	PostgresType            = "postgres"
	DB           *gorm.DB
)

type Config struct {
	Type        string
	DSN         string
	TablePrefix string
	MaxIdleConn int
	MaxOpenConn int
	Logger      logrus.FieldLogger
}

func Init(c *Config) {
	dialector, err := getDialector(c.Type, c.DSN)
	if err != nil {
		c.Logger.Fatalf("get dialector err, error:%v", err)
	}
	DB, err = gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   c.TablePrefix,
		},
		Logger: logger.New(c.Logger, logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		}),
	})
	if err != nil {
		c.Logger.Fatalf("init db error, err: %v", err)
	}
	c.Logger.Infof("init db successfully!")
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	c.Logger.Infof("db stats %v", sqlDB.Stats())
	if err != nil {
		c.Logger.Fatalf("init db error, err: %v", err)
	}
	c.Logger.Infof("init db successfully!")

	Migrate(context.Background(), "test", "shin", "new")
}

func SessionDB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionLogger := logx.FromCtx(r.Context())
		sessionDB := DB.Session(&gorm.Session{
			NewDB:   true,
			Context: r.Context(),
		})
		sessionDB.Logger = logger.New(sessionLogger, logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
		r = r.WithContext(context.WithValue(r.Context(), DBXContext, sessionDB))
		next.ServeHTTP(w, r)
	})
}

func FromCtx(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(DBXContext).(*gorm.DB)
	if !ok {
		panic("please use dbx.SessionDB to init db in context ")
	}
	return db
}

func getDialector(dbType string, DSN string) (gorm.Dialector, error) {
	switch strings.ToLower(dbType) {
	case MysqlType:
		return mysql.Open(DSN), nil
	case PostgresType:
		return postgres.Open(DSN), nil
	case SQLLiteType:
		return sqlite.Open(DSN), nil
	default:
		return nil, fmt.Errorf("unsupported db type, %v  ", dbType)
	}
}
