package mysql_tags

import (
	"fmt"
	"reflect"
	"strings"
  "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var placeholder_sep = "?,"

func SetPlaceHolderSeperator(sep string) {
	placeholder_sep = sep
}

func GetPlaceHolderSeperator() string {
	return placeholder_sep
}

var struct_tag = "sql"

func SetStructTag(tag string) {
	struct_tag = tag
}

func GetStructTag() string {
	return struct_tag
}

var db *sql.DB

func Setdb(pool *sql.DB) {
	db = pool
}

type GetFieldDesc struct {
	Numbers []int
	Columns	string
}

type SetFieldDesc struct {
	Numbers []int
	Columns	string
	Placeholders string
	Updates string
}

type TagQuery struct {
	Struct interface{}
	Value reflect.Value
	Table string

	SetFields *SetFieldDesc
	GetFields *GetFieldDesc

	qry_custom string
	qry_select string
	qry_insert string
}




/*Never forget to use a pointer as prm*/
func CreateTagQuery(v interface{}, ignore map[string]int, table string) *TagQuery {
	s := reflect.ValueOf(v).Elem()
	st := s.Type()

	var fe_sb, ft_sb, fu_sb strings.Builder
	fieldcounter := int(0)
	fieldnumbers := make([]int, st.NumField()-len(ignore))
	ignored_all := false
	ignored_cnt := int(0)
	if ignored_cnt == len(ignore) {
		ignored_all = true
	}
	for i := 0; i < st.NumField(); i++ {
		tag := st.Field(i).Tag.Get(struct_tag)
		if tag == "-" {
			fieldnumbers = fieldnumbers[:len(fieldnumbers)-1]
			continue
		}
		if !ignored_all {
			_, ignored := ignore[tag]
			if ignored {
				ignored_cnt += 1
				if ignored_cnt == len(ignore) {
					ignored_all = true
				}
				continue
			}
		}

		ft_sb.WriteString(tag)
		ft_sb.WriteString(",")
		fieldnumbers[fieldcounter] = i
		fieldcounter += 1
		fmt.Fprintf(&fe_sb, placeholder_sep)
		fmt.Fprintf(&fu_sb, tag+"="+placeholder_sep)
	}

	fieldtags := ft_sb.String()
	fieldenum := fe_sb.String()
	fieldupdate := fu_sb.String()
	fieldenum = fieldenum[:len(fieldenum)-1]
	fieldtags = fieldtags[:len(fieldtags)-1]
	fieldupdate = fieldupdate[:len(fieldupdate)-1]

	setfielddesc := &SetFieldDesc{Placeholders:fieldenum, Columns:fieldtags, Updates:fieldupdate, Numbers:fieldnumbers}
	getfielddesc := &GetFieldDesc{Columns:fieldtags, Numbers:fieldnumbers}

	return &TagQuery{Struct:s.Interface(), Value:s, Table:table, SetFields:setfielddesc, GetFields:getfielddesc}
}

/*Never forget to use a pointer as prm*/
func CreateTagQueryOfFields(v interface{}, fields map[string]int, table string) *TagQuery {
	s := reflect.ValueOf(v).Elem()
	st := s.Type()

	var fe_sb, ft_sb, fu_sb strings.Builder
	fieldcounter := int(0)
	fieldnumbers := make([]int, len(fields))

	for i := 0; i < st.NumField(); i++ {
		tag := st.Field(i).Tag.Get(struct_tag)
		if tag == "-" {
			continue
		}
		_, finding := fields[tag]
		if finding {
			ft_sb.WriteString(tag)
			ft_sb.WriteString(",")
			fieldnumbers[fieldcounter] = i
			fieldcounter += 1
			fmt.Fprintf(&fe_sb, placeholder_sep)
			fmt.Fprintf(&fu_sb, tag+"="+placeholder_sep)
			if fieldcounter == len(fields) {
				break
			}
		}
	}

	fieldtags := ft_sb.String()
	fieldenum := fe_sb.String()
	fieldupdate := fu_sb.String()
	fieldenum = fieldenum[:len(fieldenum)-1]
	fieldtags = fieldtags[:len(fieldtags)-1]
	fieldupdate = fieldupdate[:len(fieldupdate)-1]

	fielddesc := &SetFieldDesc{Placeholders:fieldenum, Columns:fieldtags, Updates:fieldupdate, Numbers:fieldnumbers}
	getfielddesc := &GetFieldDesc{Columns:fieldtags, Numbers:fieldnumbers}

	return &TagQuery{Struct:s.Interface(), Value:s, Table:table, SetFields:fielddesc, GetFields:getfielddesc}
}

/*Never forget to use a pointer as prm*/
func CreateTagQueryOfSetGetFields(v interface{}, setfields map[string]int, getfields map[string]int ,table string) *TagQuery {
	s := reflect.ValueOf(v).Elem()
	st := s.Type()

	var fe_sb, ft_sb, fu_sb, gft_sb strings.Builder
	fieldcounter := int(0)
	fieldnumbers := make([]int, len(setfields))
	gfieldcounter := int(0)
	gfieldnumbers := make([]int, len(getfields))

	for i := 0; i < st.NumField(); i++ {
		tag := st.Field(i).Tag.Get(struct_tag)
		if tag == "-" {
			continue
		}
		_, finding := setfields[tag]
		if finding {
			ft_sb.WriteString(tag)
			ft_sb.WriteString(",")
			fieldnumbers[fieldcounter] = i
			fieldcounter += 1
			fmt.Fprintf(&fe_sb, placeholder_sep)
			fmt.Fprintf(&fu_sb, tag+"="+placeholder_sep)
		}
		_, getfinding := getfields[tag]
		if getfinding {
			gft_sb.WriteString(tag)
			gft_sb.WriteString(",")
			gfieldnumbers[gfieldcounter] = i
			gfieldcounter += 1
		}
	}

	fieldtags := ft_sb.String()
	gfieldtags := gft_sb.String()
	fieldenum := fe_sb.String()
	fieldupdate := fu_sb.String()
	fieldenum = fieldenum[:len(fieldenum)-1]
	fieldtags = fieldtags[:len(fieldtags)-1]
	gfieldtags = gfieldtags[:len(gfieldtags)-1]
	fieldupdate = fieldupdate[:len(fieldupdate)-1]

	fielddesc := &SetFieldDesc{Placeholders:fieldenum, Columns:fieldtags, Updates:fieldupdate, Numbers:fieldnumbers}
	getfielddesc := &GetFieldDesc{Columns:gfieldtags, Numbers:gfieldnumbers}

	return &TagQuery{Struct:s.Interface(), Value:s, Table:table, SetFields:fielddesc, GetFields:getfielddesc}
}

/*Never forget to use a pointer as prm*/
func (tq *TagQuery) GetCopyWithStruct(v interface{}) *TagQuery {
	new_tq := tq
	s := reflect.ValueOf(v).Elem()
	new_tq.Struct = s.Interface()
	new_tq.Value = s
	return new_tq
}

func (tq *TagQuery) RebuildGetFields() {return}

func (tq *TagQuery) RebuildSetFields() {return}

func (tq *TagQuery) RebuildSetGetFields(v interface{}, setfields map[string]int, getfields map[string]int) {
	s := reflect.ValueOf(v).Elem()
	tq.Struct = s.Interface()
	tq.Value = s
	st := s.Type()
	var fe_sb, ft_sb, fu_sb, gft_sb strings.Builder
	fieldcounter := int(0)
	fieldnumbers := make([]int, len(setfields))
	gfieldcounter := int(0)
	gfieldnumbers := make([]int, len(getfields))

	for i := 0; i < st.NumField(); i++ {
		tag := st.Field(i).Tag.Get(struct_tag)
		if tag == "-" {
			continue
		}
		_, finding := setfields[tag]
		if finding {
			ft_sb.WriteString(tag)
			ft_sb.WriteString(",")
			fieldnumbers[fieldcounter] = i
			fieldcounter += 1
			fmt.Fprintf(&fe_sb, placeholder_sep)
			fmt.Fprintf(&fu_sb, tag+"="+placeholder_sep)
		}
		_, getfinding := getfields[tag]
		if getfinding {
			gft_sb.WriteString(tag)
			gft_sb.WriteString(",")
			gfieldnumbers[gfieldcounter] = i
			gfieldcounter += 1
		}
	}

	fieldtags := ft_sb.String()
	gfieldtags := gft_sb.String()
	fieldenum := fe_sb.String()
	fieldupdate := fu_sb.String()
	fieldenum = fieldenum[:len(fieldenum)-1]
	fieldtags = fieldtags[:len(fieldtags)-1]
	gfieldtags = gfieldtags[:len(gfieldtags)-1]
	fieldupdate = fieldupdate[:len(fieldupdate)-1]

	tq.SetFields = &SetFieldDesc{Placeholders:fieldenum, Columns:fieldtags, Updates:fieldupdate, Numbers:fieldnumbers}
	tq.GetFields = &GetFieldDesc{Columns:gfieldtags, Numbers:gfieldnumbers}
}

func (tq *TagQuery) GetReflectedMembersOf(v interface{}) []interface{} {
	s := reflect.ValueOf(v)
	mem_slice := make([]interface{}, len(tq.SetFields.Numbers))
	for index, key := range tq.SetFields.Numbers {
		mem_slice[index] = s.Field(key).Interface()
	}
	return mem_slice
}

func (tq *TagQuery) GetReflectedMembers() []interface{} {
	s := reflect.ValueOf(tq.Struct)
	mem_slice := make([]interface{}, len(tq.SetFields.Numbers))
	for index, key := range tq.SetFields.Numbers {
		mem_slice[index] = s.Field(key).Interface()
	}
	return mem_slice
}

/*Never forget to use this function with a pointer*/
func (tq *TagQuery) GetReflectedAddrOf(v interface{}) (reflect.Value, []interface{}) {
	s := reflect.ValueOf(v).Elem()
	addr_slice := make([]interface{}, len(tq.GetFields.Numbers))
	for index, key := range tq.GetFields.Numbers {
		field := s.Field(key)
		if field.CanAddr() {
			addr_slice[index] = field.Addr().Interface()
	  }
	}
	return s, addr_slice
}

func (tq *TagQuery) GetReflectedAddr() (reflect.Value, []interface{}) {
	s := tq.Value
	addr_slice := make([]interface{}, len(tq.GetFields.Numbers))
	for index, key := range tq.GetFields.Numbers {
		field := s.Field(key)
		if field.CanAddr() {
			addr_slice[index] = field.Addr().Interface()
	  }
	}
	return s, addr_slice
}

func (tq *TagQuery) formInsert() {
	if len(tq.qry_insert) == 0 {
		var qry_sb strings.Builder
		qry_sb.WriteString("insert into ")
		qry_sb.WriteString(tq.Table)
		qry_sb.WriteString("(")
		qry_sb.WriteString(tq.SetFields.Columns)
		qry_sb.WriteString(") values (")
		qry_sb.WriteString(tq.SetFields.Placeholders)
		qry_sb.WriteString(");")
		tq.qry_insert = qry_sb.String()
	}
}

func (tq *TagQuery) formInsertCustom(add string) {
	var qry_sb strings.Builder
	qry_sb.WriteString("insert into ")
	qry_sb.WriteString(tq.Table)
	qry_sb.WriteString("(")
	qry_sb.WriteString(tq.SetFields.Columns)
	qry_sb.WriteString(") values (")
	qry_sb.WriteString(tq.SetFields.Placeholders)
	qry_sb.WriteString(") ")
	qry_sb.WriteString(add+";")
	tq.qry_custom = qry_sb.String()
}

func (tq *TagQuery) FormInsertReturn() {
	if len(tq.qry_insert) == 0 {
		var qry_sb strings.Builder
		qry_sb.WriteString("insert into ")
		qry_sb.WriteString(tq.Table)
		qry_sb.WriteString("(")
		qry_sb.WriteString(tq.SetFields.Columns)
		qry_sb.WriteString(") values (")
		qry_sb.WriteString(tq.SetFields.Placeholders)
		qry_sb.WriteString(") returning ")
		qry_sb.WriteString(tq.GetFields.Columns)
		qry_sb.WriteString(";")
		tq.qry_insert = qry_sb.String()
	}
}

func (tq *TagQuery) Insert() error {
	tq.formInsert()
	stmtIns, err := db.Prepare(tq.qry_insert) // ? = placeholder
	if err != nil {
		return err
	}
  defer stmtIns.Close()
	_, err = stmtIns.Exec(tq.GetReflectedMembers()...)
	return err
}

/*Normalerweise umbauen InsertGetField*/
func (tq *TagQuery) InsertGetID() (uint64, error) {
	tq.formInsertCustom("returning id")
	var id uint64
	stmtOut, err := db.Prepare(tq.qry_custom)
	if err != nil {
		return 0, err
	}
  defer stmtOut.Close()
  err = stmtOut.QueryRow(tq.GetReflectedMembers()...).Scan(&id)
	return id, err
}

func (tq *TagQuery) InsertGetFields() (interface{}, error) {
	tq.FormInsertReturn()
	i, addr_slice := tq.GetReflectedAddr()

	stmtOut, err := db.Prepare(tq.qry_insert)
	if err != nil {
		return nil, err
	}
  defer stmtOut.Close()
  err = stmtOut.QueryRow(tq.GetReflectedMembers()...).Scan(addr_slice...)

	return i.Interface(), err
}

func (tq *TagQuery) formSelect() {
	if len(tq.qry_select) == 0 {
		var qry_sb strings.Builder
		qry_sb.WriteString("select ")
		qry_sb.WriteString(tq.GetFields.Columns)
		qry_sb.WriteString(" from ")
		qry_sb.WriteString(tq.Table+";")
		tq.qry_select = qry_sb.String()
	}
}

func (tq *TagQuery) formSelectCustom(add string) {
	var qry_sb strings.Builder
	qry_sb.WriteString("select ")
	qry_sb.WriteString(tq.GetFields.Columns)
	qry_sb.WriteString(" from ")
	qry_sb.WriteString(tq.Table)
	qry_sb.WriteString(" ")
	qry_sb.WriteString(add+";")
	tq.qry_custom = qry_sb.String()
}

func (tq *TagQuery) formSelectAll() {
	var qry_sb strings.Builder
	qry_sb.WriteString("select * from ")
	qry_sb.WriteString(tq.Table+";")
	tq.qry_custom = qry_sb.String()
}

func (tq *TagQuery) Select() ([]interface{}, error) {
	tq.formSelect()
	abp := []interface{}{}

	stmtOut, err := db.Prepare(tq.qry_select) // ? = placeholder
	if err != nil {
		return nil, err
	}
  defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		return abp, err
	}
	defer rows.Close()
	for rows.Next() {
		i, addr_slice := tq.GetReflectedAddr()
		err = rows.Scan(addr_slice...)
		if err != nil {
			return abp, err
		}
		abp = append(abp, i.Interface())
	}
	err = rows.Err()
	if err != nil {
		return abp, err
	}
	return abp, nil
}

