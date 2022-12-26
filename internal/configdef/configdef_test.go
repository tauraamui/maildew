package configdef_test

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/configdef"
)

func TestValidateEmptyConfigPasses(t *testing.T) {
	is := is.New(t)
	// TODO(tauraamui): return this to actually be empty again
	// once this root field has been removed
	body := `{}`
	config := configdef.Values{}
	is.NoErr(json.Unmarshal([]byte(body), &config))
	is.NoErr(config.RunValidate())
}

func TestValidatePopulatedConfigPassesValidation(t *testing.T) {
	is := is.New(t)
	body := `{
			"max_clip_age_in_days": 1,
			"cameras": [
				{
					"title": "NotBlank",
					"persist_location": "Nowhere",
					"max_clip_age_days": 15,
					"fps": 11,
					"seconds_per_clip": 1
				}
			]
		}`
	config := configdef.Values{}
	is.NoErr(json.Unmarshal([]byte(body), &config))
	is.NoErr(config.RunValidate())
}
