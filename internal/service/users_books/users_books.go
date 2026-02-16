package users_books

import (
	"context"
	"github.com/google/uuid"
)

func (ubserv *ubService) Buy(ctx context.Context, firebaseId string, bookIds []uuid.UUID) error {

	userId, err := ubserv.us.UserByFirebaseId(ctx, firebaseId)
	if err != nil {
		ubserv.l.Error("Error getting userId by firebaseId", "error", err)
		return err
	}
	if err := ubserv.ubs.Buy(ctx, userId, bookIds); err != nil {
		ubserv.l.Error("Error purchase a book", "error", err)
		return err
	}
	ubserv.l.Info("Successfully bought a book")
	return nil
}
