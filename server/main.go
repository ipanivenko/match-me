package main

import (
	"flag"
	"log"
	"matchme-server/handlers"
	"matchme-server/internal"
)

func main() {
	seed := flag.Bool("seed", false, "Seed the database with initial data")
	drop := flag.Bool("drop", false, "Drop all database tables")

	flag.Parse()

	internal.LoadConfig()
	err := internal.ConnectDB()

	if err != nil {
		log.Println(err)
	}

	defer internal.DB.Close()

	if *drop {
		log.Println("⚠️ -drop flag provided. Dropping all tables...")
		if err := internal.ActionDB("drop.sql"); err != nil {
			log.Println("Failed to drop tables:", err)
		}
		log.Println("✅ Tables dropped successfully.")
	}

	err = internal.ActionDB("schema.sql") //initialise db
	if err != nil {
		log.Println("Failed to initialize database:", err)
	}
	log.Println("✅ Database schema initialized successfully.")

	if *seed {
		log.Println("🌱 -seed flag provided. Seeding the database...")
		if err := internal.ActionDB("seed.sql"); err != nil { 
			log.Println("Failed to seed the database:", err)
		}
		log.Println("✅ Database seeded successfully.")
	}
	router := handlers.SetupRouter(internal.Cfg.IsDevMode, internal.DB)
	router.Run(":" + internal.Cfg.Port)
}
