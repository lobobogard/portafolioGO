package model

type ReqPerfil struct {
	CompanyID         uint
	SystemOperativeID uint
	Server            []uint
	BackEnd           []uint
	FrontEnd          []uint
	Mysql             bool
	Mariadb           bool
	Postgresql        bool
	Mongodb           bool
	Redis             bool
	Sqlite            bool
}
