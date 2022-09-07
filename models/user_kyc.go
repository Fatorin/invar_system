package models

type UserKYC struct {
	Model
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	IdentityCard string `json:"identity_card"`
	Photo        string `json:"photo"`
	IDPhotoFront string `json:"id_photo_front"`
	IDPhotoBack  string `json:"id_photo_back"`
	PhotoWithID  string `json:"photo_with_id"`
}
