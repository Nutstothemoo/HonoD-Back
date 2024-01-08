package sdk

type  UserDetails struct {
	ID    string
	Name  string
	Email string
}

type SafeUser struct {
	UserID      string `json:"userID"`
	Email       string `json:"email"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Avatar			string `json:"avatar"`
	Phone       string `json:"phone"`	
	Username    string `json:"username"`
}

type GoogleUserDetails struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}