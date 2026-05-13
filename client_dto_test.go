package dnssdk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestRecordTypeMX_ToContent(t *testing.T) {
	tests := []struct {
		name string
		mx   RecordTypeMX
		want []interface{}
	}{
		{
			name: "ok",
			mx:   "10 mail.server",
			want: []interface{}{int64(10), "mail.server"},
		},
		{
			name: "wrong",
			mx:   "10",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mx.ToContent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordTypeCAA_ToContent(t *testing.T) {
	tests := []struct {
		name string
		caa  RecordTypeCAA
		want []interface{}
	}{
		{
			name: "ok",
			caa:  "10 issue aaa",
			want: []interface{}{int64(10), "issue", "aaa"},
		},
		{
			name: "wrong",
			caa:  "10 aa",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.caa.ToContent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordTypeAny_ToContent(t *testing.T) {
	tests := []struct {
		name string
		any  RecordTypeAny
		want []interface{}
	}{
		{
			name: "ok",
			any:  "any any 34",
			want: []interface{}{"any any 34"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.any.ToContent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentFromValue(t *testing.T) {
	type args struct {
		recordType string
		content    string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "mx",
			args: args{
				recordType: "MX",
				content:    "10 mx.com",
			},
			want: []interface{}{int64(10), "mx.com"},
		},
		{
			name: "caa",
			args: args{
				recordType: "CAA",
				content:    "10 issue com",
			},
			want: []interface{}{int64(10), "issue", "com"},
		},
		{
			name: "any",
			args: args{
				recordType: "A",
				content:    "10 issue com",
			},
			want: []interface{}{"10 issue com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContentFromValue(tt.args.recordType, tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContentFromValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaIP(t *testing.T) {
	type args struct {
		ips []string
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "one",
			args: args{
				ips: []string{"1.1.1.1"},
			},
			want: ResourceMeta{
				name:     "ip",
				value:    []string{"1.1.1.1"},
				validErr: nil,
			},
		},
		{
			name: "many",
			args: args{
				ips: []string{"1.1.1.1", "1.1.1.2"},
			},
			want: ResourceMeta{
				name:     "ip",
				value:    []string{"1.1.1.1", "1.1.1.2"},
				validErr: nil,
			},
		},
		{
			name: "wrong",
			args: args{
				ips: []string{"1.1.1.1", "sadasd"},
			},
			want: ResourceMeta{
				name:     "",
				value:    nil,
				validErr: fmt.Errorf("wrong ip"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaIP(tt.args.ips...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaAsn(t *testing.T) {
	type args struct {
		asn []uint64
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "ok",
			args: args{
				asn: []uint64{1, 2},
			},
			want: ResourceMeta{
				name:     "asn",
				value:    []uint64{1, 2},
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaAsn(tt.args.asn...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaAsn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaLatLong(t *testing.T) {
	type args struct {
		latlong string
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "ok",
			args: args{
				latlong: "1,2",
			},
			want: ResourceMeta{
				name:     "latlong",
				value:    []float64{1, 2},
				validErr: nil,
			},
		},
		{
			name: "ok ()",
			args: args{
				latlong: "(1,2)",
			},
			want: ResourceMeta{
				name:     "latlong",
				value:    []float64{1, 2},
				validErr: nil,
			},
		},
		{
			name: "ok []",
			args: args{
				latlong: "[1,2]",
			},
			want: ResourceMeta{
				name:     "latlong",
				value:    []float64{1, 2},
				validErr: nil,
			},
		},
		{
			name: "ok {}",
			args: args{
				latlong: "{1,2}",
			},
			want: ResourceMeta{
				name:     "latlong",
				value:    []float64{1, 2},
				validErr: nil,
			},
		},
		{
			name: "invalid count",
			args: args{
				latlong: "1,2,3",
			},
			want: ResourceMeta{
				name:     "",
				value:    nil,
				validErr: fmt.Errorf("latlong invalid format"),
			},
		},
		{
			name: "invalid lat",
			args: args{
				latlong: "a,2",
			},
			want: ResourceMeta{
				name:  "",
				value: nil,
				validErr: fmt.Errorf("lat is invalid: %w",
					&strconv.NumError{Func: "ParseFloat", Num: "a", Err: strconv.ErrSyntax}),
			},
		},
		{
			name: "invalid long",
			args: args{
				latlong: "1,a",
			},
			want: ResourceMeta{
				name:  "",
				value: nil,
				validErr: fmt.Errorf("long is invalid: %w",
					&strconv.NumError{Func: "ParseFloat", Num: "a", Err: strconv.ErrSyntax}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaLatLong(tt.args.latlong); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaLatLong() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaNotes(t *testing.T) {
	type args struct {
		notes []string
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "ok",
			args: args{
				notes: []string{"a"},
			},
			want: ResourceMeta{
				name:     "notes",
				value:    []string{"a"},
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaNotes(tt.args.notes...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaCountries(t *testing.T) {
	type args struct {
		countries []string
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "",
			args: args{
				countries: []string{"a"},
			},
			want: ResourceMeta{
				name:     "countries",
				value:    []string{"a"},
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaCountries(tt.args.countries...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaCountries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaContinents(t *testing.T) {
	type args struct {
		continents []string
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "ok",
			args: args{
				continents: []string{"a"},
			},
			want: ResourceMeta{
				name:     "continents",
				value:    []string{"a"},
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaContinents(tt.args.continents...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaContinents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaDefault(t *testing.T) {
	tests := []struct {
		name string
		def  bool
		want ResourceMeta
	}{
		{
			name: "ok",
			def:  true,
			want: ResourceMeta{
				name:     "default",
				value:    true,
				validErr: nil,
			},
		},
		{
			name: "false",
			def:  false,
			want: ResourceMeta{
				name:     "default",
				value:    false,
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaDefault(tt.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaBackup(t *testing.T) {
	tests := []struct {
		name string
		want ResourceMeta
	}{
		{
			name: "ok",
			want: ResourceMeta{
				name:     "backup",
				value:    true,
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaBackup(true); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaBackup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResourceMetaWeight(t *testing.T) {
	type args struct {
		weight int
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "ok",
			args: args{
				weight: 6,
			},
			want: ResourceMeta{
				name:     "weight",
				value:    6,
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaWeight(tt.args.weight); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaWeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceRecords_AddMeta(t *testing.T) {
	type fields struct {
		Content []interface{}
		Meta    map[string]interface{}
	}
	type args struct {
		meta ResourceMeta
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ResourceRecord
	}{
		{
			name: "",
			fields: fields{
				Meta: nil,
			},
			args: args{
				meta: ResourceMeta{
					name:  "a",
					value: 1,
				},
			},
			want: ResourceRecord{
				Meta: map[string]interface{}{"a": 1},
			},
		},
		{
			name: "backup and weight helpers",
			fields: fields{
				Meta: nil,
			},
			args: args{
				meta: NewResourceMetaBackup(true),
			},
			want: ResourceRecord{
				Meta: map[string]interface{}{"backup": true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceRecord{
				Content: tt.fields.Content,
				Meta:    tt.fields.Meta,
			}
			if got := r.AddMeta(tt.args.meta); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("AddMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceRecords_AddMetaWeight(t *testing.T) {
	r := (&ResourceRecord{}).AddMeta(NewResourceMetaWeight(6))
	want := ResourceRecord{
		Meta: map[string]interface{}{"weight": 6},
	}
	if !reflect.DeepEqual(*r, want) {
		t.Errorf("AddMeta(NewResourceMetaWeight()) = %v, want %v", r, want)
	}
}

func TestResourceRecordMetaJSONMarshal(t *testing.T) {
	record := (&ResourceRecord{Enabled: true}).
		AddMeta(NewResourceMetaBackup(true)).
		AddMeta(NewResourceMetaWeight(6))

	got, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"content":null,"meta":{"backup":true,"weight":6},"enabled":true}`
	if string(got) != want {
		t.Fatalf("json.Marshal() = %s, want %s", string(got), want)
	}
}

func TestResourceRecordMetaJSONUnmarshal(t *testing.T) {
	raw := []byte(`{"content":["1.5.7.3"],"meta":{"backup":true,"weight":6},"enabled":true}`)

	var got ResourceRecord
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Meta["backup"] != true {
		t.Fatalf("meta.backup = %v, want true", got.Meta["backup"])
	}
	if got.Meta["weight"] != float64(6) {
		t.Fatalf("meta.weight = %v, want %v", got.Meta["weight"], float64(6))
	}
}

func TestNewResourceMetaRegions(t *testing.T) {
	type args struct {
		regions []string
	}
	tests := []struct {
		name string
		args args
		want ResourceMeta
	}{
		{
			name: "ok",
			args: args{
				regions: []string{"ru-spb", "ru-mow"},
			},
			want: ResourceMeta{
				name:     "regions",
				value:    []string{"ru-spb", "ru-mow"},
				validErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceMetaRegions(tt.args.regions...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceMetaRegions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceRecords_AddMetaRegions(t *testing.T) {
	r := (&ResourceRecord{}).AddMeta(NewResourceMetaRegions("ru-spb", "ru-mow"))
	want := ResourceRecord{
		Meta: map[string]interface{}{"regions": []string{"ru-spb", "ru-mow"}},
	}
	if !reflect.DeepEqual(*r, want) {
		t.Errorf("AddMeta(NewResourceMetaRegions()) = %v, want %v", r, want)
	}
}

func TestResourceRecordMetaRegionsJSONMarshal(t *testing.T) {
	record := (&ResourceRecord{Enabled: true}).
		AddMeta(NewResourceMetaRegions("ru-spb", "ru-mow"))

	got, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"content":null,"meta":{"regions":["ru-spb","ru-mow"]},"enabled":true}`
	if string(got) != want {
		t.Fatalf("json.Marshal() = %s, want %s", string(got), want)
	}
}

func TestResourceRecordMetaRegionsJSONUnmarshal(t *testing.T) {
	raw := []byte(`{"content":["1.5.7.3"],"meta":{"regions":["ru-spb","ru-mow"]},"enabled":true}`)

	var got ResourceRecord
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(got.Meta["regions"], []interface{}{"ru-spb", "ru-mow"}) {
		t.Fatalf("meta.regions = %v, want %v", got.Meta["regions"], []interface{}{"ru-spb", "ru-mow"})
	}
}
