package models

import "github.com/jinzhu/gorm"

const (
	ErrUserIDRequired modelError = "models: user ID is required"

	ErrTitleRequired modelError = "models: title is required"
)

type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not_null;index"`
	Title  string  `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

type GalleryService interface {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(userID uint) ([]Gallery, error)
	All() ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
}

type galleryGorm struct {
	db *gorm.DB
}

var _ GalleryDB = &galleryGorm{}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

// ByID accepts an unsigned integer and returns either
// the gallery requested or an error.
func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

// All returns all galleries in the database
func (gg *galleryGorm) All() ([]Gallery, error) {
	var galleries []Gallery
	// db := gg.db.Where("user_id = ?", userID)
	// Pass a Find instead of First
	// and pass a slice instead of a single gallery
	if err := gg.db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

// ByUserID accepts an unsigned integer as the userID
// and returns all galleries owned by that user
func (gg *galleryGorm) ByUserID(userID uint) ([]Gallery, error) {
	var galleries []Gallery
	db := gg.db.Where("user_id = ?", userID)
	// Pass a Find instead of First
	// and pass a slice instead of a single gallery
	if err := db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&gallery).Error
}

type galleryValFn func(*Gallery) error

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

func (gv *galleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) nonZeroID(gallery *Gallery) error {
	if gallery.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.titleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

func (gv *galleryValidator) Delete(id uint) error {
	var gallery Gallery
	gallery.ID = id
	if err := runGalleryValFns(&gallery,
		gv.nonZeroID); err != nil {
		return err
	}
	return gv.GalleryDB.Delete(gallery.ID)
}

func (g *Gallery) ImagesSplitN(n int) [][]Image {
	// Create 2D slice
	ret := make([][]Image, n)
	// Create the inner slices - we need N of them, and we will
	// start them with a size of 0.
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}
	// Iterate over our images, using the index % n to determine
	// which of the slices in ret to add the image to.
	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}
