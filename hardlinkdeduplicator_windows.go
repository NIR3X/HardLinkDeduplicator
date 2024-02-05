package hardlinkdeduplicator

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/NIR3X/logger"
)

func init() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	if kernel32 == nil {
		logger.Eprintln("kernel32.dll not found")
		return
	}

	procCreateHardLink := kernel32.NewProc("CreateHardLinkW")
	if procCreateHardLink == nil {
		logger.Eprintln("CreateHardLinkW not found")
		return
	}

	CreateHardLink := func(src, dest string) error {
		pSrc, err := syscall.UTF16PtrFromString(src)
		if err != nil {
			return err
		}

		pDest, err := syscall.UTF16PtrFromString(dest)
		if err != nil {
			return err
		}

		r1, _, err := procCreateHardLink.Call(uintptr(unsafe.Pointer(pDest)), uintptr(unsafe.Pointer(pSrc)), 0)
		if r1 == 0 {
			return err
		}

		return nil
	}

	procCreateFile := kernel32.NewProc("CreateFileW")
	if procCreateFile == nil {
		logger.Eprintln("CreateFileW not found")
		return
	}

	CreateFile := func(path string) (syscall.Handle, error) {
		pPath, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			return syscall.InvalidHandle, err
		}

		r1, _, err := procCreateFile.Call(uintptr(unsafe.Pointer(pPath)), syscall.GENERIC_READ, syscall.FILE_SHARE_READ, 0, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
		if syscall.Handle(r1) == syscall.InvalidHandle {
			return syscall.InvalidHandle, err
		}

		return syscall.Handle(r1), nil
	}

	procCloseHandle := kernel32.NewProc("CloseHandle")
	if procCloseHandle == nil {
		logger.Eprintln("CloseHandle not found")
		return
	}

	CloseHandle := func(handle syscall.Handle) error {
		r1, _, err := procCloseHandle.Call(uintptr(handle))
		if r1 == 0 {
			return err
		}

		return nil
	}

	procGetFileInformationByHandle := kernel32.NewProc("GetFileInformationByHandle")
	if procGetFileInformationByHandle == nil {
		logger.Eprintln("GetFileInformationByHandle not found")
		return
	}

	GetFileInformationByPath := func(path string) (syscall.ByHandleFileInformation, error) {
		handle, err := CreateFile(path)
		if err != nil {
			return syscall.ByHandleFileInformation{}, err
		}
		defer CloseHandle(handle)

		var data syscall.ByHandleFileInformation
		r1, _, err := procGetFileInformationByHandle.Call(uintptr(handle), uintptr(unsafe.Pointer(&data)))
		if r1 == 0 {
			return syscall.ByHandleFileInformation{}, err
		}

		return data, nil
	}

	createHardLink = func(src string, dest string) error {
		destBackup := dest + ".hldd"

		if _, err := os.Stat(dest); err == nil {
			if err := os.Rename(dest, destBackup); err != nil {
				return err
			}
		}

		if err := CreateHardLink(src, dest); err != nil {
			os.Rename(destBackup, dest)
			return err
		}

		if err := os.Remove(destBackup); err != nil {
			os.Remove(dest)
			os.Rename(destBackup, dest)
			return err
		}

		return nil
	}

	groupHardLinksByVolume = func(files []*hardLink, verbose bool) map[uint32]map[uint64][]*hardLink {
		volumes := map[uint32]map[uint64][]*hardLink{}

		for _, file := range files {
			data, err := GetFileInformationByPath(file.path)
			if err != nil {
				if verbose {
					logger.Eprintln(err)
				}
				continue
			}

			volume := volumes[data.VolumeSerialNumber]
			if volume == nil {
				volume = map[uint64][]*hardLink{}
				volumes[data.VolumeSerialNumber] = volume
			}

			index := uint64(data.FileIndexHigh)<<32 | uint64(data.FileIndexLow)
			file.index = index
			volume[index] = append(volume[index], file)
		}

		return volumes
	}
}
