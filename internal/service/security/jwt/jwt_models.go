package jwt

const (
	USER_ID   = "userId"
	USER_ROLE = "userRole"
)

type Claims struct {
	UserId   string `json:"userId"`
	UserRole string `json:"userRole"`
}

func NewClaims(userId, userRole string) *Claims {
	return &Claims{
		UserId:   userId,
		UserRole: userRole,
	}
}
