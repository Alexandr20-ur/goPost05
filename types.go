package goPost05

type User struct {
	ID       int
	Username string
}

type Userdata struct {
	ID          int
	Name        string
	Surname     string
	Description string
	Username    string
}

var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)
