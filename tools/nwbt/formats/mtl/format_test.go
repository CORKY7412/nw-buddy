package mtl_test

import (
	"encoding/json"
	"nw-buddy/tools/formats/mtl"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data, err := os.ReadFile("sample/male_alchemist_chest_matgroup.mtl")
	assert.NoError(t, err)

	doc, err := mtl.Parse(data)
	assert.NoError(t, err)

	assert.Equal(t, 524544, *doc.MtlFlags)
	assert.NotNil(t, doc.Params)
	assert.Equal(t, 3, doc.Params.Len())
	assert.Equal(t, "0", doc.Params.Get("SSSIndex"))
	assert.Equal(t, "0.25,0.25,0.25,0.25", doc.Params.Get("IndirectColor"))
}

func TestTextureSerialize(t *testing.T) {
	{
		tex := mtl.Texture{}
		data, err := json.Marshal(tex)
		assert.NoError(t, err)
		value := string(data)
		assert.Equal(t, "{}", value)
	}

	{
		tex := mtl.Texture{
			Filter: ptr[mtl.TextureFilter](0),
		}
		data, err := json.Marshal(tex)
		assert.NoError(t, err)
		value := string(data)
		assert.Equal(t, "{\"Filter\":0}", value)
	}

	{
		tex := mtl.Texture{
			IsTileU: ptr[mtl.TextureTile](0),
		}
		data, err := json.Marshal(tex)
		assert.NoError(t, err)
		value := string(data)
		assert.Equal(t, "{\"IsTileU\":0}", value)
	}

	{
		tex := mtl.Texture{
			IsTileV: ptr[mtl.TextureTile](0),
		}
		data, err := json.Marshal(tex)
		assert.NoError(t, err)
		value := string(data)
		assert.Equal(t, "{\"IsTileV\":0}", value)
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestIsShadowProxy(t *testing.T) {
	allExceptNoShadow := -1 ^ int(mtl.MTL_FLAG_NOSHADOW)

	e := mtl.Material{
		MaterialAttrs: mtl.MaterialAttrs{
			Opacity:  ptr[float32](0),
			MtlFlags: ptr(allExceptNoShadow),
			Shader:   "Illum",
		},
	}
	assert.True(t, e.IsShadowProxy())

	e.MaterialAttrs.Opacity = ptr[float32](1) // change
	e.MaterialAttrs.MtlFlags = ptr(0)
	e.MaterialAttrs.Shader = "Illum"
	assert.False(t, e.IsShadowProxy())

	e.MaterialAttrs.Opacity = ptr[float32](0)
	e.MaterialAttrs.MtlFlags = ptr(allExceptNoShadow | mtl.MTL_FLAG_NOSHADOW) // change
	e.MaterialAttrs.Shader = "Illum"
	assert.False(t, e.IsShadowProxy())

	e.MaterialAttrs.Opacity = ptr[float32](0)
	e.MaterialAttrs.MtlFlags = ptr(0)
	e.MaterialAttrs.Shader = "Other" // change

	assert.False(t, e.IsShadowProxy())
}
