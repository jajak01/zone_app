package models

// User Table mapping
type User struct {
	UserID   uint   `gorm:"primaryKey;column:User_ID"`
	UserName string `gorm:"column:user_name"`
	AnonName string `gorm:"column:anon_name"`
}

func (User) TableName() string {
	return "User"
}

// UserAuth Table mapping
type UserAuth struct {
	UserID uint   `gorm:"primaryKey;column:User_ID"`
	Email  string `gorm:"column:Email;unique;not null"`
	Pass   string `gorm:"column:Pass;not null"`
	Phone  string `gorm:"column:Phone"`
}


func (UserAuth) TableName() string {
	return "User_Auth"
}

// UserInfo Table mapping
type UserInfo struct {
	UserID       uint    `gorm:"primaryKey;column:User_ID"`
	Email        string  `gorm:"column:Email"`
	Phone        string  `gorm:"column:Phone"`
	BirthDate    string  `gorm:"column:Birth_Date"` // Use string or time.Time
	BaseLocation string  `gorm:"column:Base_Location"`
	Height       float64 `gorm:"column:Height"`
	Weight       float64 `gorm:"column:Weight"`
}

func (UserInfo) TableName() string {
	return "User_Info"
}

// UserActivity Table mapping
type UserActivity struct {
	IDActivity   uint   `gorm:"primaryKey;column:ID_Activity"`
	UserID       uint   `gorm:"column:User_ID"`
	TypeActivity string `gorm:"column:Type_Activity"`
}

func (UserActivity) TableName() string {
	return "User_Activity"
}

// ActivityDetail Table mapping
type ActivityDetail struct {
	IDActivity    uint   `gorm:"primaryKey;column:ID_Activity"`
	UserID        uint   `gorm:"column:User_ID"`
	LogGisSpatial string `gorm:"column:LOG_GIS_SPATIAL"` 
}

func (ActivityDetail) TableName() string {
	return "Activity_Detail"
}

// SummaryActivity Table mapping
type SummaryActivity struct {
	IDActivity       uint    `gorm:"primaryKey;column:ID_Activity"`
	DistanceActivity float64 `gorm:"column:Distance_Activity"`
	MovingTime       string  `gorm:"column:Moving_Time"`
	AvgSpeed         float64 `gorm:"column:Avg_Speed"`
	SimpleGis        string  `gorm:"column:Simple_GIS"`
	Calories         int     `gorm:"column:Calories"`
}

func (SummaryActivity) TableName() string {
	return "Summary_Activity"
}