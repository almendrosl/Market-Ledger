package models

type Invoice struct {
	Id          int32
	Number      string
	Description string
	FaceValue   float32
	Issuer      Issuer
}
