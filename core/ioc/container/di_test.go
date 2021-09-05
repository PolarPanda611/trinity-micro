package container

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type inject1 struct {
	name string
}
type injectInterface interface {
	A()
	B()
	C()
}
type injectInterfaceImpl struct{}

func (i injectInterfaceImpl) A() {}
func (i injectInterfaceImpl) B() {}
func (i injectInterfaceImpl) C() {}

type testInjectErr1 struct {
	test1 inject1 `container:"autowire:true"`
}
type testInjectErr2 struct {
	test1 inject1 `container:"autowire:true;resource:inject1"`
}
type testInject3 struct {
	Test1 inject1 `container:"autowire:true;resource:inject1"`
}
type testInject4 struct {
	Test1 *inject1 `container:"autowire:true;resource:inject1"`
}

type testInject5 struct {
	Test1 injectInterface `container:"autowire:true;resource:inject1"`
}

type testInject6 struct {
	Test1 interface{} `container:"autowire:true;resource:inject1"`
}

type testShared1 struct {
	T *testShared2 `container:"autowire:true;resource:shared2"`
}

type testShared2 struct {
	T *testShared1 `container:"autowire:true;resource:shared1"`
}

func TestContainer_DiFree(t *testing.T) {
	a := inject1{"1"}
	empty := inject1{}
	{
		obj := struct{}{}
		DiFree(logger, &obj)
		assert.Equal(t, struct{}{}, obj, "di free successfully")
	}
	{
		obj := struct {
			test1 inject1 `container:"autowire:true"`
		}{a}
		DiFree(logger, &obj)
		assert.Equal(t, a, obj.test1, "diFree error ")
	}
	{
		obj := struct {
			Test1 inject1 `container:"autowire:true"`
		}{a}
		DiFree(logger, &obj)
		assert.Equal(t, empty, obj.Test1, "diFree error ")
	}
	{
		obj := struct {
			Test1 inject1 `container:"autowire:false"`
		}{a}
		DiFree(logger, &obj)
		assert.Equal(t, empty, obj.Test1, "diFree error ")
	}
	{
		obj := struct {
			Test1 inject1 `container:"autowire:true"`
		}{a}
		DiFree(logger, &obj)
		assert.Equal(t, empty, obj.Test1, "diFree error ")
	}
	{
		obj := struct {
			Test1 *inject1 `container:"autowire:true"`
		}{&a}
		DiFree(logger, &obj)
		assert.Nil(t, obj.Test1, "diFree error ")
	}
	{
		obj := struct {
			Test1 string `container:"autowire:true"`
		}{"1"}
		DiFree(logger, &obj)
		assert.Equal(t, "", obj.Test1, "diFree error ")
	}
	{
		s := "1"
		obj := struct {
			Test1 *string `container:"autowire:true"`
		}{&s}
		DiFree(logger, &obj)
		assert.Nil(t, obj.Test1, "diFree error ")
	}
	{
		s := "1"
		obj := struct {
			Test1 *string `container:"autowire:true"`
			T     string
		}{&s, "1"}
		DiFree(logger, &obj)
		assert.Nil(t, obj.Test1, "diFree error ")
		assert.Equal(t, "1", obj.T, "diFree error ")
	}
}

