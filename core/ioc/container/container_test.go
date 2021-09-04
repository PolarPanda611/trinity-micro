package container

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// import (
// 	"reflect"
// 	"sync"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// type test1 struct{}
// type test2 struct{}

// var (
// 	pool1 = &sync.Pool{New: func() interface{} { return &test1{} }}
// 	pool2 = &sync.Pool{New: func() interface{} { return &test2{} }}
// )

// func TestNewContainer(t *testing.T) {
// 	c := Config{
// 		AutoWired: false,
// 	}
// 	tests := []struct {
// 		name string
// 		c    []Config
// 		want *Container
// 	}{}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewContainer(tt.c...); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewContainer() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestContainer_newInstance(t *testing.T) {

// 	type fields struct {
// 		poolMap          map[reflect.Type]*sync.Pool
// 		instanceTypeList []reflect.Type
// 		poolTags         map[string]reflect.Type
// 		instanceMapping  map[string]reflect.Type
// 	}
// 	type args struct {
// 		instanceType reflect.Type
// 		instancePool *sync.Pool
// 		instanceTag  []string
// 	}
// 	tests := []struct {
// 		name      string
// 		fields    fields
// 		args      args
// 		wantPanic bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "1",
// 			fields: fields{
// 				poolMap: map[reflect.Type]*sync.Pool{
// 					reflect.TypeOf(&test1{}): pool1,
// 				},
// 				instanceTypeList: []reflect.Type{reflect.TypeOf(&test1{})},
// 				poolTags: map[string]reflect.Type{
// 					"test1": reflect.TypeOf(&test1{}),
// 				},
// 				instanceMapping: map[string]reflect.Type{},
// 			},
// 			args: args{
// 				instanceType: reflect.TypeOf(&test1{}),
// 				instancePool: &sync.Pool{New: func() interface{} { return &test1{} }},
// 				instanceTag:  []string{"test1"},
// 			},
// 			wantPanic: false,
// 		},
// 		{
// 			name: "2",
// 			fields: fields{
// 				poolMap:          map[reflect.Type]*sync.Pool{},
// 				instanceTypeList: []reflect.Type{},
// 				poolTags:         map[string]reflect.Type{},
// 				instanceMapping:  map[string]reflect.Type{},
// 			},
// 			args: args{
// 				instanceType: reflect.TypeOf(&test1{}),
// 				instancePool: pool1,
// 				instanceTag:  []string{"test1"},
// 			},
// 			wantPanic: false,
// 		},
// 		{
// 			name: "3",
// 			fields: fields{
// 				poolMap: map[reflect.Type]*sync.Pool{
// 					reflect.TypeOf(&test1{}): pool1,
// 				},
// 				instanceTypeList: []reflect.Type{reflect.TypeOf(&test1{})},
// 				poolTags: map[string]reflect.Type{
// 					"test1": reflect.TypeOf(&test1{}),
// 				},
// 				instanceMapping: map[string]reflect.Type{},
// 			},
// 			args: args{
// 				instanceType: reflect.TypeOf(&test2{}),
// 				instancePool: &sync.Pool{New: func() interface{} { return &test1{} }},
// 				instanceTag:  []string{"test1"},
// 			},
// 			wantPanic: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Container{
// 				poolMap:          tt.fields.poolMap,
// 				instanceTypeList: tt.fields.instanceTypeList,
// 				poolTags:         tt.fields.poolTags,
// 				instanceMapping:  tt.fields.instanceMapping,
// 			}
// 			if tt.wantPanic {
// 				assert.PanicsWithError(t, "tag test1 already existed", func() { s.newInstance(tt.args.instanceType, tt.args.instancePool, tt.args.instanceTag...) })
// 			} else {
// 				s.newInstance(tt.args.instanceType, tt.args.instancePool, tt.args.instanceTag...)
// 				assert.Equal(t, &Container{
// 					poolMap: map[reflect.Type]*sync.Pool{
// 						reflect.TypeOf(&test1{}): pool1,
// 					},
// 					instanceTypeList: []reflect.Type{reflect.TypeOf(&test1{})},
// 					poolTags: map[string]reflect.Type{
// 						"test1": reflect.TypeOf(&test1{}),
// 					},
// 					instanceMapping: map[string]reflect.Type{},
// 				}, s, "instance not register correctly")
// 			}

// 		})
// 	}
// }

