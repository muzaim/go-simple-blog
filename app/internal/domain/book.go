package domain

type Book struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Title  string  `gorm:"size:255;not null" json:"title"`
	Isbn   string  `gorm:"size:50;not null" json:"isbn"`
	Author string  `gorm:"size:50;not null" json:"author"`
	Price  float64 `gorm:"not null" json:"price"`
	Stock  int     `gorm:"not null" json:"stock"`
}

type CreateNewBookRequest struct {
	Title  string  `json:"title" binding:"required"`
	Isbn   string  `json:"isbn" binding:"required"`
	Author string  `json:"author" binding:"required"`
	Price  float64 `json:"price" binding:"required,gt=0"`
	Stock  int     `json:"stock" binding:"required,gte=0"`
}

type UpdateBookRequest struct {
	Title  string  `json:"title" binding:"required"`
	Isbn   string  `json:"isbn" binding:"required"`
	Author string  `json:"author" binding:"required"`
	Price  float64 `json:"price" binding:"required,gt=0"`
	Stock  int     `json:"stock" binding:"required,gte=0"`
}
