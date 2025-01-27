package rancherdesktop

import (
	"strings"
	"testing"
)

func TestResponseConfig(t *testing.T) {

	t.Run(".Validate", func(t *testing.T) {

		t.Run("should return nil for a valid ResponseConfig", func(t *testing.T) {
			responseConfig := ResponseConfig{
				Versions: []Version{
					{
						Name:        "1.2.3",
						ReleaseDate: "2022-07-28T11:00:00Z",
					},
					{
						Name:        "2.3.4",
						ReleaseDate: "2022-07-28T11:00:00Z",
						Tags:        []string{"latest"},
					},
				},
			}
			err := responseConfig.Validate()
			if err != nil {
				t.Errorf("unexpected error %q", err)
			}
		})

		// Test error conditions
		testCases := []struct {
			Description    string
			ResponseConfig ResponseConfig
			ExpectedError  string
		}{
			{
				Description: "should return error when there are duplicate versions",
				ResponseConfig: ResponseConfig{
					Versions: []Version{
						{
							Name:        "1.2.3",
							ReleaseDate: "2022-07-28T11:00:00Z",
						},
						{
							Name:        "1.2.3",
							ReleaseDate: "2022-07-28T11:00:00Z",
						},
					},
				},
				ExpectedError: "duplicate version name",
			},
			{
				Description: "should return error when there is no version with a latest tag",
				ResponseConfig: ResponseConfig{
					Versions: []Version{
						{
							Name:        "1.2.3",
							ReleaseDate: "2022-07-28T11:00:00Z",
						},
					},
				},
				ExpectedError: "did not find exactly one latest tag",
			},
			{
				Description: "should return error when there is more than one version with a latest tag",
				ResponseConfig: ResponseConfig{
					Versions: []Version{
						{
							Name:        "1.2.3",
							ReleaseDate: "2022-07-28T11:00:00Z",
							Tags:        []string{"latest"},
						},
						{
							Name:        "2.3.4",
							ReleaseDate: "2022-07-28T11:00:00Z",
							Tags:        []string{"latest"},
						},
					},
				},
				ExpectedError: "did not find exactly one latest tag",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.Description, func(t *testing.T) {
				err := testCase.ResponseConfig.Validate()
				if err == nil {
					t.Errorf("did not return error for %#v", testCase.ResponseConfig)
				} else if !strings.Contains(err.Error(), testCase.ExpectedError) {
					t.Errorf("error %q does not contain %q", err, testCase.ExpectedError)
				}
			})
		}

	})
}

func TestReadConfig(t *testing.T) {
	t.Run("all Version.Supported fields in returned config should be true", func(t *testing.T) {
		config, err := ReadConfig("testdata/test-config.json")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		for _, version := range config.Versions {
			if !version.Supported {
				t.Errorf("version %q has Supported value %t", version.Name, version.Supported)
			}
		}
	})
}
