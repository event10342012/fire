package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		ctx     context.Context
		user    User
		wantErr error
	}{
		{
			name: "Insert success",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx:     context.Background(),
			user:    User{ID: 1},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			assert.NoError(t, err)
			userDAO := NewUserDAO(db)
			err = userDAO.Insert(tc.ctx, tc.user)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
