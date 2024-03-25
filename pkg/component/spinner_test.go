package component

import "testing"

func TestSpinner_Render(t *testing.T) {
	type fields struct {
		Id      string
		NoModal bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spn := &Spinner{
				Id:      tt.fields.Id,
				NoModal: tt.fields.NoModal,
			}
			_, err := spn.Render()
			if (err != nil) != tt.wantErr {
				t.Errorf("Spinner.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
