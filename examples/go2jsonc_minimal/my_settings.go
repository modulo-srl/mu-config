package main

type MySettings struct {
	Main  Main       // Main configuration parameters.
	Users []UserItem // Users list.
}

type Main struct {
	ParamString string
	ParamBool   bool
	ParamInt    int
	ParamFloat  float64
}

type UserItem struct {
	Name  string `json:"name"`  // User name.
	Email string `json:"email"` // User e-mail.
}

func MySettingsDefaults() *MySettings {
	return &MySettings{
		Main: Main{
			ParamString: "ParamValue \n \"test\"",
			ParamBool:   true,
			ParamInt:    12,
			ParamFloat:  1.234,
		},
		Users: []UserItem{
			{
				Name:  "John",
				Email: "john@email",
			},
			{
				Name: "Smith",
			},
		},
	}
}
