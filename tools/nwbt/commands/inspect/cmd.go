package inspect

import (
	"fmt"
	"io"
	"nw-buddy/tools/formats/cgf"
	"nw-buddy/tools/game"
	"nw-buddy/tools/nwfs"
	"nw-buddy/tools/utils"
	"nw-buddy/tools/utils/env"
	"nw-buddy/tools/utils/progress"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var flgGameDir string
var flgRegex bool
var Cmd = &cobra.Command{
	Use:           "inspect",
	Short:         "inspects a file",
	Long:          ``,
	Run:           run,
	SilenceErrors: false,
}

func init() {
	Cmd.Flags().StringVarP(&flgGameDir, "game", "g", env.GameDir(), "game root directory")
}

func run(ccmd *cobra.Command, args []string) {
	assets := utils.Must(game.InitPackedAssets(flgGameDir))
	files, err := listFiles(args, flgRegex, assets.Archive)
	if err != nil {
		panic(err)
	}
	inspectors := make(map[string]Inscpector)

	bar := progress.Bar(len(files), "Processing")
	for _, file := range files {
		ext := path.Ext(file.Path())
		switch ext {
		case ".cgf":
			if inspectors[ext] == nil {
				inspectors[ext] = &CgfInspector{}
			}
			inspectors[ext].Inspect(assets, file)
		default:
			if inspectors[ext] == nil {
				inspectors[ext] = &UnsupportedInspector{}
			}
			inspectors[ext].Inspect(assets, file)
		}
		bar.Add(1)
	}
	bar.Wait()
	for key := range inspectors {
		inspectors[key].Print(os.Stdout)
	}
}

func listFiles(args []string, regex bool, fs nwfs.Archive) ([]nwfs.File, error) {
	if len(args) == 0 {
		return fs.List()
	}
	if regex {
		for i := range args {
			args[i] = strings.ToLower(args[i])
			args[i] = strings.ReplaceAll(args[i], "\\\\", "\\")
		}
		return fs.Match(args...)
	}
	for i := range args {
		args[i] = strings.ToLower(args[i])
	}
	return fs.Glob(args...)
}

type Inscpector interface {
	Inspect(assets *game.Assets, file nwfs.File)
	Print(w io.Writer)
}

type CgfInspector struct {
	count     int
	nodeNames map[string]int
	nodeCount map[int]int
}

func (it *CgfInspector) Inspect(assets *game.Assets, file nwfs.File) {
	it.count += 1
	f, err := cgf.Load(file)
	if err != nil {
		return
	}
	nodes := cgf.SelectChunks[cgf.ChunkNode](f)
	hasdollar := false
	if it.nodeNames == nil {
		it.nodeNames = make(map[string]int)
	}
	for i := range nodes {
		if strings.Contains(nodes[i].Name, "$lod") {
			hasdollar = true
			it.nodeNames[nodes[i].Name] += 1
		}

	}

	if it.nodeCount == nil {
		it.nodeCount = make(map[int]int)
	}
	if hasdollar {
		it.nodeCount[len(nodes)] += 1
		fmt.Printf("%03d %s\n", len(nodes), file.Path())
	}
}

func (it *CgfInspector) Print(w io.Writer) {

	fmt.Fprintf(w, "CGF\t%d\n", it.count)
	fmt.Fprintf(w, "node name\tcount\n")
	for key := range it.nodeNames {
		if strings.Contains(key, "$") {
			count := it.nodeNames[key]
			fmt.Fprintf(w, "%03d\t%s\n", count, key)
		}
	}

	for c := range it.nodeCount {
		fmt.Fprintf(w, "nodes %03d\t%03d\n", c, it.nodeCount[c])
	}
}

type UnsupportedInspector struct {
	count int
}

func (it *UnsupportedInspector) Inspect(assets *game.Assets, file nwfs.File) {
	it.count += 1
}
func (it *UnsupportedInspector) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	fmt.Fprintf(tw, "other\t%d\n", it.count)
	tw.Flush()
}
