package documentstore

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestNewStoreFromDump(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc1 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Alice"}}}
	doc2 := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "2"}, "age": {Type: DocumentFieldTypeNumber, Value: 30.0}}}
	col1 := &Collection{config: config, documents: map[string]Document{"1": doc1}}
	col2 := &Collection{config: config, documents: map[string]Document{"2": doc2}}
	want := &Store{collections: map[string]*Collection{"users": col1, "profiles": col2}}

	dumpData := map[string]interface{}{
		"collections": map[string]interface{}{
			"users": map[string]interface{}{
				"config":    map[string]interface{}{"primaryKey": "id"},
				"documents": map[string]interface{}{"1": map[string]interface{}{"fields": map[string]interface{}{"id": map[string]interface{}{"type": "string", "value": "1"}, "name": map[string]interface{}{"type": "string", "value": "Alice"}}}},
			},
			"profiles": map[string]interface{}{
				"config":    map[string]interface{}{"primaryKey": "id"},
				"documents": map[string]interface{}{"2": map[string]interface{}{"fields": map[string]interface{}{"id": map[string]interface{}{"type": "string", "value": "2"}, "age": map[string]interface{}{"type": "number", "value": 30}}}},
			},
		},
	}
	dumpBytes, _ := json.Marshal(dumpData)

	type args struct {
		dump []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		{
			name:    "Valid dump",
			args:    args{dump: dumpBytes},
			want:    want,
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			args:    args{dump: []byte("{invalid json")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty dump",
			args:    args{dump: []byte("{}")},
			want:    &Store{collections: map[string]*Collection{}},
			wantErr: false,
		},
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
	config := &CollectionConfig{PrimaryKey: "id"}
	doc := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "1"}, "name": {Type: DocumentFieldTypeString, Value: "Test"}}}
	col := &Collection{config: config, documents: map[string]Document{"1": doc}}
	want := &Store{collections: map[string]*Collection{"test_collection": col}}

	dumpData := map[string]interface{}{
		"collections": map[string]interface{}{
			"test_collection": map[string]interface{}{
				"config":    map[string]interface{}{"primaryKey": "id"},
				"documents": map[string]interface{}{"1": map[string]interface{}{"fields": map[string]interface{}{"id": map[string]interface{}{"type": "string", "value": "1"}, "name": map[string]interface{}{"type": "string", "value": "Test"}}}},
			},
		},
	}
	dumpBytes, _ := json.Marshal(dumpData)

	tmpFile, err := os.CreateTemp("", "test_store_")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(dumpBytes); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		{
			name:    "Valid file",
			args:    args{filename: tmpFile.Name()},
			want:    want,
			wantErr: false,
		},
		{
			name:    "Non-existent file",
			args:    args{filename: "non_existent.json"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid JSON file",
			args:    args{filename: createInvalidJSONFile(t)},
			want:    nil,
			wantErr: true,
		},
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

func createInvalidJSONFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "invalid_json_")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte("{invalid json")); err != nil {
		t.Fatalf("Failed to write invalid JSON to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}
	return tmpFile.Name()
}

func TestStore_Dump(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "test_id"}, "value": {Type: DocumentFieldTypeNumber, Value: 123.45}}}
	collections := map[string]*Collection{
		"my_collection": {
			config:    config,
			documents: map[string]Document{"test_id": doc},
		},
	}
	wantData := map[string]interface{}{
		"collections": map[string]interface{}{
			"my_collection": map[string]interface{}{
				"config":    map[string]interface{}{"primaryKey": "id"},
				"documents": map[string]interface{}{"test_id": map[string]interface{}{"fields": map[string]interface{}{"id": map[string]interface{}{"type": "string", "value": "test_id"}, "value": map[string]interface{}{"type": "number", "value": 123.45}}}},
			},
		},
	}
	want, _ := json.Marshal(wantData)

	type fields struct {
		collections map[string]*Collection
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "Single collection with document",
			fields:  fields{collections: collections},
			want:    want,
			wantErr: false,
		},
		{
			name:    "Empty store",
			fields:  fields{collections: map[string]*Collection{}},
			want:    []byte(`{"collections":{}}`),
			wantErr: false,
		},
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
				t.Errorf("Dump() got = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestStore_DumpToFile(t *testing.T) {
	config := &CollectionConfig{PrimaryKey: "id"}
	doc := Document{Fields: map[string]DocumentField{"id": {Type: DocumentFieldTypeString, Value: "test_id"}, "value": {Type: DocumentFieldTypeNumber, Value: 123.45}}}
	collections := map[string]*Collection{
		"my_collection": {
			config:    config,
			documents: map[string]Document{"test_id": doc},
		},
	}

	tmpFile, err := os.CreateTemp("", "test_store_dump_")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	filename := tmpFile.Name()

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
		{
			name:    "Valid dump to file",
			fields:  fields{collections: collections},
			args:    args{filename: filename},
			wantErr: false,
		},
		{
			name:    "Error creating file (simulated)",
			fields:  fields{collections: collections},
			args:    args{filename: "/invalid/path/test_dump.json"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{ // Оголошення всередині циклу - коректне
				collections: tt.fields.collections,
			}
			err := s.DumpToFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				content, err := os.ReadFile(tt.args.filename)
				if err != nil {
					t.Fatalf("Failed to read dumped file: %v", err)
				}
				var got map[string]interface{}
				if err := json.Unmarshal(content, &got); err != nil {
					t.Fatalf("Failed to unmarshal dumped content: %v", err)
				}
				wantData := map[string]interface{}{
					"collections": map[string]interface{}{
						"my_collection": map[string]interface{}{
							"config":    map[string]interface{}{"primaryKey": "id"},
							"documents": map[string]interface{}{"test_id": map[string]interface{}{"fields": map[string]interface{}{"id": map[string]interface{}{"type": "string", "value": "test_id"}, "value": map[string]interface{}{"type": "number", "value": 123.45}}}},
						},
					},
				}
				var want map[string]interface{}
				wantBytes, _ := json.Marshal(wantData)
				json.Unmarshal(wantBytes, &want)

				if !reflect.DeepEqual(got, want) {
					t.Errorf("DumpToFile() dumped content mismatch: got = %v, want = %v", got, want)
				}
			}
		})
	}
}
