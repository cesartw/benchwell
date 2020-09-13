package config

import (
	"fmt"
)

type HTTPCollection struct {
	ID    int64
	Name  string
	Items []*HTTPItem
	Count int64
}

func (i *HTTPCollection) Save() error {
	if i.ID == 0 {
		sql := `INSERT INTO http_collections(name)
				VALUES(?)`
		result, err := db.Exec(sql, i.Name)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		i.ID = id
	} else {
		sql := `UPDATE http_collections
					SET name = ?
				WHERE ID = ?`
		_, err := db.Exec(sql,
			i.Name, i.ID)
		return err
	}

	return nil
}

func (c *HTTPCollection) LoadRootItems() error {
	c.Items = nil
	query := `SELECT id, name, is_folder, sort, http_collections_id, method
			FROM http_items
			WHERE http_collections_id = ? AND (parent_id IS NULL OR parent_id = 0)
			ORDER BY sort ASC`
	rows, err := db.Query(query, c.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		item := &HTTPItem{}
		err := rows.Scan(&item.ID, &item.Name, &item.IsFolder, &item.Sort, &item.HTTPCollectionID, &item.Method)
		if err != nil {
			return err
		}
		c.Items = append(c.Items, item)
	}

	return nil
}

func (i *HTTPItem) LoadFull() error {
	if i.Loaded {
		return nil
	}

	if i.IsFolder {
		i.Items = nil
		query := `SELECT id, name, parent_id, is_folder, sort, http_collections_id, method
				FROM http_items
				WHERE http_collections_id = ? AND parent_id = ?
				ORDER BY sort ASC`
		rows, err := db.Query(query, i.HTTPCollectionID, i.ID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			item := &HTTPItem{}
			err := rows.Scan(&item.ID, &item.Name, &item.ParentID,
				&item.IsFolder, &item.Sort, &item.HTTPCollectionID, &item.Method)
			if err != nil {
				return err
			}
			i.Items = append(i.Items, item)
		}

		return nil
	}

	query := `SELECT ifnull(method,""), ifnull(url,""), ifnull(body, ""), ifnull(mime,"")
				FROM http_items
				WHERE id = ?`
	rows, err := db.Query(query, i.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&i.Method, &i.URL, &i.Body, &i.Mime)
		if err != nil {
			return err
		}
		break
	}

	i.Params = nil
	i.Headers = nil
	query = `SELECT id, ifnull(key,""), ifnull(value,""), type, sort
				FROM http_kvs
				WHERE http_items_id = ?
				ORDER BY sort ASC`
	rows, err = db.Query(query, i.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		kv := &HTTPKV{HTTPItemID: i.ID}
		err := rows.Scan(&kv.ID, &kv.Key, &kv.Value, &kv.Type, &kv.Sort)
		if err != nil {
			return err
		}
		if kv.Type == "param" {
			i.Params = append(i.Params, kv)
		}
		if kv.Type == "header" {
			i.Headers = append(i.Headers, kv)
		}
	}

	return nil
}

type HTTPItem struct {
	ID          int64
	ParentID    int64
	Name        string
	Description string
	// Not pretty but makes little sense
	// to separate them just for normalization sake
	IsFolder         bool
	HTTPCollectionID int64

	Items []*HTTPItem
	Sort  int
	HTTPRequest

	Loaded bool
}

func (i *HTTPItem) UIName() string {
	if i.IsFolder {
		return i.Name
	}
	return fmt.Sprintf("% -6s %s", i.Method, i.Name)
}

func (i *HTTPItem) SearchID(id int64) *HTTPItem {
	if i.ID == id {
		return i
	}

	for _, ii := range i.Items {
		if ii.ID == id {
			return ii
		}
	}

	return nil
}

func (i *HTTPItem) Save() error {
	if i.ID == 0 {
		sql := `INSERT INTO http_is(name, description, parent_id, is_folder, sort, http_collections_id, external_data, method, url, body, mime)
				VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := db.Exec(sql, i.Name,
			i.Description, i.ParentID, i.IsFolder,
			i.Sort, i.HTTPCollectionID, "", i.Method,
			i.URL, i.Body, i.Mime)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		i.ID = id
	} else {
		sql := `UPDATE http_is
					SET name = ?, description = ?, parent_id = ?, is_folder = ?,
					sort = ?, http_collections_id = ?, external_data = ?,
					method = ?, url = ?, body = ?, mime = ?
				WHERE ID = ?`
		_, err := db.Exec(sql,
			i.Name, i.Description, i.ParentID,
			i.IsFolder, i.Sort, i.HTTPCollectionID,
			"", i.Method, i.URL, i.Body, i.Mime,
			i.ID)
		if err != nil {
			return err
		}
	}

	for _, kv := range i.Params {
		kv.HTTPItemID = i.ID
		err := kv.Save()
		if err != nil {
			return err
		}
	}
	for _, kv := range i.Headers {
		kv.HTTPItemID = i.ID
		err := kv.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *HTTPItem) Delete() error {
	if i.ID == 0 {
		return nil
	}

	sql := `delete from http_items where id = ?`

	_, err := db.Exec(sql, i.ID)
	if err != nil {
		return err
	}

	for _, kv := range i.Params {
		err := kv.Delete()
		if err != nil {
			return err
		}
	}

	for _, kv := range i.Headers {
		err := kv.Delete()
		if err != nil {
			return err
		}
	}

	for _, subitem := range i.Items {
		err = subitem.Delete()
		if err != nil {
			return err
		}
	}

	return err
}

type HTTPRequest struct {
	Method  string
	URL     string
	Body    string
	Mime    string
	Headers []*HTTPKV
	Params  []*HTTPKV
}

type HTTPKV struct {
	Var
	Type string // header | param
	Sort int

	HTTPItemID int64
}

func (i *HTTPKV) Save() error {
	if i.ID == 0 {
		sql := `INSERT INTO http_kvs(key, value,  type, sort, enabled, http_items_id)
				VALUES(?, ?, ?, ?, ?, ?)`
		result, err := db.Exec(sql, i.Key, i.Value, i.Type, i.Sort, i.Enabled, i.HTTPItemID)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		i.ID = id
	} else {
		sql := `UPDATE http_kvs
				SET key = ?, value = ?, type = ?, sort = ?, enabled = ?, http_items_id = ?
				WHERE id = ?`
		_, err := db.Exec(sql, i.Key, i.Value, i.Type, i.Sort, i.Enabled, i.HTTPItemID, i.ID)
		return err
	}

	return nil
}

func (i *HTTPKV) Delete() error {
	if i.ID == 0 {
		return nil
	}

	sql := `delete from http_kvs where id = ?`
	_, err := db.Exec(sql, i.ID)

	return err
}
