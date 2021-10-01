package model

import "gorm.io/gorm"

type CatFrontEnd struct {
	gorm.Model
	FrontEnd string `gorm:"type:varchar(20);not null;unique" json:"frontEnd" validate:"required"`
}

// insert into cat_front_ends (front_end) values('Vue');
// insert into cat_front_ends (front_end) values('React');
// insert into cat_front_ends (front_end) values('Angular');
// insert into cat_front_ends (front_end) values('Svelte');
// insert into cat_front_ends (front_end) values('Typescript');
// insert into cat_front_ends (front_end) values('jquery');
