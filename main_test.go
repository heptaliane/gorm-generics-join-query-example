package main

import (
	"context"
	"log"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestJoinQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}

	cases := map[string]struct {
		g        gorm.ChainInterface[JoinedOrder]
		expected JoinedOrder
	}{
		"not_null": {
			g: gorm.G[JoinedOrder](db).
				Joins(clause.LeftJoin.Association("User"), nil).
				Joins(clause.LeftJoin.Association("Product"), nil),
			expected: JoinedOrder{
				Order: Order{
					Id:        1,
					UserId:    1,
					ProductId: 1,
				},
				User:    &User{Id: 1, Name: "Alice"},
				Product: &Product{Id: 1, Name: "Apple", Price: 100},
			},
		},
		"null": {
			g: gorm.G[JoinedOrder](db).
				Joins(clause.LeftJoin.Association("User"), nil).
				Joins(clause.LeftJoin.Association("Product"), nil).
				Where("?.id = ?", clause.Table{Name: "joined_orders"}, 4),
			expected: JoinedOrder{
				Order: Order{
					Id: 4,
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			j, err := tc.g.First(context.Background())
			if err != nil {
				t.Errorf("Failed to query sql: %v", err)
				return
			}
			if tc.expected.Id != j.Id {
				t.Errorf(
					"Mismatched order.id. Expected %d, but got %d",
					tc.expected.Id,
					j.Id,
				)
				return
			}
			if tc.expected.UserId != j.UserId {
				t.Errorf(
					"Mismatched user_id. Expected %d, but got %d",
					tc.expected.UserId,
					j.UserId,
				)
				return
			}
			if tc.expected.ProductId != j.ProductId {
				t.Errorf(
					"Mismatched product.id. Expected %d, but got %d",
					tc.expected.ProductId,
					j.ProductId,
				)
				return
			}
			if tc.expected.User == nil {
				if j.User != nil {
					t.Errorf("`user` is expected to be nil, but got %v.", j.User)
					return
				}
			} else {
				if tc.expected.User.Id != j.User.Id {
					t.Errorf(
						"Mismatched user.id. Expected %d, but got %d",
						tc.expected.User.Id,
						j.User.Id,
					)
					return
				}
				if tc.expected.User.Name != j.User.Name {
					t.Errorf(
						"Mismatched user.name. Expected '%s', but got '%s'",
						tc.expected.User.Name,
						j.User.Name,
					)
					return
				}
			}
			if tc.expected.Product == nil {
				if j.Product != nil {
					t.Errorf("`product` is expected to be nil, but got %v.", j.Product)
					return
				}
			} else {
				if tc.expected.Product.Id != j.Product.Id {
					t.Errorf(
						"Mismatched product.id. Expected %d, but got %d",
						tc.expected.Product.Id,
						j.Product.Id,
					)
					return
				}
				if tc.expected.Product.Name != j.Product.Name {
					t.Errorf(
						"Mismatched product.name. Expected '%s', but got '%s'",
						tc.expected.Product.Name,
						j.Product.Name,
					)
					return
				}
				if tc.expected.Product.Price != j.Product.Price {
					t.Errorf(
						"Mismatched product.price. Expected %d, but got %d",
						tc.expected.Product.Price,
						j.Product.Price,
					)
					return
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the test database: %v", err)
	}

	err = db.AutoMigrate(
		&User{},
		&Product{},
		&JoinedOrder{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate the schema: %v", err)
	}

	users := []User{
		{Id: 1, Name: "Alice"},
		{Id: 2, Name: "Bob"},
	}
	for _, u := range users {
		if err := gorm.G[User](db).Create(context.Background(), &u); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
	}
	products := []Product{
		{Id: 1, Name: "Apple", Price: 100},
		{Id: 2, Name: "Banana", Price: 200},
	}
	for _, p := range products {
		if err := gorm.G[Product](db).Create(context.Background(), &p); err != nil {
			log.Fatalf("Failed to create product: %v", err)
		}
	}
	orders := []Order{
		{Id: 1, UserId: 1, ProductId: 1},
		{Id: 2, UserId: 1, ProductId: 2},
		{Id: 3, UserId: 2, ProductId: 1},
		{Id: 4},
	}
	for _, o := range orders {
		j := JoinedOrder{
			Order: o,
		}
		if err := gorm.G[JoinedOrder](db).Create(context.Background(), &j); err != nil {
			log.Fatalf("Failed to create order: %v", err)
		}
	}

	code := m.Run()

	// Clean up
	os.Remove("test.db")

	os.Exit(code)
}
