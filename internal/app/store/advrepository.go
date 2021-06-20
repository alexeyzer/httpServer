package store

import (
	"github.com/alexeyzer/httpServer/internal/app/model"
)

type AdvRepository struct{
	store *Store
}

func (a * AdvRepository) Create(adv *model.Adv) (*model.Adv, error){
	if err := a.store.db.QueryRow("INSERT INTO adv(name, description, price) VALUES ($1,$2, $3) RETURNING ID",
		adv.Name,
		adv.Description,
		adv.Price).Scan(&adv.ID); err != nil{
		return nil,err
	}
	return adv,nil
}

func (a *AdvRepository) FindById(ID int, optional bool) (*model.Adv, error){
	adv := &model.Adv{}
	if optional == false{
		if err := a.store.db.QueryRow("SELECT ID, name, price, date_create where ID = $1", ID).Scan(
			&adv.ID, &adv.Name, &adv.Price, &adv.Date); err != nil {
			return nil, err
		}
	} else {
		if err := a.store.db.QueryRow("SELECT ID, name, description, price, date_create where ID = $1", ID).Scan(
			&adv.ID, &adv.Name, &adv.Description, &adv.Price, &adv.Date); err != nil {
			return nil, err
		}
	}

	return adv,nil
}