package mongoke

import (
	"testing"

	"github.com/PaesslerAG/gval"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBasicEvaluation(t *testing.T) {
	// notice gval does not support single quotes
	vars := Map{"document": Map{"name": "World"}}

	t.Run("using dots", func(t *testing.T) {

		value, err := gval.Evaluate(`"Hello " + document.name + "!"`, vars)
		if err != nil {
			t.Error(err)
		}
		t.Log(value)
	})
	t.Run("using subscription", func(t *testing.T) {

		value, err := gval.Evaluate(`"Hello " + document.name + "!"`, vars)
		if err != nil {
			t.Error(err)
		}
		t.Log(value)
	})
	t.Run("using nested ObjectId", func(t *testing.T) {
		hex := "000000000000000000000000"
		objID, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			t.Error(err)
		}
		vars := Map{
			"document":        Map{"user": Map{"_id": objID}},
			"jwt":             Map{"user_id": hex},
			"ObjectIDFromHex": primitive.ObjectIDFromHex,
		}
		value, err := gval.Evaluate(`document.user._id == ObjectIDFromHex(jwt.user_id)`, vars)
		if err != nil {
			t.Error(err)
		}
		t.Log(value)
		require.Equal(t, true, value)
	})

	// Output:
	// Hello World!
}
