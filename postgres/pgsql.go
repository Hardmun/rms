package pgsql

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"rms/logs"
	"rms/utils"
	"strings"
	"sync"
)

type CollectionRemaining struct {
	Collection  string  `json:"collection"`
	Uuid        string  `json:"uuid"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Length      string  `json:"length"`
	Width       string  `json:"width"`
	Count       float32 `json:"count"`
}

type DetailRemaining struct {
	CollectionRemaining
	UuidDetails     string  `json:"uuidDetails"`
	CodeDetails     string  `json:"codeDetails"`
	Picture         string  `json:"picture"`
	Form            string  `json:"form"`
	Color           string  `json:"color"`
	Brand           string  `json:"brand"`
	Barcode         string  `json:"barcode"`
	CountBalance    float32 `json:"countBalance"`
	ReservedBalance float32 `json:"reservedBalance"`
}

type Images struct {
	LogoURL    string   `json:"logoURL"`
	ImagesURL  []string `json:"imagesURL"`
	Collection string   `json:"collection"`
	Picture    string   `json:"picture"`
	Form       string   `json:"form"`
	Color      string   `json:"color"`
	Brand      string   `json:"brand"`
}

type queryMap struct {
	sync.Map
}

func (q *queryMap) add(key, value string) {
	q.Store(key, value)
}

func (q *queryMap) get(key string) (string, bool) {
	value, ok := q.Load(key)
	if ok {
		return value.(string), true
	}

	return "", false
}

func (q *queryMap) clear() {
	q.Clear()
}

type collectionFilter queryMap

func (c *collectionFilter) add(key, value string) {
	c.Store(key, value)
}

func (c *collectionFilter) get(key string) (string, bool) {
	value, ok := c.Load(key)
	if ok {
		return value.(string), true
	}

	return "", false
}

func (c *collectionFilter) clear() {
	c.Clear()
}

func (c *collectionFilter) updateFilter() {
	utils.Params.MU.RLock()
	collection := utils.Params.Tables.Collection
	utils.Params.MU.RUnlock()

	c.clear()

	if v, ok := collection["products.collection"]; ok {
		c.add("collection", v)
	}
	if v, ok := collection["products.uuid"]; ok {
		c.add("uuid", v)
	}
	if v, ok := collection["products.code"]; ok {
		c.add("code", v)
	}
	if v, ok := collection["products.description"]; ok {
		c.add("description", v)
	}
	if v, ok := collection["products.length"]; ok {
		c.add("length", v)
	}
	if v, ok := collection["products.width"]; ok {
		c.add("width", v)
	}
}

func (_ *collectionFilter) applyFilter(queryText *string, r *http.Request) []any {
	var (
		values []any
	)
	filters := r.URL.Query()
	if len(filters) == 0 {
		return values
	}
	var whereClauses []string
	i := 1

	for key, value := range filters {
		if len(value) == 0 || value[0] == "" {
			continue
		}

		if k, ok := collFilter.get(key); ok {
			if key == "uuid" {
				k = fmt.Sprintf("encode(%s, 'hex')", k)
			}
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", k, i))
			values = append(values, value[0])
		}
		i++
	}

	if lenClosures, lenFilters := len(whereClauses), len(filters); lenClosures > 0 {
		*queryText = strings.ReplaceAll(*queryText, "true", strings.Join(whereClauses, " AND "))
	} else if lenFilters > lenClosures {
		*queryText = strings.ReplaceAll(*queryText, "true", "false")
	}

	return values
}

type imageFilter queryMap

func (i *imageFilter) add(key, value string) {
	i.Store(key, value)
}

func (i *imageFilter) get(key string) (string, bool) {
	value, ok := i.Load(key)
	if ok {
		return value.(string), true
	}

	return "", false
}

func (i *imageFilter) clear() {
	i.Clear()
}

func (i *imageFilter) updateFilter() {
	utils.Params.MU.RLock()
	images := utils.Params.Tables.Images
	utils.Params.MU.RUnlock()

	i.clear()

	if v, ok := images["img.imageID"]; ok {
		i.add("imageID", v)
	}
	if v, ok := images["img.collection"]; ok {
		i.add("collection", v)
	}
	if v, ok := images["img.picture"]; ok {
		i.add("picture", v)
	}
	if v, ok := images["img.form"]; ok {
		i.add("form", v)
	}
	if v, ok := images["img.color"]; ok {
		i.add("color", v)
	}
	if v, ok := images["brands.brandCode"]; ok {
		i.add("brandCode", v)
	}
	if v, ok := images["brands.brand"]; ok {
		i.add("brand", v)
	}
}

func (_ *imageFilter) applyFilter(queryText *string, r *http.Request) []any {
	var (
		values []any
	)
	filters := r.URL.Query()
	if len(filters) == 0 {
		return values
	}
	var whereClauses []string
	l := 1

	for key, value := range filters {
		if len(value) == 0 || value[0] == "" {
			continue
		}

		if k, ok := imgFilter.get(key); ok {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", k, l))
			values = append(values, value[0])
		}
		l++
	}

	if lenClosures, lenFilters := len(whereClauses), len(filters); lenFilters > lenClosures {
		*queryText = strings.Replace(*queryText, "true", "false", 1)
	} else if lenClosures > 0 {
		*queryText = strings.Replace(*queryText, "true", strings.Join(whereClauses, " AND "), 1)
	}

	return values
}

type detailFilter queryMap

func (d *detailFilter) add(key, value string) {
	d.Store(key, value)
}

func (d *detailFilter) get(key string) (string, bool) {
	value, ok := d.Load(key)
	if ok {
		return value.(string), true
	}

	return "", false
}

func (d *detailFilter) clear() {
	d.Clear()
}

func (d *detailFilter) updateFilter() {
	utils.Params.MU.RLock()
	details := utils.Params.Tables.Details
	utils.Params.MU.RUnlock()

	d.clear()

	if v, ok := details["products.collection"]; ok {
		d.add("collection", v)
	}
	if v, ok := details["products.uuid"]; ok {
		d.add("uuid", v)
	}
	if v, ok := details["products.code"]; ok {
		d.add("code", v)
	}
	if v, ok := details["products.description"]; ok {
		d.add("description", v)
	}
	if v, ok := details["products.length"]; ok {
		d.add("length", v)
	}
	if v, ok := details["products.width"]; ok {
		d.add("width", v)
	}
	if v, ok := details["details.uuidDetails"]; ok {
		d.add("uuidDetails", v)
	}
	if v, ok := details["details.codeDetails"]; ok {
		d.add("codeDetails", v)
	}
	if v, ok := details["details.color"]; ok {
		d.add("color", v)
	}
	if v, ok := details["details.picture"]; ok {
		d.add("picture", v)
	}
	if v, ok := details["details.form"]; ok {
		d.add("form", v)
	}
	if v, ok := details["details.barcode"]; ok {
		d.add("barcode", v)
	}
	if v, ok := details["brands.brandName"]; ok {
		d.add("brand", v)
	}
}

func (_ *detailFilter) applyFilter(queryText *string, r *http.Request) []any {
	var (
		values []any
	)
	filters := r.URL.Query()
	if len(filters) == 0 {
		return values
	}
	var whereClauses []string
	l := 1

	for key, value := range filters {
		if len(value) == 0 || value[0] == "" {
			continue
		}

		if k, ok := dtlFilter.get(key); ok {
			if key == "uuid" || key == "uuidDetails" {
				k = fmt.Sprintf("encode(%s, 'hex')", k)
			}
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", k, l))
			values = append(values, value[0])
		}
		l++
	}

	if lenClosures, lenFilters := len(whereClauses), len(filters); lenFilters > lenClosures {
		*queryText = strings.ReplaceAll(*queryText, "true", "false")
	} else if lenClosures > 0 {
		*queryText = strings.ReplaceAll(*queryText, "true", strings.Join(whereClauses, " AND "))
	}

	return values
}

//go:embed collection.sql images.sql details.sql
var embedFiles embed.FS
var (
	DbPG       = dbConnection()
	qMap       = getQueryMap()
	collFilter = getCollectionFilter()
	imgFilter  = getImageFilter()
	dtlFilter  = getDetailFilter()
)

func dbConnection() *sql.DB {
	p := utils.Params
	connStr := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=disable",
		p.LoginPGSQL, p.PasswordPGSQL, p.DatabasePGSQL, p.ServerPGSQL, p.PortPGSQL)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logs.ErrLog.Write("Error opening database connection:%v", err.Error())
	}
	return db
}

func getDbPG() *sql.DB {
	if err := DbPG.Ping(); err != nil {
		DbPG = dbConnection()
	}
	return DbPG
}

func getQueryMap() *queryMap {
	qM := queryMap{}

	return &qM
}

func getCollectionFilter() *collectionFilter {
	cF := collectionFilter{}
	cF.updateFilter()

	return &cF
}

func getImageFilter() *imageFilter {
	iF := imageFilter{}
	iF.updateFilter()
	return &iF
}

func getDetailFilter() *detailFilter {
	dF := detailFilter{}
	dF.updateFilter()
	return &dF
}

func UpdateSettings() {
	qMap.clear()
	collFilter.updateFilter()
	imgFilter.updateFilter()
	dtlFilter.updateFilter()
}

func collectionFieldsMapping(queryText *string) {
	params := utils.Params
	params.MU.RLock()
	for k, v := range params.Tables.Collection {
		*queryText = strings.ReplaceAll(*queryText, k, v)
	}
	params.MU.RUnlock()
}

func imageFieldsMapping(queryText *string) {
	params := utils.Params
	params.MU.RLock()
	for k, v := range params.Tables.Images {
		*queryText = strings.ReplaceAll(*queryText, k, v)
	}
	params.MU.RUnlock()
}

func detailsFieldsMapping(queryText *string) {
	params := utils.Params
	params.MU.RLock()
	for k, v := range params.Tables.Details {
		*queryText = strings.ReplaceAll(*queryText, k, v)
	}
	params.MU.RUnlock()
}

func GetCollectionRemaining(r *http.Request) ([]CollectionRemaining, error) {
	db := getDbPG()

	queryText, ok := qMap.get(r.URL.Path)
	if !ok {
		queryByte, err := embedFiles.ReadFile("collection.sql")
		if err != nil {
			return nil, err
		}
		queryText = string(queryByte)
		collectionFieldsMapping(&queryText)
		qMap.add(r.URL.Path, queryText)
	}

	args := collFilter.applyFilter(&queryText, r)

	rows, errRows := db.Query(queryText, args...)
	if errRows != nil {
		return nil, errRows
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logs.ErrLog.Write(err)
		}
	}()

	var results []CollectionRemaining

	for rows.Next() {
		var row CollectionRemaining
		err := rows.Scan(&row.Collection, &row.Uuid, &row.Code, &row.Description, &row.Length, &row.Width, &row.Count)
		if err != nil {
			logs.ErrLog.Write(err)
			return nil, err
		}
		results = append(results, row)
	}

	return results, nil
}

func GetCollectionImage(r *http.Request) ([]Images, error) {
	db := getDbPG()
	params := utils.Params
	params.MU.RLock()
	urlDir := fmt.Sprintf("%s%s%s/files", params.HTTPServer, params.UrlPort, params.HttpRedirectPrefix)
	params.MU.RUnlock()

	queryText, ok := qMap.get(r.URL.Path)
	if !ok {
		queryByte, err := embedFiles.ReadFile("images.sql")
		if err != nil {
			return nil, err
		}
		queryText = string(queryByte)
		imageFieldsMapping(&queryText)
		qMap.add(r.URL.Path, queryText)
	}

	args := imgFilter.applyFilter(&queryText, r)

	rows, errRows := db.Query(queryText, args...)
	if errRows != nil {
		return nil, errRows
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logs.ErrLog.Write(err)
		}
	}()

	var results []Images
	for rows.Next() {
		var (
			row     = Images{}
			imageID string
		)

		err := rows.Scan(&imageID, &row.Collection, &row.Picture, &row.Form, &row.Color, &row.Brand)
		if err != nil {
			logs.ErrLog.Write(err)
			return nil, err
		}
		row.LogoURL = fmt.Sprintf("%s/%s/logo/%s.png", urlDir, imageID, imageID)
		results = append(results, row)
	}

	return results, nil
}

func GetDetailRemaining(r *http.Request) ([]DetailRemaining, error) {
	db := getDbPG()

	queryText, ok := qMap.get(r.URL.Path)
	if !ok {
		queryByte, err := embedFiles.ReadFile("details.sql")
		if err != nil {
			return nil, err
		}
		queryText = string(queryByte)
		detailsFieldsMapping(&queryText)
		qMap.add(r.URL.Path, queryText)
	}

	args := dtlFilter.applyFilter(&queryText, r)

	rows, errRows := db.Query(queryText, args...)
	if errRows != nil {
		return nil, errRows
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logs.ErrLog.Write(err)
		}
	}()

	var results []DetailRemaining

	for rows.Next() {
		var row DetailRemaining
		err := rows.Scan(&row.Collection, &row.Uuid, &row.Code, &row.Description, &row.Length, &row.Width,
			&row.UuidDetails, &row.CodeDetails, &row.Color, &row.Picture, &row.Form, &row.Barcode,
			&row.Brand, &row.CountBalance, &row.ReservedBalance, &row.Count)
		if err != nil {
			logs.ErrLog.Write(err)
			return nil, err
		}
		results = append(results, row)
	}

	return results, nil
}
