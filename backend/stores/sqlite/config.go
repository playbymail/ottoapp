// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

func (db *DB) GetKeyValue(key string) (value string, err error) {
	value, err = db.q.GetConfigKeyValue(db.ctx, key)
	if err != nil {
		value = ""
	}
	return value, err
}
