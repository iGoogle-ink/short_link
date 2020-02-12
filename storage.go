package main

type Storage interface {
	Shorten(url string, expire int64) (string, error)
	ShortlinkInfo(short string) (interface{}, error)
	Unshorten(short string) (string, error)
}
