package service

import (
	"payslip-system/internal/models"
	"payslip-system/internal/repository"
	mock_repository "payslip-system/internal/repository/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_adminService_CreateAttendancePeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminID := uuid.New() // Use a valid UUID for adminID
	type args struct {
		startDate time.Time
		endDate   time.Time
		adminID   uuid.UUID
		ipAddress string
		requestID string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AttendancePeriod
		wantErr bool
	}{
		{
			name: "success - valid period",
			args: args{
				startDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
				adminID:   adminID,
				ipAddress: "127.0.0.1",
				requestID: "req-123",
			},
			want: &models.AttendancePeriod{
				BaseModel: models.BaseModel{
					CreatedBy: &adminID, // will be replaced below
					IPAddress: "127.0.0.1",
					RequestID: "req-123",
				},
				StartDate:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
				IsProcessed: false,
			},
			wantErr: false,
		},
		{
			name: "error - end date before start date",
			args: args{
				startDate: time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
				endDate:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				adminID:   adminID,
				ipAddress: "127.0.0.1",
				requestID: "req-456",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAttendancePeriodRepo := mock_repository.NewMockIAttendancePeriodRepository(ctrl)
			mockAuditLogRepo := mock_repository.NewMockIAuditLogRepository(ctrl)

			if tt.args.endDate.After(tt.args.startDate) {
				mockAttendancePeriodRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
				mockAuditLogRepo.EXPECT().Create(gomock.Any()).Return(nil).AnyTimes()
			}

			repos := &repository.Repositories{
				AttendancePeriod: mockAttendancePeriodRepo,
				AuditLog:         mockAuditLogRepo,
			}

			s := NewAdminService(repos)
			got, err := s.CreateAttendancePeriod(tt.args.startDate, tt.args.endDate, tt.args.adminID, tt.args.ipAddress, tt.args.requestID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
