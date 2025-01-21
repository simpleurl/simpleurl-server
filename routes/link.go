package routes

type Link struct {
	id   int
	url  string
	name string
}

func CreateLink(data any) (*Link, error) {
	return &Link{
		id:   1,
		url:  "",
		name: "Google",
	}, nil
}
