package documentstore

import (
	"reflect"
	"testing"
)

func TestCollection_GetAll(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc1 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}}
	doc2 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "2"}, "age": {Type: DocumentFieldTypeNumber, Value: 30.0}}}

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
		{
			name: "Empty collection",
			fields: fields{
				config:    config,
				documents: map[string]Document{},
			},
			want:    []Document{},
			wantErr: true,
		},
		{
			name: "Single document",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1},
			},
			want:    []Document{doc1},
			wantErr: false,
		},
		{
			name: "Multiple documents",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1, "2": doc2},
			},
			want:    []Document{doc1, doc2},
			wantErr: false,
		},
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
		{
			name: "Valid struct",
			args: args{input: struct {
				ID   string  `json:"id"`
				Name string  `json:"name"`
				Age  float64 `json:"age"`
			}{ID: "1", Name: "Alice", Age: 30}},
			want: &Document{Fields: map[string]DocumentField{
				"id":   {Type: DocumentFieldTypeString, Value: "1"},
				"name": {Type: DocumentFieldTypeString, Value: "Alice"},
				"age":  {Type: DocumentFieldTypeNumber, Value: 30.0},
			}},
			wantErr: false,
		},
		{
			name:    "Unsupported type",
			args:    args{input: map[int]string{1: "one"}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Nested object",
			args: args{input: struct {
				ID   string `json:"id"`
				Info struct {
					City    string `json:"city"`
					Country string `json:"country"`
				} `json:"info"`
			}{ID: "1", Info: struct {
				City    string `json:"city"`
				Country string `json:"country"`
			}{City: "Kharkiv", Country: "Ukraine"}}},
			want: &Document{Fields: map[string]DocumentField{
				"id": {Type: DocumentFieldTypeString, Value: "1"},
				"info": {Type: DocumentFieldTypeObject, Value: map[string]any{
					"city":    "Kharkiv",
					"country": "Ukraine",
				}},
			}},
			wantErr: false,
		},
		{
			name: "Array of strings",
			args: args{input: struct {
				ID      string   `json:"id"`
				Hobbies []string `json:"hobbies"`
			}{ID: "1", Hobbies: []string{"reading", "hiking"}}},
			want: &Document{Fields: map[string]DocumentField{
				"id":      {Type: DocumentFieldTypeString, Value: "1"},
				"hobbies": {Type: DocumentFieldTypeArray, Value: []any{"reading", "hiking"}},
			}},
			wantErr: false,
		},
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
		{
			name: "Create new store",
			want: &Store{collections: make(map[string]*Collection)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_CreateCollection(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
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
		{
			name: "Create new collection",
			fields: fields{
				collections: make(map[string]*Collection),
			},
			args: args{
				name: "users",
				cfg:  config,
			},
			wantErr: false,
		},
		{
			name: "Collection already exists",
			fields: fields{
				collections: map[string]*Collection{"users": {config: config, documents: make(map[string]Document)}},
			},
			args: args{
				name: "users",
				cfg:  config,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.collections,
			}
			if err := s.CreateCollection(tt.args.name, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("CreateCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, exists := s.collections[tt.args.name]; !exists {
					t.Errorf("CreateCollection() collection '%s' not created", tt.args.name)
				}
			}
		})
	}
}

func TestStore_DeleteCollection(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
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
		{
			name: "Delete existing collection",
			fields: fields{
				collections: map[string]*Collection{"users": {config: config, documents: make(map[string]Document)}},
			},
			args:    args{name: "users"},
			wantErr: false,
		},
		{
			name: "Delete non-existent collection",
			fields: fields{
				collections: make(map[string]*Collection),
			},
			args:    args{name: "users"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.collections,
			}
			err := s.DeleteCollection(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, exists := s.collections[tt.args.name]; exists {
					t.Errorf("DeleteCollection() collection '%s' not deleted", tt.args.name)
				}
			}
		})
	}
}

func TestStore_GetCollection(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	collection := &Collection{config: config, documents: make(map[string]Document)}
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
		{
			name: "Get existing collection",
			fields: fields{
				collections: map[string]*Collection{"users": collection},
			},
			args:    args{name: "users"},
			want:    collection,
			wantErr: false,
		},
		{
			name: "Get non-existent collection",
			fields: fields{
				collections: make(map[string]*Collection),
			},
			args:    args{name: "users"},
			want:    nil,
			wantErr: true,
		},
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
	config := &CollectionConfig{PrimaryKey: "id"}
	type fields struct {
		collections map[string]*Collection
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Empty store",
			fields: fields{
				collections: make(map[string]*Collection),
			},
			want: 0,
		},
		{
			name: "Single collection",
			fields: fields{
				collections: map[string]*Collection{"users": {config: config, documents: make(map[string]Document)}},
			},
			want: 1,
		},
		{
			name: "Multiple collections",
			fields: fields{
				collections: map[string]*Collection{
					"users":    {config: config, documents: make(map[string]Document)},
					"profiles": {config: config, documents: make(map[string]Document)},
				},
			},
			want: 2,
		},
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
		name       string
		args       args
		wantErr    bool
		wantOutput any
	}{
		{
			name: "Valid document",
			args: args{
				doc: &Document{Fields: map[string]DocumentField{
					"data": {Type: DocumentFieldTypeString, Value: `{"id": "1", "name": "Alice"}`},
				}},
				output: &struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{},
			},
			wantErr: false,
			wantOutput: &struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{ID: "1", Name: "Alice"},
		},
		{
			name: "Missing data field",
			args: args{
				doc: &Document{Fields: map[string]DocumentField{
					"other": {Type: DocumentFieldTypeString, Value: `{"id": "1"}`},
				}},
				output: &struct {
					ID string `json:"id"`
				}{},
			},
			wantErr: true,
		},
		{
			name: "Invalid data field type",
			args: args{
				doc: &Document{Fields: map[string]DocumentField{
					"data": {Type: DocumentFieldTypeNumber, Value: 123},
				}},
				output: &struct {
					ID string `json:"id"`
				}{},
			},
			wantErr: true,
			wantOutput: &struct {
				ID string `json:"id"`
			}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalDocument(tt.args.doc, tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.output, tt.wantOutput) {
				t.Errorf("UnmarshalDocument() output = %v, want %v", tt.args.output, tt.wantOutput)
			}
		})
	}
}
