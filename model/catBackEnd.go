package model

import "gorm.io/gorm"

type CatBackEnd struct {
	gorm.Model
	BackEnd string `gorm:"type:varchar(20);not null;unique" json:"backEnd" validate:"required"`
}

// insert into cat_back_ends (back_end) values('golang');
// insert into cat_back_ends (back_end) values('PHP');
// insert into cat_back_ends (back_end) values('C# net');
// insert into cat_back_ends (back_end) values('node.js');
// insert into cat_back_ends (back_end) values('rust');

/*todo*/
// insert into cat_system_operatives(system_operative) values('Linux');
// insert into cat_system_operatives(system_operative) values('Windows');
// insert into cat_system_operatives(system_operative) values('Mac OS');

// insert into cat_countries (country) values('méxico');
// insert into cat_countries (country) values('united state');
// insert into cat_countries (country) values('england');
// insert into cat_countries (country) values('brazil');
// insert into cat_countries (country) values('indian');
// insert into cat_countries (country) values('china');
// insert into cat_countries (country) values('germany');
// insert into cat_countries (country) values('australia');

// insert into cat_front_ends (front_end) values('Vue');
// insert into cat_front_ends (front_end) values('React');
// insert into cat_front_ends (front_end) values('Angular');
// insert into cat_front_ends (front_end) values('Svelte');
// insert into cat_front_ends (front_end) values('Typescript');
// insert into cat_front_ends (front_end) values('jquery');

// insert into cat_back_ends (back_end) values('golang');
// insert into cat_back_ends (back_end) values('PHP');
// insert into cat_back_ends (back_end) values('visual.net C#');
// insert into cat_back_ends (back_end) values('scala');
// insert into cat_back_ends (back_end) values('node.js');
// insert into cat_back_ends (back_end) values('rust');

// insert into cat_servers(server) values('AWS');
// insert into cat_servers(server) values('Google cloud');
// insert into cat_servers(server) values('Azure microsoft');
// insert into cat_servers(server) values('Vultr');
// insert into cat_servers(server) values('Digital ocean');
// insert into cat_servers(server) values('Linode');
// insert into cat_servers(server) values('Heroku');
