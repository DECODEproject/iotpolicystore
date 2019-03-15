// Code generated by go-bindata. DO NOT EDIT.
// sources:
// sql/20181009155147_add_policy_table.down.sql (38B)
// sql/20181009155147_add_policy_table.up.sql (386B)
// sql/20190306110030_add_cert_cache_table.down.sql (34B)
// sql/20190306110030_add_cert_cache_table.up.sql (106B)
// sql/20190308154458_add_coconut_fields.down.sql (107B)
// sql/20190308154458_add_coconut_fields.up.sql (134B)
// sql/20190315114713_add_uuid_column.down.sql (127B)
// sql/20190315114713_add_uuid_column.up.sql (175B)

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

	info := bindataFileInfo{name: "20181009155147_add_policy_table.down.sql", size: 38, mode: os.FileMode(420), modTime: time.Unix(1541168694, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xca, 0xcd, 0x26, 0xf6, 0x69, 0xe7, 0x14, 0x6, 0x4f, 0xc6, 0x27, 0x96, 0x99, 0x7d, 0x34, 0xb0, 0xcb, 0xd0, 0xb4, 0xb2, 0x61, 0x8c, 0x2a, 0xf7, 0xb2, 0xb1, 0x53, 0xaa, 0x6e, 0xe9, 0x29, 0x68}}
	return a, nil
}

var __20181009155147_add_policy_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x90\x51\x6b\x83\x40\x10\x84\xdf\xef\x57\xcc\x5b\x14\xfa\x0f\x52\x0a\xa7\xd9\x90\x6b\xf4\x4c\xbd\x95\x68\x5f\xc4\xe8\x51\x8e\x48\x14\x63\xa1\xf9\xf7\x45\x5b\x6a\x28\xed\xe3\xee\x7c\x33\xcb\x4e\x98\x92\x64\x02\xe5\x4c\xda\xa8\x44\x43\x6d\xa1\x13\x06\xe5\xca\xb0\x41\xff\x56\x0f\xb7\x7e\xec\xd6\x42\x7c\x93\x2c\x83\x88\x7e\x53\x5d\xeb\x6a\x67\xaf\xf0\x04\xe0\x1a\x18\x4a\x95\x8c\x70\x48\x55\x2c\xd3\x02\x7b\x2a\x1e\x04\xd0\xbf\x9f\x5a\x57\x97\x67\x7b\x03\x53\xce\x73\x80\xce\xa2\x08\xe1\x8e\xc2\x3d\xbc\x3b\xfd\xf1\x09\xab\x95\x3f\x99\xda\xea\x64\xdb\xbf\xf9\x2f\x69\x41\xeb\xc1\x56\xa3\x6d\xca\x6a\x04\xab\x98\x0c\xcb\xf8\x80\xa3\xe2\xdd\x3c\xe2\x35\xd1\x84\x0d\x6d\x65\x16\x4d\x59\x47\x6f\x36\x75\xbd\x1d\xaa\xd1\x75\x97\x2b\x9e\x4d\xa2\x83\x69\x37\x76\x67\x7b\x41\x50\x30\xc9\x9f\xa3\xc2\x5f\x2a\xc8\xb4\x7a\xc9\x08\x4a\x6f\x28\xff\xa7\x89\x72\x79\xa6\x74\xcd\x87\x00\x12\x7d\x57\xd3\xa2\xfa\xeb\xcf\x00\x00\x00\xff\xff\xd5\xd7\x33\xdd\x82\x01\x00\x00")

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

	info := bindataFileInfo{name: "20181009155147_add_policy_table.up.sql", size: 386, mode: os.FileMode(420), modTime: time.Unix(1552653817, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xb2, 0x31, 0xf0, 0x42, 0x61, 0xc8, 0xbc, 0x89, 0xc5, 0x29, 0xe6, 0xb8, 0xbb, 0x59, 0xd3, 0x2c, 0x76, 0xe4, 0xa7, 0xfa, 0x71, 0x16, 0xa, 0x36, 0xa8, 0xc7, 0xbc, 0x9f, 0x6c, 0x57, 0x24, 0x50}}
	return a, nil
}

