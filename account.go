package arloclient

// Account is the account data.
type Account struct {
	UserId        string  `json:"userId"`
	Email         string  `json:"email"`
	Token         string  `json:"token"`
	PaymentId     string  `json:"paymentId"`
	Authenticated uint32  `json:"authenticated"`
	AccountStatus string  `json:"accountStatus"`
	SerialNumber  string  `json:"serialNumber"`
	CountryCode   string  `json:"countryCode"`
	TocUpdate     bool    `json:"tocUpdate"`
	PolicyUpdate  bool    `json:"policyUpdate"`
	ValidEmail    bool    `json:"validEmail"`
	Arlo          bool    `json:"arlo"`
	DateCreated   float64 `json:"dateCreated"`
}
