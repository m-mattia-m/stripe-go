package stripe

import (
	"encoding/json"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestGetObjectValue(t *testing.T) {
	event := &Event{
		Data: &EventData{
			Object: map[string]interface{}{
				"top_level_key": "top_level",
				"integer_key":   123,
				"map": map[string]interface{}{
					"nested_key": "nested",
				},
				"slice": []interface{}{
					"index-0",
					"index-1",
					"index-2",
				},
				"slice_of_maps": []interface{}{
					map[string]interface{}{
						"slice_nested_key": "slice_nested",
					},
				},
			},
		},
	}

	assert.Equal(t, "top_level", event.GetObjectValue("top_level_key"))

	// Check that it coerces non-string values into strings (this behavior is
	// somewhat questionable, but I'm going with how it already works)
	assert.Equal(t, "123", event.GetObjectValue("integer_key"))

	assert.Equal(t, "nested", event.GetObjectValue("map", "nested_key"))
	assert.Equal(t, "index-1", event.GetObjectValue("slice", "1"))
	assert.Equal(t, "slice_nested",
		event.GetObjectValue("slice_of_maps", "0", "slice_nested_key"))

	// By design a `nil` just returns an empty string
	assert.Equal(t, "", event.GetObjectValue("bad_key"))

	// Panic conditions. Usually the function tries to just return a value is
	// fairly forgiving, but it does panic under certain obviously impossible
	// cases.
	assert.PanicsWithValue(t, "Cannot access nested slice element with non-integer key: string_key", func() {
		event.GetObjectValue("slice", "string_key")
	})
	assert.PanicsWithValue(t, "Cannot descend into non-map non-slice object with key: bad_key", func() {
		event.GetObjectValue("top_level_key", "bad_key")
	})
}

func TestEmptyOptionalFields(t *testing.T) {

	// Marshals from a JSON object - without optional attributes
	{
		v := Event{}
		data, err := json.Marshal(&v)
		assert.NoError(t, err)

		expected := `{"api_version":"","created":0,"data":null,"id":"","livemode":false,"object":"","pending_webhooks":0,"request":null,"type":""}`
		assert.Equal(t, expected, string(data))
	}

	// Marshals from a JSON object - with empty optional attributes
	{
		v := Event{Account: "", Data: &EventData{PreviousAttributes: map[string]interface{}{}}}
		data, err := json.Marshal(&v)
		assert.NoError(t, err)

		expected := `{"api_version":"","created":0,"data":{"object":null},"id":"","livemode":false,"object":"","pending_webhooks":0,"request":null,"type":""}`
		assert.Equal(t, expected, string(data))
	}

	// Marshals from a JSON object - with filled account attribute
	{
		v := Event{
			Account: "",
			Data: &EventData{
				PreviousAttributes: map[string]interface{}{
					"amount_paid": 0,
				},
			},
		}
		data, err := json.Marshal(&v)
		assert.NoError(t, err)

		expected := `{"api_version":"","created":0,"data":{"previous_attributes":{"amount_paid":0},"object":null},"id":"","livemode":false,"object":"","pending_webhooks":0,"request":null,"type":""}`
		assert.Equal(t, expected, string(data))
	}

}