var __20190306110030_add_cert_cache_tableDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4e\x2d\x2a\xc9\x4c\xcb\x4c\x4e\x2c\x49\x2d\xb6\x06\x04\x00\x00\xff\xff\x9b\x6a\xf7\x60\x22\x00\x00\x00")

func _20190306110030_add_cert_cache_tableDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__20190306110030_add_cert_cache_tableDownSql,
		"20190306110030_add_cert_cache_table.down.sql",
	)
}

func _20190306110030_add_cert_cache_tableDownSql() (*asset, error) {
	bytes, err := _20190306110030_add_cert_cache_tableDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20190306110030_add_cert_cache_table.down.sql", size: 34, mode: os.FileMode(420), modTime: time.Unix(1551877755, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x3d, 0x85, 0xef, 0x15, 0xa1, 0x51, 0x74, 0x22, 0x6b, 0x2f, 0xde, 0x28, 0x99, 0xb5, 0x60, 0xd6, 0xe8, 0x10, 0x23, 0xa7, 0x48, 0x63, 0xf2, 0xc4, 0x3c, 0xca, 0x83, 0x1f, 0xb4, 0x65, 0xad, 0x98}}
	return a, nil
}

var __20190306110030_add_cert_cache_tableUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x0e\x72\x75\x0c\x71\x55\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\xf0\xf3\x0f\x51\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x4e\x2d\x2a\xc9\x4c\xcb\x4c\x4e\x2c\x49\x2d\x56\xd0\xe0\x52\x50\xc8\x4e\xad\x54\x08\x71\x8d\x08\x01\x2b\xf2\x0b\xf5\xf1\x51\x08\x08\xf2\xf4\x75\x0c\x8a\x54\xf0\x76\x8d\xd4\xe1\x52\x40\xd6\xa1\xe0\x14\x19\xe2\xea\x08\x57\xc9\xa5\x69\x0d\x08\x00\x00\xff\xff\x2d\x4d\xb2\x71\x6a\x00\x00\x00")

func _20190306110030_add_cert_cache_tableUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__20190306110030_add_cert_cache_tableUpSql,
		"20190306110030_add_cert_cache_table.up.sql",
	)
}

func _20190306110030_add_cert_cache_tableUpSql() (*asset, error) {
	bytes, err := _20190306110030_add_cert_cache_tableUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20190306110030_add_cert_cache_table.up.sql", size: 106, mode: os.FileMode(420), modTime: time.Unix(1551877755, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x66, 0x3c, 0x3, 0x6a, 0x8c, 0x5a, 0x0, 0xe2, 0xca, 0x24, 0x4b, 0xf0, 0x4b, 0x55, 0xb2, 0xc4, 0x3f, 0x19, 0x75, 0x20, 0x4f, 0xd3, 0x4d, 0xc6, 0xa6, 0x9b, 0xbb, 0xc1, 0x94, 0x70, 0xbc, 0x38}}
	return a, nil
}

var __20190308154458_add_coconut_fieldsDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\xca\x31\x0e\xc2\x30\x0c\x05\xd0\x9d\x53\xf8\x00\xdc\x80\xa9\x40\xb7\x40\x51\x55\x66\x2b\x6d\x2c\xf1\x25\x2b\x89\x1c\x7b\xe1\xf4\xec\xec\x6f\x4a\xdb\xbc\xd2\x36\x5d\xd3\x4c\xbd\x29\x0e\xc8\x38\x11\xdd\xd7\xe5\x45\xb7\x25\xbd\x1f\x4f\xca\xe1\x9f\x66\xf8\xe6\x5d\x85\xb3\xbb\x61\x0f\x17\x46\x39\xff\xc1\xc3\xa4\x48\x75\x64\x65\x8c\x11\x62\x2c\xb5\xf4\x86\xea\x1c\xa6\x97\x5f\x00\x00\x00\xff\xff\x2c\xbe\x2a\xac\x6b\x00\x00\x00")

func _20190308154458_add_coconut_fieldsDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__20190308154458_add_coconut_fieldsDownSql,
		"20190308154458_add_coconut_fields.down.sql",
	)
}