// func TestContainer_GetInstanceTypeByTag(t *testing.T) {
// 	p1 := NewContainer()
// 	p1.newInstance(reflect.TypeOf(&test1{}), pool1, "test1")
// 	p1.newInstance(reflect.TypeOf(&test1{}), pool1, "test2")
// 	p1.newInstance(reflect.TypeOf(&test2{}), pool2, "test3")
// 	p1.newInstance(reflect.TypeOf(&test2{}), pool2, "test4")
// 	type fields struct {
// 		poolMap          map[reflect.Type]*sync.Pool
// 		instanceTypeList []reflect.Type
// 		poolTags         map[string]reflect.Type
// 		instanceMapping  map[string]reflect.Type
// 	}
// 	type args struct {
// 		tags []string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields *Container
// 		args   args
// 		want   []reflect.Type
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name:   "1",
// 			fields: p1,
// 			args: args{
// 				tags: []string{"test1"},
// 			},
// 			want: []reflect.Type{reflect.TypeOf(&test1{})},
// 		},
// 		{
// 			name:   "2",
// 			fields: p1,
// 			args: args{
// 				tags: []string{},
// 			},
// 			want: []reflect.Type{reflect.TypeOf(&test1{}), reflect.TypeOf(&test2{})},
// 		},
// 		{
// 			name:   "3",
// 			fields: p1,
// 			args: args{
// 				tags: []string{"test1", "test3"},
// 			},
// 			want: []reflect.Type{reflect.TypeOf(&test1{}), reflect.TypeOf(&test2{})},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.fields.GetInstanceTypeByTag(tt.args.tags...); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Container.GetInstanceTypeByTag() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestContainer_CheckInstanceNameIfExist(t *testing.T) {
// 	p1 := NewContainer()
// 	p1.newInstance(reflect.TypeOf(&test1{}), pool1, "test1")
// 	type args struct {
// 		instanceName reflect.Type
// 	}
// 	tests := []struct {
// 		name   string
// 		fields *Container
// 		args   args
// 		want   bool
// 	}{
// 		{
// 			name:   "1",
// 			fields: p1,
// 			args: args{
// 				instanceName: reflect.TypeOf(&test1{}),
// 			},
// 			want: true,
// 		},
// 		{
// 			name:   "2",
// 			fields: p1,
// 			args: args{
// 				instanceName: reflect.TypeOf(&test2{}),
// 			},
// 			want: false,
// 		},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &Container{
// 				poolMap:          tt.fields.poolMap,
// 				instanceTypeList: tt.fields.instanceTypeList,
// 				poolTags:         tt.fields.poolTags,
// 				instanceMapping:  tt.fields.instanceMapping,
// 			}
// 			if got := s.CheckInstanceNameIfExist(tt.args.instanceName); got != tt.want {
// 				t.Errorf("Container.CheckInstanceNameIfExist() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
			want: true,
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

// func TestContainer_RegisterInstance(t *testing.T) {
// 	type fields struct {
// 		c *Container
// 	}
// 	type args struct {
// 		instance    interface{}
// 		instanceTag []string
// 	}
// 	type want struct {
// 		typeList []reflect.Type
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want
// 	}{
// 		{
// 			name: "1",
// 			fields: fields{
// 				c: NewContainer(Config{
// 					AutoFree:  true,
// 					AutoWired: false,
// 				}),
// 			},
// 			args: args{
// 				instance:    test1{},
// 				instanceTag: []string{"test1"},
// 			},
// 			want: want{
// 				typeList: []reflect.Type{reflect.TypeOf(test1{})},
// 			},
// 		},
// 		{
// 			name: "2",
// 			fields: fields{
// 				c: NewContainer(Config{
// 					AutoFree:  true,
// 					AutoWired: false,
// 				}),
// 			},
// 			args: args{
// 				instance:    &test1{},
// 				instanceTag: []string{"test1"},
// 			},
// 			want: want{
// 				typeList: []reflect.Type{reflect.TypeOf(test1{})},
// 			},
// 		},
// 		{
// 			name: "3",
// 			fields: fields{
// 				c: NewContainer(Config{
// 					AutoFree:  true,
// 					AutoWired: false,
// 				}),
// 			},
// 			args: args{
// 				instance:    func() interface{} { return test1{} },
// 				instanceTag: []string{"test1"},
// 			},
// 			want: want{
// 				typeList: []reflect.Type{reflect.TypeOf(test1{})},
// 			},
// 		},
// 		{
// 			name: "4",
// 			fields: fields{
// 				c: NewContainer(Config{
// 					AutoFree:  true,
// 					AutoWired: false,
// 				}),
// 			},
// 			args: args{
// 				instance:    func() interface{} { return "123" },
// 				instanceTag: []string{"test1"},
// 			},
// 			want: want{
// 				typeList: []reflect.Type{reflect.TypeOf("123")},
// 			},
// 		},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.fields.c.RegisterInstance(tt.args.instance, tt.args.instanceTag...)
// 			assert.Equal(t, tt.want.typeList, tt.fields.c.instanceTypeList, "wrong type ")
// 		})
// 	}
// }

