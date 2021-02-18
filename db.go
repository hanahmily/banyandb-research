package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/dustin/go-humanize"
	"github.com/golang/snappy"

	"github.com/hanahmily/banyandb-research/api"
)

type DB struct {
	db *badger.DB
	path string
	writtenKeySize uint64
	writtenValueSize uint64
}


type mockLogger struct {
	output string
}

func (l *mockLogger) Errorf(f string, v ...interface{}) {
	l.output = fmt.Sprintf("ERROR: "+f, v...)
}

func (l *mockLogger) Infof(f string, v ...interface{}) {
	l.output = fmt.Sprintf("INFO: "+f, v...)
}

func (l *mockLogger) Warningf(f string, v ...interface{}) {
	l.output = fmt.Sprintf("WARNING: "+f, v...)
}

func (l *mockLogger) Debugf(f string, v ...interface{}) {
	l.output = fmt.Sprintf("DEBUG: "+f, v...)
}

func NewDB(blockSize int) DB {
	path, err := ioutil.TempDir("", "banyandb")
	if err != nil {
		log.Fatalf("failed to create tmp dir: %v", err)
	}
	log.Printf("database dir:%s", path)
	opts := badger.DefaultOptions(path)
	opts.BlockSize = blockSize * 1024
	//opts.Compression = options.None
	opts.CompactL0OnClose = true
	opts.BaseTableSize = 20 << 20
	opts.Logger = &mockLogger{}

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("failed to open badger database: %v", err)
	}
	return DB{db:db, path: path}
}

func (db *DB) Write(key, val []byte) {
	db.writtenKeySize = db.writtenKeySize + uint64(len(key))
	if val != nil {
		db.writtenValueSize = db.writtenValueSize + uint64(len(val))
		val = snappy.Encode(nil, val)
	}
	//compressedSpans := make([]byte, len(val))
	//l, err := lz4.CompressBlock(val, compressedSpans, nil)
	//if err != nil {
	//	panic(err)
	//}
	
	err := db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, val)
	})
	if err != nil {
		log.Fatalf("failed to write: %v", err)
	}
}

func (db *DB) Close() {
	_ = db.db.Close()

	log.Printf("uncompressed key size: %s", humanize.Bytes(db.writtenKeySize))
	log.Printf("uncompressed vaule size: %s", humanize.Bytes(db.writtenValueSize))
	files, err := ioutil.ReadDir(db.path)
	if err != nil {
		log.Fatal(err)
	}

	var sstSizeOnDisk uint64
	var vlogSizeOnDisk uint64
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "sst") {
			log.Printf("%s: %s", f.Name(), humanize.Bytes(uint64(f.Size())))
			sstSizeOnDisk = sstSizeOnDisk + uint64(f.Size())
		}
		if strings.HasSuffix(f.Name(), "vlog") {
			log.Printf("%s: %s", f.Name(), humanize.Bytes(uint64(f.Size())))
			vlogSizeOnDisk = vlogSizeOnDisk + uint64(f.Size())
		}
	}
	log.Printf("key uncompress:%s compressed:%s ratio:%f%%", humanize.Bytes(db.writtenKeySize), 
		humanize.Bytes(sstSizeOnDisk), float32((db.writtenKeySize - sstSizeOnDisk) * 100) / float32(db.writtenKeySize))
	log.Printf("value uncompress:%s compressed:%s ratio:%f%%", humanize.Bytes(db.writtenValueSize),
		humanize.Bytes(vlogSizeOnDisk), float32((db.writtenValueSize - vlogSizeOnDisk) * 100.0) / float32(db.writtenValueSize))
}

type ValueExtractor func(key, val []byte)

func (db *DB) Read(prefix string, keyOnly bool, extractor ValueExtractor) {
	err := db.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = api.SegmentVariants
		if keyOnly {
			opts.PrefetchValues = false
		}
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if keyOnly {
				if extractor != nil {
					extractor(k, nil)
				}
				continue
			}
			err := item.Value(func(v []byte) error {
				vv, _ := snappy.Decode(nil, v)
				if extractor != nil {
					extractor(k, vv)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed to get: %v", err)
	}
}

