package pager

import (
	"os"
	"reflect"
	"testing"
)

func TestOpenPageManager(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *PageManager
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenPageManager(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenPageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenPageManager() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageManager_AllocatePage(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	tests := []struct {
		name   string
		fields fields
		want   *Page
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if got := f.AllocatePage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllocatePage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageManager_Close(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if err := f.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPageManager_DeletePage(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	type args struct {
		pid uint32
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if err := f.DeletePage(tt.args.pid); (err != nil) != tt.wantErr {
				t.Errorf("DeletePage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPageManager_GetFreeOrAllocate(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	tests := []struct {
		name   string
		fields fields
		want   *Page
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if got := f.GetFreeOrAllocate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFreeOrAllocate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageManager_GetFreePageIDs(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if got := f.GetFreePageIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFreePageIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageManager_PageCount(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if got := f.PageCount(); got != tt.want {
				t.Errorf("PageCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageManager_Range(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	type args struct {
		start uint32
		fn    func(rid *RecordID) bool
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			f.Range(tt.args.start, tt.args.fn)
		})
	}
}

func TestPageManager_ReadPage(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	type args struct {
		pid uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Page
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			got, err := f.ReadPage(tt.args.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadPage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageManager_ReadPages(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	type args struct {
		pid uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Page
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			got, err := f.ReadPages(tt.args.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadPages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPageManager_WritePage(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	type args struct {
		p *Page
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if err := f.WritePage(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("WritePage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPageManager_WritePages(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	type args struct {
		ps []*Page
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if err := f.WritePages(tt.args.ps); (err != nil) != tt.wantErr {
				t.Errorf("WritePages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPageManager_grow(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
	}
	type args struct {
		sizeToGrow int64
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			f.grow(tt.args.sizeToGrow)
		})
	}
}

func TestPageManager_load(t *testing.T) {
	type fields struct {
		name        string
		fp          *os.File
		pageHeaders []*pageHeader
		pageCache   *Page
		freePages   int
		pids        *autoPageID
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
			f := &PageManager{
				name:        tt.fields.name,
				fp:          tt.fields.fp,
				pageHeaders: tt.fields.pageHeaders,
				pageCache:   tt.fields.pageCache,
				freePages:   tt.fields.freePages,
				pids:        tt.fields.pids,
			}
			if err := f.load(); (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Remove(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getPagePosition(t *testing.T) {
	type args struct {
		pid uint32
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPagePosition(tt.args.pid); got != tt.want {
				t.Errorf("getPagePosition() = %v, want %v", got, tt.want)
			}
		})
	}
}
