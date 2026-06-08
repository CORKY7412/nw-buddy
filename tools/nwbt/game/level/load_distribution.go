package level

// func (dir *RegionDirectory) LoadRegionDistribution(assets *game.Assets) *RegionCapitalsData {

// }

// func (r *regionLoader) Distribution() *DistributionInfo {
// 	if r.distribution != nil {
// 		return r.distribution
// 	}

// 	if r.Info().Meta.Distribution == nil {
// 		return nil

// 	}
// 	file := r.Info().Meta.Distribution
// 	doc, err := distribution.Load(file)
// 	if err != nil {
// 		slog.Error("distribution not loaded", "file", file.Path(), "error", err)
// 		return nil
// 	}

// 	slices := maps.NewSafeDict[[]EntityInfo]()
// 	segments := maps.NewSafeDict[map[string]DistributionSlice]()

// 	progress.Concurrent(10, doc.Positions, func(position [2]uint16, i int) error {
// 		index := doc.Indices[i]
// 		sliceName := doc.Slices[index]

// 		slices.LoadOrStoreFn(sliceName, func() []EntityInfo {
// 			sliceFile := r.assets.ResolveDynamicSliceByName(sliceName)
// 			if sliceFile == nil {
// 				return nil
// 			}
// 			return LoadEntities(r.assets, sliceFile.Path(), mat4.Identity())
// 		})

// 		worldX, worldY, localX, localY := distribution.ConvertPosition(doc.Region, position)
// 		segmentX := int(m.Floor(float64(localX) / 128))
// 		segmentY := int(m.Floor(float64(localY) / 128))

// 		key := fmt.Sprintf("%d_%d", segmentX, segmentY)

// 		segments.Write(func() {
// 			bucket, _ := segments.LoadOrStoreFn(key, func() map[string]DistributionSlice {
// 				return make(map[string]DistributionSlice)
// 			})
// 			slice, ok := bucket[sliceName]
// 			if !ok {
// 				slice = DistributionSlice{
// 					Slice:     sliceName,
// 					Positions: make([][2]float32, 0),
// 				}
// 			}
// 			slice.Positions = append(slice.Positions, [2]float32{float32(worldX), float32(worldY)})
// 			bucket[sliceName] = slice
// 		})

// 		return nil
// 	})

// 	r.distribution = &DistributionInfo{
// 		Slices:   slices.ToMap(),
// 		Segments: make(map[string][]DistributionSlice),
// 	}

// 	return r.distribution
// }