func (tq *TagQuery) SelectByID(id uint64) (interface{}, error) {
	tq.formSelectCustom("where id=$1")
	i, addr_slice := tq.GetReflectedAddr()

	stmtOut, err := db.Prepare(tq.qry_custom)
	if err != nil {
		return nil, err
	}
  defer stmtOut.Close()
  err = stmtOut.QueryRow(id).Scan(addr_slice...)

	return i.Interface(), err
}

func (tq *TagQuery) SelectCustom(custom string, args ...interface{}) ([]interface{}, error) {
	tq.formSelectCustom(custom)
	return tq.SelectCommon(tq.qry_custom, args)
}

/*produces error, when QueryTag has ignored fields*/
func (tq *TagQuery) SelectAll() ([]interface{}, error) {
	tq.formSelectAll()
	return tq.SelectCommon(tq.qry_custom)
}

func (tq *TagQuery) SelectCommon(qry string, args ...interface{}) ([]interface{}, error) {
	abp := []interface{}{}

	stmtOut, err := db.Prepare(qry)
	if err != nil {
		return nil, err
	}
  defer stmtOut.Close()

	rows, err := stmtOut.Query(args)
	if err != nil {
		return abp, err
	}
	defer rows.Close()
	for rows.Next() {
		i, addr_slice := tq.GetReflectedAddr()
		err = rows.Scan(addr_slice...)
		if err != nil {
			return abp, err
		}
		abp = append(abp, i.Interface())
	}
	err = rows.Err()
	if err != nil {
		return abp, err
	}
	return abp, nil
}

