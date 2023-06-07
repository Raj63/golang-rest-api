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

	orderService "github.com/Raj63/golang-rest-api/pkg/app/usecases/order"
	appErr "github.com/Raj63/golang-rest-api/pkg/domain/errors"
	domainOrder "github.com/Raj63/golang-rest-api/pkg/domain/order"
	mockRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/mocks/repository"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/repository"
	orderController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/order"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/routes"
	"github.com/brianvoe/gofakeit"
	"github.com/golang/mock/gomock"
)

func TestOrderRoutes(t *testing.T) {

	dinerID := gofakeit.Int64()
	dinerName := gofakeit.Name()
	menuID := gofakeit.Int64()
	menuName := gofakeit.BeerHop()
	menuDesc := gofakeit.BeerName()
	quantity := gofakeit.Number(1, 20)

	dinerName2 := gofakeit.Name()
	menuName2 := gofakeit.BeerHop()
	menuDesc2 := gofakeit.BeerName()
	quantity2 := gofakeit.Number(1, 20)
	type args struct {
		method       string
		endpoint     string
		body         interface{}
		mockrepoFn   func() repository.Orders
		outputStatus int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add new Order successfully",
			args: args{
				method:   "POST",
				endpoint: "/v1/orders/",
				body: orderController.NewOrderRequest{
					DinnerID: dinerID,
					MenuID:   menuID,
					Quantity: quantity,
				},
				outputStatus: http.StatusCreated,
				mockrepoFn: func() repository.Orders {
					mRepository := mockRepository.NewMockOrders(gomock.NewController(t))
					mRepository.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().Return(&domainOrder.Request{
						ID:        gofakeit.Int64(),
						DinnerID:  dinerID,
						MenuID:    menuID,
						Quantity:  quantity,
						CreatedAt: time.Now(),
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "Add new Order failed due to missing order menuid validation error",
			args: args{
				method:   "POST",
				endpoint: "/v1/orders/",
				body: orderController.NewOrderRequest{
					DinnerID: dinerID,
					Quantity: quantity,
				},
				outputStatus: http.StatusBadRequest,
				mockrepoFn: func() repository.Orders {
					mRepository := mockRepository.NewMockOrders(gomock.NewController(t))
					return mRepository
				},
			},
		},
		{
			name: "Fetch Order by ID successfully",
			args: args{
				method:       "GET",
				endpoint:     "/v1/orders/1",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Orders {
					mRepository := mockRepository.NewMockOrders(gomock.NewController(t))
					mRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).AnyTimes().Return([]domainOrder.Response{
						{
							ID:              gofakeit.Int64(),
							DinnerName:      dinerName,
							MenuName:        menuName,
							MenuDescription: menuDesc,
							Quantity:        quantity,
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
						{
							ID:              gofakeit.Int64(),
							DinnerName:      dinerName2,
							MenuName:        menuName2,
							MenuDescription: menuDesc2,
							Quantity:        quantity2,
							CreatedAt:       time.Now(),
							UpdatedAt:       time.Now(),
						},
					}, nil)
					return mRepository
				},
			},
		},
		{
			name: "Failed to fetch Order by ID due to repository error",
			args: args{
				method:       "GET",
				endpoint:     "/v1/orders/1",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Orders {
					mRepository := mockRepository.NewMockOrders(gomock.NewController(t))
					mRepository.EXPECT().GetByID(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, appErr.NewAppErrorWithType(appErr.RepositoryError))
					return mRepository
				},
			},
		},
		{
			name: "Deleted Order by ID successfully",
			args: args{
				method:       "DELETE",
				endpoint:     "/v1/orders/1",
				outputStatus: http.StatusOK,
				mockrepoFn: func() repository.Orders {
					mRepository := mockRepository.NewMockOrders(gomock.NewController(t))
					mRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
					return mRepository
				},
			},
		},
		{
			name: "Faiuled to delete Order by ID due to repository error",
			args: args{
				method:       "DELETE",
				endpoint:     "/v1/orders/1",
				outputStatus: http.StatusInternalServerError,
				mockrepoFn: func() repository.Orders {
					mRepository := mockRepository.NewMockOrders(gomock.NewController(t))
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
			routes.OrderRoutes(routerV1, &orderController.Controller{OrderService: orderService.Service{OrderRepository: tt.args.mockrepoFn()}})
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.args.outputStatus {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", tt.args.outputStatus, status)
			}
		})
	}
}