type userRep struct{}
type addressRepo struct{}
type orderRepo struct{}
type shipRepo struct{}

type userSrv struct {
	UserRep    *userRep    `container:"autowire:true;resource:userRep"`
	AddressSrv *addressSrv `container:"autowire:true;resource:addressSrv"`
	OrderSrv   *orderSrv   `container:"autowire:true;resource:orderSrv"`
}
type addressSrv struct {
	UserSrv     *userSrv     `container:"autowire:true;resource:userSrv"`
	AddressRepo *addressRepo `container:"autowire:true;resource:addressRepo"`
}
type shipSrv struct {
	UserSrv    *userSrv    `container:"autowire:true;resource:userSrv"`
	AddressSrv *addressSrv `container:"autowire:true;resource:addressSrv"`
	OrderSrv   *orderSrv   `container:"autowire:true;resource:orderSrv"`
	ShipRepo   *shipRepo   `container:"autowire:true;resource:shipRepo"`
}
type orderSrv struct {
	UserSrv    *userSrv    `container:"autowire:true;resource:userSrv"`
	AddressSrv *addressSrv `container:"autowire:true;resource:addressSrv"`
	ShipSrv    *shipSrv    `container:"autowire:true;resource:shipSrv"`
	OrderRepo  *orderRepo  `container:"autowire:true;resource:orderRepo"`
}

type userCtl struct {
	UserSrv *userSrv `container:"autowire:true;resource:userSrv"`
}
type addressCtl struct {
	AddressSrv *addressSrv `container:"autowire:true;resource:addressSrv"`
}
type orderCtl struct {
	OrderSrv *orderSrv `container:"autowire:true;resource:orderSrv"`
}
type shipCtl struct {
	ShipSrv *shipSrv `container:"autowire:true;resource:shipSrv"`
}

