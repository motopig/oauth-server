package model

type Client struct {
	*ClientBase `xorm:"extends"`
	BaseModel   `xorm:"extends"`
}

type ClientBase struct {
	ClientId     string `xorm:"client_id unique notnull  default '' comment(客户端id)" json:"client_id"`
	ClientSecret string `xorm:"client_secret unique notnull  default '' comment(客户端秘钥)" json:"client_secret"`
	ClientName   string `xorm:"client_name unique notnull comment(客户端名称)" json:"client_name"`
	ClientDoamin string `xorm:"client_domain unique notnull comment(客户端主域名)" json:"client_domain"`
}

func (c *Client) TableName() string {
	return "clients"
}

func GetAllClients() []Client {
	var clients = make([]Client, 0)
	_ = DB().Table("clients").Find(&clients)
	return clients
}

func GetClientByClientId(id string) (bool, error) {
	var client Client
	has, err := DB().Table("clients").Where("client_id = ?", id).Get(&client)
	return has, err
}
