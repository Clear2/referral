package response

import "github.com/keep/sunny/ent"

// UserAPIResponse -.
type UserAPIResponse struct {
	APIResponse[ent.UserEntity]
}
