package storage

type Model struct {
	ID uint64 `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
}
