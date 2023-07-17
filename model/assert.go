package model

type AssertionMsg struct {
	Type      string `json:"type"`
	Code      int64  `json:"code" bson:"code"`
	IsSucceed bool   `json:"is_succeed" bson:"is_succeed"`
	Msg       string `json:"msg" bson:"msg"`
}
