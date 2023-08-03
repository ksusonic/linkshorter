package controller

type Shorten struct {
	Link string `form:"link"`
}

type Redirect struct {
	ID string `uri:"id"`
}
