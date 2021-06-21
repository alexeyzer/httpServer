package model

type Adv struct{
	ID int
	Name string
	Description string
	Price int
	Date string
	Ref []Ref
	NextPage bool
}