func _20190308154458_add_coconut_fieldsDownSql() (*asset, error) {
	bytes, err := _20190308154458_add_coconut_fieldsDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20190308154458_add_coconut_fields.down.sql", size: 107, mode: os.FileMode(420), modTime: time.Unix(1552444813, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd2, 0x30, 0x7b, 0x66, 0x8c, 0x36, 0x59, 0x67, 0x70, 0x5b, 0x52, 0x71, 0x53, 0x0, 0x68, 0x82, 0x3, 0xd8, 0x93, 0xe5, 0x55, 0x36, 0xee, 0x23, 0xfb, 0x74, 0x2e, 0x64, 0x9b, 0x3, 0x43, 0xba}}
	return a, nil
}

var __20190308154458_add_coconut_fieldsUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\xcb\x31\x0a\xc3\x30\x0c\x00\xc0\x3d\xaf\xd0\x03\xfa\x83\x4e\x69\xe3\x4d\x75\xa0\x28\xd0\x4d\x38\xb1\xa0\x02\x63\x07\x59\x5e\xfa\xfa\xce\x79\xc0\xcd\x48\xe1\x0d\x34\x3f\x30\xc0\xd9\x8a\x1e\x2a\x7d\x02\x98\x97\x05\x9e\x2b\x6e\xaf\x08\x69\xf8\xb7\x99\xfe\xd2\x5e\x84\x93\xbb\xe9\x3e\x5c\x58\x33\x50\xf8\x10\xc4\x95\x20\x6e\x88\xb7\xab\x3a\x4c\xb2\x54\xd7\x54\x58\x7b\x1f\x62\x2c\x35\x9f\x4d\xab\xf3\xb0\x72\xa5\xf7\xe9\x1f\x00\x00\xff\xff\x10\x6f\x78\xad\x86\x00\x00\x00")

func _20190308154458_add_coconut_fieldsUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__20190308154458_add_coconut_fieldsUpSql,
		"20190308154458_add_coconut_fields.up.sql",
	)
}

func _20190308154458_add_coconut_fieldsUpSql() (*asset, error) {
	bytes, err := _20190308154458_add_coconut_fieldsUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20190308154458_add_coconut_fields.up.sql", size: 134, mode: os.FileMode(420), modTime: time.Unix(1552444813, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x8f, 0x6f, 0xd2, 0x93, 0x5a, 0xc7, 0xd9, 0x97, 0x64, 0xd8, 0x4e, 0xd9, 0x4b, 0x47, 0xd8, 0x5d, 0xe2, 0x3c, 0x58, 0x22, 0x68, 0x73, 0xe7, 0xb3, 0x10, 0x20, 0x98, 0xc7, 0xb1, 0x8f, 0x4a, 0x61}}
	return a, nil
}

var __20190315114713_add_uuid_columnDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\xcb\xb1\x0a\xc2\x30\x10\x06\xe0\x3d\x4f\xf1\x8f\xfa\x0c\x9d\x62\x7b\x42\x20\xde\x69\x7a\x81\x6e\x01\xdb\x0e\x87\x05\x0b\x12\xd0\xb7\x77\x4b\xe7\x8f\xcf\x47\xa5\x04\xf5\x97\x48\xd8\xdf\x9b\xcd\xb6\x7e\x1c\x30\x24\xb9\xa3\x97\x98\x6f\x8c\x5a\x6d\xe9\x9c\xeb\x13\x79\x25\x64\x0e\x8f\x4c\x08\x3c\xd0\x84\x70\x05\x8b\x82\xa6\x30\xea\xd8\x7e\xd9\xeb\x73\xb3\xb9\xbc\xd6\x5f\xb1\xe5\xeb\x00\xe1\x86\x38\x1d\x7a\xee\xfe\x01\x00\x00\xff\xff\xe6\x25\xe8\x45\x7f\x00\x00\x00")

func _20190315114713_add_uuid_columnDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__20190315114713_add_uuid_columnDownSql,
		"20190315114713_add_uuid_column.down.sql",
	)
}

func _20190315114713_add_uuid_columnDownSql() (*asset, error) {
	bytes, err := _20190315114713_add_uuid_columnDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20190315114713_add_uuid_column.down.sql", size: 127, mode: os.FileMode(420), modTime: time.Unix(1552650592, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xcf, 0x8f, 0x8b, 0xe8, 0x50, 0x44, 0xad, 0xed, 0x47, 0x4e, 0xf1, 0x54, 0xd3, 0x3c, 0xaf, 0xaf, 0x5b, 0xee, 0x9f, 0xd9, 0x9f, 0xda, 0x1a, 0x6c, 0x6b, 0x15, 0x3d, 0x5b, 0xce, 0xcf, 0xc8, 0x52}}
	return a, nil
}

