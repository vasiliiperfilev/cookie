package data

import "database/sql"

type Repositories struct {
	Order OrderRepository
}

func NewRepositories(db *sql.DB, models Models) Repositories {
	return Repositories{
		Order: NewPsqlOrderRepository(db, models.Order, models.Message),
	}
}
