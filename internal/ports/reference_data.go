package ports

import (
	backgroundModel "modules/dndcharactersheet/internal/domain/background"
	classModel "modules/dndcharactersheet/internal/domain/class"
)

// BackgroundRepository provides access to background reference data.
type BackgroundRepository interface {
	LoadBackgrounds() ([]backgroundModel.Background, error)
	FindByName(name string) (*backgroundModel.Background, error)
}

// ClassRepository provides access to class reference data.
type ClassRepository interface {
	LoadClasses() ([]classModel.Class, error)
	FindByName(name string) (*classModel.Class, error)
}
