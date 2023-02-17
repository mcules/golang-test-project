package models

import (
	"bufio"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"reflect"
	"time"
	"unsafe"
)

// Model database model
type Model struct {
	ID        int `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// DB model
type DB struct {
	Con *gorm.DB
}

// Init initializes database connection
func (db DB) Init() {
	var err error

	fmt.Println("Database: Init")
	db.Con, err = gorm.Open(sqlite.Open(ramData.Config.Sqlite.Filename), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
}

/*func (db DB) Init() {
	var err error

	fmt.Println("Database: Init")

	dsn := ramData.Config.Mysql.Username + ":" + ramData.Config.Mysql.Password + "@tcp(" + ramData.Config.Mysql.Hostname + ":" + strconv.Itoa(ramData.Config.Mysql.Port) + ")/" + ramData.Config.Mysql.DBName
	ramData.DB.Con, err = gorm.Open(mysql.Open(dsn+"?charset=utf8mb4&maxAllowedPacket=0&parseTime=True"), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}*/

// DBImportVersion model
type DBImportVersion struct {
	gorm.Model
	FileType      string
	FileName      string
	FileTimestamp time.Time
}

// StartUpdates gets all sql update files and import them
func (db DB) StartUpdates() {
	path := "updates/mysql/"

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, file := range files {
		if !db.checkIfImported(file.Name(), "update_mysql") {
			fmt.Println("Importing: ", file.Name())

			db.importSQLFile(path + file.Name())

			info, _ := file.Info()
			db.Con.Create(&DBImportVersion{FileType: "update_mysql", FileName: file.Name(), FileTimestamp: info.ModTime()})
		}
	}
}

func (db DB) importSQLFile(filename string) {
	queries, err := db.getQueriesFromSQLFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, query := range queries {
		if err := db.Con.Exec(query).Error; err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (db DB) getQueriesFromSQLFile(filename string) ([]string, error) {
	var result []string

	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return result, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(f)

	var line string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line += " " + sc.Text()

		if len(sc.Text()) > 0 && sc.Text()[len(sc.Text())-1:] == ";" {
			result = append(result, MyString(line).CombineWhitespaces())
			line = ""
		}
	}
	if err := sc.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func (db DB) checkIfImported(file string, fileType string) bool {
	var importVersions []DBImportVersion
	db.Con.Where("file_type=? AND file_name=?", fileType, file).Find(&importVersions)

	for _, fileImport := range importVersions {
		if fileImport.FileName == file {
			return true
		}
	}

	return false
}

func (db DB) Migrate() {
	fmt.Println("DB Migrate")

	sections, offsets := typelinks()

	for i, base := range sections {
		for _, offset := range offsets[i] {
			typeAddr := add(base, uintptr(offset), "")
			typ := reflect.TypeOf(*(*interface{})(unsafe.Pointer(&typeAddr)))
			fmt.Println(typ)

			/*g, ok := typ.(AutoMigrate)
			if ok {
				fmt.Println(g)
			}*/
		}
	}
}

func typelinks() (sections []unsafe.Pointer, offset [][]int32) {
	return
}

func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return nil
}
