/*
 * FrankyGo - UCI chess engine in GO for learning purposes
 *
 * MIT License
 *
 * Copyright (c) 2018-2020 Frank Kopp
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const debug = false

// ResolveFile is resolving a path to a file and try to find the file
// in a specific set of places and then will return an absolute path to
// it.
// Path needs to be a file or a not found error will be returned.
// The order will be check like this:
//  - if path is absolute it will return a os specific path and
//    an error if the file does not exist
// 	- if path is not absolute we will try first
// 	  - relative to working directory
//	  - relative to executable
//    - relative to user home directory
func ResolveFile(file string) (string, error) {

	fileNotFoundErr := errors.New(fmt.Sprintf("File could not be found: %s", file))

	file = filepath.Clean(file)

	if debug {
		log.Println("Searching folder", file)
	}

	// file is a absolute path
	if filepath.IsAbs(file) {
		if fileExists(file) {
			return file, nil
		}
		return file, fileNotFoundErr
	}

	// file is a relative path
	dir, err := os.Getwd()
	if err == nil {
		if fileExists(filepath.Join(dir, file)) {
			if debug {
				log.Println("Found relative to CWD")
			}
			return filepath.Clean(filepath.Join(dir, file)), nil
		}
	}

	dir, err = os.Executable()
	if err == nil {
		if fileExists(filepath.Join(dir, file)) {
			if debug {
				log.Println("Found relative to EXE")
			}
			return filepath.Clean(filepath.Join(dir, file)), nil
		}
	}

	dir, err = os.UserHomeDir()
	if err == nil {
		if fileExists(filepath.Join(dir, file)) {
			if debug {
				log.Println("Found relative to USER HOME")
			}
			return filepath.Clean(filepath.Join(dir, file)), nil
		}
	}

	if debug {
		log.Println("File not found", file)
	}
	return file, fileNotFoundErr
}

// ResolveFolder is resolving a path to a folder and try to find the folder
// in a specific set of places and then will return an absolute path to
// it.
// Path needs to be a folder or a not found error will be returned.
// The folder will not be created.
// The order will be check like this:
//  - if path is absolute it will return a os specific path and
//    an error if the file does not exist
// 	- if path is not absolute we will try first
// 	  - relative to working directory
//	  - relative to executable
//    - relative to user home directory
func ResolveFolder(folder string) (string, error) {

	folderNotFoundErr := errors.New(fmt.Sprintf("Folder could not be found: %s", folder))

	folder = filepath.Clean(folder)

	if debug {
		log.Println("Searching folder", folder)
	}

	// folder is a absolute path
	if filepath.IsAbs(folder) {
		if fileExists(folder) {
			return folder, nil
		}
		return folder, folderNotFoundErr
	}

	// folder is a relative path
	dir, err := os.Getwd()
	if err == nil {
		if folderExists(filepath.Join(dir, folder)) {
			if debug {
				log.Println("Found relative to CWD")
			}
			return filepath.Clean(filepath.Join(dir, folder)), nil
		}
	}

	dir, err = os.Executable()
	if err == nil {
		if folderExists(filepath.Join(dir, folder)) {
			if debug {
				log.Println("Found relative to EXE")
			}
			return filepath.Clean(filepath.Join(dir, folder)), nil
		}
	}

	log.Println("Testing home")
	dir, err = os.UserHomeDir()
	if folderExists(filepath.Join(dir, folder)) {
		if debug {
			log.Println("Found relative to USER HOME")
		}
		return filepath.Clean(filepath.Join(dir, folder)), nil
	}

	if debug {
		log.Println("Folder not found", folder)
	}
	return folder, folderNotFoundErr
}

// ResolveCreateFolder is resolving a path to a folder and try to find the
// folder in a specific set of places.
// If no folder can be found it will try to create a folder from the last
// part of the given folder path in the working directory.
// If a folder can't be created in the working directory it will be created
// in the os's temp directory.
// The order will be check like this:
//  - if path is absolute it will test if the folder exists or try to create
//    the folder there
// 	- if path is not absolute we will try finding the folder relative to
// 	  working directory. If not found we try to create the folder there
//  - last we create the folder in temp
func ResolveCreateFolder(folderPath string) (string, error) {
	folderPath = filepath.Clean(folderPath)

	// file is a absolute path
	if filepath.IsAbs(folderPath) {
		if folderExists(folderPath) {
			return folderPath, nil
		}
		errDir := os.Mkdir(folderPath, 0755)
		return folderPath, errDir
	}

	// file is a relative path

	// try working directory
	dir, _ := os.Getwd()
	folderPath = filepath.Join(dir, filepath.Base(folderPath))
	if folderExists(folderPath) {
		return folderPath, nil
	}
	errDir := os.Mkdir(folderPath, 0755)
	if errDir == nil {
		return folderPath, nil
	}

	// try temp
	dir = os.TempDir()
	folderPath = filepath.Join(dir, filepath.Base(folderPath))
	if folderExists(folderPath) {
		return folderPath, nil
	}
	errDir = os.Mkdir(folderPath, 0755)
	return folderPath, errDir
}

func fileExists(filename string) bool {
	// log.Println("File", filename)
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			if debug {
				log.Println("Error trying to find file", filename)
				log.Println("Stat", info, "Err: ", err)
			}
			return false
		}
	}
	if info == nil {
		if debug {
			log.Println("Stat for file is NIL when trying to find file", filename)
			log.Println("Stat", info, "Err: ", err)
		}
		return false
	}
	// log.Println("Info", info)
	// log.Println("File", info.Mode().IsRegular())
	return info.Mode().IsRegular()
}

func folderExists(foldername string) bool {
	// log.Println("Folder", foldername)
	info, err := os.Stat(foldername)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			if debug {
				log.Println("Error trying to find folder", foldername)
				log.Println("Stat", info, "Err: ", err)
			}
			return false
		}
	}
	if info == nil {
		if debug {
			log.Println("Stat for folder is NIL when trying to find folder", foldername)
			log.Println("Stat", info, "Err: ", err)
		}
		return false
	}
	// log.Println("Info", info)
	// log.Println("Dir", info.Mode().IsDir())
	return info.Mode().IsDir()
}

// 	tmp, err := os.TempDir()
//	userConfig, err := os.UserConfigDir()
//	userConfig, err := os.UserCacheDir()
