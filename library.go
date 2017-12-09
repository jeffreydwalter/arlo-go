package arloclient

// LibraryMetaData is the library meta data.
type LibraryMetaData struct {
	DateTo   string                         `json:"dateTo"`
	DateFrom string                         `json:"dateFrom"`
	Meta     map[string]map[string]Favorite `json:"meta"`
}

// presignedContentUrl is a link to the actual video in Amazon AWS.
// presignedThumbnailUrl is a link to the thumbnail .jpg of the actual video in Amazon AWS.
type Recording struct {
	MediaDurationSecond   int    `json:"mediaDurationSecond"`
	ContentType           string `json:"contentType"`
	Name                  string `json:"name"`
	PresignedContentUrl   string `json:"presignedContentUrl"`
	LastModified          int64  `json:"lastModified"`
	LocalCreatedDate      int64  `json:"localCreatedDate"`
	PresignedThumbnailUrl string `json:"presignedThumbnailUrl"`
	Reason                string `json:"reason"`
	DeviceId              string `json:"deviceId"`
	CreatedBy             string `json:"createdBy"`
	CreatedDate           string `json:"createdDate"`
	TimeZone              string `json:"timeZone"`
	OwnerId               string `json:"ownerId"`
	UtcCreatedDate        int64  `json:"utcCreatedDate"`
	CurrentState          string `json:"currentState"`
	MediaDuration         string `json:"mediaDuration"`
}

type Library []Recording
