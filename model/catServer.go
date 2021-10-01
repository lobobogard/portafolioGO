package model

import "gorm.io/gorm"

type CatServer struct {
	gorm.Model
	Server string `gorm:"type:varchar(20);not null;unique" json:"server"`
}

// insert into cat_servers(server) values('AWS');
// insert into cat_servers(server) values('Google cloud');
// insert into cat_servers(server) values('Azure microsoft');
// insert into cat_servers(server) values('Vultr');
// insert into cat_servers(server) values('Digital ocean');
// insert into cat_servers(server) values('Linode');
// insert into cat_servers(server) values('Heroku');
