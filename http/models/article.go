package model

type Article struct {
	Id            int    `gorm:"primary_key,AUTO_INCREMENT" json:"id"`
	ApplicationId int    `gorm:"unique_index:idx_application_id_target_id;not null" json:"application_id"`
	TargetId      string `gorm:"type:varchar(20);unique_index:idx_application_id_target_id;not null" json:"target_id" `
	Title         string `gorm:"type:varchar(100);not null" json:"title" `
	Cover         string `gorm:"type:varchar(255);default:null" json:"cover" `
	Json          string `gorm:"type:text;not null" json:"json"`
	Hit           int    `gorm:"default:0" json:"hit"`
	BaseModel
}
