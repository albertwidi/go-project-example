package authentication

// list of authentication action
const (
	ActionRegister        string = "register"
	ActionLogin           string = "login"
	ActionVerifyPayment   string = "verify-payment"
	ActionVerifyChangePin string = "verify-changepin"
)

// list of authentication provider
const (
	ProviderOTP      Provider = "otp"
	ProviderPassword Provider = "password"
	ProviderPin      Provider = "pin"
)
