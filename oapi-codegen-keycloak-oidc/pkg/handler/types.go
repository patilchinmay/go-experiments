package handler

import (
	"log/slog"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/patilchinmay/go-experiments/oapi-codegen-keycloak-oidc/api/spec"
)

type HandlerIntf interface {
	// Returns all pets
	// (GET /pets)
	FindPets(ctx echo.Context, params spec.FindPetsParams) error
	// Creates a new pet
	// (POST /pets)
	AddPet(ctx echo.Context) error
	// Deletes a pet by ID
	// (DELETE /pets/{id})
	DeletePet(ctx echo.Context, id int64) error
	// Returns a pet by ID
	// (GET /pets/{id})
	FindPetByID(ctx echo.Context, id int64) error
}

type PetStore struct {
	Pets   map[int64]spec.Pet
	NextId int64
	Lock   sync.Mutex
	logger *slog.Logger
}

var _ HandlerIntf = &PetStore{}

var p *PetStore

func NewPetStore(logger *slog.Logger) *PetStore {
	if p == nil {
		p = &PetStore{
			Pets:   make(map[int64]spec.Pet),
			NextId: 1000,
			logger: logger,
		}
	}

	logger.Info("petstore handler initialized")

	return p
}
