package models

type WhiteList struct {
	Model
	UserID    uint   `json:"user_id"`
	ChainName string `json:"coin_type"`
	ChainID   uint   `json:"chain_id"`
	Address   string `json:"address"`
	NickName  string `json:"nick_name"`
	Comment   string `json:"comment"`
}
