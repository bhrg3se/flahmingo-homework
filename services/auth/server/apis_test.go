package server

import (
	"context"
	"database/sql"
	pb "github.com/bhrg3se/flahmingo-homework/services/auth/pb/proto"
	"github.com/bhrg3se/flahmingo-homework/services/auth/store"
	"github.com/bhrg3se/flahmingo-homework/services/auth/testutils"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
	"reflect"
	"testing"
)

func TestServer_GetProfile(t *testing.T) {
	mockStore := new(store.MockStore)
	privateKey := testutils.GetMockPrivateKey1()

	mockStore.On("GetJWTPublicKey").Return(&privateKey.PublicKey)
	mockStore.On("GetUser", "someNumber").Return(&testutils.MockUser1, nil)

	type fields struct {
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		store                          store.GenericStore
	}
	type args struct {
		ctx     context.Context
		request *pb.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.User
		wantErr bool
	}{{
		name: "should get a profile with valid token",
		fields: fields{
			UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
			store:                          mockStore,
		},
		args: args{
			ctx:     context.Background(),
			request: &pb.Token{Token: testutils.MockToken1},
		},
		want: &pb.User{
			Id:          "someID",
			Name:        "Some User",
			PhoneNumber: "someNumber",
		},
		wantErr: false,
	},
		{
			name: "should return error if token is invalid",
			fields: fields{
				UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
				store:                          mockStore,
			},
			args: args{
				ctx:     context.Background(),
				request: &pb.Token{Token: "asdsadasdasdas"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error if token is signed with different private key",
			fields: fields{
				UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
				store:                          mockStore,
			},
			args: args{
				ctx:     context.Background(),
				request: &pb.Token{Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjYwNDgwMDAwMDAwMDAwMCwiaWF0IjoxNjMwNDc0ODYzLCJwaG9uZU51bWJlciI6InNvbWVOdW1iZXIifQ.DrvVIz_L7gw-rr6gkCQ0TPzXW70lPFvFGs2a7g_BkVeRV94MaRIRrc2aYSl9BdAXR5DDYjzlbD9ViJyt0fmXlDsApA-wG-D3WJhKg-x1fUoTfgCeq5wQAibmuCtoY_TJNYGwfWQQ_eEI0-wHKsTXujM4hNtvNUGRswxX8fP90_t9mIMCAy4HkaAm2Zpfjj2ECh_ZUKv8vzq8wLixICkpieZsQl9DjvwQuYSYFv7u5FNF0D2pbWLB6iqSzp_-YVAwDsvKFd8ScGxcwMuzWdP1RF7dPWbmWrcfPX2Z4NZS28GreCFKDJd7HYpyfAtx-56iYLSuNHg_D2EL0XCj5Bx9gA"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				store:                          tt.fields.store,
			}
			got, err := s.GetProfile(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProfile() got = %v, want %v", got, tt.want)
			}

			mockStore.AssertCalled(t, "GetJWTPublicKey")
			mockStore.AssertCalled(t, "GetUser", "someNumber")
		})
	}
}

func TestServer_LoginWithPhoneNumber(t *testing.T) {
	mockStore := new(store.MockStore)

	mockStore.On("SaveOTP", mock.AnythingOfType("string"), testutils.MockUser2.PhoneNumber).Return(nil)
	mockStore.On("GetUser", "").Return(nil, sql.ErrNoRows)

	mockStore.On("GetUser", testutils.MockUser2.PhoneNumber).Return(&testutils.MockUser2, nil)
	mockStore.On("SendOTP", context.Background(), mock.AnythingOfType("string"), testutils.MockUser2.PhoneNumber).Return(nil)

	type fields struct {
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		store                          store.GenericStore
	}
	type args struct {
		ctx     context.Context
		request *pb.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}{
		{
			name: "should fail when phone number is empty",
			fields: fields{
				UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
				store:                          mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &pb.User{
					PhoneNumber: "",
				},
			},
			want:    &emptypb.Empty{},
			wantErr: true,
		},

		{
			name: "should pass when phone number is valid",
			fields: fields{
				UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
				store:                          mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &pb.User{
					PhoneNumber: testutils.MockUser2.PhoneNumber,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				store:                          tt.fields.store,
			}
			got, err := s.LoginWithPhoneNumber(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginWithPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginWithPhoneNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_SignupWithPhoneNumber(t *testing.T) {
	mockStore := new(store.MockStore)

	mockStore.On("SaveOTP", mock.AnythingOfType("string"), testutils.MockUser2.PhoneNumber).Return(nil)

	mockStore.On("CreateUser", mock.Anything).Return(nil)
	mockStore.On("SendOTP", context.Background(), mock.AnythingOfType("string"), testutils.MockUser2.PhoneNumber).Return(nil)

	type fields struct {
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		store                          store.GenericStore
	}
	type args struct {
		ctx     context.Context
		request *pb.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}{
		{
			name: "should fail when phone number is empty",
			fields: fields{
				UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
				store:                          mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &pb.User{
					PhoneNumber: "",
				},
			},
			want:    &emptypb.Empty{},
			wantErr: true,
		},

		{
			name: "should pass when phone number is valid",
			fields: fields{
				UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
				store:                          mockStore,
			},
			args: args{
				ctx: context.Background(),
				request: &pb.User{
					PhoneNumber: testutils.MockUser2.PhoneNumber,
					Name:        testutils.MockUser2.Name,
				},
			},
			want:    &emptypb.Empty{},
			wantErr: false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				store:                          tt.fields.store,
			}
			got, err := s.SignupWithPhoneNumber(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignupWithPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignupWithPhoneNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ValidatePhoneNumberLogin(t *testing.T) {
	type fields struct {
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		store                          store.GenericStore
	}
	type args struct {
		ctx     context.Context
		request *pb.VerifyPhoneNumberRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				store:                          tt.fields.store,
			}
			got, err := s.ValidatePhoneNumberLogin(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneNumberLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidatePhoneNumberLogin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_VerifyPhoneNumber(t *testing.T) {
	type fields struct {
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		store                          store.GenericStore
	}
	type args struct {
		ctx     context.Context
		request *pb.VerifyPhoneNumberRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				store:                          tt.fields.store,
			}
			got, err := s.VerifyPhoneNumber(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyPhoneNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}
