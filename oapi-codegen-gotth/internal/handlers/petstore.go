//revive:disable
package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/internal/views"
	"github.com/patilchinmay/go-experiments/oapi-codegen-gotth/pkg/spec/generated"
)

type PetStore struct {
	Pets   map[int64]generated.Pet
	NextId int64
	Lock   sync.Mutex
}

func NewPetStore() *PetStore {
	return &PetStore{
		Pets:   make(map[int64]generated.Pet),
		NextId: 1000,
	}
}

// sendPetStoreError wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendPetStoreError(ctx echo.Context, code int, message string) error {
	petErr := generated.Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}

// PetStore implements all the handlers in the ServerInterface
// GET /pets
func (p *PetStore) FindPets(ctx echo.Context, params generated.FindPetsParams) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	var result []generated.Pet

	for _, pet := range p.Pets {
		if params.Tags != nil {
			// If we have tags,  filter pets by tag
			for _, t := range *params.Tags {
				if pet.Tag != nil && (*pet.Tag == t) {
					result = append(result, pet)
				}
			}
		} else {
			// Add all pets if we're not filtering
			result = append(result, pet)
		}

		if params.Limit != nil {
			l := int(*params.Limit)
			if len(result) >= l {
				// We're at the limit
				break
			}
		}
	}

	// Check if request is from HTMX
	if ctx.Request().Header.Get("HX-Request") == "true" {
		return views.PetList(result).Render(ctx.Request().Context(), ctx.Response().Writer)
	}

	// Return JSON for API clients
	return ctx.JSON(http.StatusOK, result)
}

func (p *PetStore) AddPet(ctx echo.Context) error {
	// We expect a NewPet object in the request body.
	var newPet generated.NewPet
	err := ctx.Bind(&newPet)
	if err != nil {
		return sendPetStoreError(ctx, http.StatusBadRequest, "Invalid format for NewPet")
	}
	// We now have a pet, let's add it to our "database".

	// We're always asynchronous, so lock unsafe operations below
	p.Lock.Lock()
	defer p.Lock.Unlock()

	// We handle pets, not NewPets, which have an additional ID field
	var pet generated.Pet
	pet.Name = newPet.Name
	pet.Tag = newPet.Tag
	pet.Id = p.NextId
	p.NextId++

	// Insert into map
	p.Pets[pet.Id] = pet

	// Handle HTMX request
	if ctx.Request().Header.Get("HX-Request") == "true" {
		return views.PetCard(pet).Render(ctx.Request().Context(), ctx.Response().Writer)
	}

	// Now, we have to return the NewPet
	err = ctx.JSON(http.StatusCreated, pet)
	if err != nil {
		// Something really bad happened, tell Echo that our handler failed
		return err
	}

	// Return no error. This refers to the handler. Even if we return an HTTP
	// error, but everything else is working properly, tell Echo that we serviced
	// the error. We should only return errors from Echo handlers if the actual
	// servicing of the error on the infrastructure level failed. Returning an
	// HTTP/400 or HTTP/500 from here means Echo/HTTP are still working, so
	// return nil.
	return nil
}

func (p *PetStore) FindPetByID(ctx echo.Context, petId int64) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	pet, found := p.Pets[petId]
	if !found {
		return sendPetStoreError(ctx, http.StatusNotFound,
			fmt.Sprintf("Could not find pet with ID %d", petId))
	}

	// Handle HTMX request
	if ctx.Request().Header.Get("HX-Request") == "true" {
		return views.PetDetail(pet).Render(ctx.Request().Context(), ctx.Response().Writer)
	}

	return ctx.JSON(http.StatusOK, pet)
}

func (p *PetStore) DeletePet(ctx echo.Context, id int64) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, found := p.Pets[id]
	if !found {
		return sendPetStoreError(ctx, http.StatusNotFound,
			fmt.Sprintf("Could not find pet with ID %d", id))
	}
	delete(p.Pets, id)

	// Handle HTMX request
	if ctx.Request().Header.Get("HX-Request") == "true" {
		// Send empty response with HX-Trigger
		ctx.Response().Header().Set("HX-Trigger", fmt.Sprintf("petDeleted-%d", id))
		return ctx.NoContent(http.StatusOK)
	}

	return ctx.NoContent(http.StatusNoContent)
}
