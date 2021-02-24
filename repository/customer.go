package repository

import "AREX-Market-Ledger/models"

func (db Database) GetAllInvestor() ([]models.Investor, error) {
	var investors []models.Investor

	rows, err := db.Conn.Query("SELECT c.id, c.name FROM customer c inner join investor i on c.id = i.customer_id ORDER BY ID DESC")
	if err != nil {
		return investors, err
	}

	for rows.Next() {
		var investor models.Investor
		err := rows.Scan(&investor.Id, &investor.Name)
		if err != nil {
			return investors, err
		}
		investors = append(investors, investor)
	}
	return investors, nil
}

