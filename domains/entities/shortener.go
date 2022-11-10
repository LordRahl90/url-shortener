package entities

import "shortener/domains/shortener/storage"

// Shortener service entity as DTO from server to repo
type Shortener struct {
	LongText  string `json:"long"`
	ShortText string `json:"short"`
}

// ToDBEntity converts service entities to Db entities
func (s *Shortener) ToDBEntity() *storage.Record {
	return &storage.Record{
		LongText:  s.LongText,
		ShortText: s.ShortText,
	}
}

// ShortenerFromDBEntity converts from repo entities to service entities
func ShortenerFromDBEntity(r *storage.Record) *Shortener {
	return &Shortener{
		LongText:  r.LongText,
		ShortText: r.ShortText,
	}
}
