package netrc

type machine struct {
	name     string
	login    string
	password string
}

func newMachine(
	name string,
	login string,
	password string,
) *machine {
	return &machine{
		name:     name,
		login:    login,
		password: password,
	}
}

func (m *machine) Name() string {
	return m.name
}

func (m *machine) Login() string {
	return m.login
}

func (m *machine) Password() string {
	return m.password
}
