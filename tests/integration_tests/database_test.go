package integrationTests

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/innovay-software/famapp-main/app/models"
	"github.com/innovay-software/famapp-main/app/repositories"
	"github.com/innovay-software/famapp-main/app/services"
	"github.com/innovay-software/famapp-main/app/utils"
	"gorm.io/gorm/schema"
)

// Resets database to initial state
func resetDatabase() {
	if err := ensureTablesExist(); err != nil {
		panic(err)
	}
	if err := truncateAllTables(); err != nil {
		panic(err)
	}
	utils.LogSuccess("Reset database succeeded")
}

// Make sure all tables are present
func ensureTablesExist() error {
	db := services.GetMainDBConnection()
	tables := getTableNames()
	for _, tableName := range tables {
		q := fmt.Sprintf("SELECT * FROM %s;", tableName)
		if err := db.Exec(q).Error; err != nil {
			// if was able to create table, then the table exists
			utils.LogError("Table missing:", q)
			return err
		}
	}
	// Insert admin
	return nil
}

// Purge all tables
func truncateAllTables() error {
	db := services.GetMainDBConnection()
	tables := getTableNames()
	for _, tableName := range tables {
		q1 := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
		if err := db.Exec(q1).Error; err != nil {
			return err
		}
		q2 := fmt.Sprintf("ALTER SEQUENCE %s_id_seq restart with 1", tableName)
		if err := db.Exec(q2).Error; err != nil {
			return err
		}
	}

	// Create the first super admin user account
	user1 := models.User{
		UUID:           uuid.New(),
		Name:           superAdminName,
		Mobile:         superAdminMobile,
		LockerPasscode: superAdminPassword,
		Role:           "admin",
		Status:			true,
	}
	user1.SetPassword(superAdminPassword)
	repositories.SaveDbModel(&user1)
	log.Println("Purge table succeeded")
	return nil
}

func getTableNames() []string {
	// Manually typing all models.
	// Todo: use AST or reflections to dynamically get a list of model structs that implements the Tabler interface

	modelList := []schema.Tabler{
		models.AppVersion{},
		models.Config{},
		models.FolderFileUpload{},
		models.FolderFile{},
		models.FolderInvitee{},
		models.Folder{},
		models.LockerNote{},
		models.LockerNoteInvitee{},
		models.LockerNoteVersion{},
		models.Upload{},
		models.Traffic{},
		models.User{},
	}

	tableNames := []string{}
	for _, model := range modelList {
		tableNames = append(tableNames, model.TableName())
	}

	return tableNames
}

// Insert a config record to database
func insertModel(model schema.Tabler) error {
	err := repositories.SaveDbModel(model)
	if err != nil {
		utils.LogError("Error inserting config: ", err)
	}
	return err
}
