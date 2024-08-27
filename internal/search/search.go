package search

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"sync"
	"time"
)

type Summary struct {
	Path      string
	TotalSize uint64
}

func GetInspectionSummary(fs afero.Fs, rootPaths *[]string) *map[string]uint64 {
	summary := make(map[string]uint64)
	for _, firstDirPath := range *rootPaths {
		summary[firstDirPath] = uint64(0)
	}

	ch := make(chan *Summary, len(*rootPaths))
	wg := new(sync.WaitGroup)
	wg.Add(len(*rootPaths))

	for _, root := range *rootPaths {
		go GetSummary(fs, root, ch, wg)
	}

	wg.Wait()
	close(ch)

	for dirSummary := range ch {
		summary[dirSummary.Path] = dirSummary.TotalSize
	}

	return &summary
}

func GetSummary(fs afero.Fs, root string, ch chan<- *Summary, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()

	var totalSize uint64
	err := afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += uint64(info.Size())
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}

	duration := time.Since(start)
	log.Info(fmt.Sprintf(
		"Search %s duration: %d Sec",
		root,
		duration/time.Second))

	ch <- &Summary{
		Path:      root,
		TotalSize: totalSize,
	}
}
