package internal

import (
	"context" 
	"fmt"
	"os"
)

func ActionDB(fileName string) error {
	sqlBytes, err := os.ReadFile("sqlfiles/" + fileName)
	if err != nil {
		return fmt.Errorf("could not read %s file: %w", fileName, err)
	}


	_, err = DB.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		return fmt.Errorf("could not execute schema.sql script: %w", err)
	}

	return nil
}