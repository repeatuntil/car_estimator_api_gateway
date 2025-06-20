package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"log/slog"

	"github.com/gorilla/mux"
	feed "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/feed_v1"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/domain"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/utils"
	"google.golang.org/grpc"
)

type FeedHandler struct {
	r      *mux.Router
	logger *slog.Logger
	client feed.FeedServiceClient
}

func (h *FeedHandler) setupgRPC(conn *grpc.ClientConn) {
	h.client = feed.NewFeedServiceClient(conn)
}

func (h *FeedHandler) setupRoutes() {
	h.r.HandleFunc("/listings", h.ListListings).Methods("GET")
	h.r.HandleFunc("/listings/search", h.SearchListings).Methods("GET")
	h.r.HandleFunc("/listings/{listingId}", h.GetListing).Methods("GET")
	h.r.HandleFunc("/listings", h.CreateListing).Methods("POST")
	h.r.HandleFunc("/listings/{listingId}", h.UpdateListing).Methods("PUT")
	h.r.HandleFunc("/listings/{listingId}", h.DeleteListing).Methods("DELETE")
	h.r.HandleFunc("/users/{userId}/favorites", h.AddToFavorites).Methods("POST")
}

func (h *FeedHandler) ListListings(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "ListListings"))

	q := r.URL.Query()
	pageNum, _ := strconv.Atoi(q.Get("page_number"))
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	sortByStr := strings.ToUpper(q.Get("sort_by"))
	sortBy := feed.SortBy_SORT_UNSPECIFIED
	if v, ok := feed.SortBy_value[sortByStr]; ok {
		sortBy = feed.SortBy(v)
	}

	grpcReq := &feed.ListListingsRequest{
		Page: &feed.PageRequest{
			PageNumber: int32(pageNum),
			PageSize:   int32(pageSize),
		},
		SortBy: sortBy,
	}

	log.Info("→ gRPC ListListings", slog.Any("req", grpcReq))
	grpcResp, err := h.client.ListListings(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "ListListings failed: ", err)
		return
	}

	out := domain.ListListingsResponse{
		PageMetadata: domain.PageResponseMetadata{
			TotalItems:  grpcResp.PageMetadata.TotalItems,
			TotalPages:  grpcResp.PageMetadata.TotalPages,
			CurrentPage: grpcResp.PageMetadata.CurrentPage,
		},
		Listings: make([]domain.CarListing, len(grpcResp.Listings)),
	}
	for i, c := range grpcResp.Listings {
		out.Listings[i] = domain.CarListing{
			ListingId:        c.ListingId,
			SellerId:         c.SellerId,
			Description:      c.Description,
			PostedAt:         c.PostedAt.AsTime(),
			Status:           c.Status,
			DealType:         c.DealType,
			Price:            c.Price,
			CarId:            c.CarId,
			Mileage:          c.Mileage,
			OwnersCount:      c.OwnersCount,
			AccidentsCount:   c.AccidentsCount,
			Condition:        c.Condition,
			Color:            c.Color,
			ConfigId:         c.ConfigId,
			EngineType:       c.EngineType,
			EngineVolume:     c.EngineVolume,
			EnginePower:      c.EnginePower,
			Cylinders:        c.Cylinders,
			Transmission:     c.Transmission,
			Drivetrain:       c.Drivetrain,
			ModelId:          c.ModelId,
			ModelName:        c.ModelName,
			Make:             c.Make,
			Year:             c.Year,
			BodyType:         c.BodyType,
			Generation:       c.Generation,
			WeightKg:         c.WeightKg,
			SellerName:       c.SellerName,
			SellerRating:     c.SellerRating,
			SellerSalesCount: c.SellerSalesCount,
			SellerIsBusiness: c.SellerIsBusiness,
		}
	}

	utils.RenderJson(w, out)
}

