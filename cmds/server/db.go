package server

// MapStateDB ...
type MapStateDB struct {
	db       map[string][]byte
	cachePut map[string]map[string][]byte
	cacheDel map[string]map[string]struct{}
}

// NewDB ...
func NewDB() *MapStateDB {
	return &MapStateDB{
		db:       make(map[string][]byte),
		cachePut: make(map[string]map[string][]byte),
		cacheDel: make(map[string]map[string]struct{}),
	}
}

// Put PutState
func (db *MapStateDB) Put(txid string, key string, value []byte) error {
	if cdb, ok := db.cachePut[txid]; ok {
		cdb[key] = value
	} else {
		cdb := make(map[string][]byte)
		cdb[key] = value
		db.cachePut[txid] = cdb
	}
	return nil
}

// Get GetState
func (db *MapStateDB) Get(txid string, key string) ([]byte, error) {
	return db.db[key], nil
}

// Del DelState
func (db *MapStateDB) Del(txid string, key string) error {
	if cdb, ok := db.cacheDel[txid]; ok {
		cdb[key] = struct{}{}
	} else {
		cdb := make(map[string]struct{})
		cdb[key] = struct{}{}
		db.cacheDel[txid] = cdb
	}
	return nil
}

// Commit 交易成功入库
func (db *MapStateDB) Commit(txid string) error {
	if pc, ok := db.cachePut[txid]; ok {
		for k, v := range pc {
			db.db[k] = v
		}
		delete(db.cachePut, txid)
	}
	if dc, ok := db.cacheDel[txid]; ok {
		for k := range dc {
			delete(db.db, k)
		}
		delete(db.cacheDel, txid)
	}
	return nil
}

// Clear 交易失败清理缓存
func (db *MapStateDB) Clear(txid string) error {
	delete(db.cachePut, txid)
	delete(db.cacheDel, txid)
	return nil
}
