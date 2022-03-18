package impl

import (
	"testing"

	"github.com/gndw/gank/services/env"
)

func TestService_Get(t *testing.T) {
	type fields struct {
		env            string
		isReleaseLevel bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantEnv string
	}{
		{
			name:    "success",
			fields:  fields{env: "my-env"},
			wantEnv: "my-env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				env:            tt.fields.env,
				isReleaseLevel: tt.fields.isReleaseLevel,
			}
			if gotEnv := s.Get(); gotEnv != tt.wantEnv {
				t.Errorf("Service.Get() = %v, want %v", gotEnv, tt.wantEnv)
			}
		})
	}
}

func TestService_IsDevelopment(t *testing.T) {
	type fields struct {
		env            string
		isReleaseLevel bool
	}
	tests := []struct {
		name              string
		fields            fields
		wantIsDevelopment bool
	}{
		{
			name:              "success",
			fields:            fields{env: env.DEFAULT_ENV_NAME_ENV_DEVELOPMENT},
			wantIsDevelopment: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				env:            tt.fields.env,
				isReleaseLevel: tt.fields.isReleaseLevel,
			}
			if gotEnv := s.IsDevelopment(); gotEnv != tt.wantIsDevelopment {
				t.Errorf("Service.Get() = %v, want %v", gotEnv, tt.wantIsDevelopment)
			}
		})
	}
}

func TestService_IsStaging(t *testing.T) {
	type fields struct {
		env            string
		isReleaseLevel bool
	}
	tests := []struct {
		name          string
		fields        fields
		wantIsStaging bool
	}{
		{
			name:          "success",
			fields:        fields{env: env.DEFAULT_ENV_NAME_ENV_STAGING},
			wantIsStaging: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				env:            tt.fields.env,
				isReleaseLevel: tt.fields.isReleaseLevel,
			}
			if gotEnv := s.IsStaging(); gotEnv != tt.wantIsStaging {
				t.Errorf("Service.Get() = %v, want %v", gotEnv, tt.wantIsStaging)
			}
		})
	}
}

func TestService_IsProduction(t *testing.T) {
	type fields struct {
		env            string
		isReleaseLevel bool
	}
	tests := []struct {
		name             string
		fields           fields
		wantIsProduction bool
	}{
		{
			name:             "success",
			fields:           fields{env: env.DEFAULT_ENV_NAME_ENV_PRODUCTION},
			wantIsProduction: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				env:            tt.fields.env,
				isReleaseLevel: tt.fields.isReleaseLevel,
			}
			if gotEnv := s.IsProduction(); gotEnv != tt.wantIsProduction {
				t.Errorf("Service.Get() = %v, want %v", gotEnv, tt.wantIsProduction)
			}
		})
	}
}

func TestService_IsReleaseLevel(t *testing.T) {
	type fields struct {
		env            string
		isReleaseLevel bool
	}
	tests := []struct {
		name               string
		fields             fields
		wantIsReleaseLevel bool
	}{
		{
			name:               "success",
			fields:             fields{isReleaseLevel: true},
			wantIsReleaseLevel: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				env:            tt.fields.env,
				isReleaseLevel: tt.fields.isReleaseLevel,
			}
			if gotIsReleaseLevel := s.IsReleaseLevel(); gotIsReleaseLevel != tt.wantIsReleaseLevel {
				t.Errorf("Service.IsReleaseLevel() = %v, want %v", gotIsReleaseLevel, tt.wantIsReleaseLevel)
			}
		})
	}
}
