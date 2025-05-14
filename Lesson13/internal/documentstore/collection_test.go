package documentstore

import (
	"reflect"
	"testing"
)

func TestCollection_Delete(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc1 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}}
	doc2 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "2"}, "age": {Type: DocumentFieldTypeNumber, Value: 30.0}}}

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
		{
			name: "Delete existing document",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1, "2": doc2},
			},
			args:    args{key: "1"},
			wantErr: false,
		},
		{
			name: "Delete non-existent document",
			fields: fields{
				config:    config,
				documents: map[string]Document{"2": doc2},
			},
			args:    args{key: "1"},
			wantErr: false,
		},
		{
			name: "Empty collection",
			fields: fields{
				config:    config,
				documents: map[string]Document{},
			},
			args:    args{key: "1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			initialLen := len(c.documents)
			err := c.Delete(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, exists := c.documents[tt.args.key] // Оголошення 'exists' тут
			if exists && !tt.wantErr {
				t.Errorf("Delete() document with key '%s' not deleted", tt.args.key)
			}
			if !exists && !tt.wantErr && len(c.documents) != initialLen-1 {
				t.Errorf("Delete() document count mismatch after deletion")
			}
		})
	}
}

func TestCollection_Get(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc1 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}}
	doc2 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "2"}, "age": {Type: DocumentFieldTypeNumber, Value: 30.0}}}

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
		{
			name: "Get existing document",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1, "2": doc2},
			},
			args:    args{key: "1"},
			want:    &doc1,
			wantErr: false,
		},
		{
			name: "Get non-existent document",
			fields: fields{
				config:    config,
				documents: map[string]Document{"2": doc2},
			},
			args:    args{key: "1"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty collection",
			fields: fields{
				config:    config,
				documents: map[string]Document{},
			},
			args:    args{key: "1"},
			want:    nil,
			wantErr: true,
		},
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

func TestCollection_List(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc1 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}}
	doc2 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "2"}, "age": {Type: DocumentFieldTypeNumber, Value: 30.0}}}

	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
	}
	tests := []struct {
		name   string
		fields fields
		want   []Document
	}{
		{
			name: "Empty collection",
			fields: fields{
				config:    config,
				documents: map[string]Document{},
			},
			want: []Document{},
		},
		{
			name: "Single document",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1},
			},
			want: []Document{doc1},
		},
		{
			name: "Multiple documents",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1, "2": doc2},
			},
			want: []Document{doc1, doc2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			got := c.List()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_NumDocuments(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc1 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}}
	doc2 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "2"}, "age": {Type: DocumentFieldTypeNumber, Value: 30.0}}}

	type fields struct {
		config    *CollectionConfig
		documents map[string]Document
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Empty collection",
			fields: fields{
				config:    config,
				documents: map[string]Document{},
			},
			want: 0,
		},
		{
			name: "Single document",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1},
			},
			want: 1,
		},
		{
			name: "Multiple documents",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1, "2": doc2},
			},
			want: 2,
		},
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
	config := &CollectionConfig{PrimaryKey: "id"}
	doc1 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}}
	doc2 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "2"}, "age": {Type: DocumentFieldTypeNumber, Value: 30.0}}}
	doc3 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Bob"}}} // Same ID as doc1

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
		{
			name: "Put new document",
			fields: fields{
				config:    config,
				documents: map[string]Document{},
			},
			args:    args{doc: doc1},
			wantErr: false,
		},
		{
			name: "Put existing document (update)",
			fields: fields{
				config:    config,
				documents: map[string]Document{"1": doc1, "2": doc2},
			},
			args:    args{doc: doc3},
			wantErr: false,
		},
		{
			name: "Put with missing primary key",
			fields: fields{
				config:    &CollectionConfig{PrimaryKey: "nonexistent_id"},
				documents: map[string]Document{},
			},
			args:    args{doc: doc1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collection{
				config:    tt.fields.config,
				documents: tt.fields.documents,
			}
			initialLen := len(c.documents)
			err := c.Put(tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, exists := c.documents[tt.args.doc.Fields[c.config.PrimaryKey].Value.(string)]; !exists {
					t.Errorf("Put() document not added")
				}
				if tt.name == "Put new document" && len(c.documents) != initialLen+1 {
					t.Errorf("Put() document count mismatch after adding")
				}
				if tt.name == "Put existing document (update)" && !reflect.DeepEqual(c.documents["1"], doc3) {
					t.Errorf("Put() document not updated correctly")
				}
			}
		})
	}
}
