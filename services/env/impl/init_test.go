package impl

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gndw/gank/services/env"
	"github.com/gndw/gank/services/flag"
	"github.com/stretchr/testify/mock"

	mocksLog "github.com/gndw/gank/services/utils/log/mocks"
	mocksMachinevar "github.com/gndw/gank/services/utils/machinevar/mocks"
)

func TestNew(t *testing.T) {
	type args struct {
		Log        *mocksLog.Service
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
				Log:        new(mocksLog.Service),
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return("", errors.New("not-found"))
				args.Log.On("Debugf", mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			want: &Service{
				env:            env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
				isReleaseLevel: false,
			},
		},
		{
			name: "success - using default env - from user preference",
			args: args{
				Log:        new(mocksLog.Service),
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
				Preference: &env.Preference{DefaultEnv: "my-custom-env"},
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return("", errors.New("not-found"))
				args.Log.On("Debugf", mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			want: &Service{
				env:            "my-custom-env",
				isReleaseLevel: false,
			},
		},
		{
			name: "success - using machine environment variable",
			args: args{
				Log:        new(mocksLog.Service),
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return(env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT, nil)
				args.Log.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			want: &Service{
				env:            env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
				isReleaseLevel: false,
			},
		},
		{
			name: "success - using custom machine environment variable name - from user preference",
			args: args{
				Log:        new(mocksLog.Service),
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
				Preference: &env.Preference{MachineEnvName: "CUSTOM_ENV"},
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", "CUSTOM_ENV").Return(env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT, nil)
				args.Log.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			want: &Service{
				env:            env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
				isReleaseLevel: false,
			},
		},
		{
			name: "success - using flag value",
			args: args{
				Log: new(mocksLog.Service),
				Flag: flag.Service{
					Env: &env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
				},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Log.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			want: &Service{
				env:            env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT,
				isReleaseLevel: false,
			},
		},
		{
			name: "success - using flag value - from user preference additional env",
			args: args{
				Log: new(mocksLog.Service),
				Flag: flag.Service{
					Env: func() *string { v := "custom-env"; return &v }(),
				},
				Machinevar: new(mocksMachinevar.Service),
				Preference: &env.Preference{AdditionalEnvs: []env.EnvLevel{{EnvName: "custom-env", IsReleaseLevel: true}}},
			},
			mock: func(args *args) {
				args.Log.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			want: &Service{
				env:            "custom-env",
				isReleaseLevel: true,
			},
		},
		{
			name: "error - using flag value - not allowed env",
			args: args{
				Log: new(mocksLog.Service),
				Flag: flag.Service{
					Env: func() *string { v := "custom-env"; return &v }(),
				},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Log.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			wantErr: true,
		},
		{
			name: "error - using machine var value - not allowed env",
			args: args{
				Log:        new(mocksLog.Service),
				Flag:       flag.Service{},
				Machinevar: new(mocksMachinevar.Service),
			},
			mock: func(args *args) {
				args.Machinevar.On("GetVar", env.DEFAULT_MACHINE_ENV_NAME).Return("custom-env", nil)
				args.Log.On("Debugf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
				args.Log.On("Infof", mock.Anything, mock.Anything)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(&tt.args)
			got, err := New(Parameters{Log: tt.args.Log, Flag: tt.args.Flag, Machinevar: tt.args.Machinevar, Preference: tt.args.Preference})
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
