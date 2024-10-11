package utils

// MessageTemplates holds the message format for each action type in both Arabic and English
var MessageTemplates = map[string]map[string]string{
	"like": {
		"ar": "أبدى %s إعجاباً بمشاركتك: %s",
		"en": "%s liked your post: %s",
	},
	"comment": {
		"ar": "علق %s على مشاركتك: %s",
		"en": "%s commented on your post: %s",
	},
	"mention": {
		"ar": "أشار إليك %s: %s",
		"en": "%s mentioned you in a comment: %s",
	},
}
