package pager

import (
	"reflect"
	"testing"
)

func TestNewPageBuffer(t *testing.T) {
	type args struct {
		pm *PageManager
	}
	tests := []struct {
		name    string
		args    args
		want    *PageBuffer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageBuffer(tt.args.pm)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPageBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPageBuffer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPageBufferSize(t *testing.T) {
	type args struct {
		pm *PageManager
		np int
	}
	tests := []struct {
		name    string
		args    args
		want    *PageBuffer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageBufferSize(tt.args.pm, tt.args.np)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPageBufferSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPageBufferSize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageBuffer_AddRecord(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	type args struct {
		r []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RecordID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			got, err := pb.AddRecord(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageBuffer_Close(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			if err := pb.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPageBuffer_DelRecord(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	type args struct {
		rid *RecordID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			if err := pb.DelRecord(tt.args.rid); (err != nil) != tt.wantErr {
				t.Errorf("DelRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPageBuffer_DirtyPages(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			if got := pb.DirtyPages(); got != tt.want {
				t.Errorf("DirtyPages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageBuffer_Flush(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			if err := pb.Flush(); (err != nil) != tt.wantErr {
				t.Errorf("Flush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPageBuffer_FreeSpace(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			if got := pb.FreeSpace(); got != tt.want {
				t.Errorf("FreeSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageBuffer_GetRecord(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	type args struct {
		rid *RecordID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			got, err := pb.GetRecord(tt.args.rid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageBuffer_Range(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	type args struct {
		fn func(rid *RecordID) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			pb.Range(tt.args.fn)
		})
	}
}

func TestPageBuffer_load(t *testing.T) {
	type fields struct {
		manager *PageManager
		buffer  []*Page
		metas   []pageMeta
		pinned  int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBuffer{
				manager: tt.fields.manager,
				buffer:  tt.fields.buffer,
				metas:   tt.fields.metas,
				pinned:  tt.fields.pinned,
			}
			if err := pb.load(); (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
