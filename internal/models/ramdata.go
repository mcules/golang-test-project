package models

var ramData *RAMData

// RAMData model
type RAMData struct {
	DB     DB
	Config struct {
		Mysql struct {
			Username string
			Password string
			Hostname string
			Port     int
			DBName   string
		}
		Sqlite struct {
			Filename string
		}
	}
}

// Init initializes ramdata object
func (rd RAMData) Init() {
	ramData = &rd

	ramData.Config.Sqlite.Filename = "gorm.db"
}
