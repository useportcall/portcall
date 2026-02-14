package df

// IncrementInput is the input for incrementing a dogfood feature.
type IncrementInput struct {
	UserID  string `json:"user_id"`
	Feature string `json:"feature"`
	IsTest  bool   `json:"is_test"`
}

// DecrementInput is the input for decrementing a dogfood feature.
type DecrementInput struct {
	UserID  string `json:"user_id"`
	Feature string `json:"feature"`
	IsTest  bool   `json:"is_test"`
}
