package main

import (
	"database/sql"
	"fmt"
	"time"
)



type DBManager struct {
	db *sql.DB
}

func NewDBManager(db *sql.DB) DBManager{
	return DBManager{db}
}

type Cars struct {
	id int
	car_name string
	color string
	price float64
	image_url string
	created_at time.Time
	Images []*Cars_image
}

type Cars_image struct {
	id int 
	image_url string
	sequence_number int 
}

type GetProductsParams struct {
	Limit int32
	Page int32
	Search string 
}

type GetProductsResponse struct {
	Car []*Cars
	count int32 
}

func (c *DBManager) CreateCar(car *Cars) (int64, error) {
	var carID int64
	query := `
		INSERT INTO cars_db (	
			car_name,
			color,
			price,
			image_url
		) VALUES ($1,$2,$3,$4)
		RETURNING id
		`
	row := c.db.QueryRow(
		query,
		car.car_name,
		car.color,
		car.price,
		car.image_url,
	)	

	err := row.Scan(&car.id)
	if err != nil {
		return 0, err 
	}

	queryInsertImage := `
		INSERT INTO cars_image (
			car_id,
			image_url,
			sequence_number
		) VALUES ($1,$2,$3)`

	for _, image := range car.Images {
		_, err := c.db.Exec(
			queryInsertImage,
			carID,
			image.image_url,
			image.sequence_number,
		)
		if err != nil {
			return 0, nil 
		}
	}	
	return carID, nil 
}

func (c *DBManager) GetCar(id int64) (*Cars, error) {
	var car Cars

	car.Images = make([]*Cars_image, 0)

	query := `
		SELECT
		c.id,
		c.car_name,
		c.color,
		c.price,
		c.created_at,
		c.image_url
	FROM cars_db c`


	row := c.db.QueryRow(query,id)

	err := row.Scan(
		&car.id,
		&car.car_name,
		&car.color,
		&car.price,
		&car.created_at,
		&car.image_url,
	)

	if err != nil {
		return nil, err 
	}

	queryImages := `
		SELECT 
		id
		image_url,
		sequence_number
		FROM cars_image
		WHERE car_id=$1
	`

	rows, err := c.db.Query(queryImages,id)
	if err != nil {
		return nil, err 
	}
	defer rows.Close()

	for rows.Next() {
		var image Cars_image

		err := rows.Scan(
			&image.id,
			&image.image_url,
			&image.sequence_number,
		)
		if err != nil {
			return nil, err 
		}
		car.Images = append(car.Images, &image)
	}

	return &car, nil 

}


func (c *DBManager) GetAllCars(params *GetProductsParams) (*GetProductsResponse, error) {
	var result GetProductsResponse

	result.Car = make([]*Cars, 0)

	filter := ""
	if params.Search != "" {
		filter = fmt.Sprintf("WHERE car_name ilike '%s'",
	"%"+ params.Search +"%")
	}

	query := `
		SELECT
			c.id,
			c.car_name,
			c.color,
			c.price,
			c.image_url,
			c.created_at
		FROM cars_db c 
		` + filter + `
		ORDER BY created_at DESC
		LIMIT $1 OFSSET $2`

	offset := (params.Page - 1) * params.Limit
	rows,err := c.db.Query(query,params.Limit,offset)
	if err != nil {
		return nil, err 
	}	
	defer rows.Close()

	for rows.Next() {
		var car Cars

		err := rows.Scan(
			&car.id,
			&car.car_name,
			&car.color,
			&car.price,
			&car.created_at,
			&car.image_url,
		)
		if err != nil {
			return nil, err 
		}

		result.Car = append(result.Car, &car)
	}
	return &result, nil 
}

func (c *DBManager) UpdateCar(car *Cars) error {
	query := `
		UPDATE cars_db SET
			car_name=$1
			price=$2,
			image_url=$3
		WHERE id=$5`

	result, err := c.db.Exec(
		query,
		car.car_name,
		car.price,
		car.image_url,
		car.id,
	)	
	if err != nil {
		return err 
	}
	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err 
	}
	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	queryDeleteImages := `DELETE FROM cars_image WHERE 
		car_id=$1`
		_,err = c.db.Exec(queryDeleteImages,car.id)
		if err != nil {
			return err 
		}
	queryInsertImage := `
		INSERT INTO cars_image (
			car_id,
			image_url,
			sequence_number
		) VALUES ($1,$2,$3)`
	
	for _, image := range car.Images {
		_,err := c.db.Exec(
			queryInsertImage,
			car.id,
			image.image_url,
			image.sequence_number,
		)
		if err != nil {
			return err 
		}
	}

	return nil 
}

func (c *DBManager) DeleteCar(id int) error {
	queryDeleteImages := `DELETE FROM cars_image WHERE car_id=$1`
	_, err := c.db.Exec(queryDeleteImages,id)
	if err != nil {
		return err 
	}
	queryDelete := `DELETE FROM cars_db WHERE id=$1`
	result,err := c.db.Exec(queryDelete,id)
	if err != nil {
		return err 
	}
	rowsCount,err := result.RowsAffected()
	if err != nil {
		return err 
	}
	if rowsCount == 0 {
		return sql.ErrNoRows
	}
	return nil 
}

