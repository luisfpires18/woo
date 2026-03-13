package model

// CampTemplate defines an admin-configurable camp composition.
type CampTemplate struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Tier          int    `json:"tier"`
	SpriteKey     string `json:"sprite_key"`
	Description   string `json:"description"`
	RewardTableID *int64 `json:"reward_table_id,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     *int64 `json:"updated_by,omitempty"`
}

// CampBeastSlot links a beast template to a camp template with spawn counts.
type CampBeastSlot struct {
	ID              int64 `json:"id"`
	CampTemplateID  int64 `json:"camp_template_id"`
	BeastTemplateID int64 `json:"beast_template_id"`
	MinCount        int   `json:"min_count"`
	MaxCount        int   `json:"max_count"`
}
