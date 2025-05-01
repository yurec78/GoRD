package documentstore

import (
	"reflect"
	"testing"
)

func TestCollection_Delete(t *testing.T) {
	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
	}
	type args struct {
		key string
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
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			if err := c.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollection_Get(t *testing.T) {
	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_GetAll(t *testing.T) {
	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			got, err := c.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_List(t *testing.T) {
	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
	}
	tests := []struct {
		name   string
		fields fields
		want   []Document
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			if got := c.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_NumDocuments(t *testing.T) {
	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
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
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			if got := c.NumDocuments(); got != tt.want {
				t.Errorf("NumDocuments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_Put(t *testing.T) {
	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
	}
	type args struct {
		doc Document
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
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			if err := c.Put(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarshalDocument(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    *Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalDocument(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalDocument() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStore(t *testing.T) {
	tests := []struct {
		name string
		want *Store
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStoreFromDump(t *testing.T) {
	type args struct {
		dump []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStoreFromDump(tt.args.dump)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStoreFromDump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreFromDump() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStoreFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStoreFromFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStoreFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_CreateCollection(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	type args struct {
		name string
		cfg  *CollectionConfig
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
			s := &Store{
				collections: tt.fields.collections,
			}
			if err := s.CreateCollection(tt.args.name, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("CreateCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_DeleteCollection(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	type args struct {
		name string
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
			s := &Store{
				collections: tt.fields.collections,
			}
			if err := s.DeleteCollection(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_Dump(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.collections,
			}
			got, err := s.Dump()
			if (err != nil) != tt.wantErr {
				t.Errorf("Dump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dump() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_DumpToFile(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	type args struct {
		filename string
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
			s := &Store{
				collections: tt.fields.collections,
			}
			if err := s.DumpToFile(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("DumpToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_GetCollection(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.collections,
			}
			got, err := s.GetCollection(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCollection() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_NumCollections(t *testing.T) {
	type fields struct {
		collections map[string]*Collection
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
			s := &Store{
				collections: tt.fields.collections,
			}
			if got := s.NumCollections(); got != tt.want {
				t.Errorf("NumCollections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalDocument(t *testing.T) {
	type args struct {
		doc    *Document
		output any
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
			if err := UnmarshalDocument(tt.args.doc, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_toDocumentField(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		args    args
		want    DocumentField
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toDocumentField(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDocumentField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDocumentField() got = %v, want %v", got, tt.want)
			}
		})
	}
}