func (h *FeedHandler) SearchListings(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "SearchListings"))

	q := r.URL.Query()
	query := q.Get("query")
	pageNum, _ := strconv.Atoi(q.Get("page_number"))
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	sortByStr := strings.ToUpper(q.Get("sort_by"))
	sortBy := feed.SortBy_SORT_UNSPECIFIED
	if v, ok := feed.SortBy_value[sortByStr]; ok {
		sortBy = feed.SortBy(v)
	}

	grpcReq := &feed.SearchListingsRequest{
		Query: query,
		Page: &feed.PageRequest{
			PageNumber: int32(pageNum),
			PageSize:   int32(pageSize),
		},
		SortBy: sortBy,
	}

	log.Info("→ gRPC SearchListings", slog.Any("req", grpcReq))
	grpcResp, err := h.client.SearchListings(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "SearchListings failed: ", err)
		return
	}

	out := domain.SearchListingsResponse{
		PageMetadata: domain.PageResponseMetadata{
			TotalItems:  grpcResp.PageMetadata.TotalItems,
			TotalPages:  grpcResp.PageMetadata.TotalPages,
			CurrentPage: grpcResp.PageMetadata.CurrentPage,
		},
		Listings: make([]domain.CarListing, len(grpcResp.Listings)),
	}
	for i, c := range grpcResp.Listings {
		out.Listings[i] = domain.CarListing{
			ListingId:        c.ListingId,
			SellerId:         c.SellerId,
			Description:      c.Description,
			PostedAt:         c.PostedAt.AsTime(),
			Status:           c.Status,
			DealType:         c.DealType,
			Price:            c.Price,
			CarId:            c.CarId,
			Mileage:          c.Mileage,
			OwnersCount:      c.OwnersCount,
			AccidentsCount:   c.AccidentsCount,
			Condition:        c.Condition,
			Color:            c.Color,
			ConfigId:         c.ConfigId,
			EngineType:       c.EngineType,
			EngineVolume:     c.EngineVolume,
			EnginePower:      c.EnginePower,
			Cylinders:        c.Cylinders,
			Transmission:     c.Transmission,
			Drivetrain:       c.Drivetrain,
			ModelId:          c.ModelId,
			ModelName:        c.ModelName,
			Make:             c.Make,
			Year:             c.Year,
			BodyType:         c.BodyType,
			Generation:       c.Generation,
			WeightKg:         c.WeightKg,
			SellerName:       c.SellerName,
			SellerRating:     c.SellerRating,
			SellerSalesCount: c.SellerSalesCount,
			SellerIsBusiness: c.SellerIsBusiness,
		}
	}

	utils.RenderJson(w, out)
}

func (h *FeedHandler) GetListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "GetListing"))

	id := mux.Vars(r)["listingId"]
	grpcReq := &feed.GetListingRequest{ListingId: id}

	log.Info("→ gRPC GetListing", slog.String("id", id))
	grpcResp, err := h.client.GetListing(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "GetListing failed: ", err)
		return
	}

	out := domain.GetListingResponse{
		Listing: domain.CarListing{
			ListingId:        grpcResp.Listing.ListingId,
			SellerId:         grpcResp.Listing.SellerId,
			Description:      grpcResp.Listing.Description,
			PostedAt:         grpcResp.Listing.PostedAt.AsTime(),
			Status:           grpcResp.Listing.Status,
			DealType:         grpcResp.Listing.DealType,
			Price:            grpcResp.Listing.Price,
			CarId:            grpcResp.Listing.CarId,
			Mileage:          grpcResp.Listing.Mileage,
			OwnersCount:      grpcResp.Listing.OwnersCount,
			AccidentsCount:   grpcResp.Listing.AccidentsCount,
			Condition:        grpcResp.Listing.Condition,
			Color:            grpcResp.Listing.Color,
			ConfigId:         grpcResp.Listing.ConfigId,
			EngineType:       grpcResp.Listing.EngineType,
			EngineVolume:     grpcResp.Listing.EngineVolume,
			EnginePower:      grpcResp.Listing.EnginePower,
			Cylinders:        grpcResp.Listing.Cylinders,
			Transmission:     grpcResp.Listing.Transmission,
			Drivetrain:       grpcResp.Listing.Drivetrain,
			ModelId:          grpcResp.Listing.ModelId,
			ModelName:        grpcResp.Listing.ModelName,
			Make:             grpcResp.Listing.Make,
			Year:             grpcResp.Listing.Year,
			BodyType:         grpcResp.Listing.BodyType,
			Generation:       grpcResp.Listing.Generation,
			WeightKg:         grpcResp.Listing.WeightKg,
			SellerName:       grpcResp.Listing.SellerName,
			SellerRating:     grpcResp.Listing.SellerRating,
			SellerSalesCount: grpcResp.Listing.SellerSalesCount,
			SellerIsBusiness: grpcResp.Listing.SellerIsBusiness,
		},
	}

	utils.RenderJson(w, out)
}

