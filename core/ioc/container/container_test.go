package container

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test1 struct{}
type test2 struct{}

var (
	pool1 = &sync.Pool{New: func() interface{} { return &test1{} }}
	pool2 = &sync.Pool{New: func() interface{} { return &test2{} }}
)

func TestNewContainer(t *testing.T) {
	c := Config{
		AutoFree:  true,
		AutoWired: false,
	}
	tests := []struct {
		name string
		c    []Config
		want *Container
	}{
		// TODO: Add test cases.
		{
			name: "1",
			c:    []Config{},
			want: &Container{
				c:               DefaultConfig,
				poolMap:         make(map[reflect.Type]*sync.Pool),
				poolTags:        make(map[string]reflect.Type),
				instanceMapping: make(map[string]reflect.Type),
			},
		},
		{
			name: "2",
			c:    []Config{c},
			want: &Container{
				c:               &c,
				poolMap:         make(map[reflect.Type]*sync.Pool),
				poolTags:        make(map[string]reflect.Type),
				instanceMapping: make(map[string]reflect.Type),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContainer(tt.c...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_newInstance(t *testing.T) {

	type fields struct {
		poolMap          map[reflect.Type]*sync.Pool
		instanceTypeList []reflect.Type
		poolTags         map[string]reflect.Type
		instanceMapping  map[string]reflect.Type
	}
	type args struct {
		instanceType reflect.Type
		instancePool *sync.Pool
		instanceTag  []string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantPanic bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			fields: fields{
				poolMap: map[reflect.Type]*sync.Pool{
					reflect.TypeOf(&test1{}): pool1,
				},
				instanceTypeList: []reflect.Type{reflect.TypeOf(&test1{})},
				poolTags: map[string]reflect.Type{
					"test1": reflect.TypeOf(&test1{}),
				},
				instanceMapping: map[string]reflect.Type{},
			},
			args: args{
				instanceType: reflect.TypeOf(&test1{}),
				instancePool: &sync.Pool{New: func() interface{} { return &test1{} }},
				instanceTag:  []string{"test1"},
			},
			wantPanic: false,
		},
		{
			name: "2",
			fields: fields{
				poolMap:          map[reflect.Type]*sync.Pool{},
				instanceTypeList: []reflect.Type{},
				poolTags:         map[string]reflect.Type{},
				instanceMapping:  map[string]reflect.Type{},
			},
			args: args{
				instanceType: reflect.TypeOf(&test1{}),
				instancePool: pool1,
				instanceTag:  []string{"test1"},
			},
			wantPanic: false,
		},
		{
			name: "3",
			fields: fields{
				poolMap: map[reflect.Type]*sync.Pool{
					reflect.TypeOf(&test1{}): pool1,
				},
				instanceTypeList: []reflect.Type{reflect.TypeOf(&test1{})},
				poolTags: map[string]reflect.Type{
					"test1": reflect.TypeOf(&test1{}),
				},
				instanceMapping: map[string]reflect.Type{},
			},
			args: args{
				instanceType: reflect.TypeOf(&test2{}),
				instancePool: &sync.Pool{New: func() interface{} { return &test1{} }},
				instanceTag:  []string{"test1"},
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Container{
				poolMap:          tt.fields.poolMap,
				instanceTypeList: tt.fields.instanceTypeList,
				poolTags:         tt.fields.poolTags,
				instanceMapping:  tt.fields.instanceMapping,
			}
			if tt.wantPanic {
				assert.PanicsWithError(t, "tag test1 already existed", func() { s.newInstance(tt.args.instanceType, tt.args.instancePool, tt.args.instanceTag...) })
			} else {
				s.newInstance(tt.args.instanceType, tt.args.instancePool, tt.args.instanceTag...)
				assert.Equal(t, &Container{
					poolMap: map[reflect.Type]*sync.Pool{
						reflect.TypeOf(&test1{}): pool1,
					},
					instanceTypeList: []reflect.Type{reflect.TypeOf(&test1{})},
					poolTags: map[string]reflect.Type{
						"test1": reflect.TypeOf(&test1{}),
					},
					instanceMapping: map[string]reflect.Type{},
				}, s, "instance not register correctly")
			}

		})
	}
}

func TestContainer_GetInstanceTypeByTag(t *testing.T) {
	p1 := NewContainer()
	p1.newInstance(reflect.TypeOf(&test1{}), pool1, "test1")
	p1.newInstance(reflect.TypeOf(&test1{}), pool1, "test2")
	p1.newInstance(reflect.TypeOf(&test2{}), pool2, "test3")
	p1.newInstance(reflect.TypeOf(&test2{}), pool2, "test4")
	type fields struct {
		poolMap          map[reflect.Type]*sync.Pool
		instanceTypeList []reflect.Type
		poolTags         map[string]reflect.Type
		instanceMapping  map[string]reflect.Type
	}
	type args struct {
		tags []string
	}
	tests := []struct {
		name   string
		fields *Container
		args   args
		want   []reflect.Type
	}{
		// TODO: Add test cases.
		{
			name:   "1",
			fields: p1,
			args: args{
				tags: []string{"test1"},
			},
			want: []reflect.Type{reflect.TypeOf(&test1{})},
		},
		{
			name:   "2",
			fields: p1,
			args: args{
				tags: []string{},
			},
			want: []reflect.Type{reflect.TypeOf(&test1{}), reflect.TypeOf(&test2{})},
		},
		{
			name:   "3",
			fields: p1,
			args: args{
				tags: []string{"test1", "test3"},
			},
			want: []reflect.Type{reflect.TypeOf(&test1{}), reflect.TypeOf(&test2{})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.GetInstanceTypeByTag(tt.args.tags...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Container.GetInstanceTypeByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_CheckInstanceNameIfExist(t *testing.T) {
	p1 := NewContainer()
	p1.newInstance(reflect.TypeOf(&test1{}), pool1, "test1")
	type args struct {
		instanceName reflect.Type
	}
	tests := []struct {
		name   string
		fields *Container
		args   args
		want   bool
	}{
		{
			name:   "1",
			fields: p1,
			args: args{
				instanceName: reflect.TypeOf(&test1{}),
			},
			want: true,
		},
		{
			name:   "2",
			fields: p1,
			args: args{
				instanceName: reflect.TypeOf(&test2{}),
			},
			want: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Container{
				poolMap:          tt.fields.poolMap,
				instanceTypeList: tt.fields.instanceTypeList,
				poolTags:         tt.fields.poolTags,
				instanceMapping:  tt.fields.instanceMapping,
			}
			if got := s.CheckInstanceNameIfExist(tt.args.instanceName); got != tt.want {
				t.Errorf("Container.CheckInstanceNameIfExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_getAutoWireTag(t *testing.T) {
	type test3 struct {
		test1 string `container:"autowire:true"`
		test2 string `container:"autowire:false"`
		test3 string `container:"autowire"`
		test4 string `container:""`
	}
	type args struct {
		obj   interface{}
		index int
	}
	tests := []struct {
		name   string
		fields *Container
		args   args
		want   bool
	}{
		{
			name:   "1",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 0,
			},
			want: true,
		},
		{
			name:   "2",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 1,
			},
			want: false,
		},
		{
			name:   "3",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 2,
			},
			want: true,
		},
		{
			name:   "4",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 3,
			},
			want: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.getAutoWireTag(tt.args.obj, tt.args.index); got != tt.want {
				t.Errorf("Container.getAutoWireTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_getResourceTag(t *testing.T) {
	type test3 struct {
		test1 string `container:"resource:1234"`
		test2 string `container:"resource:asdf"`
		test3 string `container:"resource:12fs"`
		test4 string `container:"resource:213"`
	}
	type args struct {
		obj   interface{}
		index int
	}
	tests := []struct {
		name   string
		fields *Container
		args   args
		want   string
	}{
		{
			name:   "1",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 0,
			},
			want: "1234",
		},
		{
			name:   "2",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 1,
			},
			want: "asdf",
		},
		{
			name:   "3",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 2,
			},
			want: "12fs",
		},
		{
			name:   "4",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 3,
			},
			want: "213",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.getResourceTag(tt.args.obj, tt.args.index); got != tt.want {
				t.Errorf("Container.getAutoFreeTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_getAutoFreeTag(t *testing.T) {
	type test3 struct {
		test1 string `container:"autofree:true"`
		test2 string `container:"autofree:false"`
		test3 string `container:"autofree"`
		test4 string `container:""`
	}
	type args struct {
		obj   interface{}
		index int
	}
	tests := []struct {
		name   string
		fields *Container
		args   args
		want   bool
	}{
		{
			name:   "1",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 0,
			},
			want: true,
		},
		{
			name:   "2",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 1,
			},
			want: false,
		},
		{
			name:   "3",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 2,
			},
			want: true,
		},
		{
			name:   "4",
			fields: NewContainer(),
			args: args{
				obj:   &test3{},
				index: 3,
			},
			want: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.getAutoFreeTag(tt.args.obj, tt.args.index); got != tt.want {
				t.Errorf("Container.getAutoFreeTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeTag(t *testing.T) {
	type args struct {
		value string
		key   Keyword
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantBool bool
	}{
		{
			name: "1",
			args: args{
				value: "autofree:true",
				key:   AUTOFREE,
			},
			want:     "true",
			wantBool: true,
		},
		{
			name: "2",
			args: args{
				value: "autofree:true;autofree:false;",
				key:   AUTOFREE,
			},
			want:     "false",
			wantBool: true,
		},
		{
			name: "3",
			args: args{
				value: "autofree;",
				key:   AUTOFREE,
			},
			want:     "",
			wantBool: true,
		},
		{
			name: "4",
			args: args{
				value: ":;",
				key:   AUTOFREE,
			},
			want:     "",
			wantBool: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := decodeTag(tt.args.value, tt.args.key)
			if got != tt.want {
				t.Errorf("decodeTag() value = %v, want %v", got, tt.want)
			}
			if ok != tt.wantBool {
				t.Errorf("decodeTag() isExist= %v, want %v", ok, tt.wantBool)
			}
		})
	}
}

func TestContainer_RegisterInstance(t *testing.T) {
	type fields struct {
		c *Container
	}
	type args struct {
		instance    interface{}
		instanceTag []string
	}
	type want struct {
		typeList []reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want
	}{
		{
			name: "1",
			fields: fields{
				c: NewContainer(Config{
					AutoFree:  true,
					AutoWired: false,
				}),
			},
			args: args{
				instance:    test1{},
				instanceTag: []string{"test1"},
			},
			want: want{
				typeList: []reflect.Type{reflect.TypeOf(test1{})},
			},
		},
		{
			name: "2",
			fields: fields{
				c: NewContainer(Config{
					AutoFree:  true,
					AutoWired: false,
				}),
			},
			args: args{
				instance:    &test1{},
				instanceTag: []string{"test1"},
			},
			want: want{
				typeList: []reflect.Type{reflect.TypeOf(test1{})},
			},
		},
		{
			name: "3",
			fields: fields{
				c: NewContainer(Config{
					AutoFree:  true,
					AutoWired: false,
				}),
			},
			args: args{
				instance:    func() interface{} { return test1{} },
				instanceTag: []string{"test1"},
			},
			want: want{
				typeList: []reflect.Type{reflect.TypeOf(test1{})},
			},
		},
		{
			name: "4",
			fields: fields{
				c: NewContainer(Config{
					AutoFree:  true,
					AutoWired: false,
				}),
			},
			args: args{
				instance:    func() interface{} { return "123" },
				instanceTag: []string{"test1"},
			},
			want: want{
				typeList: []reflect.Type{reflect.TypeOf("123")},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.c.RegisterInstance(tt.args.instance, tt.args.instanceTag...)
			assert.Equal(t, tt.want.typeList, tt.fields.c.instanceTypeList, "wrong type ")
		})
	}
}
