package config

type HTTPEnvironment struct {
	ID        int64
	Variables []*HTTPVariable
}

type HTTPVariable struct {
	ID    int64
	Name  string
	Value string
}

type HTTPCollection struct {
	ID    int64
	Name  string
	Items []*HTTPItem
	Count int64
}

func (i *HTTPCollection) Save() error {
	return SaveHTTPCollection(i)
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
	return SaveHTTPItem(i)
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
	ID      int64
	Key     string
	Value   string
	Type    string // header | param
	Sort    int
	Enabled bool

	HTTPItemID int64
}

func (i *HTTPKV) Save() error {
	return SaveHTTPKV(i)
}
