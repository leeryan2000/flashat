package wire

type AvatarUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectURL string `json:"object_url"`
}

type SetAvatarURLInput struct {
	AvatarURL string `json:"avatar_url"`
}
