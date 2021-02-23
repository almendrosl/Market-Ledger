package models

type Invoice struct {
	Id          int
	Number      string
	Description string
	FaceValue   float32
	Issuer      Issuer
}
