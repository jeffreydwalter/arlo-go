package arloclient

/*
// Credentials is the login credential data.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Duration holds two dates used when you need to specify a date range in the format "20060102".
type Duration struct {
	DateFrom string `json:"dateFrom""`
	DateTo   string `json:"dateTo"`
}

// PasswordPair is used when updating the account password.
type PasswordPair struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// FullName is used when updating the account username.
type FullName struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
*/

// Owner is the owner of a Device data.
type Owner struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	OwnerId   string `json:"ownerId"`
}

// Properties is the Device properties data.
type Properties struct {
	ModelId       string `json:"modelId"`
	OlsonTimeZone string `json:"olsonTimeZone"`
	HwVersion     string `json:"hwVersion"`
}

type Favorite struct {
	NonFavorite uint8 `json:"nonFavorite"`
	Favorite    uint8 `json:"Favorite"`
}

type Friend struct {
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Devices      DeviceOrder `json:"devices"`
	LastModified float64     `json:"lastModified"`
	AdminUser    bool        `json:"adminUser"`
	Email        string      `json:"email"`
	Id           string      `json:"id"`
}
