package physical_books

import (
	"context"

	"github.com/bookshop/internal/domain"
	"github.com/google/uuid"
)

func (pbserv *physicalBooksService) All(ctx context.Context) ([]domain.PhysicalBook, error) {
	physicalBooks, err := pbserv.pbs.All(ctx)
	if err != nil {
		pbserv.l.Error("Error getting physical books", "error", err)
		return nil, err
	}

	pbserv.l.Info("Successfully got physical books")
	return physicalBooks, nil
}

func (pbserv *physicalBooksService) ById(ctx context.Context, id uuid.UUID) (domain.PhysicalBook, error) {
	physicalBook, err := pbserv.pbs.ById(ctx, id)
	if err != nil {
		pbserv.l.Error("Error getting physical book", "error", err)
		return domain.PhysicalBook{}, err
	}

	pbserv.l.Info("Successfully got physical book", "id", id)
	return physicalBook, nil
}
