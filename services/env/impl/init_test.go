package impl

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/flag"

	mocksMachinevar "github.com/gndw/gank/services/utils/machinevar/mocks"
)

func TestNew(t *testing.T) {
	type args struct {
		Flag       flag.Service
		Machinevar *mocksMachinevar.Service
		Preference *env.Preference `optional:"true"`
	}
	tests := []struct {
		name    string
		args    args
		mock    func(args *args)
		want    env.Service
		wantErr bool
	}{
		{
			name: "success - using default env",
			args: args{
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return("", errors.New("not-found"))
			},
			want: &Service{
				env: env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
			},
		},
		{
			name: "success - using default env - from user preference",
			args: args{
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
				Preference: &env.Preference{DefaultEnv: "my-custom-env"},
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return("", errors.New("not-found"))
			},
			want: &Service{
				env: "my-custom-env",
			},
		},
		{
			name: "success - using machine environment variable",
			args: args{
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return(env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT, nil)
			},
			want: &Service{
				env: env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
			},
		},
		{
			name: "success - using custom machine environment variable name - from user preference",
			args: args{
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
				Preference: &env.Preference{MachineEnvName: "CUSTOM_ENV"},
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", "CUSTOM_ENV").Return(env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT, nil)
			},
			want: &Service{
				env: env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
			},
		},
		{
			name: "success - using flag value",
			args: args{
				Flag: flag.Service{
					Env: &env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
				},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
			},
			want: &Service{
				env: env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
			},
		},
		{
			name: "success - using flag value - from user preference additional env",
			args: args{
				Flag: flag.Service{
					Env: func() *string { v := "custom-env"; return &v }(),
				},
				Machinevar: new(mocksMachinevar.Service),
				Preference: &env.Preference{AdditionalEnvs: []string{"custom-env"}},
			},
			mock: func(args *args) {
			},
			want: &Service{
				env: "custom-env",
			},
		},
		{
			name: "error - using flag value - not allowed env",
			args: args{
				Flag: flag.Service{
					Env: func() *string { v := "custom-env"; return &v }(),
				},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
			},
			wantErr: true,
		},
		{
			name: "error - using machine var value - not allowed env",
			args: args{
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return("custom-env", nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(&tt.args)
			got, err := New(Parameters{Flag: tt.args.Flag, Machinevar: tt.args.Machinevar, Preference: tt.args.Preference})
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