func TestContainer_GetInstance(t *testing.T) {

	{
		// shared instance
		s1 := &testShared1{}
		s2 := &testShared2{}
		s1.T = s2
		s2.T = s1
		s := NewContainer()
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testShared1{} }}, "shared1")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testShared2{} }}, "shared2")
		injectMap := make(map[string]interface{})
		instance := s.GetInstance("shared1", injectMap)
		assert.Equal(t, s1, instance, "wrong instance ")
		t1 := instance.(*testShared1)
		assert.Equal(t, fmt.Sprintf("%p", t1), fmt.Sprintf("%p", t1.T.T), "wrong ptr")
		assert.Equal(t, fmt.Sprintf("%p", t1.T), fmt.Sprintf("%p", t1.T.T.T), "wrong ptr")
		assert.NotPanics(t, func() { fmt.Printf("%p", t1.T.T.T.T.T.T.T.T.T.T.T.T.T.T.T.T.T.T.T) }, "inject failed")
	}
	{
		// ptr register
		type testRep struct{}
		type testSrv struct {
			TestRep *testRep `container:"autowire:true;resource:testRep"`
		}
		type testCtl struct {
			TestSrv *testSrv `container:"autowire:true;resource:testSrv"`
		}

		res := &testCtl{
			TestSrv: &testSrv{
				TestRep: &testRep{},
			},
		}
		s := NewContainer()
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testRep{} }}, "testRep")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testSrv{} }}, "testSrv")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testCtl{} }}, "testCtl")
		injectMap := make(map[string]interface{})
		instance := s.GetInstance("testCtl", injectMap)
		assert.Equal(t, res, instance, "wrong instance ")
		t1 := instance.(*testCtl)
		assert.NotPanics(t, func() { fmt.Printf("%p", t1.TestSrv.TestRep) }, "inject failed")
	}
	{
		// interface register
		type testRepI interface{}
		type testRep struct{}
		type testSrvI interface{}
		type testSrv struct {
			TestRep testRepI `container:"autowire:true;resource:testRepI"`
		}
		type testCtl struct {
			TestSrv testSrvI `container:"autowire:true;resource:testSrvI"`
		}

		res := &testCtl{
			TestSrv: &testSrv{
				TestRep: &testRep{},
			},
		}
		s := NewContainer()
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testRep{} }}, "testRepI")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testSrv{} }}, "testSrvI")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &testCtl{} }}, "testCtl")
		injectMap := make(map[string]interface{})
		instance := s.GetInstance("testCtl", injectMap)
		assert.Equal(t, res, instance, "wrong instance ")
		t1 := instance.(*testCtl)
		assert.NotPanics(t, func() { fmt.Printf("%p", &t1.TestSrv) }, "inject failed")
		t2 := t1.TestSrv.(*testSrv)
		assert.NotPanics(t, func() { fmt.Printf("%p", &t2.TestRep) }, "inject failed")
	}

	{
		// blackbox
		uSrv := &userSrv{}
		aSrv := &addressSrv{}
		oSrv := &orderSrv{}
		sSrv := &shipSrv{}
		uSrv.AddressSrv = aSrv
		uSrv.OrderSrv = oSrv
		uSrv.UserRep = &userRep{}
		aSrv.UserSrv = uSrv
		aSrv.AddressRepo = &addressRepo{}
		oSrv.AddressSrv = aSrv
		oSrv.UserSrv = uSrv
		oSrv.ShipSrv = sSrv
		oSrv.OrderRepo = &orderRepo{}
		sSrv.AddressSrv = aSrv
		sSrv.ShipRepo = &shipRepo{}
		sSrv.OrderSrv = oSrv
		sSrv.UserSrv = uSrv
		userCtlRes := &userCtl{
			UserSrv: uSrv,
		}
		addressCtlRes := &addressCtl{
			AddressSrv: aSrv,
		}
		orderCtlRes := &orderCtl{
			OrderSrv: oSrv,
		}
		shipCtlRes := &shipCtl{
			ShipSrv: sSrv,
		}
		s := NewContainer()
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &userRep{} }}, "userRep")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &addressRepo{} }}, "addressRepo")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &orderRepo{} }}, "orderRepo")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &shipRepo{} }}, "shipRepo")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &userSrv{} }}, "userSrv")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &addressSrv{} }}, "addressSrv")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &orderSrv{} }}, "orderSrv")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &shipSrv{} }}, "shipSrv")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &userCtl{} }}, "userCtl")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &addressCtl{} }}, "addressCtl")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &orderCtl{} }}, "orderCtl")
		s.RegisterInstance(&sync.Pool{New: func() interface{} { return &shipCtl{} }}, "shipCtl")
		if err := s.InstanceDISelfCheck(); err != nil {
			assert.Error(t, err, "self check error ")
		}
		{
			injectMap := make(map[string]interface{})
			instance := s.GetInstance("userCtl", injectMap)
			assert.Equal(t, userCtlRes, instance, "wrong instance ")
			t1 := instance.(*userCtl)
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.UserSrv.UserRep) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.UserSrv.AddressSrv) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.UserSrv.OrderSrv) }, "inject failed")
			for k, v := range injectMap {
				s.Release(k, v)
			}
			assert.Panics(t, func() { fmt.Printf("%p", t1.UserSrv.UserRep) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.UserSrv.AddressSrv) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.UserSrv.OrderSrv) }, "release failed")
		}
		{
			injectMap := make(map[string]interface{})
			instance := s.GetInstance("addressCtl", injectMap)
			assert.Equal(t, addressCtlRes, instance, "wrong instance ")
			t1 := instance.(*addressCtl)
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.AddressSrv.AddressRepo) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.AddressSrv.UserSrv) }, "inject failed")
			for k, v := range injectMap {
				s.Release(k, v)
			}
			assert.Panics(t, func() { fmt.Printf("%p", t1.AddressSrv.AddressRepo) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.AddressSrv.UserSrv) }, "release failed")
		}
		{
			injectMap := make(map[string]interface{})
			instance := s.GetInstance("orderCtl", injectMap)
			assert.Equal(t, orderCtlRes, instance, "wrong instance ")
			t1 := instance.(*orderCtl)
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.OrderSrv.UserSrv) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.OrderSrv.AddressSrv) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.OrderSrv.ShipSrv) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.OrderSrv.OrderRepo) }, "inject failed")
			for k, v := range injectMap {
				s.Release(k, v)
			}
			assert.Panics(t, func() { fmt.Printf("%p", t1.OrderSrv.UserSrv) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.OrderSrv.AddressSrv) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.OrderSrv.ShipSrv) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.OrderSrv.OrderRepo) }, "release failed")
		}
		{
			injectMap := make(map[string]interface{})
			instance := s.GetInstance("shipCtl", injectMap)
			assert.Equal(t, shipCtlRes, instance, "wrong instance ")
			t1 := instance.(*shipCtl)
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.ShipSrv.UserSrv) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.ShipSrv.AddressSrv) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.ShipSrv.ShipRepo) }, "inject failed")
			assert.NotPanics(t, func() { fmt.Printf("%p", t1.ShipSrv.OrderSrv) }, "inject failed")
			for k, v := range injectMap {
				s.Release(k, v)
			}
			assert.Panics(t, func() { fmt.Printf("%p", t1.ShipSrv.UserSrv) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.ShipSrv.AddressSrv) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.ShipSrv.ShipRepo) }, "release failed")
			assert.Panics(t, func() { fmt.Printf("%p", t1.ShipSrv.OrderSrv) }, "release failed")
		}

	}
}
