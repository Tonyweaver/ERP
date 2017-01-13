package models

import (
	"errors"
	"fmt"
	"strings"

	"goERP/utils"

	"github.com/astaxie/beego/orm"
)

// SaleConfig 销售设置
type SaleConfig struct {
	Base
	Name    string   `orm:"unique"`
	Company *Company `orm:"rel(fk)"` //公司
}

func init() {
	orm.RegisterModel(new(SaleConfig))
}

// AddSaleConfig insert a new SaleConfig into database and returns
// last inserted ID on success.
func AddSaleConfig(obj *SaleConfig) (id int64, err error) {
	o := orm.NewOrm()

	id, err = o.Insert(obj)
	return id, err
}

// GetSaleConfigByID retrieves SaleConfig by ID. Returns error if
// ID doesn't exist
func GetSaleConfigByID(id int64) (obj *SaleConfig, err error) {
	o := orm.NewOrm()
	obj = &SaleConfig{Base: Base{ID: id}}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetSaleConfigByName retrieves SaleOrder by Name. Returns error if
// Name doesn't exist
func GetSaleConfigByName(name string) (obj *SaleConfig, err error) {
	o := orm.NewOrm()
	obj = &SaleConfig{Name: name}
	if err = o.Read(obj); err == nil {
		return obj, nil
	}
	return nil, err
}

// GetAllSaleConfig retrieves all SaleConfig matches certain condition. Returns empty list if
// no records exist
func GetAllSaleConfig(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (utils.Paginator, []SaleConfig, error) {
	var (
		objArrs   []SaleConfig
		paginator utils.Paginator
		num       int64
		err       error
	)
	o := orm.NewOrm()
	qs := o.QueryTable(new(SaleConfig))
	qs = qs.RelatedSel()

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return paginator, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return paginator, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return paginator, nil, errors.New("Error: unused 'order' fields")
		}
	}

	qs = qs.OrderBy(sortFields...)
	if cnt, err := qs.Count(); err == nil {
		paginator = utils.GenPaginator(limit, offset, cnt)
	}
	if num, err = qs.Limit(limit, offset).All(&objArrs, fields...); err == nil {
		paginator.CurrentPageSize = num
	}
	return paginator, objArrs, err
}

// UpdateSaleConfigByID updates SaleConfig by ID and returns error if
// the record to be updated doesn't exist
func UpdateSaleConfigByID(m *SaleConfig) (err error) {
	o := orm.NewOrm()
	v := SaleConfig{Base: Base{ID: m.ID}}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSaleConfig deletes SaleConfig by ID and returns error if
// the record to be deleted doesn't exist
func DeleteSaleConfig(id int64) (err error) {
	o := orm.NewOrm()
	v := SaleConfig{Base: Base{ID: id}}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SaleConfig{Base: Base{ID: id}}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
