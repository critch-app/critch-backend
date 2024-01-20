package api

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createServerRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	OwnerID     string `json:"owner_id" binding:"required"`
	Photo       string `json:"photo"`
}
