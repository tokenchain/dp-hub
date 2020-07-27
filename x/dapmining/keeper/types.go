package keeper

type (
	LoginRequest struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	}

	MeResponse struct {
		Address string `json:"address" yaml:"address"`
	}
)
