package folder

type CreateFolderParams struct {
	Title string `json:"title" validate:"required" example:"日本史"`
}