func (tq *TagQuery) Update(where string) error {
	var cnt_sb strings.Builder
	cnt_sb.WriteString("update ")
	cnt_sb.WriteString(tq.Table)
	cnt_sb.WriteString(" set ")
	cnt_sb.WriteString(tq.SetFields.Updates)
	if len(where) > 0 {
		cnt_sb.WriteString(where)
	}
	cnt_sb.WriteString(";")

	stmtOut, err := db.Prepare(cnt_sb.String())
	if err != nil {
		return err
	}
  defer stmtOut.Close()

	_, err = stmtOut.Exec(tq.GetReflectedMembers()...)
	return err
}

func (tq *TagQuery) UpdateFieldWith(field string, where string, args ...interface{},) error {
	var cnt_sb strings.Builder
	cnt_sb.WriteString("update ")
	cnt_sb.WriteString(tq.Table)
	cnt_sb.WriteString(" set ")
	cnt_sb.WriteString(field)
	cnt_sb.WriteString("=$1 ")
	if len(where) > 0 {
		cnt_sb.WriteString(where)
	}
	cnt_sb.WriteString(";")

	stmtOut, err := db.Prepare(cnt_sb.String())
	if err != nil {
		return err
	}
  defer stmtOut.Close()

	_, err = stmtOut.Exec(args...)
	return err
}

func (tq *TagQuery) Count(where string) (uint64, error) {
	var cnt_sb strings.Builder
	cnt_sb.WriteString("select count(*) from ")
	cnt_sb.WriteString(tq.Table)
	if len(where) > 0 {
		cnt_sb.WriteString(where)
	}
	cnt_sb.WriteString(";")
	var count uint64

	stmtOut, err := db.Prepare(cnt_sb.String())
	if err != nil {
		return 0, err
	}
  defer stmtOut.Close()

	err = stmtOut.QueryRow().Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (tq *TagQuery) FieldExists(field string, arg interface{}) (bool, error) {
	var cnt_sb strings.Builder
	cnt_sb.WriteString("select exists(select 1 from ")
	cnt_sb.WriteString(tq.Table)
	cnt_sb.WriteString(" where ")
	cnt_sb.WriteString(field)
	cnt_sb.WriteString("=$1) as \"exists\";")
	var exist bool

	stmtOut, err := db.Prepare(cnt_sb.String())
	if err != nil {
		return false, err
	}
  defer stmtOut.Close()

	err = stmtOut.QueryRow(arg).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}
