package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/devsujalpatel/olx-api/internal/httpx"
	"github.com/devsujalpatel/olx-api/internal/middleware"
)

type listing struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Price string `json:"price"`
	City string `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type ListingHandler struct {
	db *sql.DB
	logger *slog.Logger
}

// Contructor function
func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler {
		db: db,
		logger: logger,
	}
}

// GetAllListings
func (lh ListingHandler) List(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context() // request scoped context
		rows, err := lh.db.QueryContext(ctx,
				`SELECT id, title, description, price, city, created_at
			 	FROM listings
			 	ORDER BY created_at DESC
				LIMIT 100`)

		if err != nil {
			lh.logger.Error("listings query error", "err", err)
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
			return
		}

		defer rows.Close()
		listings := []listing{}
		for rows.Next() {
			var l listing
			if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
				lh.logger.Error("rows scan error", "err", err)
				httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
				return
			}

			lh.logger.Info("listings fetched", "total", len(listings))
			listings = append(listings, l)
		}

		if err := rows.Err(); err != nil {
			lh.logger.Error("rows error", "err", err)
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(listings)
}

// Delete Listing
func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context() // request scoped context
		requestId := middleware.RequestIDFromContext(ctx)
		id := r.PathValue("id")
		_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
		if err != nil {
			lh.logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)
			// http.Error(w, "internal error", http.StatusInternalServerError)
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
}
