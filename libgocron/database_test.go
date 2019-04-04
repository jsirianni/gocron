package libgocron

import (
    "testing"
)

func TestqueryDatabase(t *testing.T) {
    g := getTestConfig()
    q := "SELECR * FROM gocron;"

    _, err := queryDatabase(g, q)
    if err == nil {
        t.Errorf("Expected queryDatabase() to return an error, due to no database to connect to")
    }
}

func TestupdateDatabase(t *testing.T) {
    g := getTestConfig()
    c := getTestCron()

    if g.updateDatabase(c) == true {
        t.Errorf("Expected updateDatabase() to return false, due to no database to connect to")
    }
}

func TestcreateGocronTable(t *testing.T) {
    g := getTestConfig()

    if err := g.createGocronTable(); err == nil {
        t.Errorf("Expected createGocronTable() to return an error, due to no database to connect to")
    }
}

func TesttestDatabaseConnection(t *testing.T) {
    g := getTestConfig()

    if err := g.testDatabaseConnection(); err == nil {
        t.Errorf("Expected testDatabaseConnection() to return an error, due to no database to connect to")
    }
}

func TestgetDatabaseVersion(t *testing.T) {
    g := getTestConfig()

    if _, err := g.getDatabaseVersion(); err == nil {
        t.Errorf("Expected getDatabaseVersion() to return an error, due to no database to connect to")
    }
}

func TestgetCronsByAccount(t *testing.T) {
    g := getTestConfig()

    if _, err := g.getCronsByAccount("test"); err == nil {
        t.Errorf("Expected getCronsByAccount() to return an error, due to no database to connect to")
    }
}
