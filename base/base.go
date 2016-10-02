package base

// UserID is id of user
type UserID uint64

// BallID is consist of userID and id. id is a value from 1 - 2^16.
// After user creating a ball, id add to one. 0 is user s airplane.
type BallID struct {
	UserID UserID
	ID     uint16
}

// ImageID is id of image
type ImageID uint8
