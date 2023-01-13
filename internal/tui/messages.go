package tui

type (
	createAccountMsg struct {
		nick, email, password string
	}
	updateFocusedInputsMsg struct {
		index int
	}
	switchModeMsg struct {
		mode mode
	}
	resetFormMsg        struct{}
	authenticateUserMsg struct{}
)
