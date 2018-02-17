package ws

type Message struct {
	User
	Content string `json:"content"`
}

type User struct {
	Username string `json:"username"`
}

type PrivateMessage struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Message *Message `json:"message"`
}
