package schema

type Response struct {
	Status  string         `json:"status,omitempty"`
	Payload DrosophilaTray `json:"drosphila"`
}

type DrosophilaTray struct {
	Vials []struct {
		GenomeLabel  string `json:"GenomeLabel"`
		ID           string `json:"id"`
		Color        string `json:"Color"`
		Description  string `json:"Description"`
		LastTaskDate string `json:"LastTaskDate"`
		NextTaskData string `json:"NextTaskData"`
		Temperature  int    `json:"Temperature"`
		Task         string `json:"Task"`
		State        string `json:"State"`
		Status       string `json:"Status"`
	} `json:"Vials"`
}
