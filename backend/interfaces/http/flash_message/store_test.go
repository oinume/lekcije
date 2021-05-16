package flash_message

import (
	"os"
	"reflect"
	"testing"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/model"
)

var (
	storeMySQL *StoreMySQL
)

func TestMain(m *testing.M) {
	helper := model.NewTestHelper()
	config.MustProcessDefault()
	db := helper.DB(nil)
	defer func() { _ = db.Close() }()
	storeMySQL = NewStoreMySQL(db)

	os.Exit(m.Run())
}

func TestStoreMySQL_Save_Load(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		want *FlashMessage
	}{
		"ok": {
			want: New(KindInfo, "データの削除に成功しました"),
		},
	}
	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if err := storeMySQL.Save(test.want); err != nil {
				t.Fatal(err)
			}
			got, err := storeMySQL.Load(test.want.Key)
			if err != nil {
				t.Fatalf("Load() failed: %got", err)
			}
			if !reflect.DeepEqual(test.want.Messages, got.Messages) {
				t.Fatalf("unexpected flash message: want=%+v, got=%+v", test.want, got)
			}
		})
	}
}
