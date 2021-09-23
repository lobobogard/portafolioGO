package model

import "gorm.io/gorm"

type CatCountry struct {
	gorm.Model
	Country string `gorm:"type:varchar(20);not null;unique" json:"country" validate:"required"`
}

// insert into cat_countries (country) values('méxico');
// insert into cat_countries (country) values('united state');
// insert into cat_countries (country) values('england');
// insert into cat_countries (country) values('brazil');
// insert into cat_countries (country) values('indian');
// insert into cat_countries (country) values('china');
// insert into cat_countries (country) values('germany');
// insert into cat_countries (country) values('australia');
