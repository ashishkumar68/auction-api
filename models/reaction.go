package models

const (
	ReactionTypeLike = iota
	ReactionTypeDislike
)

const (
	ReactionTypeLikeString    = "Like"
	ReactionTypeDislikeString = "Dislike"
)

type Reaction struct {
	IdentityAuditableModel

	ItemId uint  `gorm:"column:item_id;index" json:"-"`
	Item   *Item `gorm:"foreignKey:ItemId" json:"item"`
	Type   uint8 `gorm:"column:type;type:smallint" json:"type"`
}

func (Reaction) TableName() string {
	return "reactions"
}

func NewReactionFromValues(item *Item, reactionType uint8, reactionBy *User) *Reaction {
	newReaction := Reaction{
		ItemId: item.ID,
		Type:   reactionType,
	}
	newReaction.UserCreatedBy = reactionBy.ID
	newReaction.UserUpdatedBy = reactionBy.ID

	return &newReaction
}

func (r Reaction) IsTypeSameAs(reaction Reaction) bool {
	if r.Type == reaction.Type {
		return true
	}

	return false
}

func FindReactionTypeString(reactionType uint8) string {
	var rType string
	switch reactionType {
	case ReactionTypeLike:
		rType = ReactionTypeLikeString
	case ReactionTypeDislike:
		rType = ReactionTypeDislikeString
	}

	return rType
}
