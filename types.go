package main

// Config records
type Config struct {
	TelegramAPIKey string `json:"telegramapikey"`
	TMDBAPIKey     string `json:"tmdbapikey"`
	FlickrAPIKey   string `json:"flickrapikey"`
	DBHost         string `json:"dbhost"`
	DBName         string `json:"dbname"`
}

// Object is the reference for the DB object
type Object struct {
	ID      string `gorethink:"id,omitempty"`
	URL     string `gorethink:"url,omitempty"`
	Content []byte `gorethink:"content,omitempty"`
	FileID  string `gorethink:"fileid,omitempty"`
}

// Result is the one update from Telegram API
type Result struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// Updates is a collection of updates from Telegram API
type Updates struct {
	Ok     bool     `json:"ok"`
	Result []Result `json:"result"`
}

// ResultMessage is the feedback message from sending a message from Telegram API
type ResultMessage struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"chat"`
		Date           int `json:"date"`
		ReplyToMessage struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"reply_to_message"`
		Text string `json:"text"`
	} `json:"result"`
}

// ResultPhoto is the feedback message from sending a photo from Telegram API
type ResultPhoto struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"chat"`
		Date           int `json:"date"`
		ReplyToMessage struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"reply_to_message"`
		Photo []struct {
			FileID   string `json:"file_id"`
			FileSize int    `json:"file_size"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
		} `json:"photo"`
	} `json:"result"`
}

// ResultDocument is the feedback message from sending a document from Telegram API
type ResultDocument struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
		} `json:"chat"`
		Date           int `json:"date"`
		ReplyToMessage struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"reply_to_message"`
		Document struct {
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				FileSize int    `json:"file_size"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
			} `json:"thumb"`
			FileID   string `json:"file_id"`
			FileSize int    `json:"file_size"`
		} `json:"document"`
	} `json:"result"`
}

// ResultAction is the feedback message from sending a chat action from Telegram API
type ResultAction struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}
