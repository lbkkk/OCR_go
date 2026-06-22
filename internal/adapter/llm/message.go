package llm

// Role identifies the author of a chat message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Image is an image attached to a message for vision-capable models.
type Image struct {
	// Data holds the raw image bytes.
	Data []byte
	// Format is the image encoding without the "image/" prefix, e.g. "png".
	Format string
}

// Message is a single chat message. When Images is non-empty the message is
// sent as a multi-part (vision) message; otherwise it is sent as plain text.
type Message struct {
	Role   Role
	Text   string
	Images []Image
}

// System returns a system message with the given text.
func System(text string) Message {
	return Message{Role: RoleSystem, Text: text}
}

// User returns a user message with the given text.
func User(text string) Message {
	return Message{Role: RoleUser, Text: text}
}

// UserWithImages returns a vision user message with text and images.
func UserWithImages(text string, images []Image) Message {
	return Message{Role: RoleUser, Text: text, Images: images}
}
