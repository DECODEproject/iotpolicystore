// Code generated by go-bindata. DO NOT EDIT.
// sources:
// sql/20181009155147_add_policy_table.down.sql (38B)
// sql/20181009155147_add_policy_table.up.sql (341B)

package migrations

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __20181009155147_add_policy_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\xc8\xcf\xc9\x4c\xce\x4c\x2d\x56\x70\x76\x0c\x76\x76\x74\x71\xb5\x06\x04\x00\x00\xff\xff\x7a\x0e\x77\x6e\x26\x00\x00\x00")

func _20181009155147_add_policy_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__20181009155147_add_policy_tableDownSql,
		"20181009155147_add_policy_table.down.sql",
	)
}

func _20181009155147_add_policy_tableDownSql() (*asset, error) {
	bytes, err := _20181009155147_add_policy_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20181009155147_add_policy_table.down.sql", size: 38, mode: os.FileMode(420), modTime: time.Unix(1539100797, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xca, 0xcd, 0x26, 0xf6, 0x69, 0xe7, 0x14, 0x6, 0x4f, 0xc6, 0x27, 0x96, 0x99, 0x7d, 0x34, 0xb0, 0xcb, 0xd0, 0xb4, 0xb2, 0x61, 0x8c, 0x2a, 0xf7, 0xb2, 0xb1, 0x53, 0xaa, 0x6e, 0xe9, 0x29, 0x68}}
	return a, nil
}

var __20181009155147_add_policy_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8f\xcd\x6a\xeb\x30\x10\x85\xf7\x7a\x8a\xb3\x8c\xe1\xbe\x41\x56\xf2\xcd\x84\xaa\x95\xa5\xd4\x1a\x13\xbb\x1b\xe3\xd8\xa2\x88\x98\xc8\x38\x2a\x34\x6f\x5f\x1c\x4a\x5d\x5a\xba\x9c\xf3\x33\x33\xdf\xff\x92\x24\x13\xa8\x66\x32\x4e\x59\x03\xb5\x87\xb1\x0c\xaa\x95\x63\x87\xe9\xb5\x9f\x6f\x53\x8a\x5b\x21\x3e\x93\x2c\x73\x4d\x3f\x53\x71\x0c\x7d\xf0\x57\x6c\x04\x10\x06\x38\x2a\x95\xd4\x38\x94\xaa\x90\x65\x83\x27\x6a\xfe\x09\x60\x7a\x3b\x8d\xa1\x6f\xcf\xfe\x06\xa6\x9a\xef\x0b\x4c\xa5\xf5\xe2\x8d\xdd\xc9\x8f\xbf\xe5\x7e\xf6\x5d\xf2\x43\xdb\x25\xb0\x2a\xc8\xb1\x2c\x0e\x38\x2a\x7e\xb8\x8f\x78\xb1\x86\xb0\xa3\xbd\xac\xf4\xd2\x3b\x6e\xb2\xa5\x14\x27\x3f\x77\x29\xc4\xcb\x15\x8f\xce\x9a\x7c\xd1\x52\x3c\xfb\x0b\xf2\x86\x49\x7e\x1d\x10\xd9\x4a\x55\x19\xf5\x5c\x11\x94\xd9\x51\xfd\x07\x5c\xbb\xfe\xdf\x86\xe1\x5d\x00\xd6\x7c\x23\x5f\xdd\x6c\xfb\x11\x00\x00\xff\xff\x83\xff\x4e\xe9\x55\x01\x00\x00")

func _20181009155147_add_policy_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__20181009155147_add_policy_tableUpSql,
		"20181009155147_add_policy_table.up.sql",
	)
}

func _20181009155147_add_policy_tableUpSql() (*asset, error) {
	bytes, err := _20181009155147_add_policy_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20181009155147_add_policy_table.up.sql", size: 341, mode: os.FileMode(420), modTime: time.Unix(1539188091, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd5, 0x3, 0x95, 0x32, 0x6f, 0xea, 0xfc, 0x6a, 0xb4, 0x63, 0x9b, 0x6f, 0x2e, 0xa0, 0xe, 0xca, 0x89, 0x6d, 0xfd, 0x4b, 0x21, 0x40, 0x68, 0x80, 0xdc, 0xb2, 0xfa, 0x21, 0x1a, 0x6f, 0x3f, 0x9a}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"20181009155147_add_policy_table.down.sql": _20181009155147_add_policy_tableDownSql,

	"20181009155147_add_policy_table.up.sql": _20181009155147_add_policy_tableUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"20181009155147_add_policy_table.down.sql": &bintree{_20181009155147_add_policy_tableDownSql, map[string]*bintree{}},
	"20181009155147_add_policy_table.up.sql":   &bintree{_20181009155147_add_policy_tableUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
