package model

import "gorm.io/gorm"

type CatSystemOperative struct {
	gorm.Model
	SystemOperative string `gorm:"type:varchar(20);not null;unique" json:"systemOperative"`
}

// insert into cat_system_operatives(system_operative) values('Linux');
// insert into cat_system_operatives(system_operative) values('Windows');
// insert into cat_system_operatives(system_operative) values('Mac OS');
