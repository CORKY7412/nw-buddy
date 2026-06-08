package gltf

func (doc *Document) ImportMaterials(file string) {
	materials, _ := doc.Options.LoadMaterials(file)

	for _, mtl := range materials {
		doc.FindOrAddMaterial(mtl)
	}

}
