package model

type ChatType string

const (
	ChatPrivate        ChatType = "private"
	ChatGroup          ChatType = "group"
	ChatSuperGroup     ChatType = "supergroup"
	ChatChannel        ChatType = "channel"
	ChatChannelPrivate ChatType = "privatechannel"
)

type ChatDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Bio         string `json:"bio,omitempty"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"has_private_forwards,omitempty"`
	Protected   bool   `json:"has_protected_content,omitempty"`
}
