package model

type Result struct {
	Recommend Recommend
	Delete    func() error
}
