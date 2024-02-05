package hardlinkdeduplicator

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/NIR3X/logger"
)

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hashSha512 := sha512.New()
	if _, err := io.Copy(hashSha512, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hashSha512.Sum(nil)), nil
}

type hardLink struct {
	index uint64
	path  string
	size  int64
}

func newHardLink(index uint64, path string, size int64) *hardLink {
	return &hardLink{index: index, path: path, size: size}
}

type hardLinks struct {
	mains []*hardLink
	links []*hardLink
}

var createHardLink func(src, dest string) error
var groupHardLinksByVolume func(files []*hardLink, verbose bool) map[uint32]map[uint64][]*hardLink

func Deduplicate(path string, all, deduplicate bool, minSize int64, verbose bool) {
	if createHardLink == nil || groupHardLinksByVolume == nil {
		logger.Eprintln("unsupported OS")
		return
	}

	filesBySize := map[int64][]*hardLink{}

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if verbose {
				logger.Eprintln(err)
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		size := info.Size()
		if size < minSize {
			return nil
		}

		if strings.HasSuffix(path, ".hldd") {
			newPath := strings.TrimSuffix(path, ".hldd")
			if err := os.Rename(path, newPath); err != nil {
				if verbose {
					logger.Eprintln(err)
				}
				return nil
			}
			path = newPath
		}

		filesBySize[size] = append(filesBySize[size], newHardLink(0, path, size))
		return nil
	})

	for _, files := range filesBySize {
		if len(files) < 2 {
			continue
		}

		volumes := groupHardLinksByVolume(files, verbose)
		for _, volume := range volumes {
			hardLinksByHash := map[string]*hardLinks{}

			for _, links := range volume {
				mainLink := links[0]
				hash, err := hashFile(mainLink.path)
				if err != nil {
					if verbose {
						logger.Eprintln(err)
					}
					continue
				}

				if hardLinksByHash[hash] == nil {
					hardLinksByHash[hash] = &hardLinks{
						mains: []*hardLink{mainLink},
						links: links[1:],
					}
				} else {
					hardLinks := hardLinksByHash[hash]
					if all || len(hardLinks.mains) >= 2 {
						hardLinks.links = append(hardLinks.links, links...)
					} else {
						hardLinks.links = append(hardLinks.links, links[1:]...)
					}
					hardLinks.mains = append(hardLinks.mains, mainLink)
				}
			}

			for _, hardLinks := range hardLinksByHash {
				mains := hardLinks.mains

				if all {
					if len(mains) < 2 {
						continue
					}
				} else {
					if len(mains) < 3 {
						continue
					}
				}

				if !deduplicate || verbose {
					fmt.Println("Duplicates found:")
					for _, mainLink := range mains {
						fmt.Printf("%s (%d bytes)\n", mainLink.path, mainLink.size)
					}
					fmt.Println()
				}

				files := hardLinks.links
				if deduplicate {
					if all {
						src := mains[0]
						for _, dest := range files {
							if src.index == dest.index || dest.index == 0 {
								continue
							}

							if verbose {
								fmt.Printf("Creating hard link for \"%s\" to \"%s\"\n", src.path, dest.path)
							}

							if err := createHardLink(src.path, dest.path); err != nil {
								if verbose {
									logger.Eprintln(err)
								}
							}
						}
					} else {
						src, src2 := mains[0], mains[1]
						for i, dest := range files {
							if i%2 == 0 {
								if src.index == dest.index || dest.index == 0 {
									continue
								}

								if verbose {
									fmt.Printf("Creating hard link for \"%s\" to \"%s\"\n", src.path, dest.path)
								}

								if err := createHardLink(src.path, dest.path); err != nil {
									if verbose {
										logger.Eprintln(err)
									}
								}
							} else {
								if src2.index == dest.index || dest.index == 0 {
									continue
								}

								if verbose {
									fmt.Printf("Creating hard link for \"%s\" to \"%s\"\n", src2.path, dest.path)
								}

								if err := createHardLink(src2.path, dest.path); err != nil {
									if verbose {
										logger.Eprintln(err)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
