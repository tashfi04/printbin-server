package dtos

type GoogleUserData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type MinimalUserInfo struct {
	Username   string `json:"username"`
	Role       uint   `json:"role"`
	TeamName   string `json:"team_name"`
	RoomNumber string `json:"room_number"`
	AvailablePrintPageCount int `json:"available_print_page_count,omitempty"`
}