var __20190315114713_add_uuid_columnUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\xcd\xb1\x0a\xc2\x30\x14\x85\xe1\xfd\x3e\xc5\x19\xf5\x19\x3a\xc5\xe6\x0a\x81\x78\xa3\x69\x2e\x74\x0b\xd8\x76\x08\x16\x2c\x48\x40\xdf\x5e\xb2\x54\x70\xfd\xe1\x3b\xc7\xf8\xc4\x11\xc9\x9c\x3c\x63\x7b\xae\x65\x2a\xcb\x8b\x00\x63\x2d\xfa\xe0\xf5\x22\xa8\xb5\xcc\x50\x75\x16\x12\x12\x44\xbd\xef\x88\x6c\x0c\x57\x38\xb1\x3c\xc2\x9d\xc1\xa3\x1b\xd2\xb0\xfb\xbc\xd5\xfb\x5a\xa6\xfc\x58\x3e\xb9\xcc\xef\x8e\xa8\x8f\x6c\x12\x43\xc5\xdd\x94\x7f\xae\x0d\xfe\xdb\x76\xd7\x14\x01\x41\xf6\x8c\x43\xeb\xc7\xee\x1b\x00\x00\xff\xff\x90\x02\xc8\x7f\xaf\x00\x00\x00")

func _20190315114713_add_uuid_columnUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__20190315114713_add_uuid_columnUpSql,
		"20190315114713_add_uuid_column.up.sql",
	)
}

func _20190315114713_add_uuid_columnUpSql() (*asset, error) {
	bytes, err := _20190315114713_add_uuid_columnUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "20190315114713_add_uuid_column.up.sql", size: 175, mode: os.FileMode(420), modTime: time.Unix(1552650540, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xdf, 0x33, 0x53, 0xb0, 0xfd, 0xc5, 0x0, 0xfa, 0xe7, 0x89, 0x52, 0x83, 0x19, 0x82, 0x32, 0x67, 0x38, 0x63, 0xb6, 0x9b, 0xbb, 0x59, 0xd8, 0x75, 0x13, 0x41, 0x3e, 0xee, 0xa7, 0x34, 0xb8, 0xac}}
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

	"20190306110030_add_cert_cache_table.down.sql": _20190306110030_add_cert_cache_tableDownSql,

	"20190306110030_add_cert_cache_table.up.sql": _20190306110030_add_cert_cache_tableUpSql,

	"20190308154458_add_coconut_fields.down.sql": _20190308154458_add_coconut_fieldsDownSql,

	"20190308154458_add_coconut_fields.up.sql": _20190308154458_add_coconut_fieldsUpSql,

	"20190315114713_add_uuid_column.down.sql": _20190315114713_add_uuid_columnDownSql,

	"20190315114713_add_uuid_column.up.sql": _20190315114713_add_uuid_columnUpSql,
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
	"20181009155147_add_policy_table.down.sql":     &bintree{_20181009155147_add_policy_tableDownSql, map[string]*bintree{}},
	"20181009155147_add_policy_table.up.sql":       &bintree{_20181009155147_add_policy_tableUpSql, map[string]*bintree{}},
	"20190306110030_add_cert_cache_table.down.sql": &bintree{_20190306110030_add_cert_cache_tableDownSql, map[string]*bintree{}},
	"20190306110030_add_cert_cache_table.up.sql":   &bintree{_20190306110030_add_cert_cache_tableUpSql, map[string]*bintree{}},
	"20190308154458_add_coconut_fields.down.sql":   &bintree{_20190308154458_add_coconut_fieldsDownSql, map[string]*bintree{}},
	"20190308154458_add_coconut_fields.up.sql":     &bintree{_20190308154458_add_coconut_fieldsUpSql, map[string]*bintree{}},
	"20190315114713_add_uuid_column.down.sql":      &bintree{_20190315114713_add_uuid_columnDownSql, map[string]*bintree{}},
	"20190315114713_add_uuid_column.up.sql":        &bintree{_20190315114713_add_uuid_columnUpSql, map[string]*bintree{}},
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
