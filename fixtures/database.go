package fixtures

import (
	"conformity-core/config"
	database "conformity-core/frameworks/database/gorm"
	"conformity-core/frameworks/database/gorm/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testConfig = config.DatabaseConfig{
	Host:         testEnv("TEST_DATABASE_HOST", "localhost"),
	User:         testEnv("TEST_DATABASE_USER", "postgres"),
	Password:     testEnv("TEST_DATABASE_PASSWORD", "postgres"),
	DatabaseName: testEnv("TEST_DATABASE_NAME", "conformityprotest"),
	Port:         testEnv("TEST_DATABASE_PORT", "5432"),
}

func testEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

var DB_Test = GetConnectionTestDatabase()

func GetConnectionTestDatabase() *database.Database {
	time.Local = time.UTC
	CreateTestDatabase()
	db, err := gorm.Open(postgres.Open(database.ParseDSN(&testConfig)))

	if err != nil {
		panic(err)
	}

	if err := database.Migrate(db); err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)

	return &database.Database{Conn: db}
}

var tables = []string{
	models.User{}.TableName(),
	models.Company{}.TableName(),
	models.Department{}.TableName(),
	models.UserCompanyDepartment{}.TableName(),
	models.Challenge{}.TableName(),
	models.Campaign{}.TableName(),
	models.CampaignResult{}.TableName(),
	models.ScheduledChallenge{}.TableName(),
}

func CleanTestDatabase() {
	for _, table := range tables {
		if err := DB_Test.Conn.Exec("TRUNCATE " + table + " CASCADE").Error; err != nil {
			panic(err)
		}
	}
}

func CreateTestDatabase() {
	withoutDb := config.DatabaseConfig{
		DatabaseName: "",
		Host:         testConfig.Host,
		User:         testConfig.User,
		Password:     testConfig.Password,
		Port:         testConfig.Port,
	}

	dsnWithoutDb := database.ParseDSN(&withoutDb)

	db, err := gorm.Open(postgres.Open(dsnWithoutDb))
	if err != nil {
		panic(err)
	}

	result := 0
	if err := db.Raw("SELECT 1 from pg_database WHERE datname=?", testConfig.DatabaseName).Scan(&result).Error; err != nil {
		panic(err)
	}

	hasDatabase := result > 0
	if hasDatabase {
		return
	}

	if err := db.Exec("CREATE DATABASE " + testConfig.DatabaseName).Error; err != nil {
		panic(err)
	}
}