func (h *FeedHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "CreateListing"))

	var body domain.CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := &feed.CreateListingRequest{Listing: &feed.CarListing{
		SellerId:    body.Listing.SellerId,
		Description: body.Listing.Description,
		Status:      body.Listing.Status,
		DealType:    body.Listing.DealType,
		Price:       body.Listing.Price,
		Tags:		 body.Listing.Tags,
		CarId: 		 body.Listing.CarId,
		Mileage: 	 body.Listing.Mileage,
		OwnersCount: body.Listing.OwnersCount,
		AccidentsCount: body.Listing.AccidentsCount,
		Condition: 	 body.Listing.Condition,
		Color: 		 body.Listing.Color,
		ConfigId: 	 body.Listing.ConfigId,
		EngineType:       body.Listing.EngineType,
		EngineVolume:     body.Listing.EngineVolume,
		EnginePower:      body.Listing.EnginePower,
		Cylinders:        body.Listing.Cylinders,
		Transmission:     body.Listing.Transmission,
		Drivetrain:       body.Listing.Drivetrain,
		ModelId:          body.Listing.ModelId,
		ModelName:        body.Listing.ModelName,
		Make:             body.Listing.Make,
		Year:             body.Listing.Year,
		BodyType:         body.Listing.BodyType,
		Generation:       body.Listing.Generation,
		WeightKg:         body.Listing.WeightKg,
		SellerName:       body.Listing.SellerName,
		SellerRating:     body.Listing.SellerRating,
		SellerSalesCount: body.Listing.SellerSalesCount,
		SellerIsBusiness: body.Listing.SellerIsBusiness,
	}}
	log.Info("→ gRPC CreateListing", slog.Any("req", grpcReq))

	grpcResp, err := h.client.CreateListing(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "CreateListing failed: ", err)
		return
	}

	out := domain.CreateListingResponse{
		Listing: domain.CarListing{
			ListingId:        grpcResp.Listing.ListingId,
			SellerId:         grpcResp.Listing.SellerId,
			Description:      grpcResp.Listing.Description,
			PostedAt:         grpcResp.Listing.PostedAt.AsTime(),
			Status:           grpcResp.Listing.Status,
			DealType:         grpcResp.Listing.DealType,
			Price:            grpcResp.Listing.Price,
			CarId:            grpcResp.Listing.CarId,
			Mileage:          grpcResp.Listing.Mileage,
			OwnersCount:      grpcResp.Listing.OwnersCount,
			AccidentsCount:   grpcResp.Listing.AccidentsCount,
			Condition:        grpcResp.Listing.Condition,
			Color:            grpcResp.Listing.Color,
			ConfigId:         grpcResp.Listing.ConfigId,
			EngineType:       grpcResp.Listing.EngineType,
			EngineVolume:     grpcResp.Listing.EngineVolume,
			EnginePower:      grpcResp.Listing.EnginePower,
			Cylinders:        grpcResp.Listing.Cylinders,
			Transmission:     grpcResp.Listing.Transmission,
			Drivetrain:       grpcResp.Listing.Drivetrain,
			ModelId:          grpcResp.Listing.ModelId,
			ModelName:        grpcResp.Listing.ModelName,
			Make:             grpcResp.Listing.Make,
			Year:             grpcResp.Listing.Year,
			BodyType:         grpcResp.Listing.BodyType,
			Generation:       grpcResp.Listing.Generation,
			WeightKg:         grpcResp.Listing.WeightKg,
			SellerName:       grpcResp.Listing.SellerName,
			SellerRating:     grpcResp.Listing.SellerRating,
			SellerSalesCount: grpcResp.Listing.SellerSalesCount,
			SellerIsBusiness: grpcResp.Listing.SellerIsBusiness,
		},
	}
	utils.RenderJson(w, out)
}

