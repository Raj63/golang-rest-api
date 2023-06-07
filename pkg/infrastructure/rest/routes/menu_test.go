// Package routes contains all routes of the application
package routes_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	menuService "github.com/Raj63/golang-rest-api/pkg/app/usecases/menu"
	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"

	appErr "github.com/Raj63/golang-rest-api/pkg/domain/errors"
	mockRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/mocks/repository"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/repository"
	menuController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/menu"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/routes"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/mock/gomock"
)

func TestMenuRoutes(t *testing.T) {

	menuName := gofakeit.BeerHop()
	menuDesc1 := gofakeit.BeerName()
	menuPrice1 := gofakeit.Float64()
	menuName2 := gofakeit.BeerHop()
	menuDesc2 := gofakeit.BeerName()
	menuPrice2 := gofakeit.Float64()
	menuName3 := gofakeit.BeerHop()
	menuDesc3 := gofakeit.BeerName()
	menuPrice3 := gofakeit.Float64()
	type args struct {
		method       string
		endpoint     string
		body         interface{}
		mockrepoFn   func() repository.Menus
		outputStatus int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add new Menu successfully",
			args: args{
				method:   "POST",
				endpoint: "/v1/menus/",
				body: menuController.NewMenuRequest{
					Name:        menuName,
					Description: menuDesc1,
					Price:       float64(menuPrice1),
				},
				outputStatus: http.StatusCreated,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(&domainMenu.Menu{
						ID:          gofakeit.Int64(),
						Name:        menuName,
						Description: menuDesc1,
						Price:       menuPrice1,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "Add new Menu failed due to missing menu desc validation error",
			args: args{
				method:   "POST",
				endpoint: "/v1/menus/",
				body: menuController.NewMenuRequest{
					Name:  menuName,
					Price: float64(menuPrice1),
				},
				outputStatus: http.StatusBadRequest,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					return mRepository
				},
			},
		},
		{
			name: "Fetch Menu List successfully",
			args: args{
				method:       "GET",
				endpoint:     "/v1/menus/?page=1&limit=10",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().GetTotalCount(gomock.Any()).AnyTimes().Return(int64(2), nil)
					mRepository.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&repository.PaginationResultMenu{
						Data: &[]domainMenu.Menu{
							{
								ID:          gofakeit.Int64(),
								Name:        menuName,
								Description: menuDesc1,
								Price:       menuPrice1,
								CreatedAt:   time.Now(),
								UpdatedAt:   time.Now(),
							},
							{
								ID:          gofakeit.Int64(),
								Name:        menuName2,
								Description: menuDesc2,
								Price:       menuPrice2,
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
			name: "failed to fetch Menu List due to repository error",
			args: args{
				method:       "GET",
				endpoint:     "/v1/menus/?page=1&limit=10",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, appErr.NewAppErrorWithType(appErr.RepositoryError))
					return mRepository
				},
			},
		},
		{
			name: "Fetch Menu top 3 List successfully",
			args: args{
				method:       "GET",
				endpoint:     "/v1/menus/top?count=3",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().GetTotalCount(gomock.Any()).AnyTimes().Return(int64(2), nil)
					mRepository.EXPECT().GetByTopCount(gomock.Any(), gomock.Any()).AnyTimes().Return([]domainMenu.Menu{
						{
							ID:          gofakeit.Int64(),
							Name:        menuName,
							Description: menuDesc1,
							Price:       menuPrice1,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
						{
							ID:          gofakeit.Int64(),
							Name:        menuName2,
							Description: menuDesc2,
							Price:       menuPrice2,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
						{
							ID:          gofakeit.Int64(),
							Name:        menuName3,
							Description: menuDesc3,
							Price:       menuPrice3,
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						},
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "Failed to fetch Menu top 3 List due to repository error",
			args: args{
				method:       "GET",
				endpoint:     "/v1/menus/top?count=3",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().GetTotalCount(gomock.Any()).AnyTimes().Return(int64(2), nil)
					mRepository.EXPECT().GetByTopCount(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, appErr.NewAppErrorWithType(appErr.RepositoryError))
					return mRepository
				},
			},
		},
		{
			name: "Fetch Menu by ID successfully",
			args: args{
				method:       "GET",
				endpoint:     "/v1/menus/1",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).AnyTimes().Return(&domainMenu.Menu{
						ID:          gofakeit.Int64(),
						Name:        menuName,
						Description: menuDesc1,
						Price:       menuPrice1,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "Failed to fetch Menu by ID due to repository error",
			args: args{
				method:       "GET",
				endpoint:     "/v1/menus/1",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, appErr.NewAppErrorWithType(appErr.RepositoryError))
					return mRepository
				},
			},
		},
		{
			name: "Deleted Menu by ID successfully",
			args: args{
				method:       "DELETE",
				endpoint:     "/v1/menus/1",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
					mRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
					return mRepository
				},
			},
		},
		{
			name: "Faiuled to delete Menu by ID due to repository error",
			args: args{
				method:       "DELETE",
				endpoint:     "/v1/menus/1",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Menus {
					mRepository := mockRepository.NewMockMenus(gomock.NewController(t))
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
			routes.MenuRoutes(routerV1, &menuController.Controller{MenuService: menuService.Service{MenuRepository: tt.args.mockrepoFn()}})
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.args.outputStatus {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", tt.args.outputStatus, status)
			}
		})
	}
}
