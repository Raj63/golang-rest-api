package routes_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dinerService "github.com/Raj63/golang-rest-api/pkg/app/usecases/diner"
	domainDiner "github.com/Raj63/golang-rest-api/pkg/domain/diner"
	appErr "github.com/Raj63/golang-rest-api/pkg/domain/errors"
	mockRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/mocks/repository"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/repository"
	dinerController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/diner"
	errorsController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/errors"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/routes"
	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func getTestRouter() (*gin.Engine, *gin.RouterGroup) {
	// initialize the router
	router := gin.Default()
	// the application errors will be processed here before returning to the caller
	router.Use(errorsController.Handler)

	return router, router.Group("/v1")
}

func TestDinerRoutes(t *testing.T) {

	dinerName := gofakeit.Name()
	dinertable1 := gofakeit.Number(1, 20)
	dinerName2 := gofakeit.Name()
	dinertable2 := gofakeit.Number(1, 20)
	type args struct {
		method       string
		endpoint     string
		body         interface{}
		mockrepoFn   func() repository.Diners
		outputStatus int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add new Diner successfully",
			args: args{
				method:   "POST",
				endpoint: "/v1/diners/",
				body: dinerController.NewDinerRequest{
					Name:        dinerName,
					TableNumber: dinertable1,
				},
				outputStatus: http.StatusCreated,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					mRepository.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(&domainDiner.Diner{
						ID:          gofakeit.Int64(),
						Name:        dinerName,
						TableNumber: dinertable1,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "Add new Diner failed due to missing diner name validation error",
			args: args{
				method:   "POST",
				endpoint: "/v1/diners/",
				body: dinerController.NewDinerRequest{
					TableNumber: dinertable1,
				},
				outputStatus: http.StatusBadRequest,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					return mRepository
				},
			},
		},
		{
			name: "Fetch Diner List successfully",
			args: args{
				method:       "GET",
				endpoint:     "/v1/diners/?page=1&limit=10",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					mRepository.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&repository.PaginationResultDiner{
						Data: &[]domainDiner.Diner{
							{
								ID:          gofakeit.Int64(),
								Name:        dinerName,
								TableNumber: dinertable1,
								CreatedAt:   time.Now(),
								UpdatedAt:   time.Now(),
							},
							{
								ID:          gofakeit.Int64(),
								Name:        dinerName2,
								TableNumber: dinertable2,
								CreatedAt:   time.Now(),
								UpdatedAt:   time.Now(),
							},
						},
						Total:   2,
						Limit:   10,
						Current: 1,
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "failed to fetch Diner List due to repository error",
			args: args{
				method:       "GET",
				endpoint:     "/v1/diners/?page=1&limit=10",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					mRepository.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, appErr.NewAppErrorWithType(appErr.RepositoryError))
					return mRepository
				},
			},
		},
		{
			name: "Fetch Diner by ID successfully",
			args: args{
				method:       "GET",
				endpoint:     "/v1/diners/1",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					mRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).AnyTimes().Return(&domainDiner.Diner{
						ID:          gofakeit.Int64(),
						Name:        dinerName,
						TableNumber: dinertable1,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "Failed to fetch Diner by ID due to repository error",
			args: args{
				method:       "GET",
				endpoint:     "/v1/diners/1",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					mRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, appErr.NewAppErrorWithType(appErr.RepositoryError))
					return mRepository
				},
			},
		},
		{
			name: "Deleted Diner by ID successfully",
			args: args{
				method:       "DELETE",
				endpoint:     "/v1/diners/1",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					mRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
					return mRepository
				},
			},
		},
		{
			name: "Faiuled to delete Diner by ID due to repository error",
			args: args{
				method:       "DELETE",
				endpoint:     "/v1/diners/1",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Diners {
					mRepository := mockRepository.NewMockDiners(gomock.NewController(t))
					mRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().Return(appErr.NewAppErrorWithType(appErr.RepositoryError))
					return mRepository
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if tt.args.body != nil {
				err := json.NewEncoder(&buf).Encode(tt.args.body)
				if err != nil {
					log.Fatal(err)
				}
			}

			req, err := http.NewRequest(tt.args.method, tt.args.endpoint, &buf)
			if err != nil {
				t.Errorf("Error creating a new request: %v", err)
			}
			rr := httptest.NewRecorder()
			router, routerV1 := getTestRouter()
			routes.DinerRoutes(routerV1, &dinerController.Controller{DinerService: dinerService.Service{DinerRepository: tt.args.mockrepoFn()}})
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.args.outputStatus {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", tt.args.outputStatus, status)
			}
		})
	}
}