func TestContainer_DiSelfCheck(t *testing.T) {
	type fields struct {
		c         *Container
		instances map[string]*sync.Pool
	}
	type args struct {
		instanceName string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "instance not in map ",
			fields: fields{
				c:         NewContainer(),
				instances: map[string]*sync.Pool{},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "instance `instance1` not exist in pool map",
		},
		{
			name: "instance cannot be addressable",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} { return testInjectErr1{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "the object to be injected container.testInjectErr1 should be addressable",
		},
		{
			name: "resource tag not set ",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} { return &testInjectErr1{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "self check error: instanceName: instance1 index: 0 objectName: *container.testInjectErr1.test1.(container.inject1), the resource tag not exist in container",
		},
		{
			name: "resource name not register ",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} { return &testInjectErr2{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "self check error: instanceName: instance1 index: 0 objectName: *container.testInjectErr2.test1.(container.inject1), resource name: inject1 not register in container ",
		},
		{
			name: "private param",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} { return &testInjectErr2{} },
					},
					"inject1": {
						New: func() interface{} { return &testInjectErr2{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "self check error: instanceName: instance1 index: 0 objectName: *container.testInjectErr2.test1.(container.inject1), private param",
		},
		{
			name: "object is not null ",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} {
							return &testInject3{
								Test1: inject1{
									name: "1",
								},
							}
						},
					},
					"inject1": {
						New: func() interface{} { return &testInjectErr2{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "self check error: instanceName: instance1 index: 0 objectName: *container.testInject3.Test1.(container.inject1), the param to be injected is not null",
		},
		{
			name: "struct inject type not equal",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} { return &testInject3{} },
					},
					"inject1": {
						New: func() interface{} { return &testInject3{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "self check error: instanceName: instance1 index: 0 objectName: *container.testInject3.Test1.(container.inject1), resource name: inject1 type not same, expected: container.inject1 actual: *container.testInject3",
		},
		{
			name: "ptr inject type not equal",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} { return &testInject4{} },
					},
					"inject1": {
						New: func() interface{} { return &testInject3{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "self check error: instanceName: instance1 index: 0 objectName: *container.testInject4.Test1.(*container.inject1), resource name: inject1 type not same, expected: *container.inject1 actual: *container.testInject3",
		},
		{
			name: "interface inject type not implement",
			fields: fields{
				c: NewContainer(),
				instances: map[string]*sync.Pool{
					"instance1": {
						New: func() interface{} { return &testInject5{} },
					},
					"inject1": {
						New: func() interface{} { return &testInject3{} },
					},
				},
			},
			args: args{
				instanceName: "instance1",
			},
			wantErr:    true,
			wantErrMsg: "self check error: instanceName: instance1 index: 0 objectName: *container.testInject5.Test1.(container.injectInterface), resource name: inject1 type:  not implement the interface injectInterface",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.fields.instances {
				tt.fields.c.RegisterInstance(k, v)
			}
			if err := tt.fields.c.DiSelfCheck(tt.args.instanceName); err != nil {
				if tt.wantErr {
					assert.Equal(t, tt.wantErrMsg, err.Error())
				} else {
					t.Error("unexpected error ")
					t.FailNow()
				}
			}

		})
	}
}

func TestContainer_DiAllFields(t *testing.T) {
	p1 := &inject1{}
	s1 := &testShared1{}
	s2 := &testShared2{}
	s1.T = s2
	s2.T = s1
	type fields struct {
		c       *Config
		poolMap map[string]*sync.Pool
	}
	type args struct {
		dest         interface{}
		injectingMap map[string]interface{}
	}
	type want struct {
		instance  interface{}
		injectMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "di ptr",
			fields: fields{
				c: DefaultConfig,
				poolMap: map[string]*sync.Pool{
					"inject1": {
						New: func() interface{} { return &inject1{} },
					},
				},
			},
			args: args{
				dest:         &testInject4{},
				injectingMap: make(map[string]interface{}),
			},
			want: want{
				instance: &testInject4{
					Test1: &inject1{},
				},
				injectMap: map[string]interface{}{
					"inject1": &inject1{},
				},
			},
		},
		{
			name: "di with inject map instance",
			fields: fields{
				c: DefaultConfig,
				poolMap: map[string]*sync.Pool{
					"inject1": {
						New: func() interface{} { return &inject1{} },
					},
				},
			},
			args: args{
				dest: &testInject4{},
				injectingMap: map[string]interface{}{
					"inject1": p1,
				},
			},
			want: want{
				instance: &testInject4{
					Test1: p1,
				},
				injectMap: map[string]interface{}{
					"inject1": p1,
				},
			},
		},
		{
			name: "di with inject map instance",
			fields: fields{
				c: DefaultConfig,
				poolMap: map[string]*sync.Pool{
					"shared1": {
						New: func() interface{} { return &testShared1{} },
					},
					"shared2": {
						New: func() interface{} { return &testShared2{} },
					},
				},
			},
			args: args{
				dest: &testShared1{},
				injectingMap: map[string]interface{}{
					"shared1": s1,
				},
			},
			want: want{
				instance: s1,
				injectMap: map[string]interface{}{
					"shared1": s1,
					"shared2": s2,
				},
			},
		},
		{
			name: "di with interface",
			fields: fields{
				c: DefaultConfig,
				poolMap: map[string]*sync.Pool{
					"shared1": {
						New: func() interface{} { return &testShared1{} },
					},
					"shared2": {
						New: func() interface{} { return &testShared2{} },
					},
					"inject1": {
						New: func() interface{} {
							return inject1{
								name: "1",
							}
						},
					},
				},
			},
			args: args{
				dest:         &testInject3{},
				injectingMap: map[string]interface{}{},
			},
			want: want{
				instance: &testInject3{
					Test1: inject1{
						name: "1",
					},
				},
				injectMap: map[string]interface{}{
					"inject1": inject1{
						name: "1",
					},
				},
			},
		},
		{
			name: "di with interface",
			fields: fields{
				c: DefaultConfig,
				poolMap: map[string]*sync.Pool{
					"shared1": {
						New: func() interface{} { return &testShared1{} },
					},
					"shared2": {
						New: func() interface{} { return &testShared2{} },
					},
					"inject1": {
						New: func() interface{} {
							return &inject1{
								name: "1",
							}
						},
					},
				},
			},
			args: args{
				dest:         &testInject4{},
				injectingMap: map[string]interface{}{},
			},
			want: want{
				instance: &testInject4{
					Test1: &inject1{
						name: "1",
					},
				},
				injectMap: map[string]interface{}{
					"inject1": &inject1{
						name: "1",
					},
				},
			},
		},
		{
			name: "di with empty interface",
			fields: fields{
				c: DefaultConfig,
				poolMap: map[string]*sync.Pool{
					"shared1": {
						New: func() interface{} { return &testShared1{} },
					},
					"shared2": {
						New: func() interface{} { return &testShared2{} },
					},
					"inject1": {
						New: func() interface{} {
							return &injectInterfaceImpl{}
						},
					},
				},
			},
			args: args{
				dest:         &testInject6{},
				injectingMap: map[string]interface{}{},
			},
			want: want{
				instance: &testInject6{
					Test1: &injectInterfaceImpl{},
				},
				injectMap: map[string]interface{}{
					"inject1": &injectInterfaceImpl{},
				},
			},
		},
		{
			name: "di with empty interface 2 ",
			fields: fields{
				c: DefaultConfig,
				poolMap: map[string]*sync.Pool{
					"shared1": {
						New: func() interface{} { return &testShared1{} },
					},
					"shared2": {
						New: func() interface{} { return &testShared2{} },
					},
					"inject1": {
						New: func() interface{} {
							return inject1{}
						},
					},
				},
			},
			args: args{
				dest:         &testInject6{},
				injectingMap: map[string]interface{}{},
			},
			want: want{
				instance: &testInject6{
					Test1: inject1{},
				},
				injectMap: map[string]interface{}{
					"inject1": inject1{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Container{
				c:       tt.fields.c,
				poolMap: tt.fields.poolMap,
			}
			s.DiAllFields(tt.args.dest, tt.args.injectingMap)
			assert.Equal(t, tt.want.instance, tt.args.dest, "wrong instance")
			assert.Equal(t, tt.want.injectMap, tt.args.injectingMap, "wrong inject map ")
		})
	}
}
