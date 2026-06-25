package database

type Contact struct {
	FullName       string `json:"full_name"`
	Passport       string `json:"passport"`
	Address        string `json:"address"`
	PhoneNumber    string `json:"phone_number"`
	Pin            string `json:"pin"`
	SipUsername    string `json:"sip_username"`
	SipSecret      string `json:"sip_secret"`
	Email          string `json:"email"`
	ConnectionType string `json:"connection_type"`
	RegTimeout     int    `json:"reg_timeout"`
	AsteriskNotes  string `json:"asterisk_notes"`
	CreatedAt      string `json:"created_at,omitempty"`
}

type Extension struct {
	ID         string `json:"id"`
	Transport  string `json:"transport"`
	Aors       string `json:"aors"`
	Auth       string `json:"auth"`
	Context    string `json:"context"`
	Disallow   string `json:"disallow"`
	Allow      string `json:"allow"`
	FullName   string `json:"full_name,omitempty"`
	Phone      string `json:"phone_number,omitempty"`
	MaxContact int    `json:"max_contacts,omitempty"`
}

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

type ConfigUpdate struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name,omitempty"`
	User     string `json:"user"`
	Password string `json:"password,omitempty"`
}