func (h *FeedHandler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "UpdateListing"))

	listingId := mux.Vars(r)["listingId"]
	var body domain.UpdateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	body.Listing.ListingId = listingId

	grpcReq := &feed.UpdateListingRequest{
		Listing: &feed.CarListing{
			SellerId:    body.Listing.SellerId,
			Description: body.Listing.Description,
			Status:      body.Listing.Status,
			DealType:    body.Listing.DealType,
			Price:       body.Listing.Price,
			Tags:		 body.Listing.Tags,
			CarId: 		 body.Listing.CarId,
			Mileage: 	 body.Listing.Mileage,
			OwnersCount: body.Listing.OwnersCount,
			AccidentsCount: body.Listing.AccidentsCount,
			Condition: 	 body.Listing.Condition,
			Color: 		 body.Listing.Color,
			ConfigId: 	 body.Listing.ConfigId,
			EngineType:       body.Listing.EngineType,
			EngineVolume:     body.Listing.EngineVolume,
			EnginePower:      body.Listing.EnginePower,
			Cylinders:        body.Listing.Cylinders,
			Transmission:     body.Listing.Transmission,
			Drivetrain:       body.Listing.Drivetrain,
			ModelId:          body.Listing.ModelId,
			ModelName:        body.Listing.ModelName,
			Make:             body.Listing.Make,
			Year:             body.Listing.Year,
			BodyType:         body.Listing.BodyType,
			Generation:       body.Listing.Generation,
			WeightKg:         body.Listing.WeightKg,
			SellerName:       body.Listing.SellerName,
			SellerRating:     body.Listing.SellerRating,
			SellerSalesCount: body.Listing.SellerSalesCount,
			SellerIsBusiness: body.Listing.SellerIsBusiness,
		},
	}

	log.Info("→ gRPC UpdateListing", slog.Any("req", grpcReq))
	grpcResp, err := h.client.UpdateListing(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "UpdateListing failed: ", err)
		return
	}

	out := domain.UpdateListingResponse{
		Listing: domain.CarListing{
			ListingId:        grpcResp.Listing.ListingId,
			SellerId:         grpcResp.Listing.SellerId,
			Description:      grpcResp.Listing.Description,
			PostedAt:         grpcResp.Listing.PostedAt.AsTime(),
			Status:           grpcResp.Listing.Status,
			DealType:         grpcResp.Listing.DealType,
			Price:            grpcResp.Listing.Price,
			CarId:            grpcResp.Listing.CarId,
			Mileage:          grpcResp.Listing.Mileage,
			OwnersCount:      grpcResp.Listing.OwnersCount,
			AccidentsCount:   grpcResp.Listing.AccidentsCount,
			Condition:        grpcResp.Listing.Condition,
			Color:            grpcResp.Listing.Color,
			ConfigId:         grpcResp.Listing.ConfigId,
			EngineType:       grpcResp.Listing.EngineType,
			EngineVolume:     grpcResp.Listing.EngineVolume,
			EnginePower:      grpcResp.Listing.EnginePower,
			Cylinders:        grpcResp.Listing.Cylinders,
			Transmission:     grpcResp.Listing.Transmission,
			Drivetrain:       grpcResp.Listing.Drivetrain,
			ModelId:          grpcResp.Listing.ModelId,
			ModelName:        grpcResp.Listing.ModelName,
			Make:             grpcResp.Listing.Make,
			Year:             grpcResp.Listing.Year,
			BodyType:         grpcResp.Listing.BodyType,
			Generation:       grpcResp.Listing.Generation,
			WeightKg:         grpcResp.Listing.WeightKg,
			SellerName:       grpcResp.Listing.SellerName,
			SellerRating:     grpcResp.Listing.SellerRating,
			SellerSalesCount: grpcResp.Listing.SellerSalesCount,
			SellerIsBusiness: grpcResp.Listing.SellerIsBusiness,
		},
	}
	utils.RenderJson(w, out)
}

func (h *FeedHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "DeleteListing"))

	listingId := mux.Vars(r)["listingId"]
	log.Info("→ gRPC DeleteListing", slog.String("Id", listingId))

	grpcResp, err := h.client.DeleteListing(r.Context(), &feed.DeleteListingRequest{
		ListingId: listingId,
	})
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "DeleteListing failed: ", err)
		return
	}

	out := domain.DeleteListingResponse{Success: grpcResp.Success}
	utils.RenderJson(w, out)
}

func (h *FeedHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "AddToFavorites"))

	userId := mux.Vars(r)["userId"]
	var body domain.AddToFavoritesRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := &feed.AddToFavoritesRequest{
		UserId:    userId,
		ListingId: body.ListingId,
	}

	log.Info("→ gRPC AddToFavorites", slog.Any("req", grpcReq))
	grpcResp, err := h.client.AddToFavorites(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "AddToFavorites failed: ", err)
		return
	}

	out := domain.AddToFavoritesResponse{Success: grpcResp.Success}
	utils.RenderJson(w, out)
}