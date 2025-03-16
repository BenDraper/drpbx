package file_manager

import (
	mock_transfer "drpbx/file-monitor/transfer/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"time"

	"testing"
)

func TestFileManager_diff(t *testing.T) {

	tests := map[string]struct {
		oldFiles   map[string]time.Time
		entries    map[string]time.Time
		wantCreate []string
		wantDelete []string
		wantUpdate []string
	}{
		"One of each": {
			oldFiles: map[string]time.Time{
				"update": time.Now().Add(-10 * time.Second),
				"delete": time.Now().Add(-10 * time.Second),
			},
			entries: map[string]time.Time{
				"update": time.Now(),
				"create": time.Now(),
			},
			wantCreate: []string{"create"},
			wantDelete: []string{"delete"},
			wantUpdate: []string{"update"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			fm := &FileManager{
				oldFiles: tt.oldFiles,
			}
			gotCreate, gotUpdate, gotDelete := fm.diff(tt.entries)
			assert.Equal(t, tt.wantCreate, gotCreate)
			assert.Equal(t, tt.wantDelete, gotDelete)
			assert.Equal(t, tt.wantUpdate, gotUpdate)
		})
	}
}

func TestFileManager_sendDiffs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTransfer := mock_transfer.NewMockTransfer(mockCtrl)

	tests := map[string]struct {
		mockTransferOutcomes func(transferMock *mock_transfer.MockTransfer)
		folder               string
		creates              []string
		updates              []string
		deletes              []string
	}{
		"Send One Of Each": {
			folder: "folder/",
			mockTransferOutcomes: func(transferMock *mock_transfer.MockTransfer) {
				transferMock.EXPECT().Create("folder/create1").Times(1).Return(nil)
				transferMock.EXPECT().Update("folder/update1").Times(1).Return(nil)
				transferMock.EXPECT().Delete("delete1").Times(1).Return(nil)
			},
			creates: []string{"create1"},
			updates: []string{"update1"},
			deletes: []string{"delete1"},
		},
		"send others if one fails": {
			folder: "folder/",
			mockTransferOutcomes: func(transferMock *mock_transfer.MockTransfer) {
				transferMock.EXPECT().Create("folder/create1").Times(1).Return(fmt.Errorf("Error!"))
				transferMock.EXPECT().Update("folder/update1").Times(1).Return(nil)
				transferMock.EXPECT().Delete("delete1").Times(1).Return(nil)
			},
			creates: []string{"create1"},
			updates: []string{"update1"},
			deletes: []string{"delete1"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.mockTransferOutcomes(mockTransfer)

			fm := &FileManager{
				folder:   tt.folder,
				transfer: mockTransfer,
			}
			fm.sendDiffs(tt.creates, tt.updates, tt.deletes)
		})
	}
}
