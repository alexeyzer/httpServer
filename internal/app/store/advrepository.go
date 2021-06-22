package store

import (
	"database/sql"
	"github.com/alexeyzer/httpServer/internal/app/model"
)

type AdvRepository struct {
	store *Store
}

func (a *AdvRepository) Create(adv *model.Adv) (*model.Adv, error) {
	if err := a.store.db.QueryRow("INSERT INTO adv(name, description, price) VALUES ($1,$2, $3) RETURNING ID",
		adv.Name,
		adv.Description,
		adv.Price).Scan(&adv.ID); err != nil {
		return nil, err
	}
	return adv, nil
}

func (a *AdvRepository) FindById(ID int, optional bool) (*model.Adv, error) {

	adv := &model.Adv{}
	if optional == false {
		if err := a.store.db.QueryRow("SELECT id, name, price, date_create from adv where ID = $1", ID).Scan(
			&adv.ID, &adv.Name, &adv.Price, &adv.Date); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		}
	} else {
		if err := a.store.db.QueryRow("SELECT id, name, description, price, date_create from adv where id = $1", ID).Scan(
			&adv.ID, &adv.Name, &adv.Description, &adv.Price, &adv.Date); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		}
	}

	return adv, nil
}

func (a *AdvRepository) List(sort string) ([]model.Adv, error) {

	var res *sql.Rows = nil
	var err error = nil
	switch sort {
	case "-price":
		res, err = a.store.db.Query("SELECT id, name, description, price, date_create from adv order by price desc")
		if err != nil {
			return nil, err
		}
	case "price":
		res, err = a.store.db.Query("SELECT id, name, description, price, date_create from adv order by price asc")
		if err != nil {
			return nil, err
		}
	case "-date":
		res, err = a.store.db.Query("SELECT id, name, description, price, date_create from adv order by date_create desc")
		if err != nil {
			return nil, err
		}
	case "date":
		res, err = a.store.db.Query("SELECT id, name, description, price, date_create from adv order by date_create asc")
		if err != nil {
			return nil, err
		}
	default:
		res, err = a.store.db.Query("SELECT id, name, description, price, date_create from adv")
		if err != nil {
			return nil, err
		}
	}
	if res == nil {
		return nil, nil
	}
	defer res.Close()
	list := []model.Adv{}
	for res.Next() {
		newadv := model.Adv{}
		err = res.Scan(&newadv.ID, &newadv.Name, &newadv.Description, &newadv.Price, &newadv.Date)
		if err != nil {
			return nil, err
		}
		list = append(list, newadv)
	}
	return list, nil
}
