package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	PostgresUser     = "postgres"
	PostgresPassword = "12345"
	PostgresHost     = "localhost"
	PostgresPort     = 5432
	PostgresDatabase = "cars"
)

func main() {
	connStr := fmt.Sprintf(
		"host = %s port = %d user = %s password = %s dbname = %s sslmode=disable",
		PostgresHost,
		PostgresPort,
		PostgresUser,
		PostgresPassword,
		PostgresDatabase,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
	}

	m := NewDBManager(db)

	id, err := m.CreateCar(&Cars{
		car_name:  "Malibu",
		color:     "black",
		price:     320000.00,
		image_url: "test_url",
		Images: []*Cars_image{
			{
				image_url:       "test_url_1",
				sequence_number: 1,
			},
			{
				image_url:       "test_url_2",
				sequence_number: 2,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to create car: %v", err)
	}

	car, err := m.GetCar(id)
	if err != nil {
		log.Fatalf("failed to get car: %v", err)
	}
	fmt.Println(car)

	resc, err := m.GetAllCars(&GetProductsParams{
		Limit: 10,
		Page:  1,
	})
	if err != nil {
		log.Fatalf("failed to get cars: %v", err)
	}
	fmt.Println(resc)

	err = m.UpdateCar(&Cars{
		id:        2,
		car_name:  "gentra",
		color:     "grey",
		price:     15000.00,
		image_url: "test_url_2",
		Images: []*Cars_image{
			{
				image_url:       "test_url_10",
				sequence_number: 10,
			},
			{
				image_url:       "test_url_3",
				sequence_number: 3,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to update car: %v", err)
	}

	err = m.DeleteCar(car.id)
	if err != nil {
		log.Fatalf("failed to delete car: %v", err)
	}
}
