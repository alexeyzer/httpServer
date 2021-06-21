package store

import (
	"database/sql"
	"fmt"
	"github.com/alexeyzer/httpServer/internal/app/model"
)

type refrepository struct{
	store *Store
}

func (r * refrepository) Create(ref *model.Ref) (*model.Ref, error){
	if err := r.store.db.QueryRow("INSERT INTO ref(adv_id, ref) VALUES ($1,$2, $3) RETURNING ID",
		ref.AdvId,
		ref.Ref).Scan(&ref.ID); err != nil{
		return nil, err
	}
	return ref, nil
}

func (r * refrepository) GetList(AdvId int) ([]model.Ref, error){
	var res *sql.Rows = nil
	var err error = nil
	res, err = r.store.db.Query("SELECT id, adv_id, ref from ref where adv_id = $1", AdvId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer res.Close()
	list := []model.Ref{}
	for res.Next(){
		newref := model.Ref{}
		err = res.Scan(&newref.ID,&newref.AdvId ,&newref.Ref)
		if err != nil {
			return nil, err
		}
		fmt.Println(newref)
		list = append(list, newref)
	}
	return list, nil
}