package tui

type (
	createAccountMsg struct {
		nick, email, password string
	}
	updateFocusedInputsMsg struct {
		index int
	}
	resetFormMsg        struct{}
	authenticateUserMsg struct{}
)
